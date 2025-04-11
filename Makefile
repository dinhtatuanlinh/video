
sqlc:
	sqlc generate
test:
	go test -v -cover -coverprofile=.coverage.html ./...

mock:
	mockgen -package mockdb -destination=db/mock/store.go simplebank/db/sqlc Store
	mockgen -package mockwk -destination worker/mock/distributor.go simplebank/worker TaskDistributor
proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb \
	--go-grpc_opt=paths=source_relative proto/*.proto

redis:
	docker run --name redis -p 6379:6379 --restart unless-stopped  -d redis:7-alpine

.PHONY: postgres createdb dropdb migrateup migratedown new_migration sqlc test mock proto redis