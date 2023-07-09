package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Book struct {
	ID    string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year int32  `json:"year"`
}


func HomepageHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to Golang LiveCode Jenkins Daffa"})
}

func NewBookHandler(c *gin.Context) {
	var newBook Book
	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	newBook.ID = uuid.New().String()
	db, err := DbConn()
	if err != nil {
		panic(err.Error())
	}
	db.Create(&newBook)
	c.JSON(http.StatusCreated, newBook)
}

func GetBookHandler(c *gin.Context) {
	db, err := DbConn()
	if err != nil {
		panic(err.Error())
	}
	var Books []Book
	db.Find(&Books)
	c.JSON(http.StatusOK, Books)
}

func UpdateBookHandler(c *gin.Context) {
	id := c.Param("id")
	var Book Book
	if err := c.ShouldBindJSON(&Book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	db, err := DbConn()
	if err != nil {
		panic(err.Error())
	}
	err = db.Where("id=?", id).First(&Book).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Book not found",
		})
		return
	}
	db.Save(&Book)
	c.JSON(http.StatusOK, Book)
}

func DeleteBookHandler(c *gin.Context) {
	id := c.Param("id")
	db, err := DbConn()
	if err != nil {
		panic(err.Error())
	}
	db.Where("id=?", id).Delete(&Book{})
	c.JSON(http.StatusOK, gin.H{
		"message": "Book has been deleted",
	})
}

func InitHandler(c *gin.Context) {
	var Books = []Book{
		{ID: "B001", Title: "Book A", Author: "Daffa", Year: 2002},
		{ID: "B002", Title: "Book B", Author: "Daffa", Year: 2002},
		{ID: "B003", Title: "Book C", Author: "Daffa", Year: 2002},
	}

	db, err := DbConn()
	if err != nil {
		panic(err.Error())
	}
	err = db.AutoMigrate(&Book{})
	if err != nil {
		panic(err.Error())
	}

	db.Create(&Books)
	c.JSON(http.StatusOK, gin.H{
		"message": "Init Book",
	})
}

func DbConn() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", host, user, password, dbName, port)
	// dsn := "host=localhost user=postgres password=admin dbname=livecode-cicd port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database")
	}
	return db, nil
}

func main() {

	router := gin.Default()
	router.GET("/", HomepageHandler)
	router.GET("/init", InitHandler)
	router.GET("/books", GetBookHandler)
	router.POST("/books", NewBookHandler)
	router.PUT("/books/:id", UpdateBookHandler)
	router.DELETE("/books/:id", DeleteBookHandler)
	err := router.Run(":8888")
	if err != nil {
		panic("Failed run server")
	}
}