package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type NASAimg struct{
	ID int  `json:"id"`
	Title string
 	URL string `json:"url"`         // Image link
    Author    string `json:"author"`
    CreatedAt string `json:"created_at"`
}

func handlerhome(c *fiber.Ctx)error{
	var images [] NASAimg
	if err := db.Find(&images).Error; err!=nil{
		return c.Status(500).SendString("Failed to fetch images!")
	}

	return c.Status(fiber.StatusCreated).JSON(images)

}

func main(){
	app := fiber.New()
	app.Get("/",handlerhome)
	fmt.Println("Hello world this is for NASA image library")

	host := "localhost"
	user := "postgres"
	password := "yipikaye2123"
	port := 5432
	dbname := "valinenasa" //using database this
	sslmode := "disable"


	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
    host, user, password, dbname, port, sslmode,)

	var errConnectDb error

	db, errConnectDb = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if errConnectDb != nil {
		fmt.Println("Failed to connect to database !")
	}

	
	db.AutoMigrate(&NASAimg{})
	fmt.Println("Connected to database!")


	errorlisten := app.Listen(":8282")
	if errorlisten != nil {
		fmt.Println("Failed to connect to route localhost8282")
	}

	



}

