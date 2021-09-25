
service: FORCE
	go build -o service ./cmd

docker:
	docker build -t checker -f service.dockerfile .

up:
	docker-compose up -d

down:
	docker-compose down

FORCE: ;
