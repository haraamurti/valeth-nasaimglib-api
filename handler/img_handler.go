package handler

import (
	"fmt"
	"io"
	"net/http"
	"valeth/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var DB *gorm.DB

func HandlerHome(c *fiber.Ctx)error{
	var images []model.NasaolineImage
	if err := DB.Find(&images).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch images!")
	}

	// Build slice of URLs only
	var urls []string
	for _, img := range images {
		urls = append(urls, " "+img.URL+"                  ")
	}

	return c.Status(fiber.StatusOK).JSON(urls)

}

func HandlerImgDetails(c *fiber.Ctx)error{
	var img model.NasaolineImage;
	idimg := c.Params("id");
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



func Handlerdownloadimg(c *fiber.Ctx) error {
    idimg := c.Params("id")
    var img model.NasaolineImage
    result := DB.First(&img, idimg)
    if result.Error == gorm.ErrRecordNotFound {
        return c.Status(fiber.StatusNotFound).SendString("img Not Found")
    }

    // Download the image from its URL
    resp, err := http.Get(img.URL)
    if err != nil || resp.StatusCode != 200 {
        return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch image from URL")
    }
    defer resp.Body.Close()

    // Set headers for file download
    c.Set("Content-Type", "application/octet-stream")
    c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"image_%s.jpg\"", idimg))

    // Stream the image to the client
    _, err = io.Copy(c, resp.Body)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString("Failed to send file")
    }

    return nil
}