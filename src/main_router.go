package main

import (
	"fmt"
	"os"

	"valeth/handler"
	"valeth/model"
	"valeth/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB


func main(){

	app := fiber.New()

	fmt.Println("Hello world this is for NASA image library")

	dsn := "postgresql://postgres.xrsnptveunsdsnfcvxjz:yipikaye2123@aws-1-ap-southeast-1.pooler.supabase.com:5432/postgres"
	var errConnectDb error //eeror handler for making db handler eerro and the type is error in golang, error is a type that if error occur it will be filledwith the interface , but if therea re no tthe value will be nil

	db, errConnectDb = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if errConnectDb != nil {
		fmt.Println("Failed to connect to database !")
		return
	}

	db.AutoMigrate(&model.NasaolineImage{},&model.User{})
	fmt.Println("Connected to database!")
	handler.DB = db//this will asign the gorm db connection in handler to use the detials and spesification in the handler.
	//adding details for review

	app.Get("/", handler.HandlerHome)
	app.Get("/imgs/:id",handler.HandlerImgDetails)
	app.Get("/imgs/:id/download",utils.AuthMiddleware,handler.Handlerdownloadimg)
	app.Get("/imgs/:id/view",utils.AuthMiddleware,handler.HandlerImgView)

	app.Post("/register",handler.Register)
	app.Post("/login",handler.Login)


	port := os.Getenv("PORT")
	if port == "" {
    port = "8282" // fallback for local dev
}
	errorlisten := app.Listen(":" + port)
	if errorlisten != nil {
        fmt.Println("Failed to connect to route on port " + port)
    }
}
//new update so it can run on railway