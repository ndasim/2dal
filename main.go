package main

import (
	"2dal/core"
	"2dal/shortener"
	"2dal/shortener/db"
	"2dal/shortener/middleware"
	"2dal/shortener/models"
	"fmt"
	"math/big"
	"time"

	"github.com/catinello/base62"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/kamva/mgm/v3"
	"golang.org/x/exp/slices"
)

type Book struct {
	// DefaultModel adds _id, created_at and updated_at fields to the Model
	mgm.DefaultModel `bson:",inline"`
	Name             string `json:"name" bson:"name"`
	Pages            int    `json:"pages" bson:"pages"`
}

func NewBook(name string, pages int) *Book {
	return &Book{
		Name:  name,
		Pages: pages,
	}
}

func main() {
	router := gin.Default()

	router.Use(middleware.Errors())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Static("/assets", "./assets")
	router.GET("/api/create", gin.Bind(shortener.CreateLinkStruct{}), shortener.CreateLink)
	router.GET("/api/qr", shortener.CreateQR)

	// Register custom validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("url", core.UrlValidator())
		v.RegisterValidation("isUnique", db.IsUniqueValidator())
	}

	router.Run("localhost:8080")

	//mi := i.Uint64()
	//fmt.Println(len(uniques))
	//fmt.Println(len(uniques))
	//fmt.Println(len(uniques))
}

func someTest() {
	a := big.NewInt(2103)

	m := big.NewInt(100000)

	i := big.Int{}
	i.ModInverse(a, m)

	list := map[int]int{1: 0, 2: 0, 3: 0, 4: 0, 5: 0, 6: 0, 7: 0}
	uniques := []string{}
	start_time := time.Now().UnixMilli()

	for i := 0; i < 90000; i++ {
		char := base62.Encode(i * 1881 % 100000)

		if slices.Contains(uniques, char) {
			list[7]++
		}

		uniques = append(uniques, char)
		list[len(char)]++

		if i%1000 == 0 {
			fmt.Println("took ", time.Now().UnixMilli()-start_time, " ms for, ", len(uniques), " items")
			start_time = time.Now().UnixMilli()
		}
	}

	fmt.Println(list)
}

func linkTest() {
	user := models.User{Username: "ndasim", Subscription: "sd"}

	link := models.Link{}
	link.Create("https://2d.al/", "", &user)

	link2 := &models.Link{}
	_ = mgm.Coll(link2).First(bson.M{}, link2)
}
