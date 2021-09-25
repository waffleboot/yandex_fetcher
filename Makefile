
service: FORCE
	go build -o service ./cmd

docker: docker/service.dockerfile docker/checker.dockerfile
	docker build -t service -f docker/service.dockerfile .
	docker build -t checker -f docker/checker.dockerfile .

up:
	docker-compose up -d

down:
	docker-compose down

FORCE: ;
