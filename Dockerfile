# syntax=docker/dockerfile:1

FROM golang:1.22 as build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /go-example-api

FROM gcr.io/distroless/static-debian12
COPY --from=build /go-example-api /
CMD ["/go-example-api"]

# CMD ["/go-example-api"]