package api

import (
	"BookAPI/pkg/database"
	"bookAPI"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strconv"
)

func HealthCheck(repo *database.Repository) gin.HandlerFunc {
	return func(c *gin.Context){
		err := repo.Ping()
		if err != nil {
			fmt.Println("Error connecting to database: ", err)
			c.String(http.StatusInternalServerError, "Error connecting to the database")
			return
		}
		c.String(http.StatusOK, "ok")
	}
}

func Get(repo bookAPI.Repository) gin.HandlerFunc {
	return func(c *gin.Context){
		var bookId int
		var year int
		var err error
		if bookIdString := c.Query("bookId"); bookIdString != ""{
			bookId, err = strconv.Atoi(bookIdString)
			if err != nil {
				fmt.Println(err)
				c.String(http.StatusBadRequest, "bookId query must be a nonzero positive integer")
				return
			}
		}

		if yearString := c.Query("year"); yearString != ""{
			year, err = strconv.Atoi(yearString)
			if err != nil {
				fmt.Println(err)
				c.String(http.StatusBadRequest, "year query must be a nonzero positive integer")
				return
			}
		}

		search := bookAPI.Book{
			BookId: int32(bookId),
			Title: c.Query("title"),
			Author: c.Query("author"),
			Year: int32(year),
		}

		output, err := repo.GetBooks(search)
		if err != nil{
			fmt.Println(err)
			c.String(http.StatusInternalServerError, "An error occurred processing this request")
			return
		}

		if len(output.Books) == 0 {
			c.String(http.StatusNotFound, "No books found")
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

func GetById(repo bookAPI.Repository) gin.HandlerFunc {
	return func(c *gin.Context){
		id := c.Param("id")

		output, err := repo.GetBookById(id)
		if err != nil {
			fmt.Println(err)
			c.String(http.StatusInternalServerError, "An error occurred processing this request")
			return
		}
		if output == nil {
			c.String(http.StatusNotFound, "No book found with that id")
			return
		}

		c.Header("Location", "/" + url.PathEscape(id))
		c.JSON(http.StatusOK, output)
	}
}

func Post(repo bookAPI.Repository) gin.HandlerFunc {
	return func(c *gin.Context){
		var u = bookAPI.Book{}
		err := json.NewDecoder(c.Request.Body).Decode(&u)
		if err != nil{
			fmt.Println(err)
			c.String(http.StatusBadRequest, "An error occurred decoding this book")
			return
		}

		id, err := repo.PostBook(&u)
		if err != nil {
			fmt.Println(err)
			c.String(http.StatusInternalServerError, "An error occurred processing this request")
			return
		}

		c.Header("Location", "/" + url.PathEscape(id))
		c.JSON(http.StatusCreated, u)
	}
}

func Put(repo bookAPI.Repository) gin.HandlerFunc {
	return func(c *gin.Context){
		id := c.Param("id")

		var u = bookAPI.Book{}
		err := json.NewDecoder(c.Request.Body).Decode(&u)
		if err != nil {
			fmt.Println(err)
			c.String(http.StatusBadRequest, "An error occurred decoding this book")
			return
		}

		updated, updatedBook, err := repo.PutBook(id, &u)
		if err != nil {
			fmt.Println(err)
			if err.Error() == "invalid id given" {
				c.String(http.StatusBadRequest, "Invalid id given")
				return
			}
			c.String(http.StatusInternalServerError, "An error occurred processing this request")
			return
		}

		c.Header("Location", "/" + url.PathEscape(id))

		if updated {
			c.JSON(http.StatusOK, updatedBook)
		} else {
			c.JSON(http.StatusCreated, updatedBook)
		}
	}
}

func Patch(repo bookAPI.Repository) gin.HandlerFunc {
	return func(c *gin.Context){
		id := c.Param("id")

		var update map[string]interface{}

		err := json.NewDecoder(c.Request.Body).Decode(&update)
		if err != nil {
			fmt.Println(err)
			c.String(http.StatusBadRequest, "An error occurred decoding this update")
			return
		}

		err = repo.PatchBook(id, update)
		if err != nil {
			fmt.Println(err)
			if err.Error() == "not found" {
				c.String(http.StatusNotFound, "No book with that id found to update")
				return
			} else if err.Error() == "invalid id given" {
				c.String(http.StatusBadRequest, "Invalid id given")
				return
			}
			c.String(http.StatusInternalServerError, "An error occurred processing this request")
			return
		}

		c.Header("Location", "/" + url.PathEscape(id))
		c.Writer.WriteHeader(http.StatusOK)
	}
}

func Delete(repo bookAPI.Repository) gin.HandlerFunc {
	return func(c *gin.Context){
		id := c.Param("id")

		err := repo.DeleteBook(id)
		if err != nil {
			if err.Error() == "invalid id given" {
				c.String(http.StatusBadRequest, "Invalid id given")
				return
			} else {
				c.String(http.StatusInternalServerError, "An error occurred processing this request")
			}
		}
	}
}
