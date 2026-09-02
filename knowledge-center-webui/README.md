# Knowledge Center WebUI

Vue 3 + Vite + TypeScript + Tailwind + shadcn-vue PWA for the Knowledge Center API.

## Pages

- `/import` — submit a website URL
- `/list` — imports and scrape status; open detail for content
- `/scraping` — per-domain CSS scrape profiles

## Setup

```powershell
npm install
npm run dev
```

Proxies `/api` to `http://127.0.0.1:8091`. Optional `VITE_API_BASE_URL` in `.env`.

Production build enables `vite-plugin-pwa` (manifest + Workbox).
