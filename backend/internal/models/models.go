package models

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Content represents a generated content item
type Content struct {
	ID          string                 `json:"id" db:"id"`
	Type        string                 `json:"type" db:"type"`
	Title       string                 `json:"title" db:"title"`
	Body        string                 `json:"body" db:"body"`
	Status      string                 `json:"status" db:"status"`
	Provider    string                 `json:"provider" db:"provider"`
	Model       string                 `json:"model" db:"model"`
	Metadata    map[string]interface{} `json:"metadata" db:"metadata"`
	Tags        []string               `json:"tags" db:"tags"`
	ScheduledAt *time.Time             `json:"scheduledAt,omitempty" db:"scheduled_at"`
	PublishedAt *time.Time             `json:"publishedAt,omitempty" db:"published_at"`
	CreatedAt   time.Time              `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time              `json:"updatedAt" db:"updated_at"`
}

// ContentType enum
const (
	ContentTypeScript      = "script"
	ContentTypeThumbnail   = "thumbnail"
	ContentTypeXPost       = "x_post"
	ContentTypeBlogPost    = "blog_post"
	ContentTypeNewsletter  = "newsletter"
	ContentTypeOutreach    = "outreach"
	ContentTypeClipScript  = "clip_script"
	ContentTypeTrendReport = "trend_report"
	ContentTypeIdea        = "idea"
)

// ContentStatus enum
const (
	ContentStatusDraft      = "draft"
	ContentStatusScheduled  = "scheduled"
	ContentStatusPublished  = "published"
	ContentStatusFailed     = "failed"
	ContentStatusProcessing = "processing"
	ContentStatusArchived   = "archived"
)

// NewContent creates a new content item
func NewContent(contentType, title, body string) *Content {
	now := time.Now()
	return &Content{
		ID:        uuid.New().String(),
		Type:      contentType,
		Title:     title,
		Body:      body,
		Status:    ContentStatusDraft,
		Metadata:  make(map[string]interface{}),
		Tags:      []string{},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Execute creates the content tables
func Execute(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS contents (
			id TEXT PRIMARY KEY,
			type TEXT NOT NULL,
			title TEXT NOT NULL,
			body TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'draft',
			provider TEXT,
			model TEXT,
			metadata TEXT,
			tags TEXT,
			scheduled_at DATETIME,
			published_at DATETIME,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_content_type ON contents(type)`,
		`CREATE INDEX IF NOT EXISTS idx_content_status ON contents(status)`,
		`CREATE INDEX IF NOT EXISTS idx_content_created ON contents(created_at)`,
		
		`CREATE TABLE IF NOT EXISTS ideas (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			category TEXT NOT NULL,
			sources TEXT,
			virality_score REAL,
			created_at DATETIME NOT NULL,
			used_at DATETIME
		)`,
		
		`CREATE TABLE IF NOT EXISTS analytics (
			id TEXT PRIMARY KEY,
			content_id TEXT,
			content_type TEXT,
			platform TEXT,
			views INTEGER DEFAULT 0,
			engagement INTEGER DEFAULT 0,
			ctr REAL DEFAULT 0,
			retention REAL DEFAULT 0,
			date DATE NOT NULL
		)`,
		
		`CREATE TABLE IF NOT EXISTS connections (
			id TEXT PRIMARY KEY,
			type TEXT NOT NULL,
			name TEXT NOT NULL,
			token TEXT,
			settings TEXT,
			is_active BOOLEAN DEFAULT 1,
			last_used DATETIME,
			created_at DATETIME NOT NULL
		)`,
	}
	
	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}
	
	return nil
}

// ContentStore provides database operations for content
type ContentStore struct {
	db *sql.DB
}

// NewContentStore creates a new content store
func NewContentStore(db *sql.DB) *ContentStore {
	return &ContentStore{db: db}
}

// Create inserts a new content item
func (s *ContentStore) Create(c *Content) error {
	metadata, _ := json.Marshal(c.Metadata)
	tags, _ := json.Marshal(c.Tags)
	
	_, err := s.db.Exec(
		`INSERT INTO contents (id, type, title, body, status, provider, model, metadata, tags, scheduled_at, published_at, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		c.ID, c.Type, c.Title, c.Body, c.Status, c.Provider, c.Model,
		string(metadata), string(tags), c.ScheduledAt, c.PublishedAt, c.CreatedAt, c.UpdatedAt,
	)
	return err
}

// Get retrieves a content item by ID
func (s *ContentStore) Get(id string) (*Content, error) {
	row := s.db.QueryRow(
		`SELECT id, type, title, body, status, provider, model, metadata, tags, scheduled_at, published_at, created_at, updated_at
		 FROM contents WHERE id = ?`, id,
	)
	
	var c Content
	var metadataJSON, tagsJSON string
	
	err := row.Scan(&c.ID, &c.Type, &c.Title, &c.Body, &c.Status, &c.Provider, &c.Model,
		&metadataJSON, &tagsJSON, &c.ScheduledAt, &c.PublishedAt, &c.CreatedAt, &c.UpdatedAt)
	
	if err != nil {
		return nil, err
	}
	
	json.Unmarshal([]byte(metadataJSON), &c.Metadata)
	json.Unmarshal([]byte(tagsJSON), &c.Tags)
	
	return &c, nil
}

// List returns a list of content items with optional filtering
func (s *ContentStore) List(contentType, status string, limit, offset int) ([]Content, error) {
	query := `SELECT id, type, title, body, status, provider, model, metadata, tags, scheduled_at, published_at, created_at, updated_at FROM contents WHERE 1=1`
	args := []interface{}{}
	
	if contentType != "" {
		query += " AND type = ?"
		args = append(args, contentType)
	}
	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}
	
	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)
	
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var contents []Content
	for rows.Next() {
		var c Content
		var metadataJSON, tagsJSON string
		
		err := rows.Scan(&c.ID, &c.Type, &c.Title, &c.Body, &c.Status, &c.Provider, &c.Model,
			&metadataJSON, &tagsJSON, &c.ScheduledAt, &c.PublishedAt, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			continue
		}
		
		json.Unmarshal([]byte(metadataJSON), &c.Metadata)
		json.Unmarshal([]byte(tagsJSON), &c.Tags)
		
		contents = append(contents, c)
	}
	
	return contents, nil
}

// Update updates a content item
func (s *ContentStore) Update(c *Content) error {
	c.UpdatedAt = time.Now()
	metadata, _ := json.Marshal(c.Metadata)
	tags, _ := json.Marshal(c.Tags)
	
	_, err := s.db.Exec(
		`UPDATE contents SET title = ?, body = ?, status = ?, provider = ?, model = ?, 
		 metadata = ?, tags = ?, scheduled_at = ?, published_at = ?, updated_at = ?
		 WHERE id = ?`,
		c.Title, c.Body, c.Status, c.Provider, c.Model,
		string(metadata), string(tags), c.ScheduledAt, c.PublishedAt, c.UpdatedAt, c.ID,
	)
	return err
}

// Delete removes a content item
func (s *ContentStore) Delete(id string) error {
	_, err := s.db.Exec("DELETE FROM contents WHERE id = ?", id)
	return err
}

// Count returns the total count of content items
func (s *ContentStore) Count(contentType, status string) (int, error) {
	query := "SELECT COUNT(*) FROM contents WHERE 1=1"
	args := []interface{}{}
	
	if contentType != "" {
		query += " AND type = ?"
		args = append(args, contentType)
	}
	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}
	
	var count int
	err := s.db.QueryRow(query, args...).Scan(&count)
	return count, err
}

// Search searches content by title or body
func (s *ContentStore) Search(query string, limit int) ([]Content, error) {
	likeQuery := "%" + query + "%"
	
	rows, err := s.db.Query(
		`SELECT id, type, title, body, status, provider, model, metadata, tags, scheduled_at, published_at, created_at, updated_at 
		 FROM contents WHERE title LIKE ? OR body LIKE ? ORDER BY created_at DESC LIMIT ?`,
		likeQuery, likeQuery, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var contents []Content
	for rows.Next() {
		var c Content
		var metadataJSON, tagsJSON string
		
		err := rows.Scan(&c.ID, &c.Type, &c.Title, &c.Body, &c.Status, &c.Provider, &c.Model,
			&metadataJSON, &tagsJSON, &c.ScheduledAt, &c.PublishedAt, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			continue
		}
		
		json.Unmarshal([]byte(metadataJSON), &c.Metadata)
		json.Unmarshal([]byte(tagsJSON), &c.Tags)
		
		contents = append(contents, c)
	}
	
	return contents, nil
}
