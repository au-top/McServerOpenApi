SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build ./src/main.go 
copy .\main .\output\fromMcServerGetModConfig
copy .\config.json .\output\config.json