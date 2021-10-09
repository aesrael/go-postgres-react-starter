# Build go binary.
build-go:
	cd ./server && make build

run: build-go
	cd server/cmd && ./main && cd ../../ui/ && yarn && yarn start

# Build images and run
docker-build:
	docker-compose build

docker-run:
	docker-compose --env-file ./server/.env up 