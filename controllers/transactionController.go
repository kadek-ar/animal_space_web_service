package controllers

import (
	"net/http"
	"project/web-service-gin/initializers"
	"project/web-service-gin/models"

	"github.com/gin-gonic/gin"
)

func PostCart(c *gin.Context) {

	var body struct {
		AnimalID int
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	user, _ := c.Get("user")

	var cartTmp models.Cart
	resultCart := initializers.DB.Where("animal_id = ?", body.AnimalID).Where("user_id = ?", int(user.(models.User).ID)).First(&cartTmp)

	if resultCart.RowsAffected > 0 {
		initializers.DB.Model(&cartTmp).Where("animal_id = ?", body.AnimalID).Update("quantity", cartTmp.Quantity+1)
		c.JSON(http.StatusOK, gin.H{
			"messege": "Success to add to cart",
		})
	} else {
		cart := models.Cart{
			Note:     "",
			Image:    "",
			Quantity: 1,
			UserID:   int(user.(models.User).ID),
			AnimalID: body.AnimalID,
		}

		result := initializers.DB.Create(&cart)

		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to create cart",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"messege": "Success to add to cart",
		})
	}

}

func GetCartByUser(c *gin.Context) {
	var cart []models.Cart
	user, _ := c.Get("user")
	result := initializers.DB.Preload("Animal.Shelter").Where("user_id = ?", int(user.(models.User).ID)).Find(&cart)

	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"messege": "Failed to retrieve data cart",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messege": "success",
		"data":    cart,
	})

}

func PostCheckout(c *gin.Context) {
	var body []struct {
		AnimalID int
		Quantity int
		Price    int
	}
	user, _ := c.Get("user")
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	transaction := models.Transaction{
		Note:   "",
		Status: "active",
		UserID: int(user.(models.User).ID),
	}

	result := initializers.DB.Create(&transaction)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create transaction",
		})
		return
	}

	var trxFind models.Transaction

	resultLast := initializers.DB.Last(&trxFind)

	if resultLast.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to find new transaction ",
		})
		return
	}

	for _, b := range body {
		transactionDetail := models.TransactionAnimal{
			AnimalID:      b.AnimalID,
			TransactionID: trxFind.ID,
			Note:          "",
			Quantity:      b.Quantity,
			Price:         b.Price,
			Images:        "",
		}

		resultTrx := initializers.DB.Create(&transactionDetail)

		if resultTrx.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to create transaction detail animal id " + string(b.AnimalID),
			})
			return
		}

		// var cart models.Cart

		// resultCart := initializers.DB.Where("user_id = ?", int(user.(models.User).ID)).Where("animal_id = ? ", b.AnimalID).Delete(&cart)
		resultUpdate := initializers.DB.Exec(`
			DELETE FROM carts 
			WHERE user_id = ? AND animal_id = ?`,
			int(user.(models.User).ID), b.AnimalID)
		if resultUpdate.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to delete animal",
			})
			return
		}

	}

	c.JSON(http.StatusOK, gin.H{
		"messege": "success create new transaction",
	})
}
