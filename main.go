package main

// import "fmt"
// import "net/http"

import (
	"encoding/json"
	// "fmt"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	
	"github.com/go-sql-driver/mysql"
	"github.com/jinzu/gorm"
	// "strconv"
   )

type User struct {
	Nama string
	umur int
}

func PrintText(w http.ResponseWriter, r *http.Request) {
    text := "hellow wolrd"
	temp, _ := json.Marshal(text)

	w.Write(temp)
}

func PrintKata (c *gin.Context) {
	c.JSON(200, "hello world")
}

func PrintUser(c *gin.Context) {
	var orang User

	err := c.BindJSON(&orang)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "request invalid")
		return
	}

	orang.Umur = orang.Umur +5

	c.JSON(200, orang)
}

func main() {

	// fmt.Printf("Heloow Wolrd")

	// var angka int = 5
	// fmt.Println(angka)

	// var str string = "coba"
	// fmt.Println(str, "keren")

	// var flt float64 = 3.4
	// fmt.Println(flt)

	// var boolean bool = true
	// fmt.Println(boolean)

	// angka2 := 10
	// fmt.Println(angka2)

	
    // orang2 := User{
	// 	Nama: "Mely",
	// 	umur: 16,
	// }

	// fmt.Println(orang2.Nama, orang2.umur)

	http.HandleFunc("/", PrintText)
	// running di : localhost:9000
	log.Fatal(http.ListenAndServe(":9000", nil))

	app := gin.Default()
	app.GET("/", PrintKata)

	app.Run(":9000")
}