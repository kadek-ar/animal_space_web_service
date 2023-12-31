package middleware

import (
	"fmt"
	"net/http"
	"os"
	"project/web-service-gin/initializers"
	"project/web-service-gin/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {
	// Get the cookie off req
	// tokenString, err := c.Cookie(("Authorization"))
	tokenString := c.Request.Header["Authorization"][0]
	splitToken := strings.Split(tokenString, "Bearer ")
	tokenString = splitToken[1]

	fmt.Println(tokenString)
	// if err != nil {
	// 	c.AbortWithStatus(http.StatusUnauthorized)
	// }

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET")), nil
	})
	fmt.Println("---")
	fmt.Println(err)
	fmt.Println("---")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		//check exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// Find the user with toke sub
		var user models.User
		initializers.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		//attach to req

		c.Set("user", user)

		//continue
		c.Next()

		fmt.Println(claims["foo"], claims["nbf"])
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
