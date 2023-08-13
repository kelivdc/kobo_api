package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func SendServer(c *fiber.Ctx, method string, url string) (string, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println("Error in request", err.Error())
	}
	token := c.Locals("token").(string)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+token)
	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error in request", err.Error())
	}
	defer response.Body.Close()
	var jsonData map[string]interface{}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error in request", err.Error())
	}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		fmt.Println("Error when parsing", err.Error())
	}
	if response.StatusCode != 200 {
		fmt.Println("Error", response.StatusCode)
	}

	hasil, err := json.Marshal(jsonData)
	if err != nil {
		panic(err)
	}
	return string(hasil), nil
}
