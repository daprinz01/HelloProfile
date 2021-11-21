migrationCreate:
	migrate create -ext sql -dir persistence/migrations -seq init_schema
migrationUp:
	migrate -path persistence/migrations -database "postgres://iuyegkoq:Zm032Nh7TJO_A_vifLUk8gX1R49YPEMe@drona.db.elephantsql.com/iuyegkoq?sslmode=disable" -verbose up
migrationDown:
	migrate -path persistence/migrations -database "postgres://iuyegkoq:Zm032Nh7TJO_A_vifLUk8gX1R49YPEMe@drona.db.elephantsql.com/iuyegkoq?sslmode=disable" -verbose down 
migrationForce:
	migrate -path persistence/migrations -database "postgres://iuyegkoq:Zm032Nh7TJO_A_vifLUk8gX1R49YPEMe@drona.db.elephantsql.com/iuyegkoq?sslmode=disable" -verbose force 1
migrationGoto:
	migrate -path persistence/migrations -database "postgres://iuyegkoq:Zm032Nh7TJO_A_vifLUk8gX1R49YPEMe@drona.db.elephantsql.com/iuyegkoq?sslmode=disable" -verbose goto 2
installSqlc:
	go get github.com/kyleconroy/sqlc/cmd/sqlc
initialiseGoModules:
	go mod init helloprofile
sqlcGenerate:
	sqlc generate
dockerBuild:
	docker build -t helloprofile:latest .
dockerRun:
	 docker run --mount source=hello-profile-logs,destination=/usr/local/bin/log/ -p 8083:8083 --name helloprofile --env TOKEN_LIFESPAN=96h --env SESSION_LIFESPAN=120h --env DB_HOST=drona.db.elephantsql.com --env DB_PORT=5432 --env COMMUNICATION_SERVICE_ENDPOINT=http://host.docker.internal:8084 helloprofile:latest
.PHONY: migrationCreate migrationUp migrationDown migrationForce migrationGoto installSqlc initialiseGoModules dockerRun sqlcGenerate dockerBuild