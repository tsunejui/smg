package services

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"smg/pkg/models"
)

type MediaService struct {
	db *sql.DB
}

func NewMediaService(db *sql.DB) *MediaService {
	return &MediaService{db: db}
}

func (s *MediaService) GetAccounts(userID string) ([]models.MediaAccount, error) {
	rows, err := s.db.Query(`
		SELECT id, platform, account_id, account_name, expires_at, user_id, created_at, updated_at
		FROM media_accounts 
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var accounts []models.MediaAccount
	for rows.Next() {
		var account models.MediaAccount
		err := rows.Scan(
			&account.ID, &account.Platform, &account.AccountID, &account.AccountName,
			&account.ExpiresAt, &account.UserID, &account.CreatedAt, &account.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	
	return accounts, nil
}

func (s *MediaService) CreateAccount(userID string, req *models.ConnectPlatformRequest) (*models.MediaAccount, error) {
	accountID := uuid.New().String()
	now := time.Now()
	
	_, err := s.db.Exec(`
		INSERT INTO media_accounts (id, platform, account_id, account_name, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, accountID, "twitter", req.Code, req.AccountName, userID, now, now)
	
	if err != nil {
		return nil, err
	}
	
	return s.GetAccount(accountID)
}

func (s *MediaService) GetAccount(accountID string) (*models.MediaAccount, error) {
	var account models.MediaAccount
	err := s.db.QueryRow(`
		SELECT id, platform, account_id, account_name, expires_at, user_id, created_at, updated_at
		FROM media_accounts WHERE id = $1
	`, accountID).Scan(
		&account.ID, &account.Platform, &account.AccountID, &account.AccountName,
		&account.ExpiresAt, &account.UserID, &account.CreatedAt, &account.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &account, nil
}

func (s *MediaService) UpdateAccount(accountID string, req *models.ConnectPlatformRequest) (*models.MediaAccount, error) {
	now := time.Now()
	
	_, err := s.db.Exec(`
		UPDATE media_accounts 
		SET account_name = $2, updated_at = $3
		WHERE id = $1
	`, accountID, req.AccountName, now)
	
	if err != nil {
		return nil, err
	}
	
	return s.GetAccount(accountID)
}

func (s *MediaService) DeleteAccount(accountID string) error {
	_, err := s.db.Exec("DELETE FROM media_accounts WHERE id = $1", accountID)
	return err
}

func (s *MediaService) ConnectPlatform(userID, platform string, req *models.ConnectPlatformRequest) (*models.MediaAccount, error) {
	accountID := uuid.New().String()
	now := time.Now()
	
	_, err := s.db.Exec(`
		INSERT INTO media_accounts (id, platform, account_id, account_name, access_token, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, accountID, platform, req.Code, req.AccountName, &req.Code, userID, now, now)
	
	if err != nil {
		return nil, err
	}
	
	return s.GetAccount(accountID)
}

func (s *MediaService) DisconnectAccount(accountID string) error {
	return s.DeleteAccount(accountID)
}

func (s *MediaService) GetPlatformAccounts(userID, platform string) ([]models.MediaAccount, error) {
	rows, err := s.db.Query(`
		SELECT id, platform, account_id, account_name, expires_at, user_id, created_at, updated_at
		FROM media_accounts 
		WHERE user_id = $1 AND platform = $2
		ORDER BY created_at DESC
	`, userID, platform)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var accounts []models.MediaAccount
	for rows.Next() {
		var account models.MediaAccount
		err := rows.Scan(
			&account.ID, &account.Platform, &account.AccountID, &account.AccountName,
			&account.ExpiresAt, &account.UserID, &account.CreatedAt, &account.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	
	return accounts, nil
}