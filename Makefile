up:
	docker compose up --build

generate-mock:
	 mockgen -destination=mocks/auth_mock.go -package=mocks zeina/services AuthService
	 mockgen -destination=mocks/auth_repo_mock.go -package=mocks zeina/db AuthRepository
	 mockgen -destination=mocks/wallet_mock.go -package=mocks zeina/services WalletService
	 mockgen -destination=mocks/wallet_repo_mock.go -package=mocks zeina/db WalletRepository


test: generate-mock
	 ZEINA_ENV=test go test ./...

soda-down:
	soda m down -s 10

soda-up:
	soda m up