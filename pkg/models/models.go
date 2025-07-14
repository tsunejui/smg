package models

import (
	"time"
)

type User struct {
	ID            string    `json:"id" db:"id"`
	Name          *string   `json:"name" db:"name"`
	Email         string    `json:"email" db:"email"`
	EmailVerified *time.Time `json:"email_verified" db:"email_verified"`
	Password      *string   `json:"-" db:"password"`
	Image         *string   `json:"image" db:"image"`
	IsAdmin       bool      `json:"is_admin" db:"is_admin"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type Topic struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description" db:"description"`
	Keywords    []string  `json:"keywords" db:"keywords"`
	Platforms   []string  `json:"platforms" db:"platforms"`
	UserID      string    `json:"user_id" db:"user_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type MediaAccount struct {
	ID           string     `json:"id" db:"id"`
	Platform     string     `json:"platform" db:"platform"`
	AccountID    string     `json:"account_id" db:"account_id"`
	AccountName  string     `json:"account_name" db:"account_name"`
	AccessToken  *string    `json:"-" db:"access_token"`
	RefreshToken *string    `json:"-" db:"refresh_token"`
	ExpiresAt    *time.Time `json:"expires_at" db:"expires_at"`
	UserID       string     `json:"user_id" db:"user_id"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

type Article struct {
	ID          string    `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Content     string    `json:"content" db:"content"`
	OriginalURL string    `json:"original_url" db:"original_url"`
	Platform    string    `json:"platform" db:"platform"`
	AuthorName  *string   `json:"author_name" db:"author_name"`
	AuthorID    *string   `json:"author_id" db:"author_id"`
	PublishedAt time.Time `json:"published_at" db:"published_at"`
	TopicID     string    `json:"topic_id" db:"topic_id"`
	UserID      string    `json:"user_id" db:"user_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Repost struct {
	ID             string     `json:"id" db:"id"`
	ArticleID      string     `json:"article_id" db:"article_id"`
	MediaAccountID string     `json:"media_account_id" db:"media_account_id"`
	CustomCaption  *string    `json:"custom_caption" db:"custom_caption"`
	AICaption      *string    `json:"ai_caption" db:"ai_caption"`
	Status         string     `json:"status" db:"status"`
	ScheduledAt    *time.Time `json:"scheduled_at" db:"scheduled_at"`
	PostedAt       *time.Time `json:"posted_at" db:"posted_at"`
	ExternalID     *string    `json:"external_id" db:"external_id"`
	UserID         string     `json:"user_id" db:"user_id"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

type SystemSetting struct {
	ID        string    `json:"id" db:"id"`
	Key       string    `json:"key" db:"key"`
	Value     string    `json:"value" db:"value"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Platform struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Enabled     bool      `json:"enabled"`
	Config      string    `json:"config"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Request/Response models
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type QRCodeRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

type QRCodeResponse struct {
	Token   string `json:"token"`
	QRCode  string `json:"qr_code"`
	Expires int64  `json:"expires"`
}

type AuthResponse struct {
	User         User   `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type UpdateProfileRequest struct {
	Name  *string `json:"name"`
	Image *string `json:"image"`
}

type CreateTopicRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description *string  `json:"description"`
	Keywords    []string `json:"keywords" binding:"required"`
	Platforms   []string `json:"platforms" binding:"required"`
}

type ConnectPlatformRequest struct {
	Code         string `json:"code" binding:"required"`
	RedirectURI  string `json:"redirect_uri"`
	AccountName  string `json:"account_name"`
}

type RepostRequest struct {
	MediaAccountID string  `json:"media_account_id" binding:"required"`
	CustomCaption  *string `json:"custom_caption"`
	ScheduledAt    *time.Time `json:"scheduled_at"`
}

type StatsResponse struct {
	TotalUsers     int64 `json:"total_users"`
	ActiveUsers    int64 `json:"active_users"`
	TotalTopics    int64 `json:"total_topics"`
	TotalArticles  int64 `json:"total_articles"`
	TotalReposts   int64 `json:"total_reposts"`
	RepostsToday   int64 `json:"reposts_today"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}