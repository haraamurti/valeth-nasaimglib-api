package handler

import (
	"fmt"
	"io"
	"net/http"
	"valeth/model"
	"valeth/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var DB *gorm.DB //DB yang disin digunakan untuk membuat objek dari gorm.db dimana kita akan bsia menggunakan method method yang ditawarkan dari packgae libarary y tersebut sehingga naun kita akan sambiungkan dengan yan berada di main router.go karena disana kita sudah menmigratenya dengan data data dari objek dan struct terbaru yang dimikliki olleh golang. dan juga appp ini

func HandlerHome(c *fiber.Ctx)error{
	var images [] model.NasaolineImage
	if err := DB.Find(&images).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch images!")
	}//this will pritn all the elemtn and their atributes

	// Build slice of URLs only
	var urls []string
	for _, img := range images {
		urls = append(urls, " "+img.URL+"                  ")
	}

	return c.Status(fiber.StatusOK).JSON(images)

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
    fmt.Println("code is running right")
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
    hashedPassword := utils.GeneratePassword(requestData.Password)

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

//=================================================handler login
func Login (c *fiber.Ctx)error{
	var loginData struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

	if err := c.BodyParser(&loginData); err !=nil{
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request format",
            "details": "Please provide valid JSON with email and password",
	})
}
	if loginData.Email == "" || loginData.Password == "" {
		return c.Status (400).JSON(fiber.Map{
			"error": "Missing required fields",
            "details": "Both email and password are required",
		})
	}
	var user model.User
	if err := DB.Where("email = ?",loginData.Email).First(&user).Error; err !=nil{
	return c.Status(401).JSON(fiber.Map{
			"error": "Invalid credentials",
            "details": "Email or password is incorrect",
	})
}
	if !utils.ComparePassword(user.PasswordHash, loginData.Password){
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid credentials",
            "details": "Email or password is incorrect",
		})
	}

	token, err := utils.GenerateToken(user.ID, user.Email)
    if err != nil {
        fmt.Printf("Token generation error: %v\n", err)
        return c.Status(500).JSON(fiber.Map{
            "error": "Authentication failed",
            "details": "Unable to generate access token",
        })
    }

    // 6. Return successful login with token
    return c.Status(200).JSON(fiber.Map{
        "success": true,
        "message": "Login successful! Welcome to NASA Image Library",
        "data": fiber.Map{
            "token":      token,              //  JWT Token for API access
            "user_id":    user.ID,
            "email":      user.Email,
            "expires_in": "24h",
            "token_type": "Bearer",
        },
    })

}