package handler

import (
	"fmt"
	"io"
	"net/http"
	"valeth/model"
	Utils "valeth/utils"

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
//=========================================================================
//handler for register jwt tokens

func Register(c *fiber.Ctx) error {
    // 1. Get the data from the request (we still need this to receive email/password)
    var requestData struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    
    if err := c.BodyParser(&requestData); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Bad request"})
    }

    // 2. Check if user already exists
    var existingUser model.User
    if err := DB.Where("email = ?", requestData.Email).First(&existingUser).Error; err == nil {
        return c.Status(400).JSON(fiber.Map{"error": "User already exists"})
    }

    // 3. Hash the password using YOUR function
    hashedPassword := Utils.GeneratePassword(requestData.Password)

    // 4. Create new user using YOUR User struct
    user := model.User{
        Email:        requestData.Email,
        PasswordHash: hashedPassword,  // Notice: PasswordHash, not Password
    }

    // 5. Save to database
    if err := DB.Create(&user).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
    }

    // 6. Return success
    return c.JSON(fiber.Map{
        "message": "User created successfully",
        "user_id": user.ID,
        "email":   user.Email,
    })
}