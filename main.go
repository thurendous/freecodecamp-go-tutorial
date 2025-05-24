package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

type Todo struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

// `json:id`是一个标签，用于指定JSON序列化时的字段名
// 如此一来，在序列化时，字段名会变为id，而不是ID

func main() {
	fmt.Println("Hello, World!")
	app := fiber.New()

	todos := []Todo{} // 这里是创建了一个非nil的空slice

	fmt.Println("Server is running on port 4000")

	app.Get("/", func(c *fiber.Ctx) error {
		// return c.Status(200).JSON(fiber.Map{"msg": "Hello air!"})
		return c.Status(200).JSON(todos)
	})

	// 创建一个POST请求，用于创建一个todo
	app.Post("/api/todos", func(c *fiber.Ctx) error {
		todo := &Todo{}
		fmt.Println(todo)

		if err := c.BodyParser(todo); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		fmt.Println(todo)
		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Body is required"})
		}
		todo.ID = len(todos) + 1
		todo.Completed = false
		todos = append(todos, *todo)
		return c.Status(201).JSON(todo)
	})

	// Update a todo to be completed
	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos[i].Completed = true
				return c.Status(200).JSON(todos[i])
			}
		}
		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	})

	// Update a todo's content to the designated content
	app.Put("/api/todos/:id", func(c *fiber.Ctx) error {
		// 获取请求体中的内容
		id := c.Params("id")
		// 解析请求体中的内容
		var todo Todo
		// 创建一个Todo类型的变量 这是一个Todo类型的变量，不是指针
		// 初始化成了这个值
		//    Todo{
		//       ID: 0,         // int 的零值
		//       Completed: false,  // bool 的零值
		//       Body: "",      // string 的零值
		//    }
		if err := c.BodyParser(&todo); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		// 更新todo的内容
		for i, t := range todos {
			if fmt.Sprint(t.ID) == id {
				todos[i].Body = todo.Body
				return c.Status(200).JSON(todos[i])
			}
		}
		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	})

	// Delete a todo
	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos = append(todos[:i], todos[i+1:]...)
				return c.Status(200).JSON(fiber.Map{"msg": "Todo deleted"})
			}
		}
		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	})

	log.Fatal(app.Listen(":4000"))
}

// var x int = 5 // 0x00001 -> record the val	 of x
// var p *int = &x // 0x00001 -> This is the pointer point to the memory address of x
// fmt.Println(p) // 0x00001
// fmt.Println(*p) // 5
// *p = 10 // assign the val of x to 10
// fmt.Println(x) // 10
