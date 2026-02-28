// Login & JWT logic
package main

import (
	    "fmt"

	"net/http" // HTTP status codes use panna
	"time"     //token expiry set panna

	"github.com/gin-gonic/gin" //API handle panna
	"github.com/golang-jwt/jwt/v5" //JWT create & sign panna
	"golang.org/x/crypto/bcrypt" //password hash panna
	    "gorm.io/gorm"

)

var jwtSecret = []byte("supersecretkey")

func Register(c *gin.Context) {
 var input struct {
        Name     string `json:"name"`  
        Email    string `json:"email"`
        Password string `json:"password"`
    }	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Check if email exists
	var existing User
	DB.Where("email = ?", input.Email).First(&existing)
	if existing.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := User{
		Name: input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	DB.Create(&user)

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "user": gin.H{"id": user.ID, "email": user.Email}})
}

func Login(c *gin.Context) {
    var input User
    var user User

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    fmt.Println("Login attempt:", input.Email) // ✅ Debug log

    if err := DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
            return
        }
        fmt.Println("DB error:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }

    err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
    if err != nil {
        fmt.Println("Password mismatch for user:", input.Email)
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
    })

    tokenString, err := token.SignedString(jwtSecret)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
