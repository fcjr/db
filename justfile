default:
    @just --list

dev addr=":4000":
    @LISTEN_ADDR={{addr}} go run cmd/db/main.go

test:
    @go test ./...

collection:
    @posting --collection ./collection