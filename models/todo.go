package models

import (
	"github.com/gofiber/fiber/v2"
	"github.com/naufalathallah/go-todolist/database"
)

type Todo struct{
	ID uint `gorm:"primarykey" json:"id"`
	Title string `json:"title"`
	Completed bool `json:"completed"`
}

func GetTodos(c *fiber.Ctx) error {
	db := database.DBConn
	var todos []Todo
	db.Find(&todos)
	return c.JSON(&todos)
}

func GetTodoById(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var todo Todo

	result := db.First(&todo, id)

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not find todo",
			"data":    nil,
		})
	}

	return c.JSON(&todo)
}

func CreateTodo(c *fiber.Ctx) error {
	db := database.DBConn
	todo := new(Todo)
	err := c.BodyParser(todo)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Check your input", "data": err})
	}
	err = db.Create(&todo).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create todo", "data": err})
	}
	return c.JSON(&todo)
}

func UpdateTodo(c *fiber.Ctx) error {
	type UpdatedTodo struct{
		Title string `json:"title"`
		Completed bool `json:"completed"`
	}

	id := c.Params("id")
	db := database.DBConn
	var todo Todo

	result := db.First(&todo, id)

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not find todo",
			"data":    nil,
		})
	}

	var updatedTodo UpdatedTodo
	err := c.BodyParser(&updatedTodo)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Check your input", "data": err})
	}

	todo.Title = updatedTodo.Title
	todo.Completed = updatedTodo.Completed
	db.Save(&todo)

	return c.JSON(&todo)
}