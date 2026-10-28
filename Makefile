.PHONY: install run-backend run-frontend test coverage

install:
	cd frontend && npm ci

# API on :8080
run-backend:
	cd backend && go run ./cmd/server

# UI on :5173, proxying /api to the backend
run-frontend:
	cd frontend && npm run dev

test:
	cd backend && go test -race ./...
	cd frontend && npm test

# HTML reports stay local (git-ignored); text summaries go to docs/coverage/.
coverage:
	mkdir -p docs/coverage
	cd backend && go test -coverprofile=coverage.out ./... \
		&& go tool cover -html=coverage.out -o coverage.html \
		&& go tool cover -func=coverage.out > ../docs/coverage/backend.txt
	cd frontend && npm run test:coverage && cp coverage/coverage.txt ../docs/coverage/frontend.txt
