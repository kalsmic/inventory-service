FROM golang:1.17-alpine
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o /inventory-service
EXPOSE 5000
CMD ["/inventory-service"]