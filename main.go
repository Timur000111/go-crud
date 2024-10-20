package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type book struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{Id: "1", Title: "В стране чудес", Author: "Лермонтов", Quantity: 2},
	{Id: "2", Title: "Мама мимо или мюзикл пошел не так", Author: "Онегин", Quantity: 4},
	{Id: "3", Title: "33 богатыря", Author: "Пушкин", Quantity: 1},
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func createBook(c *gin.Context) {
	var newBook book
	if err := c.BindJSON(&newBook); err != nil {
		return
	}
	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	book, err := getBookbyId(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}

	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not availabel."})
		return
	}

	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)
}

func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	book, err := getBookbyId(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}

	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)
}

func BookById(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookbyId(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return

	}
	c.IndentedJSON(http.StatusOK, book)
}

func getBookbyId(id string) (*book, error) {
	for i, b := range books {
		if b.Id == id {
			return &books[i], nil
		}
	}
	return nil, errors.New("Book not found")

}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", BookById)
	router.POST("/books", createBook)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)
	router.Run("localhost:8080")
}
