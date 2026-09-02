package models

import "time"

type ImportStatus string

const (
	StatusPending ImportStatus = "pending"
	StatusRunning ImportStatus = "running"
	StatusSuccess ImportStatus = "success"
	StatusFailed  ImportStatus = "failed"
)

type ScrapeProfile struct {
	ID              string    `json:"id"`
	Host            string    `json:"host"`
	TitleSelector   string    `json:"title_selector"`
	ContentSelector string    `json:"content_selector"`
	LinkSelector    string    `json:"link_selector"`
	ExcludeSelector string    `json:"exclude_selector"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Import struct {
	ID           string       `json:"id"`
	URL          string       `json:"url"`
	Host         string       `json:"host"`
	Status       ImportStatus `json:"status"`
	ErrorMessage string       `json:"error_message"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	ScrapedAt    *time.Time   `json:"scraped_at,omitempty"`
}

type PageSummary struct {
	ID        string    `json:"id"`
	ImportID  string    `json:"import_id"`
	URL       string    `json:"url"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	ScrapedAt time.Time `json:"scraped_at"`
}

type Page struct {
	PageSummary
	ContentText string `json:"content_text"`
	ContentHTML string `json:"content_html"`
}

type ImportDetail struct {
	Import
	Pages []PageSummary `json:"pages"`
}

type CreateImportRequest struct {
	URL string `json:"url" binding:"required"`
}

type CreateProfileRequest struct {
	Host            string `json:"host" binding:"required"`
	TitleSelector   string `json:"title_selector"`
	ContentSelector string `json:"content_selector"`
	LinkSelector    string `json:"link_selector"`
	ExcludeSelector string `json:"exclude_selector"`
}

type UpdateProfileRequest struct {
	Host            string `json:"host" binding:"required"`
	TitleSelector   string `json:"title_selector"`
	ContentSelector string `json:"content_selector"`
	LinkSelector    string `json:"link_selector"`
	ExcludeSelector string `json:"exclude_selector"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
