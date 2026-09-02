package scraper

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ArminDashti/knowledge-center-api/internal/models"
	"github.com/PuerkitoBio/goquery"
)

const MaxPages = 20

type ResultPage struct {
	URL         string
	Title       string
	ContentText string
	ContentHTML string
}

type Client struct {
	http *http.Client
}

func New() *Client {
	return &Client{
		http: &http.Client{Timeout: 30 * time.Second},
	}
}

func DefaultProfile(host string) models.ScrapeProfile {
	return models.ScrapeProfile{
		Host:            host,
		TitleSelector:   "h1",
		ContentSelector: "article, main, body",
	}
}

func (c *Client) Scrape(ctx context.Context, startURL string, profile models.ScrapeProfile) ([]ResultPage, error) {
	start, err := url.Parse(startURL)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}
	if start.Scheme == "" {
		start.Scheme = "https"
	}

	visited := map[string]struct{}{}
	queue := []*url.URL{start}
	results := make([]ResultPage, 0, 1)

	for len(queue) > 0 && len(results) < MaxPages {
		current := queue[0]
		queue = queue[1:]
		key := normalizeURL(current)
		if _, ok := visited[key]; ok {
			continue
		}
		visited[key] = struct{}{}

		doc, body, err := c.fetchDocument(ctx, current.String())
		if err != nil {
			if len(results) == 0 {
				return nil, err
			}
			continue
		}

		if profile.ExcludeSelector != "" {
			doc.Find(profile.ExcludeSelector).Remove()
		}

		title := strings.TrimSpace(doc.Find(profile.TitleSelector).First().Text())
		if title == "" {
			title = strings.TrimSpace(doc.Find("title").First().Text())
		}

		contentSel := doc.Find(profile.ContentSelector).First()
		if contentSel.Length() == 0 {
			contentSel = doc.Find("body")
		}
		contentHTML, _ := contentSel.Html()
		contentText := strings.TrimSpace(contentSel.Text())

		results = append(results, ResultPage{
			URL:         current.String(),
			Title:       title,
			ContentText: contentText,
			ContentHTML: contentHTML,
		})
		_ = body

		if strings.TrimSpace(profile.LinkSelector) == "" {
			continue
		}
		doc.Find(profile.LinkSelector).Each(func(_ int, sel *goquery.Selection) {
			href, ok := sel.Attr("href")
			if !ok || strings.TrimSpace(href) == "" {
				return
			}
			next, err := current.Parse(href)
			if err != nil {
				return
			}
			if next.Host != "" && !strings.EqualFold(next.Host, start.Host) {
				return
			}
			next.Fragment = ""
			if next.Scheme != "http" && next.Scheme != "https" {
				return
			}
			nkey := normalizeURL(next)
			if _, seen := visited[nkey]; seen {
				return
			}
			queue = append(queue, next)
		})
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no pages scraped")
	}
	return results, nil
}

func (c *Client) fetchDocument(ctx context.Context, rawURL string) (*goquery.Document, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("User-Agent", "knowledge-center-api/0.1 (+local scraper)")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, "", fmt.Errorf("unexpected status %d for %s", resp.StatusCode, rawURL)
	}

	limited := io.LimitReader(resp.Body, 8<<20)
	doc, err := goquery.NewDocumentFromReader(limited)
	if err != nil {
		return nil, "", err
	}
	return doc, "", nil
}

func normalizeURL(u *url.URL) string {
	clone := *u
	clone.Fragment = ""
	return strings.ToLower(clone.Scheme) + "://" + strings.ToLower(clone.Host) + clone.Path + "?" + clone.RawQuery
}
