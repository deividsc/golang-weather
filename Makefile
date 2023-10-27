test:
	go test ./... --coverprofile cover.out

start:
	go run cmd/web/app.go