gqlgen:
	go run github.com/99designs/gqlgen generate

build:
	go build -o bin/play-with-go-graphql server.go
