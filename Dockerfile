FROM golang:1.17.8-buster

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN go build -o /imlogs

EXPOSE 8081

CMD ["/imlogs"]