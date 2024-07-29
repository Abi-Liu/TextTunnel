# Variables
SERVER_BINARY := chat_server
CLIENT_BINARY := chat_client

# Directories
SERVER_DIR := ./cmd/chat_server/
CLIENT_DIR := ./cmd/chat_client/texttunnel/

build-server:
	@echo "Building server..."
	@cd $(SERVER_DIR) && go build -o $(SERVER_BINARY)

build-linux:
	@echo "Building server for Linux..."
	@cd $(CLIENT_DIR) && GOOS=linux GOARCH=amd64 go build -o '../../../linux_binary'

build-client:
	@echo "Building client..."
	@cd $(CLIENT_DIR) && go build -o $(CLIENT_BINARY)

run-server: build-server
	@echo "Starting server..."
	@./$(SERVER_DIR)/$(SERVER_BINARY)

run-client: build-client
	@echo "Starting client..."
	@./$(CLIENT_DIR)/$(CLIENT_BINARY)

run-both: build-server build-client
	@echo "Starting server and client..."
	@./$(SERVER_DIR)/$(SERVER_BINARY) & ./$(CLIENT_DIR)/$(CLIENT_BINARY)

clean:
	@echo "Cleaning up..."
	@rm -f $(SERVER_DIR)/$(SERVER_BINARY) $(CLIENT_DIR)/$(CLIENT_BINARY)

# Default target
all: build-server build-client
