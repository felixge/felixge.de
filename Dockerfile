ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as gobuilder

WORKDIR /usr/src/app
RUN go install github.com/datadog/orchestrion@main
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /run-app .

FROM debian:bookworm

RUN apt update && apt install -y ca-certificates
COPY --from=gobuilder /run-app /usr/local/bin/
CMD ["run-app"]
