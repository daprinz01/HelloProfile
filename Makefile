migrationCreate:
	migrate create -ext sql -dir persistence/migrations -seq init_schema
migrationUp:
	migrate -path persistence/migrations -database "postgresql://postgres:Sarah4Daprinz@localhost:8669/persian_black?sslmode=disable" -verbose up
migrationDown:
	migrate -path persistence/migrations -database "postgresql://postgres:Sarah4Daprinz@localhost:8669/persian_black?sslmode=disable" -verbose down
installSqlc:
	go get github.com/kyleconroy/sqlc/cmd/sqlc
initialiseGoModules:
	go mod init authengine
.PHONY: migrationCreate migrationUp installSqlc initialiseGoModules