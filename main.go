package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// defining data types
type RectangleInput struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}
type Input struct {
	Main  RectangleInput   `json:"main"`
	Input []RectangleInput `json:"input"`
}

type Rectangle struct {
	X      int       `json:"x"`
	Y      int       `json:"y"`
	Width  int       `json:"width"`
	Height int       `json:"height"`
	Time   time.Time `json:"time"`
}
type GetAllRectanglesOutput struct {
	Rectangles []Rectangle `json:"rectangels"`
}

// data base connection instance
var db *gorm.DB
var err error

func main() {
	// initiating db instance
	db, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	r := gin.Default()
	db.AutoMigrate(&Rectangle{})

	// routes
	r.GET("/", getAllRectangles)
	r.POST("/", SendRectangles)

	r.Run()
}

func SendRectangles(c *gin.Context) {
	var receiveTime time.Time = time.Now()

	var input Input
	c.BindJSON(&input)

	main := input.Main
	for _, r := range input.Input {
		if haveOverlap(main, r) {
			rectangleToSave := Rectangle{X: r.X, Y: r.Y, Height: r.Height, Width: r.Width, Time: receiveTime}
			db.Create(&rectangleToSave)
			fmt.Println("*** Saving Rectangle => ", rectangleToSave)
		}
	}
	c.JSON(200, main)
}

func getAllRectangles(c *gin.Context) {
	var rectangles []Rectangle
	if err := db.Find(&rectangles).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, rectangles)
	}
}

// checking if two rectangles have overlap
func haveOverlap(a, b RectangleInput) bool {
	if (a.X < (b.X + b.Width)) && ((a.X + a.Width) > b.X) && ((a.Y + a.Height) > b.Y) && (a.Y < (b.Y + b.Height)) {
		return true
	}
	return false
}
