templ:
	@templ generate

tail:
	@npx tailwindcss -i ./static/css/input.css -o ./static/css/output.css

build: templ tail
	@go build -o feed-templ cmd/feed-templ/main.go

run: templ tail
	@go run cmd/feed-templ/main.go

clean:
	@echo "Cleaning..."
	@rm -f feed-templ

watch:
	@if command -v air > /dev/null; then \
	    air -c .air.toml; \
	    echo "Watching...";\
	else \
	    read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/cosmtrek/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

.PHONY: templ tail build run clean watch
