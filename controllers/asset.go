package controllers

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sahabatgolkar.com/lib"
)

type Survey struct {
	Uid                          string
	Name                         string
	Content                      map[string]interface{}
	Deployment__submission_count int16
}

type Detail struct {
	Uid     string
	Count   int16
	Results map[string]interface{}
}

func ListAssets(c *fiber.Ctx) error {
	url := os.Getenv("SERVER_URL") + "/api/v2/assets.json"
	hasil, err := lib.SendServer(c, "GET", url)
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(hasil), &data); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error parsing JSON data")
	}
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	return c.JSON(data)
}

func PullSurvey(c *fiber.Ctx) (float64, error) {
	uid := c.Params("uid")
	hasil, err := lib.SendServer(c, "GET", os.Getenv("SERVER_URL")+"/api/v2/assets/"+uid+".json")
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(hasil), &data); err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1600*time.Second)
	defer cancel()
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	coll := client.Database("kobo").Collection("surveys")
	filter := bson.D{{"uid", uid}}
	coll.DeleteMany(ctx, filter)
	coll.InsertOne(ctx, data)
	if err != nil {
		panic(err)
	}
	total := data["deployment__submission_count"]
	return total.(float64), nil
}

func PullDetail(c *fiber.Ctx) error {
	uid := c.Params("uid")
	hasil, err := lib.SendServer(c, "GET", os.Getenv("SERVER_URL")+"/api/v2/assets/"+uid+"/data.json")
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(hasil), &data); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error parsing JSON data")
	}
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1600*time.Second)
	defer cancel()
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	coll := client.Database("kobo").Collection("details")
	filter := bson.D{{"uid", uid}}
	coll.DeleteMany(ctx, filter)
	data["uid"] = uid
	updatedData, err := json.Marshal(data)
	var bsonData interface{}
	err = json.Unmarshal(updatedData, &bsonData)
	if err != nil {
		panic(err)
	}
	coll.InsertOne(ctx, bsonData)
	if err != nil {
		panic(err)
	}
	return c.JSON(updatedData)
}

func PullData(c *fiber.Ctx) error {
	// url := os.Getenv("SERVER_URL") + "/api/v2/assets/" + uid + "/data.json" ----> Data
	// url := os.Getenv("SERVER_URL") + "/api/v2/assets/" + uid + ".json" ---> List values
	uid := c.Params("uid")
	total, err := PullSurvey(c)
	if err != nil {
		panic(err)
	}
	PullDetail(c)
	return c.JSON(fiber.Map{"uid": uid, "total": total})
}

func ReadForm(c *fiber.Ctx) error {
	uid := c.Params("uid")
	uri := os.Getenv("MONGODB_URI")
	ctx := context.TODO()
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	filter := bson.D{{"uid", uid}}
	var result Survey
	coll := client.Database("kobo").Collection("surveys")
	err = coll.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		panic(err)
	}
	return c.JSON(result)
}
