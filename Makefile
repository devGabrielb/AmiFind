up:
	docker-compose -f docker-compose-dev.yml up -d

down:
	docker-compose -f docker-compose-dev.yml down -v

test:
	docker-compose run unittest

run-dev:
	go run cmd/main.go
