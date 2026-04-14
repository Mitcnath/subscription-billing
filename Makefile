deploy:
	docker compose build
	docker compose up -d

deploy-seeded:
	docker compose build
	docker compose up -d
	docker exec subscription-billing-backend-1 ./seed

down:
	docker compose down