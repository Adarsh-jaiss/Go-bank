package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"os"

	// "github.com/adarsh-jaiss/go-bank/api"
	"github.com/adarsh-jaiss/go-bank/store"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(userstore store.UserStorer) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			fmt.Println("Token not present in the header")
			c.JSON(http.StatusUnauthorized ,gin.H{"error": "Token not present in the header"})
		}
		fmt.Println("Token found in the header")

		claims, err := ValidateTokens(tokenStr)
		if err != nil {
			fmt.Println("code fat rha hai bhai!!!!!!!!!:")
			c.JSON(http.StatusUnauthorized , gin.H{"error": "Invalid token"})
			return
		}

		fmt.Println("token validated")

		expiresFloat := claims["expires"].(float64)
		expires := int64(expiresFloat)

		// check token expiration
		fmt.Println(expires)

		if time.Now().Unix() > expires {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
			return
		}

		// email, ok := claims["email"].(string)
		// if !ok {
		// 	c.JSON(401, gin.H{"error": "Invalid token"})
		// 	return
		// }

		accountNumber := claims["account_number"].(uint64)
		fmt.Println("Account Number:", accountNumber) // Debug statement

		user, err := userstore.GetUserByAccountNumber(c, accountNumber)
		if err != nil {
			fmt.Println("Failed to get user by account number:", err) // Debug statement
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// set the current authenticated user to the context
		c.Set("user", user)

		c.Next()
	}
}

func ValidateTokens(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method: %v", token.Header["alg"])
		}

		secret := os.Getenv("JWT_SECRET")
		fmt.Println("NEVER PRINT SECRET",secret)
		if secret == "" {
			return nil, errors.New("JWT secret key not found")
		}
		return []byte(secret), nil
	})

	fmt.Println("token parsed! ---------------->>>>>>>>")
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT token: %v", err)
	}

	fmt.Println("token parsed!")
	fmt.Println("token:",token)
	if !token.Valid {
		fmt.Println("token error idhar hai bhaiiiii! ---------------->>>>>>>>")
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	fmt.Println("token validated!, here's the claim :", claims)
	return claims, nil
}
