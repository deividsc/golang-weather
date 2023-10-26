test:
	go test ./... --coverprofile cover.out

start:
	go run src/app.go