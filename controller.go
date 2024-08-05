package main

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// get all todos
func GetTodos(c *fiber.Ctx) error {
	var todos []Todo

	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		return err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo Todo
		if err := cursor.Decode(&todo); err != nil {
			return err
		}

		todos = append(todos, todo)
	}

	fmt.Println()
	return c.Status(200).JSON(todos)
}

func CreateTodo(c *fiber.Ctx) error {

	// todo := &Todo{}
	todo := new(Todo) //used when no field is init

	if err := c.BodyParser(todo); err != nil {
		return err
	}
	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{
			"msg": "need todo body",
		})
	}

	insert, err := collection.InsertOne(context.Background(), todo)
	if err != nil {
		return err
	}

	todo.ID = insert.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(todo)
}
func UpdateTodo(c *fiber.Ctx) error {

	id := c.Params("id")

	ObjectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"msg": "invalid todo id",
		})
	}

	filter := bson.M{"_id": ObjectID}
	update := bson.M{"$set": bson.M{"completed": true}}

	_, err = collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
	})

}

func DeleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")

	ObjectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
		// return c.Status(400).JSON(fiber.Map{
		// 	"msg": "invalid todo id",
		// })
	}

	filter := bson.M{"_id": ObjectID}

	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
	})
}
