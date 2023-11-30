package main

import (
	"net/http"
	"project/web-service-gin/controllers"
	"project/web-service-gin/initializers"
	"project/web-service-gin/middleware"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `josn:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func postAlbums(c *gin.Context) {
	var newAlbum album

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	// DatabaseConnection()
	router := gin.Default()
	// router.Use(cors.Default())
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.Use(CORSMiddleware())
	router.GET("/albums", getAlbums)
	router.POST("/albums", postAlbums)
	router.POST("/signup", controllers.Signup)
	router.POST("/login", controllers.Login)
	router.GET("/validate", middleware.RequireAuth, controllers.Validate)
	router.GET("/user", middleware.RequireAuth, controllers.GetUserLogin)

	router.POST("/shelter/create", middleware.RequireAuth, controllers.CreateShelter)
	router.GET("/shelter", middleware.RequireAuth, controllers.GetAllShelter)
	router.GET("/shelter/:id", controllers.GetShelter)
	router.PUT("/shelter/:id", middleware.RequireAuth, controllers.EditShelter)
	router.PUT("/shelter/approval", middleware.RequireAuth, controllers.ApprovalShelter)

	router.GET("/admin/transaction", middleware.RequireAuth, controllers.GetTransactionAdmin)
	router.GET("/admin/transaction/:shelter_id/:id", middleware.RequireAuth, controllers.GetAdminDetailTransaction)
	router.GET("/admin/animal", middleware.RequireAuth, controllers.GetAllAnimalAdmin)

	router.POST("/shelter/category", middleware.RequireAuth, controllers.CreateCategory)
	router.GET("/shelter/category", controllers.GetAllCategories)
	router.PUT("/shelter/category/:id", middleware.RequireAuth, controllers.EditCategory)
	router.DELETE("/shelter/category/:id", middleware.RequireAuth, controllers.DeleteCategory)

	router.GET("/shelter/animal/:id", middleware.RequireAuth, controllers.GetShelterAnimal)
	router.GET("/shelter/transaction/:id", middleware.RequireAuth, controllers.GetShelterTransaction)
	router.GET("/shelter/transaction/detail/:shelter_id/:id", middleware.RequireAuth, controllers.GetShelterDetailTransaction)
	router.PUT("/shelter/transaction/detail/approval", middleware.RequireAuth, controllers.PostApprovalReceipt)

	router.POST("/upload", middleware.RequireAuth, controllers.UploadFile)

	router.POST("/admin/banner", middleware.RequireAuth, controllers.PostBanner)
	router.GET("/admin/banner", controllers.GetBanner)
	router.PUT("/admin/banner", middleware.RequireAuth, controllers.EditBanner)
	router.DELETE("/admin/banner/:id", middleware.RequireAuth, controllers.DeleteBanner)

	router.POST("/animal", middleware.RequireAuth, controllers.CreateAnimal)
	router.GET("/animal", controllers.GetAllAnimalByShelter)
	router.GET("/animal/:id", middleware.RequireAuth, controllers.GetAnimal)
	router.PUT("/animal/:id", middleware.RequireAuth, controllers.UpdateAnimal)
	router.DELETE("/animal/:id", middleware.RequireAuth, controllers.DeleteAnimal)

	router.GET("/animal-space", controllers.GetAllAnimal)
	router.GET("/animal-space/:id", controllers.GetSingelAnimal)
	router.POST("/animal-space/cart", middleware.RequireAuth, controllers.PostCart)
	router.DELETE("/animal-space/cart/:id", middleware.RequireAuth, controllers.DeleteCart)
	router.GET("/animal-space/cart", middleware.RequireAuth, controllers.GetCartByUser)
	router.POST("/animal-space/checkout", middleware.RequireAuth, controllers.PostCheckout)
	router.GET("/animal-space/transaction", middleware.RequireAuth, controllers.GetTransactionByUser)
	router.GET("/animal-space/transaction/:id", middleware.RequireAuth, controllers.GetDetailTransaction)
	router.POST("/animal-space/transaction/receipt", middleware.RequireAuth, controllers.PostReceipt)

	router.Run("localhost:8081")
}
