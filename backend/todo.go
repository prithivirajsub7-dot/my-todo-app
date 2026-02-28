package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateTodo creates a new todo for logged-in user
func CreateTodo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		Title string `json:"title"`
	}

	if err := c.ShouldBindJSON(&input); err != nil || input.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}

	todo := Todo{
		Title:  input.Title,
		UserID: userID.(uint),
	}

	// Save new todo
	if err := DB.Create(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
		return
	}

	// ✅ Return full todo object including id
	c.JSON(http.StatusOK, todo)
}

// GetTodos fetches todos of logged-in user
func GetTodos(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var todos []Todo
	if err := DB.Where("user_id = ? AND is_deleted = false", userID.(uint)).Find(&todos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todos"})
		return
	}

	c.JSON(http.StatusOK, todos)
}

// UpdateTodo toggles completed status or updates title
func UpdateTodo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := c.Param("id")
	var todo Todo
	if err := DB.Where("id = ? AND user_id = ?", id, userID.(uint)).First(&todo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find todo"})
		return
	}

	var input struct {
		Title     string `json:"title"`
		Completed *bool  `json:"completed"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if input.Title != "" {
		todo.Title = input.Title
	}
	if input.Completed != nil {
		todo.Completed = *input.Completed
	}

	if err := DB.Save(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// DeleteTodo marks todo as deleted
func DeleteTodo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := c.Param("id")
	var todo Todo
	if err := DB.Where("id = ? AND user_id = ?", id, userID.(uint)).First(&todo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find todo"})
		return
	}

	todo.IsDeleted = true
	if err := DB.Save(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete todo"})
		return
	}

	// ✅ Return confirmation message
	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted", "id": todo.ID})
}