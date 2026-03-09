SHELL := /bin/zsh

.PHONY: install-backend install-frontend build-backend build-frontend up down logs

install-backend:
	go mod tidy

install-frontend:
	cd ../frontend-repo && npm install

build-backend:
	go build ./...

build-frontend:
	cd ../frontend-repo && npm run build

up:
	docker compose up --build

down:
	docker compose down

logs:
	docker compose logs -f
