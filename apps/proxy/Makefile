dev:
	deno run -A dev.ts
	
run:
	go run main.go serve
	
build:
	go build -o pocketbase main.go

docker-build:
	docker build . --platform=linux/amd64 -t huakunshen/view-analytics-proxy:latest
