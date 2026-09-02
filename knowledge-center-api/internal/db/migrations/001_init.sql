CREATE TABLE IF NOT EXISTS scrape_profiles (
    id UUID PRIMARY KEY,
    host TEXT NOT NULL UNIQUE,
    title_selector TEXT NOT NULL DEFAULT 'h1',
    content_selector TEXT NOT NULL DEFAULT 'article, main, body',
    link_selector TEXT NOT NULL DEFAULT '',
    exclude_selector TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS imports (
    id UUID PRIMARY KEY,
    url TEXT NOT NULL,
    host TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    error_message TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    scraped_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_imports_created_at ON imports (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_imports_status ON imports (status);

CREATE TABLE IF NOT EXISTS pages (
    id UUID PRIMARY KEY,
    import_id UUID NOT NULL REFERENCES imports (id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    title TEXT NOT NULL DEFAULT '',
    content_text TEXT NOT NULL DEFAULT '',
    content_html TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    scraped_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_pages_import_id ON pages (import_id);
