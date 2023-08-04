# Build stage
FROM golang:1.20-alpine3.18 AS builder
ARG GH_ACCESS_TOKEN
WORKDIR /app
COPY . .
RUN apk update && apk upgrade && apk add git
RUN git config --global url.https://${GH_ACCESS_TOKEN}@github.com/.insteadOf https://github.com
RUN go env -w GOPRIVATE=github.com/nisbeyim
RUN go build -o main main.go

# Run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY db/migrations ./db/migrations

EXPOSE 8050
CMD [ "/app/main" ]