install:
	docker build . -f ./backend-api -t backend-api
	docker build . -f ./backend-exec -t backend-exec
	docker build . -f ./backend-telegram -t backend-telegram

run:
	docker-compose up -d

stop:
	docker-compose stop

down:
	docker-compose down

logs:
	docker-compose logs -f

reload:
	docker-compose down
	docker build . -f ./backend-api -t backend-api
	docker build . -f ./backend-exec -t backend-exec
	docker build . -f ./backend-telegram -t backend-telegram
	docker-compose up -d
	docker-compose logs -f