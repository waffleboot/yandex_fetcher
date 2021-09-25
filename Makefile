
service: FORCE
	go build -o service ./cmd/service
	# go build -o checker ./cmd/checker

docker: docker/service.dockerfile docker/checker.dockerfile
	docker build -t service -f docker/service.dockerfile .
	# docker build -t checker -f docker/checker.dockerfile .

up:
	docker-compose up -d

down:
	docker-compose down

FORCE: ;
