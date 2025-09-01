package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type nasaoline_image struct{

	ID int  `json:"id"`
	Title string
 	URL string `json:"url"`         // Image link
    Author    string `json:"author"`
    CreatedAt string `json:"created_at"`
}

func handlerhome(c *fiber.Ctx)error{
	var images [] nasaoline_image
	if err := db.Find(&images).Error; err!=nil{
		return c.Status(201).SendString("Failed to fetch images!")
	}

	return c.Status(fiber.StatusCreated).JSON(images)

}

func handlerimgdetails(c *fiber.Ctx)error{
	idimg := c.Params("id");
	var img nasaoline_image;
	result := db.First(&img,idimg)
	if result.Error == gorm.ErrRecordNotFound{
		return c.Status(fiber.StatusNotFound).SendString("img Not Found")
	}
	return c.JSON(img)
}

func handlerimgview(c *fiber.Ctx)error{
	idimg := c.Params("id");
	var img nasaoline_image;
	result := db.First(&img,idimg)
	if result.Error == gorm.ErrRecordNotFound{
		return c.Status(fiber.StatusNotFound).SendString("img Not Found")
	}
    return c.Redirect(img.URL, fiber.StatusFound)
}

func main(){
	app := fiber.New()
	app.Get("/",handlerhome)
	app.Get("/imgs/:id",handlerimgdetails)
	app.Get("/imgs/:id/view",handlerimgview)



	fmt.Println("Hello world this is for NASA image library")

	dsn := "postgresql://postgres.xrsnptveunsdsnfcvxjz:yipikaye2123@aws-1-ap-southeast-1.pooler.supabase.com:5432/postgres"

	var errConnectDb error

	db, errConnectDb = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if errConnectDb != nil {
		fmt.Println("Failed to connect to database !")
		return
	}


	db.AutoMigrate(&nasaoline_image{})
	fmt.Println("Connected to database!")


	errorlisten := app.Listen(":8282")
	if errorlisten != nil {
		fmt.Println("Failed to connect to route localhost8282")
	}


	


}

