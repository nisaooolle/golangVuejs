package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/gorm"
)

type Response struct {
	Success bool        `json:"success"`
	Result  interface{} `json:"result"`
}

func RequestResponseLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log the request details
		log.Printf("Request: %s %s", c.Request.Method, c.Request.URL.String())

		// Process the request
		c.Next()

		// Log the response details
		log.Printf("Response: %d %s", c.Writer.Status(), http.StatusText(c.Writer.Status()))
	}
}

func main() {

	app := gin.Default()
	app.Use(gin.Recovery())
	app.Use(RequestResponseLogger())

	domain := app.Group("api")
	domain.POST("/payment", Payment)
	domain.POST("/generate/jwt", GenerateToken)
	domain.POST("/validate/jwt", ValidateToken)

	fmt.Println(app.Run(":9000"))
}

type TokenGenerate struct {
	Username string `json:"username" binding:"required"`
}
type TokenValidate struct {
	Token string `json:"token" binding:"required"`
}

func GenerateToken(c *gin.Context) {
	var request TokenGenerate

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusOK, Response{
			Success: false,
			Result:  err.Error(),
		})
		return
	}

	secretKey := []byte("secret-key")

	// Create a new token
	token := jwt.New(jwt.SigningMethodHS256)

	temp := time.Now().Add(20 * time.Second).Unix()
	fmt.Println(temp)
	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = request.Username
	claims["exp"] = temp // Token expiration time

	// Sign the token with your secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Success: false,
			Result:  "failed generate token",
		})
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Result:  tokenString,
	})
}

func ValidateToken(c *gin.Context) {
	var request TokenValidate

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusOK, Response{
			Success: false,
			Result:  err.Error(),
		})
		return
	}

	token, err := jwt.Parse(request.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret-key"), nil
	})

	if err != nil {
		c.JSON(http.StatusOK, Response{
			Success: false,
			Result:  err.Error(),
		})
		return
	}

	// Check if the token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check token expiration
		expTime := time.Unix(int64(claims["exp"].(float64)), 0)
		if time.Now().After(expTime) {
			c.JSON(http.StatusOK, Response{
				Success: false,
				Result:  "Token is expired",
			})
			return
		}

		c.JSON(http.StatusOK, Response{
			Success: true,
			Result:  "Token is valid",
		})
		log.Println("SUCESS VALIDATION")
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: false,
		Result:  "Token is not valid",
	})
}

type Product struct {
	Index int `json:"index" binding:"required"`
	Qty   int `json:"qty" binding:"required"`
}

func (p *Product) TableName() string {
	return "history"
}

func Payment(c *gin.Context) {
	var request []Product
	token := c.GetHeader("Authorization")

	if token != "david" {
		c.JSON(http.StatusForbidden, Response{
			Success: false,
			Result:  "Bad Request, Please Try Again",
		})
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusOK, Response{
			Success: false,
			Result:  "Bad Request, Please Try Again",
		})
		return
	}

	dsn := "root:@tcp(localhost:3306)/hicolleagues"

	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	for _, v := range request {
		err = db.Omit("umur").Create(&v).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, "failed save data")
			return
		}
	}

	body, _ := json.Marshal(request)

	fmt.Println(string(body))

	c.JSON(200, Response{
		Success: true,
	})
}