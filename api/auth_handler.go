package api

import (
	"errors"
	"net/http"
	"os"
	"time"

	"fmt"

	"github.com/adarsh-jaiss/go-bank/models"
	"github.com/adarsh-jaiss/go-bank/store"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	userStore store.UserStorer
}

func NewAuthHandler(userStore store.UserStorer) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  models.User `json:"user"`
	Token string      `json:"token"`
	Error string      `json:"error,omitempty"`
}

func (h *AuthHandler) HandleAuthentication(c *gin.Context) {
	var authParams AuthParams
	if err := c.ShouldBindJSON(&authParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return //returning from the function
	}

	fmt.Println(authParams)

	user, err := h.userStore.GetUserByEmail(c, authParams.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("email found!")
	fmt.Println("auth params password :" + authParams.Password)
	fmt.Println("user password :" + user.Password)

	if user.Password != authParams.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials, password mismatch"})
		return
	}

	fmt.Println("password matched!")

	token, err := CreateTokenFromUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("token created!")

	res := AuthResponse{
		User:  *user,
		Token: token,
	}

	c.JSON(http.StatusOK, res)
	fmt.Println("response sent!")
	// Return statement is not needed here, as the function will naturally exit after this point.
}

func CreateTokenFromUser(user *models.User) (string, error) {
	now := time.Now()
	expires := now.Add(time.Minute * 20).Unix()
	claims := jwt.MapClaims{
		// "id":      user.ID,
		"email":   user.Email,
		"expires": expires,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT secret key not found")
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token with secret: %v", err)
	}

	fmt.Println("token from auth_handler.go: ", tokenString)
	return tokenString, nil
	
}
