package main

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Note struct{
	ID int `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
}

var notes = make(map[int]Note)
var nextID = 1
var db *gorm.DB

//===============response for home url route=================
// func handlerHome(c *fiber.Ctx)error {
// 	return c.JSON(notes) //-> bakal ngembaliin json yaitu hashmap dari notes ini
// }

func handlerHome(c *fiber.Ctx) error {
	var notes []Note
	if err := db.Find(&notes).Error; err != nil {
		return c.Status(500).SendString("Failed to fetch notes")
	}
	return c.JSON(notes)
}

//============getting notes with id route==================
func handlerGetNotes(c *fiber.Ctx) error {
    idStr := c.Params("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
    }

    var note Note
    // Use GORM to find by primary key
    result := db.First(&note, id)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return c.Status(fiber.StatusNotFound).SendString("Note not found")
        }
        return c.Status(fiber.StatusInternalServerError).SendString("Database error")
    }

    return c.JSON(note)
}

//r============response from posting notes or making a note============
func handlerCreateNotes(c *fiber.Ctx)error {
	note := new (Note)
	if err := c.BodyParser(note); err !=nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request")
	}

	if err := db.Create(note).Error; err !=nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid request")
	}

	note.ID = nextID
	notes[note.ID] = *note
	nextID++

	return c.Status(fiber.StatusCreated).JSON(note)
}

//=================handler updating and editing a note==============
func handlerUpdateNote(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(400)
	}

	// Step 2: Parse new content from body
    var body struct {
        Content string `json:"content"`
    }

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid Body - neval")
	}

	var updatednote Note
	if err := db.First(&updatednote, id).Error; err != nil{
		return c.Status(fiber.StatusInternalServerError).SendString("error Fetching note -neval")
	}

	updatednote.Content = body.Content

	// Save the updated note to the database
	if err := db.Save(&updatednote).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("error updating note - neval")
	}

	return c.JSON(updatednote)


}

//==============handler to delete a note============
func handlertDeleteNote(c *fiber.Ctx)error{
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.SendStatus(400)
	}
	db.Delete(&Note{}, id)
	return c.JSON("file with the id "+ strconv.Itoa(id)+"is already deleted")
}

func main(){
	//started using fiber backend framework
	app := fiber.New() // -> fiber declaration make object called "app"
	app.Get("/",handlerHome) //home
	app.Get("/notes/:id",handlerGetNotes) //getnotes
	app.Post("/notes",handlerCreateNotes) //create notes
	app.Put("/notes/:id", handlerUpdateNote) // update notes and content
	app.Delete("/delete/notes/:id", handlertDeleteNote)//delete a note

	// -------- Database connection code (add this block) ---------------
	host := "localhost"
	user := "postgres"
	password := "yipikaye2123"
	port := 5432
	dbname := "neval" //using database this
	sslmode := "disable"

	dsn := fmt.Sprintf(
    "host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
    host, user, password, dbname, port, sslmode,)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})//connecting the database using gorm or opening a database.
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	fmt.Println("Connected to database!")
	_ = db //assigning the database to ignore because i havent use any datbases yet
	//---------------------------------------------------------------------
	
	db.AutoMigrate(&Note{}) //-> migrating the db from the note strucct


	//->error handling
    err = app.Listen(":8181")
	if err != nil {
        fmt.Println("error !")
    }

}