package controller

import (
	"net/http"
	"project-mgmt-go/common"
	"project-mgmt-go/db"
	"project-mgmt-go/model"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	mongopagination "github.com/gobeam/mongo-go-pagination"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetClients(c *fiber.Ctx) error {
	collection := db.Collection("clients")

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ApiResponse{Success: false, Error: err.Error()})
	}

	var clients []model.Client = make([]model.Client, 0)
	pagination, err := mongopagination.New(collection).Limit(20).Page(int64(page)).Sort("createdAt", -1).Filter(bson.M{}).Decode(&clients).Find()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(common.ApiResponse{Success: false, Error: err.Error()})
	}

	return c.Status(200).JSON(common.ApiResponse{Success: true, Data: fiber.Map{"docs": clients, "pagination": pagination.Pagination}})
}

func GetClient(c *fiber.Ctx) error {
	rawClientId := c.Params("clientId")

	clientId, err := primitive.ObjectIDFromHex(rawClientId)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(common.ApiResponse{Success: false, Error: err.Error()})
	}

	var client model.Client
	collection := db.Collection("clients")
	err = collection.FindOne(c.Context(), bson.M{"_id": clientId}).Decode(&client)

	if err == mongo.ErrNoDocuments {
		return c.Status(http.StatusNotFound).JSON(common.ApiResponse{Success: false, Error: "Client not found"})
	} else if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(common.ApiResponse{Success: false, Error: err.Error()})
	}

	return c.Status(http.StatusOK).JSON(common.ApiResponse{Success: true, Data: client})
}

type CreateClientDTO struct {
	Name         string    `json:"name" bson:"name" validate:"required,min=3,max=24"`
	Email        string    `json:"email" bson:"email" validate:"required,email"`
	Phone        string    `json:"phone" bson:"phone" validate:"required,numeric,min=8,max=12"`
	Organization string    `json:"organization" bson:"organization" validate:"required,min=3,max=36"`
	CreatedAt    time.Time `json:"createdAt" bson:"createdAt" validate:"omitempty"`
	UpdatedAt    time.Time `json:"updatedAt" bson:"updatedAt" validate:"omitempty"`
}

func CreateClient(c *fiber.Ctx) error {
	client := new(CreateClientDTO)
	client.CreatedAt = time.Now()
	client.UpdatedAt = time.Now()

	if err := c.BodyParser(client); err != nil {
		return c.Status(http.StatusBadRequest).JSON(common.ApiResponse{Success: false, Error: err.Error()})
	}

	validate := validator.New()
	if err := validate.Struct(client); err != nil {
		return c.Status(http.StatusBadRequest).JSON(common.ApiResponse{Success: false, Error: err.Error()})
	}

	collection := db.Collection("clients")
	insertResult, err := collection.InsertOne(c.Context(), client)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(common.ApiResponse{Success: false, Error: err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(common.ApiResponse{Success: true, Data: fiber.Map{"id": insertResult.InsertedID}})
}

type UpdateClientDTO struct {
	Name         string    `json:"name" bson:"name,omitempty" validate:"omitempty,min=3,max=24"`
	Email        string    `json:"email" bson:"email,omitempty" validate:"omitempty,email"`
	Phone        string    `json:"phone" bson:"phone,omitempty" validate:"omitempty,numeric,min=8,max=12"`
	Organization string    `json:"organization" bson:"organization,omitempty" validate:"omitempty,min=3,max=36"`
	UpdatedAt    time.Time `json:"updatedAt" bson:"updatedAt" validate:"omitempty"`
}

func UpdateClient(c *fiber.Ctx) error {
	rawClientId := c.Params("clientId")
	clientId, err := primitive.ObjectIDFromHex(rawClientId)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(common.ApiResponse{Success: false, Error: err.Error()})
	}

	var client UpdateClientDTO
	client.UpdatedAt = time.Now()
	if err := c.BodyParser(&client); err != nil {
		return c.Status(http.StatusBadRequest).JSON(common.ApiResponse{Success: false, Error: err.Error()})
	}

	validate := validator.New()
	if err := validate.Struct(&client); err != nil {
		return c.Status(http.StatusBadRequest).JSON(common.ApiResponse{Success: false, Error: err.Error()})
	}

	collection := db.Collection("clients")
	updateResult, err := collection.UpdateOne(c.Context(), bson.M{"_id": clientId}, bson.M{"$set": client})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(common.ApiResponse{Success: false, Error: err.Error()})
	} else if updateResult.MatchedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(common.ApiResponse{Success: false, Error: "Client not found"})
	}

	return c.Status(http.StatusOK).JSON(common.ApiResponse{Success: true, Data: fiber.Map{"modified": updateResult.ModifiedCount == 1}})
}

func DeleteClient(c *fiber.Ctx) error {
	rawClientId := c.Params("clientId")
	clientId, err := primitive.ObjectIDFromHex(rawClientId)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(common.ApiResponse{Success: false, Error: err.Error()})
	}

	clientCollection := db.Collection("clients")
	projectCollection := db.Collection("projects")

	_, err = projectCollection.DeleteMany(c.Context(), bson.M{"clientId": clientId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(common.ApiResponse{Success: false, Error: err.Error()})
	}

	deleteResult, err := clientCollection.DeleteOne(c.Context(), bson.M{"_id": clientId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(common.ApiResponse{Success: false, Error: err.Error()})
	} else if deleteResult.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(common.ApiResponse{Success: false, Error: "Client not found"})
	}

	return c.Status(http.StatusOK).JSON(common.ApiResponse{Success: true, Data: deleteResult})
}
