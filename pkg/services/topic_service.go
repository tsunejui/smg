package services

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"smg/pkg/models"
)

type TopicService struct {
	db *sql.DB
}

func NewTopicService(db *sql.DB) *TopicService {
	return &TopicService{db: db}
}

func (s *TopicService) GetTopics(userID string, page, pageSize int) (*models.PaginatedResponse, error) {
	offset := (page - 1) * pageSize
	
	var total int64
	err := s.db.QueryRow("SELECT COUNT(*) FROM topics WHERE user_id = $1", userID).Scan(&total)
	if err != nil {
		return nil, err
	}
	
	rows, err := s.db.Query(`
		SELECT id, name, description, keywords, platforms, user_id, created_at, updated_at
		FROM topics 
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, pageSize, offset)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var topics []models.Topic
	for rows.Next() {
		var topic models.Topic
		err := rows.Scan(
			&topic.ID, &topic.Name, &topic.Description, 
			pq.Array(&topic.Keywords), pq.Array(&topic.Platforms),
			&topic.UserID, &topic.CreatedAt, &topic.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		topics = append(topics, topic)
	}
	
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	
	return &models.PaginatedResponse{
		Data:       topics,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (s *TopicService) CreateTopic(userID string, req *models.CreateTopicRequest) (*models.Topic, error) {
	topicID := uuid.New().String()
	now := time.Now()
	
	_, err := s.db.Exec(`
		INSERT INTO topics (id, name, description, keywords, platforms, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, topicID, req.Name, req.Description, pq.Array(req.Keywords), pq.Array(req.Platforms), userID, now, now)
	
	if err != nil {
		return nil, err
	}
	
	return s.GetTopic(topicID)
}

func (s *TopicService) GetTopic(topicID string) (*models.Topic, error) {
	var topic models.Topic
	err := s.db.QueryRow(`
		SELECT id, name, description, keywords, platforms, user_id, created_at, updated_at
		FROM topics WHERE id = $1
	`, topicID).Scan(
		&topic.ID, &topic.Name, &topic.Description, 
		pq.Array(&topic.Keywords), pq.Array(&topic.Platforms),
		&topic.UserID, &topic.CreatedAt, &topic.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &topic, nil
}

func (s *TopicService) UpdateTopic(topicID string, req *models.CreateTopicRequest) (*models.Topic, error) {
	now := time.Now()
	
	_, err := s.db.Exec(`
		UPDATE topics 
		SET name = $2, description = $3, keywords = $4, platforms = $5, updated_at = $6
		WHERE id = $1
	`, topicID, req.Name, req.Description, pq.Array(req.Keywords), pq.Array(req.Platforms), now)
	
	if err != nil {
		return nil, err
	}
	
	return s.GetTopic(topicID)
}

func (s *TopicService) DeleteTopic(topicID string) error {
	_, err := s.db.Exec("DELETE FROM topics WHERE id = $1", topicID)
	return err
}

func (s *TopicService) GetTopicArticles(topicID string, page, pageSize int) (*models.PaginatedResponse, error) {
	offset := (page - 1) * pageSize
	
	var total int64
	err := s.db.QueryRow("SELECT COUNT(*) FROM articles WHERE topic_id = $1", topicID).Scan(&total)
	if err != nil {
		return nil, err
	}
	
	rows, err := s.db.Query(`
		SELECT id, title, content, original_url, platform, author_name, author_id, 
			   published_at, topic_id, user_id, created_at, updated_at
		FROM articles 
		WHERE topic_id = $1
		ORDER BY published_at DESC
		LIMIT $2 OFFSET $3
	`, topicID, pageSize, offset)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var articles []models.Article
	for rows.Next() {
		var article models.Article
		err := rows.Scan(
			&article.ID, &article.Title, &article.Content, &article.OriginalURL,
			&article.Platform, &article.AuthorName, &article.AuthorID,
			&article.PublishedAt, &article.TopicID, &article.UserID,
			&article.CreatedAt, &article.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	
	return &models.PaginatedResponse{
		Data:       articles,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}