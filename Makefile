#install on widnwos use Cygwin or choco for make command 
server:
	go run cmd/server/main.go -port 8080

client:
	go run cmd/client/main.go -address 0.0.0.0:8080