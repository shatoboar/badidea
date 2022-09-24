run: main.go
	go run main.go

build: main.go
	go build -o main main.go

docker-build: Dockerfile
	docker build -t slash/badidea .

docker-run: Dockerfile
	docker run -p 8080:8080 slash/badidea
