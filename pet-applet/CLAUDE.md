# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A pet management app (宠物助手) — track pets, feeding schedules, and feeding records.

## Architecture

Two sub-projects in one monorepo:

```
pet-applet/
├── backend/          # Go API server
│   ├── main.go            — Entry point, Gin routes
│   ├── config/            — YAML config loading (config.go + config.yaml)
│   ├── database/          — MySQL connection + auto-migration on startup
│   ├── handlers/          — HTTP handlers (pets, schedules, records, meta)
│   └── models/            — Data structs (Pet, FeedingSchedule, FeedingRecord)
│
└── miniprogram/      # uni-app + Vue 3 + TypeScript
    ├── src/
    │   ├── api/            — Typed HTTP client layer wrapping uni.request()
    │   ├── config/         — API_BASE_URL constant
    │   ├── pages/
    │   │   ├── pets/       — List, detail, edit/create
    │   │   ├── schedules/  — CRUD for feeding schedules (inline form)
    │   │   ├── records/    — View feeding records
    │   │   └── index/      — Redirect shell to pets/index
    │   ├── App.vue
    │   └── main.ts         — SSR app factory
    ├── pages.json          — Route definitions
    └── vite.config.ts      — Vite + @dcloudio/vite-plugin-uni
```

## Backend (Go)

- **Framework**: Gin with raw `database/sql` + MySQL (`go-sql-driver/mysql`)
- **Port**: 3000 (set in `config/config.yaml`, override with env `CONFIG_PATH`)
- **Auto-migration**: Tables created on startup in `database/migrate()`
- **ID scheme**: 8-char random alphanumeric (lowercase + digits) from `handlers/generateID()`
- **API routes** (prefix `/api`):
  - Pets: `GET /pets`, `GET /pets/:id`, `POST /pets`, `PUT /pets/:id`, `DELETE /pets/:id`
  - Schedules: `GET /pets/schedules/:petId`, `POST /pets/schedules/:petId`, `PUT /schedules/:id`, `DELETE /schedules/:id`
  - Records: `GET /pets/records/:petId`, `GET /pets/records/today/:petId`, `POST /pets/records/:petId`, `DELETE /records/:id`
  - Meta: `GET /meta/breeds` — returns emoji list + breed options per animal type

## Frontend (uni-app / Vue 3)

- **Pages** (defined in `pages.json`):
  - `pets/index` — pet list with FAB add button; loads on `onShow`
  - `pets/detail` — pet info card + today's records + schedules list + inline record form (popup)
  - `pets/edit` — create/edit form with emoji picker, breed/dob pickers
  - `schedules/index` — list + inline CRUD form (add/edit/delete)
  - `records/index` — list with delete confirmation
  - `index/index` — empty shell, redirects to `pets/index`
- **API client**: Typed functions in `src/api/` modules (pets/schedules/records/meta) using the shared `request.ts` wrapper
- **Config**: `API_BASE_URL = 'http://localhost:3000/api'` in `src/config/index.ts`

### Commands

```bash
# Backend — build and run
cd backend && go build -o pet-applet-server . && ./pet-applet-server

# Frontend — H5 (browser dev)
cd miniprogram && npm run dev:h5

# Frontend — WeChat mini-program
npm run dev:mp-weixin

# Frontend — TypeScript type-check
npm run type-check

# Frontend — build for production (H5)
npm run build:h5
```
