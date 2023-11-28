package controllers

import (
	"net/http"
	"project/web-service-gin/initializers"
	"project/web-service-gin/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func PostBanner(c *gin.Context) {
	var body struct {
		Image string
	}

	if c.Bind((&body)) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	admin := models.Banner{Image: body.Image}
	result := initializers.DB.Create(&admin)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messege": "success upload banner slide show",
	})

}

func GetBanner(c *gin.Context) {

	var banner []models.Banner
	result := initializers.DB.Find(&banner)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messege": "success get banner slide show",
		"data":    banner,
	})

}

func EditBanner(c *gin.Context) {

	var body struct {
		image string
		id    int
	}

	if c.Bind((&body)) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	var banner models.Banner
	result := initializers.DB.Model(&banner).Where("id = ?", body.id).Update("image", body.image)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to edit banner",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messege": "success update banner slide show",
	})

}

func DeleteBanner(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))

	resultUpdate := initializers.DB.Exec(`
			DELETE FROM banners 
			WHERE id = ?`, id)
	if resultUpdate.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete animal",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messege": "success delete banner slide show",
	})

}
