SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build ./src/main.go 
copy .\main.exe .\output\fromMcServerGetModConfig.exe
copy .\config.json .\output\config.json