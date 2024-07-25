SET GO111MODULE=on
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o .\build\worker\main cmd\worker\main.go
go build -o .\build\api\main cmd\api\main.go