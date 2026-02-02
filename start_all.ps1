$env:PORT=8080
Start-Process -NoNewWindow -FilePath "go" -ArgumentList "run", "cmd/httpserver/main.go"
$env:PORT=8081
Start-Process -NoNewWindow -FilePath "go" -ArgumentList "run", "cmd/gqlserver/main.go"
$env:GRPC_PORT=8082
Start-Process -NoNewWindow -FilePath "go" -ArgumentList "run", "cmd/grpcserver/main.go"
Write-Host "All servers started!"
