package controllers

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"project/web-service-gin/initializers"
	"project/web-service-gin/models"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}

func CreateShelter(c *gin.Context) {

	var body struct {
		Name        string
		Phone       string
		Description string
		Address     string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	user, _ := c.Get("user")
	shelter := models.Shelter{
		Name:        body.Name,
		Phone:       body.Phone,
		Description: body.Description,
		Address:     body.Address,
		Status:      "pending",
		UserID:      int(user.(models.User).ID),
	}

	result := initializers.DB.Create(&shelter)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create shelter",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messege": "Success to create shelter",
	})
}

func GetAllShelter(c *gin.Context) {
	var shelters []models.Shelter
	result := initializers.DB.Find(&shelters)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"messege": "Failed to retrieve data shelter",
		})
	}
	var respone []models.GetShelter
	// initializers.DB.Raw("SELECT id, name, age FROM shelter WHERE id = ?", 3).Scan(&shelters)
	initializers.DB.Raw(` 
		SELECT 
			a.id, 
			a.name as name, 
			a.phone as phone, 
			a.description as description, 
			a.address as address, 
			a.status as status, 
			b.id as user_id, 
			b.email as email_user, 
			b.username as owner_name 
		FROM shelters a 
		JOIN users b 
		ON a.user_id = b.id`).Scan(&respone)

	// respone := shelters

	c.JSON(http.StatusOK, gin.H{
		"messege": "success",
		"data":    respone,
	})

}

func ApprovalShelter(c *gin.Context) {

	var body struct {
		Id     int
		Status string
		Note   string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	fmt.Print(body)
	// var shelter models.Shelter
	result := initializers.DB.Exec(" UPDATE shelters SET status = ?, note = ?, updated_at = ? WHERE id = ?", body.Status, body.Note, time.Now(), body.Id)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update table shelter",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messege": "Success to update shelter",
	})
}

func CreateCategory(c *gin.Context) {
	// single file

	var body struct {
		Name  string
		Image string
	}

	if c.Bind((&body)) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	categories := models.Categories{Name: body.Name, Image: body.Image}
	resultInsert := initializers.DB.Create(&categories)

	if resultInsert.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create category",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		// "image": result.Location,
		"messege": "success create category",
	})
}

func GetAllCategories(c *gin.Context) {
	var categories []models.Categories
	result := initializers.DB.Find(&categories)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"messege": "Failed to retrieve data shelter",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"messege": "success",
		"data":    categories,
	})

}

func UploadFile(c *gin.Context) {
	// single file

	file, err := c.FormFile("image")

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"messege": "Failed to upload image",
		})
		return
	}

	f, openErr := file.Open()

	if openErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"messege": "Failed to open image",
		})
		return
	}

	// s3 bucket store file
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)

	result, upoadErr := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("animal-space-img"),
		Key:    aws.String(randomString(20) + file.Filename),
		Body:   f,
		ACL:    "public-read",
	})
	// err = c.SaveUploadedFile(file, "assets/"+file.Filename)

	if upoadErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"messege": "Failed to upload image",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		// "image": result.Location,
		"messege":  "success create category",
		"file_url": result.Location,
	})
}
