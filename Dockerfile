FROM golang:1.21 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

COPY proto proto

RUN CGO_ENABLED=0 GOOS=linux go build -o /chekcer

# Deploy the application binary into a lean image
FROM alpine

COPY --from=build-stage /chekcer /chekcer

ENTRYPOINT ["/chekcer"]