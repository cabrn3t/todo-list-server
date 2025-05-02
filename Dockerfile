FROM golang:1.24.2-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG COPY_ENV=false
RUN if [ "$COPY_ENV" = "true" ]; then cp config.yaml /app/config.yaml; fi

RUN go build -o todo-list-server ./cmd/main.go

EXPOSE 8080

CMD ["./todo-list-server"]