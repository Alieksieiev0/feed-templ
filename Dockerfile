FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /feed-templ cmd/feed-templ/main.go
CMD ["/feed-templ"]
