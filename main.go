package main

import (
	"fmt"
	"log"
	"github.com/adarsh-jaiss/go-bank/db"
)

func main()  {
	err := db.Connect()
	if err!= nil{
		fmt.Errorf("Error connecting to database %s",err)
	}
	
	fmt.Println("Database connected successfully")
	
	defer db.Disconnect()

	if err := db.CreateTable(); err != nil {
		log.Fatal("Error creating table schema",err)
	}

}

