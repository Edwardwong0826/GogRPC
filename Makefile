#install on widnwos use Cygwin or choco for make command 
# if copy the command run in terminal, remember cd to C:\Git\GogRPC> to run it
server:
	go run cmd/server/main.go -port 8080

client:
	go run cmd/client/main.go -address 0.0.0.0:8080

.PHONY: gen clean server client test
