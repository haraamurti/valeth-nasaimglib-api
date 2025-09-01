package handler

import (
	"valeth/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var DB *gorm.DB

func HandlerHome(c *fiber.Ctx)error{
	var images []model.NasaolineImage
	if err := DB.Find(&images).Error; err!=nil{
		return c.Status(201).SendString("Failed to fetch images!")
	}

	return c.Status(fiber.StatusCreated).JSON(images)

}

func HandlerImgDetails(c *fiber.Ctx)error{
	idimg := c.Params("id");
	var img model.NasaolineImage;
	result := DB.First(&img,idimg)
	if result.Error == gorm.ErrRecordNotFound{
		return c.Status(fiber.StatusNotFound).SendString("img Not Found")
	}
	return c.JSON(img)
}

func HandlerImgView(c *fiber.Ctx)error{
	idimg := c.Params("id");
	var img model.NasaolineImage;
	result := DB.First(&img,idimg)
	if result.Error == gorm.ErrRecordNotFound{
		return c.Status(fiber.StatusNotFound).SendString("img Not Found")
	}
    return c.Redirect(img.URL, fiber.StatusFound)
}