package controller

import (
	"log"
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

func GetProjects(c *fiber.Ctx) error {
	collection := db.Collection("projects")

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ApiResponse{Success: false, Error: err.Error()})
	}

	var projects []model.Project = make([]model.Project, 0)
	pagination, err := mongopagination.New(collection).Limit(20).Page(int64(page)).Sort("createdAt", -1).Filter(bson.M{}).Decode(&projects).Find()
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(common.ApiResponse{Success: false, Error: err.Error()})
	}

	return c.Status(200).JSON(common.ApiResponse{Success: true, Data: fiber.Map{"docs": projects, "pagination": pagination.Pagination}})
}

func GetProject(c *fiber.Ctx) error {
	rawProjectId := c.Params("projectId")

	projectId, err := primitive.ObjectIDFromHex(rawProjectId)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(common.ApiResponse{Success: false, Error: err.Error()})
	}

	var project model.Project
	collection := db.Collection("projects")
	err = collection.FindOne(c.Context(), bson.M{"_id": projectId}).Decode(&project)

	if err == mongo.ErrNoDocuments {
		return c.Status(http.StatusNotFound).JSON(common.ApiResponse{Success: false, Error: "Project not found"})
	} else if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(common.ApiResponse{Success: false, Error: err.Error()})
	}

	return c.Status(http.StatusOK).JSON(common.ApiResponse{Success: true, Data: project})
}

type CreateProjectDTO struct {
	Name        string             `json:"name" bson:"name" validate:"required,min=3,max=36"`
	Description string             `json:"description" bson:"description" validate:"required,min=24,max=1024"`
	Status      int                `json:"status" bson:"status" validate:"required,gte=1,lte=3"`
	ClientID    primitive.ObjectID `json:"clientId" bson:"clientId" validate:"required"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt" validate:"omitempty"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt" validate:"omitempty"`
}

func CreateProject(c *fiber.Ctx) error {
	project := new(CreateProjectDTO)
	project.CreatedAt = time.Now()
	project.UpdatedAt = time.Now()

	if err := c.BodyParser(project); err != nil {
		return c.Status(http.StatusBadRequest).JSON(common.ApiResponse{Success: false, Error: err.Error()})
	}

	clientId, err := primitive.ObjectIDFromHex(project.ClientID.Hex())
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(common.ApiResponse{Success: false, Error: err.Error()})
	}

	var client model.Client
	clientCollection := db.Collection("clients")
	err = clientCollection.FindOne(c.Context(), bson.M{"_id": clientId}).Decode(&client)

	if err == mongo.ErrNoDocuments {
		return c.Status(http.StatusNotFound).JSON(common.ApiResponse{Success: false, Error: "Client not found"})
	} else if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(common.ApiResponse{Success: false, Error: err.Error()})
	}

	project.ClientID = clientId

	validate := validator.New()
	if err := validate.Struct(project); err != nil {
		return c.Status(http.StatusBadRequest).JSON(common.ApiResponse{Success: false, Error: err.Error()})
	}

	projectCollection := db.Collection("projects")
	insertResult, err := projectCollection.InsertOne(c.Context(), project)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(common.ApiResponse{Success: false, Error: err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(common.ApiResponse{Success: true, Data: fiber.Map{"id": insertResult.InsertedID}})
}

type UpdateProjectDTO struct {
	Name        string    `json:"name" bson:"name,omitempty" validate:"omitempty,min=3,max=36"`
	Description string    `json:"description" bson:"description,omitempty" validate:"omitempty,min=24,max=1024"`
	Status      int       `json:"status" bson:"status,omitempty" validate:"omitempty,gte=1,lte=3"`
	UpdatedAt   time.Time `json:"updatedAt" bson:"updatedAt" validate:"omitempty"`
}

func UpdateProject(c *fiber.Ctx) error {
	rawProjectId := c.Params("projectId")
	projectId, err := primitive.ObjectIDFromHex(rawProjectId)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(common.ApiResponse{Success: false, Error: err.Error()})
	}

	var project UpdateProjectDTO
	project.UpdatedAt = time.Now()
	if err := c.BodyParser(&project); err != nil {
		return c.Status(http.StatusBadRequest).JSON(common.ApiResponse{Success: false, Error: err.Error()})
	}

	validate := validator.New()
	if err := validate.Struct(&project); err != nil {
		return c.Status(http.StatusBadRequest).JSON(common.ApiResponse{Success: false, Error: err.Error()})
	}

	collection := db.Collection("projects")
	updateResult, err := collection.UpdateOne(c.Context(), bson.M{"_id": projectId}, bson.M{"$set": project})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(common.ApiResponse{Success: false, Error: err.Error()})
	} else if updateResult.MatchedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(common.ApiResponse{Success: false, Error: "Project not found"})
	}

	return c.Status(http.StatusOK).JSON(common.ApiResponse{Success: true, Data: fiber.Map{"modified": updateResult.ModifiedCount == 1}})
}

func DeleteProject(c *fiber.Ctx) error {
	rawProjectId := c.Params("projectId")
	projectId, err := primitive.ObjectIDFromHex(rawProjectId)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(common.ApiResponse{Success: false, Error: err.Error()})
	}

	collection := db.Collection("projects")
	deleteResult, err := collection.DeleteOne(c.Context(), bson.M{"_id": projectId})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(common.ApiResponse{Success: false, Error: err.Error()})
	} else if deleteResult.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(common.ApiResponse{Success: false, Error: "Project not found"})
	}

	return c.Status(http.StatusOK).JSON(common.ApiResponse{Success: true, Data: deleteResult})
}
