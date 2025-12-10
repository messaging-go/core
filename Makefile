generate:
	go tool mockgen -destination=./internal/test/mocks/mock_core.go -package=mocks -typed github.com/messaging-go/core MessageProcessor
	go tool mockgen -destination=./internal/test/mocks/mock_middleware.go -package=mocks -typed github.com/messaging-go/core Middleware
