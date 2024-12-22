package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/naufalathallah/go-todolist/database"
	"github.com/naufalathallah/go-todolist/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}


func helloWorld(c *fiber.Ctx)error  {
	return c.SendString("Hello World")
}

func initDatabase() {
	var err error

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password='%s' dbname=%s port=%s",
		host, user, password, dbname, port)

	fmt.Printf("DSN: %s\n", dsn)

	database.DBConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	fmt.Println("Database connected!")

	database.DBConn.AutoMigrate(&models.Todo{})
	fmt.Println("Migrated DB")
}

func closeDatabase() {
    db, err := database.DBConn.DB()
    if err != nil {
        log.Fatalf("Failed to get database connection: %v", err)
    }
    db.Close()
}

func setupRoutes(app *fiber.App)  {
	app.Get("/", helloWorld)
	app.Get("/todos", models.GetTodos)
	app.Get("/todos/:id", models.GetTodoById)
	app.Post("/todos", models.CreateTodo)
	app.Put("/todos/:id", models.UpdateTodo)
	app.Delete("/todos/:id", models.DeleteTodo)
}

func main()  {
	loadEnv()
	app := fiber.New()
	initDatabase()
    defer closeDatabase()
	
	setupRoutes(app)
	app.Listen(":8000")
}