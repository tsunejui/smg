package services

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"smg/pkg/models"
)

type ArticleService struct {
	db *sql.DB
}

func NewArticleService(db *sql.DB) *ArticleService {
	return &ArticleService{db: db}
}

func (s *ArticleService) GetArticles(userID string, page, pageSize int) (*models.PaginatedResponse, error) {
	offset := (page - 1) * pageSize
	
	var total int64
	err := s.db.QueryRow("SELECT COUNT(*) FROM articles WHERE user_id = $1", userID).Scan(&total)
	if err != nil {
		return nil, err
	}
	
	rows, err := s.db.Query(`
		SELECT id, title, content, original_url, platform, author_name, author_id, 
			   published_at, topic_id, user_id, created_at, updated_at
		FROM articles 
		WHERE user_id = $1
		ORDER BY published_at DESC
		LIMIT $2 OFFSET $3
	`, userID, pageSize, offset)
	
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

func (s *ArticleService) CreateArticle(userID string, article *models.Article) (*models.Article, error) {
	articleID := uuid.New().String()
	now := time.Now()
	
	_, err := s.db.Exec(`
		INSERT INTO articles (id, title, content, original_url, platform, author_name, author_id, 
							published_at, topic_id, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`, articleID, article.Title, article.Content, article.OriginalURL, article.Platform,
		article.AuthorName, article.AuthorID, article.PublishedAt, article.TopicID, 
		userID, now, now)
	
	if err != nil {
		return nil, err
	}
	
	return s.GetArticle(articleID)
}

func (s *ArticleService) GetArticle(articleID string) (*models.Article, error) {
	var article models.Article
	err := s.db.QueryRow(`
		SELECT id, title, content, original_url, platform, author_name, author_id, 
			   published_at, topic_id, user_id, created_at, updated_at
		FROM articles WHERE id = $1
	`, articleID).Scan(
		&article.ID, &article.Title, &article.Content, &article.OriginalURL,
		&article.Platform, &article.AuthorName, &article.AuthorID,
		&article.PublishedAt, &article.TopicID, &article.UserID,
		&article.CreatedAt, &article.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &article, nil
}

func (s *ArticleService) UpdateArticle(articleID string, article *models.Article) (*models.Article, error) {
	now := time.Now()
	
	_, err := s.db.Exec(`
		UPDATE articles 
		SET title = $2, content = $3, updated_at = $4
		WHERE id = $1
	`, articleID, article.Title, article.Content, now)
	
	if err != nil {
		return nil, err
	}
	
	return s.GetArticle(articleID)
}

func (s *ArticleService) DeleteArticle(articleID string) error {
	_, err := s.db.Exec("DELETE FROM articles WHERE id = $1", articleID)
	return err
}

func (s *ArticleService) RepostArticle(articleID, userID string, req *models.RepostRequest) (*models.Repost, error) {
	repostID := uuid.New().String()
	now := time.Now()
	
	_, err := s.db.Exec(`
		INSERT INTO reposts (id, article_id, media_account_id, custom_caption, status, 
						   scheduled_at, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, repostID, articleID, req.MediaAccountID, req.CustomCaption, "pending", 
		req.ScheduledAt, userID, now, now)
	
	if err != nil {
		return nil, err
	}
	
	return s.GetRepost(repostID)
}

func (s *ArticleService) GetRepost(repostID string) (*models.Repost, error) {
	var repost models.Repost
	err := s.db.QueryRow(`
		SELECT id, article_id, media_account_id, custom_caption, ai_caption, status, 
			   scheduled_at, posted_at, external_id, user_id, created_at, updated_at
		FROM reposts WHERE id = $1
	`, repostID).Scan(
		&repost.ID, &repost.ArticleID, &repost.MediaAccountID, &repost.CustomCaption,
		&repost.AICaption, &repost.Status, &repost.ScheduledAt, &repost.PostedAt,
		&repost.ExternalID, &repost.UserID, &repost.CreatedAt, &repost.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &repost, nil
}

func (s *ArticleService) GetReposts(userID string, page, pageSize int) (*models.PaginatedResponse, error) {
	offset := (page - 1) * pageSize
	
	var total int64
	err := s.db.QueryRow("SELECT COUNT(*) FROM reposts WHERE user_id = $1", userID).Scan(&total)
	if err != nil {
		return nil, err
	}
	
	rows, err := s.db.Query(`
		SELECT id, article_id, media_account_id, custom_caption, ai_caption, status, 
			   scheduled_at, posted_at, external_id, user_id, created_at, updated_at
		FROM reposts 
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, pageSize, offset)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var reposts []models.Repost
	for rows.Next() {
		var repost models.Repost
		err := rows.Scan(
			&repost.ID, &repost.ArticleID, &repost.MediaAccountID, &repost.CustomCaption,
			&repost.AICaption, &repost.Status, &repost.ScheduledAt, &repost.PostedAt,
			&repost.ExternalID, &repost.UserID, &repost.CreatedAt, &repost.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		reposts = append(reposts, repost)
	}
	
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	
	return &models.PaginatedResponse{
		Data:       reposts,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}