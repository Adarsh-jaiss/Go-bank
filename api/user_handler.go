package api

import (
	"github.com/adarsh-jaiss/go-bank/models"
	"github.com/adarsh-jaiss/go-bank/store"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	UserStore store.UserStorer
}

/*
Userhandler struct definition declares a struct type Userhandler with a single field userstore of type store.UserStorer.
This struct is intended to be used as a handler for user-related operations in your application.
The purpose of the constructor function NewUserHandler is to initialize a new Userhandler struct instance with the provided userstore.
This pattern allows you to encapsulate the creation of Userhandler instances and ensures that they are properly initialized before being used.
By accepting a userstore parameter of type store.UserStorer, the constructor function provides flexibility,
allowing different implementations of the UserStorer interface to be passed in, such as a database store or a mock store for testing purposes.
Overall, this design promotes modularity and testability by adhering to Go's principles of composition and dependency injection.
*/
func NewUserHandler(userstore store.UserStorer) *UserHandler {
	return &UserHandler{
		UserStore: userstore,
	}
}

func (h *UserHandler) HandlePostUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var newUser models.User
		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, err := h.UserStore.InsertUser(c,&newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
	}
}

func(h *UserHandler) HandleGetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user models.User
		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, err := h.UserStore.GetUser(ctx, &user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, res)
	}
}

