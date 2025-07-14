package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	db   *sql.DB
	cron *cron.Cron
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Connect to database
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// Create scheduler
	scheduler := &Scheduler{
		db:   db,
		cron: cron.New(cron.WithSeconds()),
	}

	// Schedule jobs
	scheduler.scheduleJobs()

	// Start scheduler
	scheduler.cron.Start()
	defer scheduler.cron.Stop()

	log.Println("Scheduler started successfully")

	// Keep the program running
	select {}
}

func (s *Scheduler) scheduleJobs() {
	// Process scheduled reposts every minute
	s.cron.AddFunc("0 * * * * *", s.processScheduledReposts)

	// Fetch articles every 10 minutes
	s.cron.AddFunc("0 */10 * * * *", s.fetchArticles)

	// Generate AI captions every 5 minutes
	s.cron.AddFunc("0 */5 * * * *", s.generateAICaptions)

	// Cleanup old data daily at 2 AM
	s.cron.AddFunc("0 0 2 * * *", s.cleanupOldData)

	log.Println("Scheduled jobs:")
	log.Println("- Process scheduled reposts: every minute")
	log.Println("- Fetch articles: every 10 minutes")
	log.Println("- Generate AI captions: every 5 minutes")
	log.Println("- Cleanup old data: daily at 2 AM")
}

func (s *Scheduler) processScheduledReposts() {
	log.Println("Processing scheduled reposts...")

	// Get reposts that are scheduled for now or earlier
	rows, err := s.db.Query(`
		SELECT id, article_id, media_account_id, custom_caption, ai_caption, user_id
		FROM reposts 
		WHERE status = 'pending' 
		AND scheduled_at <= NOW()
		LIMIT 10
	`)
	if err != nil {
		log.Printf("Error fetching scheduled reposts: %v", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var repostID, articleID, mediaAccountID, userID string
		var customCaption, aiCaption *string

		if err := rows.Scan(&repostID, &articleID, &mediaAccountID, &customCaption, &aiCaption, &userID); err != nil {
			log.Printf("Error scanning repost: %v", err)
			continue
		}

		// Process the repost
		if err := s.processRepost(repostID, articleID, mediaAccountID, customCaption, aiCaption, userID); err != nil {
			log.Printf("Error processing repost %s: %v", repostID, err)
			continue
		}

		count++
	}

	if count > 0 {
		log.Printf("Processed %d scheduled reposts", count)
	}
}

func (s *Scheduler) processRepost(repostID, articleID, mediaAccountID string, customCaption, aiCaption *string, userID string) error {
	// Get article content
	var title, content, originalURL string
	err := s.db.QueryRow(`
		SELECT title, content, original_url 
		FROM articles 
		WHERE id = $1
	`, articleID).Scan(&title, &content, &originalURL)
	if err != nil {
		return fmt.Errorf("failed to get article: %v", err)
	}

	// Get media account info
	var platform, accountName string
	err = s.db.QueryRow(`
		SELECT platform, account_name 
		FROM media_accounts 
		WHERE id = $1
	`, mediaAccountID).Scan(&platform, &accountName)
	if err != nil {
		return fmt.Errorf("failed to get media account: %v", err)
	}

	// Determine caption to use
	caption := content
	if customCaption != nil && *customCaption != "" {
		caption = *customCaption
	} else if aiCaption != nil && *aiCaption != "" {
		caption = *aiCaption
	}

	// Simulate posting to social media platform
	log.Printf("Posting to %s (@%s): %s", platform, accountName, caption)

	// Update repost status
	_, err = s.db.Exec(`
		UPDATE reposts 
		SET status = 'posted', posted_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`, repostID)
	if err != nil {
		return fmt.Errorf("failed to update repost status: %v", err)
	}

	return nil
}

func (s *Scheduler) fetchArticles() {
	log.Println("Fetching articles...")

	// Get all topics with their keywords
	rows, err := s.db.Query(`
		SELECT id, name, keywords, platforms, user_id 
		FROM topics
	`)
	if err != nil {
		log.Printf("Error fetching topics: %v", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var topicID, name, userID string
		var keywords, platforms []string

		if err := rows.Scan(&topicID, &name, &keywords, &platforms, &userID); err != nil {
			log.Printf("Error scanning topic: %v", err)
			continue
		}

		// Simulate fetching articles for this topic
		if err := s.fetchArticlesForTopic(topicID, name, keywords, platforms, userID); err != nil {
			log.Printf("Error fetching articles for topic %s: %v", topicID, err)
			continue
		}

		count++
	}

	if count > 0 {
		log.Printf("Fetched articles for %d topics", count)
	}
}

func (s *Scheduler) fetchArticlesForTopic(topicID, name string, keywords, platforms []string, userID string) error {
	// This is a simulation - in real implementation, you would:
	// 1. Use platform APIs to search for articles based on keywords
	// 2. Filter and deduplicate articles
	// 3. Store new articles in the database

	log.Printf("Fetching articles for topic '%s' with keywords: %v on platforms: %v", name, keywords, platforms)

	// Simulate finding 1-3 articles
	for i := 0; i < 2; i++ {
		articleID := fmt.Sprintf("article_%d_%s", time.Now().Unix(), topicID)
		
		// Insert simulated article
		_, err := s.db.Exec(`
			INSERT INTO articles (id, title, content, original_url, platform, published_at, topic_id, user_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
			ON CONFLICT (id) DO NOTHING
		`, articleID, 
			fmt.Sprintf("Article about %s", name),
			fmt.Sprintf("Content related to %s keywords: %v", name, keywords),
			fmt.Sprintf("https://example.com/article/%s", articleID),
			platforms[0],
			time.Now().Add(-time.Hour*time.Duration(i)),
			topicID,
			userID)
		
		if err != nil {
			log.Printf("Error inserting article: %v", err)
		}
	}

	return nil
}

func (s *Scheduler) generateAICaptions() {
	log.Println("Generating AI captions...")

	// Get reposts that need AI captions
	rows, err := s.db.Query(`
		SELECT r.id, a.title, a.content
		FROM reposts r
		JOIN articles a ON r.article_id = a.id
		WHERE r.ai_caption IS NULL
		AND r.status = 'pending'
		LIMIT 5
	`)
	if err != nil {
		log.Printf("Error fetching reposts for AI captions: %v", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var repostID, title, content string

		if err := rows.Scan(&repostID, &title, &content); err != nil {
			log.Printf("Error scanning repost for AI caption: %v", err)
			continue
		}

		// Generate AI caption (simulated)
		aiCaption := s.generateAICaption(title, content)

		// Update repost with AI caption
		_, err := s.db.Exec(`
			UPDATE reposts 
			SET ai_caption = $1, updated_at = NOW()
			WHERE id = $2
		`, aiCaption, repostID)
		if err != nil {
			log.Printf("Error updating AI caption for repost %s: %v", repostID, err)
			continue
		}

		count++
	}

	if count > 0 {
		log.Printf("Generated AI captions for %d reposts", count)
	}
}

func (s *Scheduler) generateAICaption(title, content string) string {
	// This is a simulation - in real implementation, you would:
	// 1. Use OpenAI API or similar service
	// 2. Create a prompt with the article content
	// 3. Generate an engaging caption

	return fmt.Sprintf("🚀 %s - %s... #socialmedia #growth", title, content[:min(50, len(content))])
}

func (s *Scheduler) cleanupOldData() {
	log.Println("Cleaning up old data...")

	// Delete old verification tokens (older than 1 day)
	result, err := s.db.Exec(`
		DELETE FROM verification_tokens 
		WHERE expires < NOW() - INTERVAL '1 day'
	`)
	if err != nil {
		log.Printf("Error cleaning up verification tokens: %v", err)
	} else {
		if rowsAffected, _ := result.RowsAffected(); rowsAffected > 0 {
			log.Printf("Deleted %d old verification tokens", rowsAffected)
		}
	}

	// Delete old sessions (older than 7 days)
	result, err = s.db.Exec(`
		DELETE FROM sessions 
		WHERE expires < NOW() - INTERVAL '7 days'
	`)
	if err != nil {
		log.Printf("Error cleaning up sessions: %v", err)
	} else {
		if rowsAffected, _ := result.RowsAffected(); rowsAffected > 0 {
			log.Printf("Deleted %d old sessions", rowsAffected)
		}
	}

	// Delete old articles (older than 30 days)
	result, err = s.db.Exec(`
		DELETE FROM articles 
		WHERE created_at < NOW() - INTERVAL '30 days'
		AND id NOT IN (SELECT article_id FROM reposts)
	`)
	if err != nil {
		log.Printf("Error cleaning up articles: %v", err)
	} else {
		if rowsAffected, _ := result.RowsAffected(); rowsAffected > 0 {
			log.Printf("Deleted %d old articles", rowsAffected)
		}
	}

	log.Println("Cleanup completed")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}