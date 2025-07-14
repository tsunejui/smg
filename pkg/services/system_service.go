package services

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"smg/pkg/models"
)

type SystemService struct {
	db *sql.DB
}

func NewSystemService(db *sql.DB) *SystemService {
	return &SystemService{db: db}
}

func (s *SystemService) GetSettings() (map[string]string, error) {
	rows, err := s.db.Query("SELECT key, value FROM system_settings")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	settings := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		settings[key] = value
	}
	
	return settings, nil
}

func (s *SystemService) UpdateSettings(settings map[string]string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	for key, value := range settings {
		_, err := tx.Exec(`
			INSERT INTO system_settings (id, key, value, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (key) DO UPDATE SET
			value = EXCLUDED.value,
			updated_at = EXCLUDED.updated_at
		`, uuid.New().String(), key, value, time.Now(), time.Now())
		
		if err != nil {
			return err
		}
	}
	
	return tx.Commit()
}

func (s *SystemService) GetStats() (*models.StatsResponse, error) {
	var stats models.StatsResponse
	
	err := s.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&stats.TotalUsers)
	if err != nil {
		return nil, err
	}
	
	err = s.db.QueryRow(`
		SELECT COUNT(*) FROM users 
		WHERE created_at > NOW() - INTERVAL '30 days'
	`).Scan(&stats.ActiveUsers)
	if err != nil {
		return nil, err
	}
	
	err = s.db.QueryRow("SELECT COUNT(*) FROM topics").Scan(&stats.TotalTopics)
	if err != nil {
		return nil, err
	}
	
	err = s.db.QueryRow("SELECT COUNT(*) FROM articles").Scan(&stats.TotalArticles)
	if err != nil {
		return nil, err
	}
	
	err = s.db.QueryRow("SELECT COUNT(*) FROM reposts").Scan(&stats.TotalReposts)
	if err != nil {
		return nil, err
	}
	
	err = s.db.QueryRow(`
		SELECT COUNT(*) FROM reposts 
		WHERE created_at::date = CURRENT_DATE
	`).Scan(&stats.RepostsToday)
	if err != nil {
		return nil, err
	}
	
	return &stats, nil
}

func (s *SystemService) GetPlatforms() ([]models.Platform, error) {
	rows, err := s.db.Query(`
		SELECT id, name, display_name, enabled, config, created_at, updated_at
		FROM platforms
		ORDER BY display_name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var platforms []models.Platform
	for rows.Next() {
		var platform models.Platform
		err := rows.Scan(
			&platform.ID, &platform.Name, &platform.DisplayName,
			&platform.Enabled, &platform.Config, &platform.CreatedAt, &platform.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		platforms = append(platforms, platform)
	}
	
	return platforms, nil
}

func (s *SystemService) CreatePlatform(platform *models.Platform) (*models.Platform, error) {
	platformID := uuid.New().String()
	now := time.Now()
	
	_, err := s.db.Exec(`
		INSERT INTO platforms (id, name, display_name, enabled, config, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, platformID, platform.Name, platform.DisplayName, platform.Enabled, 
		platform.Config, now, now)
	
	if err != nil {
		return nil, err
	}
	
	return s.GetPlatform(platformID)
}

func (s *SystemService) GetPlatform(platformID string) (*models.Platform, error) {
	var platform models.Platform
	err := s.db.QueryRow(`
		SELECT id, name, display_name, enabled, config, created_at, updated_at
		FROM platforms WHERE id = $1
	`, platformID).Scan(
		&platform.ID, &platform.Name, &platform.DisplayName,
		&platform.Enabled, &platform.Config, &platform.CreatedAt, &platform.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &platform, nil
}

func (s *SystemService) UpdatePlatform(platformID string, platform *models.Platform) (*models.Platform, error) {
	now := time.Now()
	
	_, err := s.db.Exec(`
		UPDATE platforms 
		SET name = $2, display_name = $3, enabled = $4, config = $5, updated_at = $6
		WHERE id = $1
	`, platformID, platform.Name, platform.DisplayName, platform.Enabled, 
		platform.Config, now)
	
	if err != nil {
		return nil, err
	}
	
	return s.GetPlatform(platformID)
}

func (s *SystemService) DeletePlatform(platformID string) error {
	_, err := s.db.Exec("DELETE FROM platforms WHERE id = $1", platformID)
	return err
}