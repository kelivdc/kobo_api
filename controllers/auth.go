package controllers

import (
	"context"
	"io"
	"net/http"
	"os"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"sahabatgolkar.com/database"
)

type User struct {
	Name     string
	Email    string
	Password string
}

type UserLogin struct {
	Username string
	Password string
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateAdmin(c *fiber.Ctx) error {
	hash, _ := HashPassword("12345678")
	client, err := database.ConnectMongoDB()
	if err != nil {
		panic(err)
	}
	coll := client.Database(os.Getenv("DATABASE")).Collection("users")
	doc := User{Email: "admin@example.com", Password: hash}
	result, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		panic(err)
	}
	return c.JSON(fiber.Map{"message": "Create Admin", "data": result.InsertedID})
}

func Login(c *fiber.Ctx) error {
	user := new(UserLogin)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	url := os.Getenv("SERVER_URL") + "/token/?format=json"
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(user.Username, user.Password)
	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return c.JSON(fiber.Map{"error": err.Error()})
	}
	defer resp.Body.Close()
	var jsonData map[string]interface{}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(fiber.Map{"error": err.Error()})
	}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		return c.Status(resp.StatusCode).JSON(jsonData)
	}
	return c.JSON(jsonData)
}

func Auth(c *fiber.Ctx) error {
	ctx := context.TODO()

	client, err := database.ConnectMongoDB()
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	DB := client.Database(os.Getenv("DATABASE"))
	collection := DB.Collection("users")
	cursor, err := collection.Find(ctx, bson.D{})
	defer cursor.Close(ctx)
	var results []User
	if err = cursor.All(ctx, &results); err != nil {
		panic(err)
	}
	return c.JSON(fiber.Map{"data": results})
}
