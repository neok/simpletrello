# SimpleTrello

A minimal Trello-like board. Go backend, React frontend, SQLite database.
<img width="1316" height="544" alt="Screenshot 2026-04-10 at 16 12 53" src="https://github.com/user-attachments/assets/5c97c8f0-e4a1-4478-936c-8fbc61df92de" />


## Stack

- **Backend** — Go, `net/http`, `html/template`
- **Frontend** — React 19, Vite, Tailwind CSS v4
- **Database** — SQLite
- **Deploy** — Docker + Docker Compose

## Features

- Create, rename, delete lists (tabs)
- Create, edit, delete cards with title and description
- Move cards between lists

## Run with Docker

```bash
docker compose up --build
```

Open [http://localhost:8080](http://localhost:8080).

## Run locally

**Prerequisites:** Go 1.26+, Node 22+

```bash
# 1. Build frontend
cd frontend && npm install && npm run build && cd ..

# 2. Start server
mkdir -p data && go run ./cmd/web
```

Open [http://localhost:8080](http://localhost:8080).

For frontend hot-reload during development:

```bash
# Terminal 1 — Go server
mkdir -p data && go run ./cmd/web

# Terminal 2 — Vite dev server (proxies /api to :8080)
cd frontend && npm run dev
```

## API

| Method | Path                  | Description              |
|--------|-----------------------|--------------------------|
| GET    | /api/v1/tabs          | All tabs with their cards |
| POST   | /api/v1/tabs          | Create tab               |
| PATCH  | /api/v1/tabs/:id      | Rename tab               |
| DELETE | /api/v1/tabs/:id      | Delete tab               |
| POST   | /api/v1/cards         | Create card              |
| PATCH  | /api/v1/cards/:id     | Edit or move card        |
| DELETE | /api/v1/cards/:id     | Delete card              |

## Project structure

```
├── cmd/web/          # HTTP server, handlers, routes
├── internal/models/  # SQLite models (tabs, cards)
├── ui/
│   ├── embed.go      # Embeds templates + built assets into binary
│   ├── html/         # Go HTML templates
│   └── static/       # Vite build output (gitignored)
├── frontend/         # React source
├── migrations/       # SQL schema
├── Dockerfile
└── docker-compose.yml
```
