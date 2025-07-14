package services

import (
	"database/sql"
	"fmt"
	"time"

	"smg/pkg/models"
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetProfile(userID string) (*models.User, error) {
	var user models.User
	err := s.db.QueryRow(`
		SELECT id, name, email, email_verified, image, is_admin, created_at, updated_at
		FROM users WHERE id = $1
	`, userID).Scan(
		&user.ID, &user.Name, &user.Email, &user.EmailVerified, 
		&user.Image, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

func (s *UserService) UpdateProfile(userID string, req *models.UpdateProfileRequest) (*models.User, error) {
	now := time.Now()
	
	_, err := s.db.Exec(`
		UPDATE users 
		SET name = COALESCE($2, name), 
			image = COALESCE($3, image),
			updated_at = $4
		WHERE id = $1
	`, userID, req.Name, req.Image, now)
	
	if err != nil {
		return nil, err
	}
	
	return s.GetProfile(userID)
}

func (s *UserService) GetUsers(page, pageSize int) (*models.PaginatedResponse, error) {
	offset := (page - 1) * pageSize
	
	var total int64
	err := s.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&total)
	if err != nil {
		return nil, err
	}
	
	rows, err := s.db.Query(`
		SELECT id, name, email, email_verified, image, is_admin, created_at, updated_at
		FROM users 
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, pageSize, offset)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID, &user.Name, &user.Email, &user.EmailVerified, 
			&user.Image, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	
	return &models.PaginatedResponse{
		Data:       users,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (s *UserService) GetUserByID(userID string) (*models.User, error) {
	return s.GetProfile(userID)
}

func (s *UserService) UpdateUser(userID string, req *models.UpdateProfileRequest) (*models.User, error) {
	return s.UpdateProfile(userID, req)
}

func (s *UserService) DeleteUser(userID string) error {
	_, err := s.db.Exec("DELETE FROM users WHERE id = $1", userID)
	return err
}

func (s *UserService) GetUserStats(userID string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	var topicsCount, mediaAccountsCount, articlesCount, repostsCount int
	
	err := s.db.QueryRow("SELECT COUNT(*) FROM topics WHERE user_id = $1", userID).Scan(&topicsCount)
	if err != nil {
		return nil, err
	}
	
	err = s.db.QueryRow("SELECT COUNT(*) FROM media_accounts WHERE user_id = $1", userID).Scan(&mediaAccountsCount)
	if err != nil {
		return nil, err
	}
	
	err = s.db.QueryRow("SELECT COUNT(*) FROM articles WHERE user_id = $1", userID).Scan(&articlesCount)
	if err != nil {
		return nil, err
	}
	
	err = s.db.QueryRow("SELECT COUNT(*) FROM reposts WHERE user_id = $1", userID).Scan(&repostsCount)
	if err != nil {
		return nil, err
	}
	
	stats["topics_count"] = topicsCount
	stats["media_accounts_count"] = mediaAccountsCount
	stats["articles_count"] = articlesCount
	stats["reposts_count"] = repostsCount
	
	return stats, nil
}