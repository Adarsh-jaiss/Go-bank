# Go-bank

# Frequently asked questions

### Why i chosed making use of interfaces if i could have done this in simple functional way?

I chosed this approach because this gives me an edge in the scenerios where i need to migrate the database, in that case, i don't need to write and every file. i'll just use the Storer and I'll just add few configurations and i am good to go. Also the functional way doesn't follow the design pattern of separating concerns and creating reusable components, which is why i started looking into using interfaces and constructors.

The code I would have written the 1st apprroach directly interacts with the database (sql.DB), executes the query, and handles the result all within the same function. This approach works fine for simple scenarios, but as your application grows, you might find it harder to maintain and test.

Using interfaces and constructors, as i'm doing in this project, allows me to decouple the database logic from the business logic. This separation makes my code more modular, easier to test, and promotes better code organization.

# Project requirements

1. Postgresql 
 
 ```
 docker run postgres
 docker exec -it <container id> psql -U postgres -d postgres
 ```


# Project outline/ TODO
- Define Schema for Database 
- Write Database connection & setup database using docker
- Look for the way, where if we create a new user with new account, its initial account balance should be $0
- Setup the gin
- Define the routes
- Start with Handling Users
- After making User Handlers , Implement User Authentication
- Then, Implement User Authorization
- 


# Future
- Create DB seeding and Add DB Dropping as well, for creating a fresh schema, if everything is okay!
- Look for ways to integrate GRPC and protobuf into the project
