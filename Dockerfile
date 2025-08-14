FROM golang:1.25-alpine as builder

RUN apk update && apk add --no-cache gcc musl-dev ca-certificates

WORKDIR /app
COPY src/go.mod src/go.sum ./
RUN go mod download
COPY . .

ARG GIT_TAG
ARG GIT_COMMIT
ARG GIT_BRANCH

RUN BUILD_DATE=$(date +'%Y-%m-%dT%H:%M:%S%z') && \
    go build -ldflags "-X 'main.version=${GIT_TAG}' -X 'main.gitCommit=${GIT_COMMIT}' -X 'main.gitBranch=${GIT_BRANCH}' -X 'main.buildDate=${BUILD_DATE}'" -o main ./src

FROM alpine:3.22.1 as prod
WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 8080

CMD ["./main"]

FROM golang:1.25-alpine as dev

RUN apk update && apk add --no-cache gcc musl-dev ca-certificates

RUN go install github.com/air-verse/air@latest

WORKDIR /app/src
COPY src/go.mod src/go.sum ./
RUN go mod download
COPY . .

CMD ["air", "-c", ".air.toml"]