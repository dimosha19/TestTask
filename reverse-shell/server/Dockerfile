FROM golang:1.19

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o server

EXPOSE 8081
CMD ["/app/server"]