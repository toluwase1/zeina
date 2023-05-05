up:
	docker compose up --build

generate-mock:
	 mockgen -destination=mocks/auth_mock.go -package=mocks zeina/services AuthService
	 mockgen -destination=mocks/auth_repo_mock.go -package=mocks zeina/db AuthRepository
	 mockgen -destination=mocks/wallet_mock.go -package=mocks zeina/services WalletService
	 mockgen -destination=mocks/wallet_repo_mock.go -package=mocks zeina/db WalletRepository


test: generate-mock
	 MEDDLE_ENV=test go test ./...
