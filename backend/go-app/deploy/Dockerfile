from golang:1.24.2 as builder

workdir /app

copy go.mod go.sum ./
run go mod download

copy . .

run CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/app ./cmd/api

from  alpine:3.20

workdir /app

run apk add --no-cache ca-certificates

copy --from=builder /out/app ./api
copy --from=builder /app/deploy ./deploy
copy --from=builder /app/migrations ./migrations

expose 8080
cmd ["./api"]