package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mu         sync.Mutex
	collection *mongo.Collection
	//client     *mongo.Client
)

func init() {
	clientoptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientoptions)
	if err != nil {
		log.Fatal(err.Error())
	}
	database := client.Database("sampleToken")
	collection = database.Collection("tokens")
	fmt.Println("Database connected successfully")
}

func main() {

	router := gin.Default()

	router.POST("/tokens", func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "token not found in the header"})
			return
		}
		go func() {
			if err := storeToken(token); err != nil {
				return
			}
		}()

		c.JSON(http.StatusOK, gin.H{"message": "token storage request received"})

	})

	router.GET("/tokens", func(c *gin.Context) {

		go func() {
			tokens, err := retrieveTokens()
			if err != nil {
				return
			}
			c.JSON(http.StatusOK, gin.H{"tokens": tokens})

		}()

	})

	router.POST("/createtokens", func(c *gin.Context) {

		header := jwt.SigningMethodHS256
		payload := jwt.MapClaims{
			"email":      "vignesh@123",
			"age":        "21",
			"customerId": "1",
			"exp":        time.Now().Add(time.Hour * 1).Unix(),
		}

		join := jwt.NewWithClaims(header, payload)
		signiture, _ := join.SignedString([]byte("SecretKey"))
		fmt.Println(signiture)

		_, err := collection.InsertOne(context.TODO(), map[string]interface{}{"token": signiture})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"message": "token inserted successfully"})

		token, _ := jwt.Parse(signiture, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Invalid Signing method")
			}
			return []byte("SecretKey"), nil

		})
		if token.Valid {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				customerId, _ := claims["customerId"].(string)
				fmt.Println(customerId)
				emailId, _ := claims["email"].(string)
				fmt.Println(emailId)
				age, _ := claims["age"].(string)
				fmt.Println(age)
			}
		}

	})

	router.Run(":4000")
}
func retrieveTokens() ([]string, error) {
	mu.Lock()
	defer mu.Unlock()

	cursor, err := collection.Find(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var tokens []string

	for cursor.Next(context.TODO()) {
		var result map[string]interface{}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		tokens = append(tokens, result["token"].(string))
	}
	return tokens, nil

}

func storeToken(token string) error {
	mu.Lock()
	defer mu.Unlock()

	_, err := collection.InsertOne(context.TODO(), map[string]interface{}{"token": token})
	return err

}
