# Project Restructure Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Reorganize the project into a Standard Go Layout (`cmd/`, `internal/`, `web/`, `scripts/`) to improve maintainability and cleanliness.

**Architecture:** Move entry point to `cmd/server`. Encapsulate core logic in `internal`. Move assets to `web`. Add a `Makefile` for task management.

**Tech Stack:** Go, Make, Docker, Bash.

### Task 1: Create Directory Structure and Move Files

**Files:**
- Create directories: `cmd/server`, `internal`, `web`, `scripts`
- Move: `main.go` -> `cmd/server/main.go`
- Move: `handlers/` -> `internal/handlers/`
- Move: `database/` -> `internal/database/`
- Move: `utils/` -> `internal/utils/`
- Move: `templates/` -> `web/templates/`
- Move: `static/` -> `web/static/`
- Move: `*.sh` -> `scripts/`

**Step 1: Create directories**
```bash
mkdir -p cmd/server internal web scripts
```

**Step 2: Move files**
```bash
mv main.go cmd/server/
mv handlers internal/
mv database internal/
mv utils internal/
mv templates web/
mv static web/
mv *.sh scripts/
```

**Step 3: Verify moves**
Run: `ls -R`
Expected: Files are in new locations.

**Step 4: Commit**
```bash
git add .
git commit -m "refactor: restructure project to standard go layout"
```

### Task 2: Update Go Imports and Paths

**Files:**
- Modify: `cmd/server/main.go`
- Modify: `internal/handlers/*.go`
- Modify: `internal/utils/render.go`

**Step 1: Fix imports in main.go**
Change `portfolio/handlers` -> `portfolio/internal/handlers`
Change `portfolio/database` -> `portfolio/internal/database`
Update static file path: `"static"` -> `"web/static"`

**Step 2: Fix imports in handlers**
Change `portfolio/database` -> `portfolio/internal/database`
Change `portfolio/utils` -> `portfolio/internal/utils`

**Step 3: Fix imports in database and utils**
Check if they have cross-dependencies and update them.

**Step 4: Fix template paths in utils/render.go**
Change `"templates/"` -> `"web/templates/"` (lines 11 and 14)

**Step 5: Verify build**
Run: `go build ./cmd/server`
Expected: Build success (or go get might be needed if mod verification fails, but local paths should resolve).

**Step 6: Commit**
```bash
git add .
git commit -m "refactor: update imports and file paths"
```

### Task 3: Create Makefile

**Files:**
- Create: `Makefile`

**Step 1: Create Makefile content**
```makefile
.PHONY: run build docker-up docker-down clean

run:
	go run cmd/server/main.go

build:
	go build -o bin/server cmd/server/main.go

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

clean:
	rm -rf bin/
```

**Step 2: Verify Makefile**
Run: `make build`
Expected: `bin/server` is created.

**Step 3: Commit**
```bash
git add Makefile
git commit -m "feat: add Makefile"
```

### Task 4: Update Docker Configuration

**Files:**
- Modify: `Dockerfile`
- Modify: `docker-compose.yml`

**Step 1: Update Dockerfile**
- Update builder stage:
    - `COPY . .` (still works, copies new structure)
    - `RUN go build -o server .` -> `RUN go build -o server ./cmd/server`
- Update runtime stage:
    - `COPY --from=builder /app/static ./static` -> `COPY --from=builder /app/web/static ./web/static`
    - `COPY --from=builder /app/templates ./templates` -> `COPY --from=builder /app/web/templates ./web/templates`
    - `CMD ["./server"]` (remains same)

**Step 2: Update docker-compose.yml**
- Ensure volumes map correctly if changed (currently maps `./sqlite:/app/sqlite`, which is still at root, so fine).

**Step 3: Verify Docker build**
Run: `make docker-up`
Expected: Container builds and runs successfully.

**Step 4: Commit**
```bash
git add Dockerfile docker-compose.yml
git commit -m "chore: update docker config for new layout"
```
