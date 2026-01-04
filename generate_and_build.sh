go run internal/main.go
cd test/out
go get golang.org/x/tools/cmd/goimports
go run golang.org/x/tools/cmd/goimports -w .
go build