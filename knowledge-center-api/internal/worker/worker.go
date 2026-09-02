package worker

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/ArminDashti/knowledge-center-api/internal/models"
	"github.com/ArminDashti/knowledge-center-api/internal/scraper"
	"github.com/ArminDashti/knowledge-center-api/internal/store"
)

type Worker struct {
	store   *store.Store
	scraper *scraper.Client
	jobs    chan string
}

func New(st *store.Store, client *scraper.Client, buffer int) *Worker {
	if buffer < 1 {
		buffer = 32
	}
	return &Worker{
		store:   st,
		scraper: client,
		jobs:    make(chan string, buffer),
	}
}

func (w *Worker) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case importID := <-w.jobs:
				w.process(ctx, importID)
			}
		}
	}()
}

func (w *Worker) Enqueue(importID string) {
	select {
	case w.jobs <- importID:
	default:
		go func() {
			w.jobs <- importID
		}()
	}
}

func (w *Worker) process(parent context.Context, importID string) {
	ctx, cancel := context.WithTimeout(parent, 5*time.Minute)
	defer cancel()

	item, err := w.store.GetImport(ctx, importID)
	if err != nil {
		log.Printf("worker: get import %s: %v", importID, err)
		return
	}

	if err := w.store.UpdateImportStatus(ctx, importID, models.StatusRunning, "", nil); err != nil {
		log.Printf("worker: set running %s: %v", importID, err)
		return
	}

	profile, err := w.store.GetProfileByHost(ctx, item.Host)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			fallback := scraper.DefaultProfile(item.Host)
			profile = &fallback
		} else {
			_ = w.store.UpdateImportStatus(ctx, importID, models.StatusFailed, err.Error(), nil)
			return
		}
	}

	pages, err := w.scraper.Scrape(ctx, item.URL, *profile)
	if err != nil {
		_ = w.store.UpdateImportStatus(ctx, importID, models.StatusFailed, err.Error(), nil)
		return
	}

	for _, page := range pages {
		if err := w.store.InsertPage(ctx, importID, page.URL, page.Title, page.ContentText, page.ContentHTML); err != nil {
			_ = w.store.UpdateImportStatus(ctx, importID, models.StatusFailed, err.Error(), nil)
			return
		}
	}

	now := time.Now().UTC()
	if err := w.store.UpdateImportStatus(ctx, importID, models.StatusSuccess, "", &now); err != nil {
		log.Printf("worker: set success %s: %v", importID, err)
	}
}
