FROM golang:1.19 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

COPY proto proto

RUN CGO_ENABLED=0 GOOS=linux go build -o /chekcer

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /chekcer /chekcer

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/chekcer"]