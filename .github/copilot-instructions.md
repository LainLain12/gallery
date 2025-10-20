<!-- .github/copilot-instructions.md - Guidance for AI coding agents working on this repo -->
# Quick orientation

This repository is a small Go (1.21) REST API that serves wallpaper images from local folders using the Gin web framework. The server discovers files under `./images/<category>/` and exposes JSON endpoints that build image URLs dynamically from the incoming request.

Key files
- `main.go` — entire application. Read this first: it contains routing, configuration (categories, extensions, resolutions), static file serving, and helper functions that generate titles/tags/resolutions.
- `go.mod` — project dependencies (Gin, cors). Use `go mod tidy` to ensure modules are present.
- `start-server.bat` — convenience script for Windows devs (runs `go run main.go`).
- `README.md` & `privacy_policy.txt` — user-facing docs; `main.go` reads `privacy_policy.txt` for HTML/JSON endpoints.

Top-level architecture notes (what to know)
- Single binary HTTP API using Gin. There are no separate packages — everything lives in `main.go`.
- Static files are served from the `images/` directory via `r.Static("/images", "./images")`. The JSON responses include absolute URLs constructed from the incoming request's Host header.
- Categories, allowed image extensions, and the list of resolutions are in top-level vars in `main.go`. Changing categories requires updating that variable and corresponding image folders.
- Randomization: Titles, tags and resolutions are generated on every request using `math/rand`. This is intentional (lightweight demo behavior) but means responses are non-deterministic in tests.

Developer workflows and commands
- Install/verify dependencies:
  - `go mod tidy`
- Run locally:
  - `go run main.go` (the app prints startup info and starts a server)
  - There is a Windows helper: `start-server.bat` (calls `go mod tidy` if needed then `go run main.go`).
- Build binary:
  - `go build -o wallpaper-api` (builds executable)
- Quick smoke tests (examples):
  - GET all wallpapers: `GET http://localhost:8664/api/v1/wallpapers`
  - GET category: `GET http://localhost:8664/api/v1/wallpapers/nature`
  - GET random: `GET http://localhost:8664/api/v1/wallpapers/nature/random`
  - Health: `GET http://localhost:8664/health`

Important project-specific gotchas
- Port mismatch: `README.md` and `start-server.bat` mention port 8080, but `main.go` starts the server on port 8664 (see `r.Run(":8664")`). When changing the port, update all doc files and `start-server.bat` for consistency.
- Host header / reverse proxy behavior: image URLs are constructed from `c.Request.Host` and the request scheme is inferred from TLS or `X-Forwarded-Proto`. When testing behind proxies, ensure `X-Forwarded-Proto` and `Host` are set correctly.
- Categories are hard-coded in `main.go` as `[]string{"nature","culture","digital"}`. Add/remove categories in that slice and create corresponding subfolders under `images/`.
- The code uses `ioutil.ReadDir` and `ioutil.ReadFile` (deprecated in newer Go but still available). Be mindful if upgrading Go or refactoring to `os`/`io` packages.
- Responses are wrapped in `APIResponse` with `Success` + `Data` fields. Keep that shape if adding new endpoints or tests.

Patterns & examples agents should follow
- Small, single-file changes are preferred for feature work unless extracting packages is explicitly requested.
- When adding new endpoints, register them under the `api := r.Group("/api/v1")` block and return JSON responses using the existing `APIResponse` shape.
- For static assets and URLs, follow the existing approach: use `r.Static("/images", "./images")` and construct `ImageURL` using the request host/scheme so URLs remain valid behind proxies.
- For tests, isolate randomness: wrap or mock `generateRandomTitle`, `generateRandomTags`, and `getRandomResolution` if deterministic output is required.

Integration points & external dependencies
- Gin (`github.com/gin-gonic/gin`) — HTTP server and routing
- Gin CORS (`github.com/gin-contrib/cors`) — configured with AllowAllOrigins = true
- No external databases or cloud services; image files are read from disk.

Suggested quick tasks an agent can do immediately
- Fix documentation: align `README.md`, `start-server.bat`, and the printed startup port in `main.go` (choose a single canonical port). When doing so, update any hard-coded examples.
- Add a small unit/integration test that stubs `loadWallpapersFromFolder` to ensure API shape remains stable and to document the expected JSON fields.
- Add a `Makefile` or update `start-server.bat` to allow configuring port via env var (e.g., `PORT`) to avoid hard-coded port mismatches.

When making PRs
- Keep diffs minimal and explain any public API changes (JSON shape, endpoint paths) in the PR description.
- Mention backwards-incompatible changes: changing category names, response schema, or the image URL format.

If something is unclear
- If you need to know why a particular choice was made (for example, port 8664 vs 8080), look for related commits or ask the maintainers. There are multiple hints in `README.md`, `start-server.bat`, and `main.go` that diverge.

— End of file
