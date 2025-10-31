package repository

import (
	"context"
	"database/sql"
	"time"

	"url_shortener/internal/models"
)

// Ensure urlRepository implements URLRepository
var _ URLRepository = (*urlRepository)(nil)

func (r *urlRepository) Create(u *models.URL) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	affected, err := r.db.ExecContext(ctx,
		"INSERT INTO url_info (ui_short_url, ui_long_url, ui_created_at, ui_status) VALUES (?, ?, COALESCE(?, CURRENT_TIMESTAMP), 1)",
		u.ShortCode, u.OriginalURL, nullableTime(u.CreatedAt),
	)
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *urlRepository) GetByShortCode(shortCode string) (*models.URL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	row := r.db.QueryRowContext(ctx,
		"SELECT ui_id, ui_short_url, ui_long_url, ui_created_at FROM url_info WHERE ui_short_url = ? AND ui_status = 1",
		shortCode,
	)
	var u models.URL
	if err := row.Scan(&u.ID, &u.ShortCode, &u.OriginalURL, &u.CreatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *urlRepository) GetByOriginalURL(originalURL string) (*models.URL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	row := r.db.QueryRowContext(ctx,
		"SELECT ui_id, ui_short_url, ui_long_url, ui_created_at FROM url_info WHERE ui_long_url = ? AND ui_status = 1",
		originalURL,
	)
	var u models.URL
	if err := row.Scan(&u.ID, &u.ShortCode, &u.OriginalURL, &u.CreatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *urlRepository) Exists(shortCode string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	row := r.db.QueryRowContext(ctx,
		"SELECT 1 FROM url_info WHERE ui_short_url = ? AND ui_status = 1 LIMIT 1",
		shortCode,
	)
	var one int
	if err := row.Scan(&one); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *urlRepository) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Soft delete: mark inactive
	affected, err := r.db.ExecContext(ctx,
		"UPDATE url_info SET ui_status = 0 WHERE ui_id = ?",
		id,
	)
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func nullableTime(t time.Time) any {
	if t.IsZero() {
		return nil
	}
	return t
}
