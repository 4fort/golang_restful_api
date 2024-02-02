package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type book struct {
	ID				string			`json:"id"`
	Title			string	`json:"title"`
	Author		string	`json:"author"`
	Quantity	int			`json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "Book 1", Author: "Author 1", Quantity: 10},
	{ID: "2", Title: "Book 2", Author: "Author 2", Quantity: 20},
	{ID: "3", Title: "Book 3", Author: "Author 3", Quantity: 30},
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func getBookById(c *gin.Context) {
	id := c.Param("id")
	book, err := bookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}

func bookById(id string) (*book, error) {
	for i, book := range books {
		if book.ID == id {
			return &books[i], nil
		}
	}

	return nil, errors.New("book not found")
}

func processBook(c *gin.Context, increaseQuantity bool) {
	id, ok := c.GetQuery("id")
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing book"})
		return
	}

	book, err := bookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	if increaseQuantity {
		book.Quantity += 1
	} else {
		if book.Quantity <= 0 {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book is out of stock"})
			return
		}
		book.Quantity -= 1
	}

	c.IndentedJSON(http.StatusOK, book)
}

func checkoutBook(c *gin.Context) {
	processBook(c, false)
}

func returnBook(c *gin.Context) {
	processBook(c, true)
}

func addBook(c *gin.Context) {
	var newBook book

	if err := c.BindJSON(&newBook); err != nil{
		return
	}

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	servHost := os.Getenv("SERVER_HOST")
	servPort := os.Getenv("SERVER_PORT")

	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", getBookById)
	router.POST("/books", addBook)
	router.PATCH("/checkout/", checkoutBook)
	router.PATCH("/return/", returnBook)
	router.Run(fmt.Sprintf("%s:%s", servHost, servPort))
}