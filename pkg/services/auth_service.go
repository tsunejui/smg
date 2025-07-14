package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"smg/pkg/models"
)

type AuthService struct {
	db          *sql.DB
	redisClient *redis.Client
	jwtSecret   string
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	IsAdmin bool  `json:"is_admin"`
	jwt.RegisteredClaims
}

func NewAuthService(db *sql.DB, redisClient *redis.Client) *AuthService {
	return &AuthService{
		db:          db,
		redisClient: redisClient,
		jwtSecret:   "your-secret-key", // Should be from config
	}
}

func (s *AuthService) Login(email, password string) (*models.AuthResponse, error) {
	var user models.User
	var hashedPassword string
	
	err := s.db.QueryRow(`
		SELECT id, name, email, email_verified, password, image, is_admin, created_at, updated_at
		FROM users WHERE email = $1
	`, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.EmailVerified, 
		&hashedPassword, &user.Image, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid credentials")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	accessToken, err := s.generateToken(user.ID, user.Email, user.IsAdmin, time.Hour*24)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateToken(user.ID, user.Email, user.IsAdmin, time.Hour*24*7)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    86400, // 24 hours
	}, nil
}

func (s *AuthService) Register(name, email, password string) (*models.AuthResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	userID := uuid.New().String()
	now := time.Now()
	
	_, err = s.db.Exec(`
		INSERT INTO users (id, name, email, password, is_admin, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, userID, name, email, string(hashedPassword), false, now, now)
	
	if err != nil {
		return nil, err
	}

	user := models.User{
		ID:        userID,
		Name:      &name,
		Email:     email,
		IsAdmin:   false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	accessToken, err := s.generateToken(user.ID, user.Email, user.IsAdmin, time.Hour*24)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateToken(user.ID, user.Email, user.IsAdmin, time.Hour*24*7)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    86400,
	}, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*models.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	var user models.User
	err = s.db.QueryRow(`
		SELECT id, name, email, email_verified, image, is_admin, created_at, updated_at
		FROM users WHERE id = $1
	`, claims.UserID).Scan(
		&user.ID, &user.Name, &user.Email, &user.EmailVerified, 
		&user.Image, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *AuthService) GenerateQRCode(userID string) (*models.QRCodeResponse, error) {
	token := uuid.New().String()
	expires := time.Now().Add(time.Minute * 5).Unix()
	
	qrData := map[string]interface{}{
		"token":   token,
		"user_id": userID,
		"expires": expires,
	}
	
	qrJSON, err := json.Marshal(qrData)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	err = s.redisClient.Set(ctx, fmt.Sprintf("qr:%s", token), userID, time.Minute*5).Err()
	if err != nil {
		return nil, err
	}

	return &models.QRCodeResponse{
		Token:   token,
		QRCode:  string(qrJSON),
		Expires: expires,
	}, nil
}

func (s *AuthService) VerifyQRCode(token string) (*models.AuthResponse, error) {
	ctx := context.Background()
	userID, err := s.redisClient.Get(ctx, fmt.Sprintf("qr:%s", token)).Result()
	if err != nil {
		return nil, fmt.Errorf("invalid or expired QR code")
	}

	var user models.User
	err = s.db.QueryRow(`
		SELECT id, name, email, email_verified, image, is_admin, created_at, updated_at
		FROM users WHERE id = $1
	`, userID).Scan(
		&user.ID, &user.Name, &user.Email, &user.EmailVerified, 
		&user.Image, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}

	accessToken, err := s.generateToken(user.ID, user.Email, user.IsAdmin, time.Hour*24)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateToken(user.ID, user.Email, user.IsAdmin, time.Hour*24*7)
	if err != nil {
		return nil, err
	}

	// Delete the QR code token after use
	s.redisClient.Del(ctx, fmt.Sprintf("qr:%s", token))

	return &models.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    86400,
	}, nil
}

func (s *AuthService) generateToken(userID, email string, isAdmin bool, duration time.Duration) (string, error) {
	claims := &Claims{
		UserID:  userID,
		Email:   email,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}