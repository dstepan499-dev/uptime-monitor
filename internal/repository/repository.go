package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Repository struct {
	db *sql.DB
}

func New(connString string) (*Repository, error) {
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(25)                 //лимит на количество одновременно открытых соединений
	db.SetMaxIdleConns(5)                  //лимит на количество неактивных соединений
	db.SetConnMaxLifetime(5 * time.Minute) //время жизни соединения с бд

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Repository{db: db}, nil
}

func (r *Repository) Close() error {
	return r.db.Close()
}

func (r *Repository) CreateWebsite(ctx context.Context, website *Website) error {
	query := `INSERT INTO websites (name, url, interval_seconds, timeout_seconds) VALUES ($1, $2, $3, $4) RETURNING id, created_at;`
	err := r.db.QueryRowContext(ctx, query, website.Name, website.URL, website.IntervalSeconds, website.TimeoutSeconds).Scan(&website.ID, &website.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create url: %w", err)
	}

	return nil
}

func (r *Repository) GetActiveWebsites(ctx context.Context) ([]Website, error) {
	query := `SELECT id, name, url, interval_seconds, timeout_seconds, is_active, last_status_code, last_checked_at, created_at, ssl_expires_at FROM websites WHERE is_active;`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query active websittes: %w", err)
	}
	defer rows.Close()

	var websites []Website
	for rows.Next() {
		var w Website
		if err := rows.Scan(&w.ID, &w.Name, &w.URL, &w.IntervalSeconds, &w.TimeoutSeconds, &w.IsActive, &w.LastStatusCode, &w.LastCheckedAt, &w.CreatedAt, &w.SSLExpiresAt); err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}

		websites = append(websites, w)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return websites, nil
}

func (r *Repository) SaveCheckLogs(ctx context.Context, log *CheckLog) error {
	query := `INSERT INTO check_logs (website_id, status_code, response_time_ms, is_successful, error_message) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, log.WebsiteID, log.StatusCode, log.ResponseTimeMS, log.IsSuccessful, log.ErrorMessage)
	if err != nil {
		return fmt.Errorf("failed to save logs: %w", err)
	}

	return nil
}

func (r *Repository) UpdateWebsiteStatus(ctx context.Context, websiteID int, statusCode *int, sslExpiresAt *time.Time) error {
	query := `UPDATE websites 
	SET last_status_code = $1, ssl_expires_at = $2, last_checked_at = CURRENT_TIMESTAMP
	WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, statusCode, sslExpiresAt, websiteID)
	if err != nil {
		return fmt.Errorf("failed to update website status: %w", err)
	}

	return nil
}
