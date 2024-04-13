templ:
	@templ generate

tail:
	@npx tailwindcss -i ./static/css/input.css -o ./static/css/output.css

run:
	@templ generate
	@go run cmd/feed-templ/main.go

air:
	@air -c .air.toml
