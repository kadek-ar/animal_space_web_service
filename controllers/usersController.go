package controllers

import (
	"bytes"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"project/web-service-gin/initializers"
	"project/web-service-gin/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

const charset2 = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.][-)(*&^%$#@!~)]?><1234567890"

func randomHash(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(charset2[rand.Intn(len(charset2))])
	}
	return sb.String()
}

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

	link := os.Getenv("BASE_URL_EMAIL") + "verify-email?idveas=" + randomHash(60)

	errEmail := sendEmail(body.Email, body.Username, "./assets/template/verifyEmail.html", link, "Verify your email to continue signup")

	if errEmail != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to send email",
		})
		return
	}

	// create the user
	user := models.User{Email: body.Email, Password: string(hash), Username: body.Username, Role: "user", Status: "pending"}
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

	if user.Status == "pending" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "this user need to verify email address, check email in this user",
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

func VerifyEmail(c *gin.Context) {
	hash := c.Query("idveas")

	err := VerifyHashUser(hash)

	if err != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messege": "email is valid",
	})
}

// func sendEmail(email string, username string, htmlTamplatePath string, link string, subjectTitle string) (err error) {

// 	var body bytes.Buffer
// 	t, err := template.ParseFiles(htmlTamplatePath)
// 	t.Execute(&body, struct{ Link string }{Link: link})

// 	from := mail.NewEmail("Animal Space Admin", "Shelterspace27@gmail.com")
// 	subject := subjectTitle
// 	to := mail.NewEmail(username, email)
// 	plainTextContent := "Verify you email"
// 	htmlContent := body.String()
// 	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
// 	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_KEY"))
// 	response, err := client.Send(message)
// 	if err != nil {
// 		return err
// 	} else {
// 		fmt.Println(response.StatusCode)
// 		fmt.Println(response.Body)
// 		fmt.Println(response.Headers)
// 	}

// 	return
// }

func sendEmail(email string, username string, htmlTamplatePath string, link string, subjectTitle string) (err error) {

	var body bytes.Buffer
	t, err := template.ParseFiles(htmlTamplatePath)
	t.Execute(&body, struct{ Link string }{Link: link})

	// from := mail.NewEmail("Animal Space Admin", "Shelterspace27@gmail.com")
	// subject := subjectTitle
	// to := mail.NewEmail(username, email)
	// plainTextContent := "Verify you email"
	htmlContent := body.String()
	// message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	// client := sendgrid.NewSendClient(os.Getenv("SENDGRID_KEY"))
	// response, err := client.Send(message)
	// if err != nil {
	// 	return err
	// } else {
	// 	fmt.Println(response.StatusCode)
	// 	fmt.Println(response.Body)
	// 	fmt.Println(response.Headers)
	// }

	from := "shelterspace27@gmail.com"

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subjectTitle)
	m.SetBody("text/html", htmlContent)
	// m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer("smtp.gmail.com", 587, from, "zkfcemkkhqbvesf")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		return err
	}

	return
}

func VerifyHashUser(hash string) (err string) {
	var user models.User
	initializers.DB.First(&user, "hash = ?", hash)

	if user.ID == 0 {
		return "hash not found"
	}

	resultUpdate := initializers.DB.Exec(`
		UPDATE users 
		SET 
			status = 'valid',
			hash = NULL
		WHERE id = ?`, user.ID,
	)
	if resultUpdate.Error != nil {
		return "error to update table users"
	}

	return ""
}

func RequestResetPassword(c *gin.Context) {

	var body struct {
		Email string
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

	codeHash := randomHash(60)

	link := os.Getenv("BASE_URL_EMAIL") + "reset-password?idveas=" + codeHash
	resultUpdate := initializers.DB.Exec(`
		UPDATE users 
		SET 
			hash = ?
		WHERE id = ?`, codeHash, user.ID,
	)
	if resultUpdate.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error to send email",
		})
		return
	}

	errEmail := sendEmail(user.Email, user.Username, "./assets/template/forgetPass.html", link, "Request to reset password")

	if errEmail != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to send email",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messege": "success send email",
	})
}

func ResetPassword(c *gin.Context) {
	var body struct {
		Password string
		Hash     string
	}

	if c.Bind((&body)) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(body.Password), 14)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	var user models.User
	initializers.DB.First(&user, "hash = ?", body.Hash)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error to get user",
		})
		return
	}

	resultUpdate := initializers.DB.Exec(`
		UPDATE users 
		SET 
			password = ?,
			hash = NULL
		WHERE id = ?`, pass, user.ID,
	)
	if resultUpdate.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error to send email",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messege": "reset password",
	})

}
