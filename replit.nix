run = "cd backend && go mod tidy && go run main.go"
language = "go"
entrypoint = "backend/main.go"
hidden = [".git", "node_modules", "vendor"]

[go]
version = "1.21"

[nix]
channel = "stable-23_05"

[deployment]
run = ["go", "run", "main.go"]
deploymentTarget = "cloudrun"
build = ["go", "build", "-o", "server", "main.go"]
port = 8080

[[ports]]
localPort = 8080
externalPort = 80

[[ports]]
localPort = 3000
externalPort = 3000