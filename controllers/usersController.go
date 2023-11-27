package controllers

import (
	"fmt"
	"net/http"
	"os"
	"project/web-service-gin/initializers"
	"project/web-service-gin/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	// Get the email/pass off req body

	var body struct {
		Email    string
		Password string
		Username string
	}

	if c.Bind((&body)) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	// hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 14)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// create the user
	user := models.User{Email: body.Email, Password: string(hash), Username: body.Username, Role: "user"}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func Login(c *gin.Context) {
	//Get the email and pass off req body
	var body struct {
		Email    string
		Password string
	}

	if c.Bind((&body)) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// look up requested user
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User Not Found",
		})
		return
	}

	//compare sent in pass with saved user pass hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	fmt.Println(err)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// generate jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	//send it back
	// c.SetSameSite(http.SameSiteLaxMode)
	// c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})

}

func GetUserLogin(c *gin.Context) {
	user, _ := c.Get("user")
	// user.(models.User).
	// var respone models.GetShelter
	fmt.Println("befor query")
	var shelter models.Shelter
	result := initializers.DB.First(&shelter, "user_id = ?", user.(models.User).ID)
	// result := initializers.DB.Raw("SELECT * FROM shelters WHERE user_id = ?", user.(models.User).ID).First(&shelter)
	fmt.Println("query")
	fmt.Println(shelter)
	if result.RowsAffected == 1 {
		// result := initializers.DB.Raw(`
		// 	SELECT
		// 		a.id,
		// 		a.name as name,
		// 		a.phone as phone,
		// 		a.description as description,
		// 		a.address as address,
		// 		a.status as status,
		// 		b.id as user_id,
		// 		b.email as email_user,
		// 		b.username as owner_name
		// 	FROM shelters a
		// 	JOIN users b
		// 	ON a.user_id = b.id`).First(&respone)

		// if result.Error != nil {
		// 	c.JSON(http.StatusOK, gin.H{
		// 		"messege": "Failed to retrieve data shelter",
		// 	})
		// }
		c.JSON(http.StatusOK, gin.H{
			"messege":        "success",
			"email":          user.(models.User).Email,
			"username":       user.(models.User).Username,
			"role":           user.(models.User).Role,
			"shelter_id":     shelter.ID,
			"shelter_status": shelter.Status,
			"shelter_name":   shelter.Name,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"messege":        "success",
			"email":          user.(models.User).Email,
			"username":       user.(models.User).Username,
			"role":           user.(models.User).Role,
			"shelter_id":     "",
			"shelter_status": "",
			"shelter_name":   "",
		})
	}

}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")
	// user.(models.User).

	c.JSON(http.StatusOK, gin.H{
		"messege": "i am login",
		"user":    user,
	})
}
