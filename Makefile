build:
	@go build -o bin/api
run : build
	@./bin/api 

drop:
	@go run ./seeding/drop/drop.go

seed:
	@go run ./seeding/seed.go