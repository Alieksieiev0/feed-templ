generate:
	@templ generate

run:
	@templ generate
	@go run cmd/feed-templ/main.go
