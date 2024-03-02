# Go-bank

# Project requirements

1. Postgresql 
 
 ```
 docker run postgres
 docker exec -it <container id> psql -U postgres -d postgres
 ```


# Project outline
- Define Schema for Database 
- Write Database connection & setup database using docker
- Look for the way, where if we create a new user with new account, its initial account balance should be $0
- Setup the gin
- Define the routes
- Start with Handling Users
- Look for ways to integrate GRPC and protobuf into the project


# Future
- Create DB seeding and Add DB Dropping as well, for creating a fresh schema, if everything is okay!
