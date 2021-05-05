install:
	docker build . -f ./backend-api -t backend-api
	docker build . -f ./backend-exec -t backend-exec

run:
	docker-compose up -d

stop:
	docker-compose stop

down:
	docker-compose down

logs:
	docker-compose logs -f