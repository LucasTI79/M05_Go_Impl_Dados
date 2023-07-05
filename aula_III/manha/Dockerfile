FROM golang:latest
WORKDIR /app
COPY go.mod go.mod
RUN go mod tidy
COPY . .
CMD ["go","run","."]