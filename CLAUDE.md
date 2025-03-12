# Swing Society Website Development Guide

## Build & Run Commands
- Run server: `cd server && go run main.go`
- Live reload: `air` (uses air.toml configuration)
- Build for production: `cd server && go build -o main .`
- Docker build: `docker build -t gcr.io/swingsociety-backend/ss-go .`
- Deploy: `./deploy.sh`

## Test Commands
- Run all tests: `go test ./...`
- Run specific test: `go test -v ./server/internal/config -run TestLoadConfig`
- API tests: `./server/tests/test_api.sh`
- Rate limit tests: `./server/tests/test_rate_limit.sh`

## Code Style Guidelines
- **Imports**: Standard library first, then external packages, then internal
- **Error Handling**: Use custom error types from `internal/errors`
- **Responses**: Use `response.JSON()` and `response.Error()` for consistent formatting
- **Validation**: Model-based with detailed field errors
- **Naming**: CamelCase for exports, handler functions prefixed with "Handle"
- **Package Structure**: Models in `api/models/`, handlers in `api/handlers/`
- **Frontend**: Use modules in `static/js/modules/` for component functionality
- **Templates**: Follow Bulma CSS framework conventions for styling
- **CSS**: Use rem measures when possible, currently used css files: style.css bulma-custom.css responsive.css, rest of the css files are deprecated