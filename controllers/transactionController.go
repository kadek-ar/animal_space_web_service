package controllers

import (
	"net/http"
	"project/web-service-gin/initializers"
	"project/web-service-gin/models"
	"strconv"

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
	total_price := c.Query("total_price")
	number_of_item := c.Query("number_of_item")
	total, _ := strconv.Atoi(total_price)
	item_total, _ := strconv.Atoi(number_of_item)
	user, _ := c.Get("user")
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	transaction := models.Transaction{
		Note:         "",
		Status:       "active",
		Total:        total,
		NumberOfItem: item_total,
		UserID:       int(user.(models.User).ID),
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

func GetTransactionByUser(c *gin.Context) {
	var transaction []models.GetUserTransaction
	user, _ := c.Get("user")
	// result := initializers.DB.Where("user_id = ?", int(user.(models.User).ID)).Order("created_at DESC").Find(&transaction)

	result := initializers.DB.Raw(`
	SELECT 
		ta.transaction_id,
		t.status,
		t.created_at,
		COUNT(a.name) as animal_count,
		SUM(CASE WHEN ta.status = 'approve' THEN 1 ELSE 0 END) as approve_count,
		SUM( ta.quantity * ta.price ) as total,
		u.id as user_id
	FROM transactions t
	JOIN transaction_animals ta
		on t.id = ta.transaction_id
	JOIN animals a 
		on ta.animal_id = a.id 
	JOIN users u
		on t.user_id = u.id
	WHERE u.id = ?
	GROUP BY ta.transaction_id
	ORDER BY t.created_at DESC`, int(user.(models.User).ID)).Scan(&transaction)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"messege": "Failed to retrieve data transaction",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messege": "success",
		"data":    transaction,
	})
}

func GetDetailTransaction(c *gin.Context) {
	// var transaction_animals []models.TransactionAnimal
	var res_transaction_detail []models.GetTransactionDetail
	id, _ := strconv.Atoi(c.Param("id"))
	// result := initializers.DB.Where("transaction_id = ?", id).Find(&transaction_animals)
	result := initializers.DB.Raw(`
		SELECT 
			ta.transaction_id, 
			ta.animal_id, 
			ta.note, 
			ta.images, 
			ta.quantity, 
			ta.status,
			a.name as animal_name, 
			a.gender as animal_gender, 
			a.type as animal_type, 
			a.description as animal_description, 
			a.image as animal_image, 
			a.price as animal_price,
			c.name as animal_category,
			s.id as shelter_id, 
			s.name as shelter_name, 
			s.phone as shelter_phone
		FROM transaction_animals ta 
		JOIN animals a 
			on ta.animal_id = a.id 
		JOIN categories c
			on a.category_id = c.id 
		JOIN shelters s 
			on a.shelter_id = s.id 
		WHERE transaction_id = ?`, id).Scan(&res_transaction_detail)

	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"messege": "Failed to retrieve data transaction detail",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messege": "success",
		"data":    res_transaction_detail,
	})
}

func PostReceipt(c *gin.Context) {
	// var transaction_animals []models.TransactionAnimal
	var body models.GetTransactionDetail

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	result := initializers.DB.Exec(` 
		UPDATE transaction_animals 
		SET 
			status = ?, 
			images = ?
		WHERE transaction_id = ? AND animal_id = ? 
		`, "pending", body.Images, body.TransactionID, body.AnimalID)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"messege": "Failed to update transaction detail",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messege": "success send receipt",
	})
}
