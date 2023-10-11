FROM golang:1.21 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

COPY proto proto

RUN CGO_ENABLED=0 GOOS=linux go build -o /chekcer

# Deploy the application binary into a lean image
FROM python:3-alpine

COPY requirements.txt .
RUN pip3 install -r requirements.txt

COPY --from=build-stage /chekcer /chekcer

COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh
ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["/chekcer"]