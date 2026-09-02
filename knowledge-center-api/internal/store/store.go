package store

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ArminDashti/knowledge-center-api/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")

type Store struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

func (s *Store) CreateImport(ctx context.Context, rawURL, host string) (*models.Import, error) {
	item := &models.Import{
		ID:        uuid.NewString(),
		URL:       rawURL,
		Host:      host,
		Status:    models.StatusPending,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_, err := s.pool.Exec(ctx, `
		INSERT INTO imports (id, url, host, status, error_message, created_at, updated_at)
		VALUES ($1, $2, $3, $4, '', $5, $6)
	`, item.ID, item.URL, item.Host, item.Status, item.CreatedAt, item.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert import: %w", err)
	}
	return item, nil
}

func (s *Store) ListImports(ctx context.Context) ([]models.Import, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, url, host, status, error_message, created_at, updated_at, scraped_at
		FROM imports
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.Import, 0)
	for rows.Next() {
		var item models.Import
		if err := rows.Scan(
			&item.ID, &item.URL, &item.Host, &item.Status, &item.ErrorMessage,
			&item.CreatedAt, &item.UpdatedAt, &item.ScrapedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) GetImport(ctx context.Context, id string) (*models.Import, error) {
	var item models.Import
	err := s.pool.QueryRow(ctx, `
		SELECT id, url, host, status, error_message, created_at, updated_at, scraped_at
		FROM imports WHERE id = $1
	`, id).Scan(
		&item.ID, &item.URL, &item.Host, &item.Status, &item.ErrorMessage,
		&item.CreatedAt, &item.UpdatedAt, &item.ScrapedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *Store) UpdateImportStatus(ctx context.Context, id string, status models.ImportStatus, errMsg string, scrapedAt *time.Time) error {
	_, err := s.pool.Exec(ctx, `
		UPDATE imports
		SET status = $2, error_message = $3, updated_at = $4, scraped_at = $5
		WHERE id = $1
	`, id, status, errMsg, time.Now().UTC(), scrapedAt)
	return err
}

func (s *Store) ListPagesByImport(ctx context.Context, importID string) ([]models.PageSummary, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, import_id, url, title, created_at, scraped_at
		FROM pages WHERE import_id = $1
		ORDER BY created_at ASC
	`, importID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.PageSummary, 0)
	for rows.Next() {
		var item models.PageSummary
		if err := rows.Scan(&item.ID, &item.ImportID, &item.URL, &item.Title, &item.CreatedAt, &item.ScrapedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) GetPage(ctx context.Context, importID, pageID string) (*models.Page, error) {
	var page models.Page
	err := s.pool.QueryRow(ctx, `
		SELECT id, import_id, url, title, content_text, content_html, created_at, scraped_at
		FROM pages WHERE import_id = $1 AND id = $2
	`, importID, pageID).Scan(
		&page.ID, &page.ImportID, &page.URL, &page.Title,
		&page.ContentText, &page.ContentHTML, &page.CreatedAt, &page.ScrapedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &page, nil
}

func (s *Store) InsertPage(ctx context.Context, importID, pageURL, title, contentText, contentHTML string) error {
	now := time.Now().UTC()
	_, err := s.pool.Exec(ctx, `
		INSERT INTO pages (id, import_id, url, title, content_text, content_html, created_at, scraped_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, uuid.NewString(), importID, pageURL, title, contentText, contentHTML, now, now)
	return err
}

func (s *Store) GetProfileByHost(ctx context.Context, host string) (*models.ScrapeProfile, error) {
	var profile models.ScrapeProfile
	err := s.pool.QueryRow(ctx, `
		SELECT id, host, title_selector, content_selector, link_selector, exclude_selector, created_at, updated_at
		FROM scrape_profiles WHERE host = $1
	`, strings.ToLower(host)).Scan(
		&profile.ID, &profile.Host, &profile.TitleSelector, &profile.ContentSelector,
		&profile.LinkSelector, &profile.ExcludeSelector, &profile.CreatedAt, &profile.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (s *Store) ListProfiles(ctx context.Context) ([]models.ScrapeProfile, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, host, title_selector, content_selector, link_selector, exclude_selector, created_at, updated_at
		FROM scrape_profiles
		ORDER BY host ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.ScrapeProfile, 0)
	for rows.Next() {
		var profile models.ScrapeProfile
		if err := rows.Scan(
			&profile.ID, &profile.Host, &profile.TitleSelector, &profile.ContentSelector,
			&profile.LinkSelector, &profile.ExcludeSelector, &profile.CreatedAt, &profile.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, profile)
	}
	return items, rows.Err()
}

func (s *Store) GetProfile(ctx context.Context, id string) (*models.ScrapeProfile, error) {
	var profile models.ScrapeProfile
	err := s.pool.QueryRow(ctx, `
		SELECT id, host, title_selector, content_selector, link_selector, exclude_selector, created_at, updated_at
		FROM scrape_profiles WHERE id = $1
	`, id).Scan(
		&profile.ID, &profile.Host, &profile.TitleSelector, &profile.ContentSelector,
		&profile.LinkSelector, &profile.ExcludeSelector, &profile.CreatedAt, &profile.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func normalizeProfileFields(req *models.CreateProfileRequest) {
	if strings.TrimSpace(req.TitleSelector) == "" {
		req.TitleSelector = "h1"
	}
	if strings.TrimSpace(req.ContentSelector) == "" {
		req.ContentSelector = "article, main, body"
	}
	req.Host = strings.ToLower(strings.TrimSpace(req.Host))
}

func (s *Store) CreateProfile(ctx context.Context, req models.CreateProfileRequest) (*models.ScrapeProfile, error) {
	normalizeProfileFields(&req)
	now := time.Now().UTC()
	profile := &models.ScrapeProfile{
		ID:              uuid.NewString(),
		Host:            req.Host,
		TitleSelector:   req.TitleSelector,
		ContentSelector: req.ContentSelector,
		LinkSelector:    req.LinkSelector,
		ExcludeSelector: req.ExcludeSelector,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	_, err := s.pool.Exec(ctx, `
		INSERT INTO scrape_profiles (
			id, host, title_selector, content_selector, link_selector, exclude_selector, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`, profile.ID, profile.Host, profile.TitleSelector, profile.ContentSelector,
		profile.LinkSelector, profile.ExcludeSelector, profile.CreatedAt, profile.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (s *Store) UpdateProfile(ctx context.Context, id string, req models.UpdateProfileRequest) (*models.ScrapeProfile, error) {
	createReq := models.CreateProfileRequest{
		Host:            req.Host,
		TitleSelector:   req.TitleSelector,
		ContentSelector: req.ContentSelector,
		LinkSelector:    req.LinkSelector,
		ExcludeSelector: req.ExcludeSelector,
	}
	normalizeProfileFields(&createReq)
	now := time.Now().UTC()
	tag, err := s.pool.Exec(ctx, `
		UPDATE scrape_profiles
		SET host = $2, title_selector = $3, content_selector = $4,
		    link_selector = $5, exclude_selector = $6, updated_at = $7
		WHERE id = $1
	`, id, createReq.Host, createReq.TitleSelector, createReq.ContentSelector,
		createReq.LinkSelector, createReq.ExcludeSelector, now)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, ErrNotFound
	}
	return s.GetProfile(ctx, id)
}

func (s *Store) DeleteProfile(ctx context.Context, id string) error {
	tag, err := s.pool.Exec(ctx, `DELETE FROM scrape_profiles WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
