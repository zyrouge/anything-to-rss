FROM golang:1.23-alpine as build

WORKDIR /usr/app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go .
COPY internal internal
RUN go build -o ./dist/anything-to-rss .

FROM alpine

WORKDIR /usr/app

COPY --from=build /usr/app/dist/anything-to-rss anything-to-rss
CMD ["./anything-to-rss"]
