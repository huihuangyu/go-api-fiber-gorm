package routes

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/huihuangyu/go-api-fiber-gorm/database"
	"github.com/huihuangyu/go-api-fiber-gorm/models"
)

type Order struct {
	ID      uint    `json:"id"`
	User    User    `json:"user"`
	Product Product `json:"product"`
}

func CreateResponseOrder(order models.Order, user User, product Product) Order {
	return Order{
		ID:      order.ID,
		User:    user,
		Product: product,
	}
}

func CreateOrder(c *fiber.Ctx) error {
	var order models.Order

	if err := c.BodyParser(&order); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	var user models.User

	if err := FindUser(order.UserRefer, &user); err != nil {
		return c.Status(404).JSON(err.Error())
	}

	var product models.Product

	if err := findProduct(order.ProductRefer, &product); err != nil {
		return c.Status(404).JSON(err.Error())
	}

	database.Database.Db.Create(&order)

	responseUser := CreateResponseUser(user)
	responseProdcut := CreateResponseProduct(product)
	responseOrder := CreateResponseOrder(order, responseUser, responseProdcut)

	return c.Status(200).JSON(responseOrder)
}

func GetOrders(c *fiber.Ctx) error {
	orders := []models.Order{}

	database.Database.Db.Find(&orders)

	responseOrders := []Order{}

	for _, order := range orders {
		var user models.User
		var product models.Product
		database.Database.Db.Find(&user, "id = ?", order.UserRefer)
		database.Database.Db.Find(&product, "id = ?", order.ProductRefer)
		responseUser := CreateResponseUser(user)
		responseProdcut := CreateResponseProduct(product)
		responseOrder := CreateResponseOrder(order, responseUser, responseProdcut)

		responseOrders = append(responseOrders, responseOrder)
	}

	return c.Status(200).JSON(responseOrders)
}

func FindOrder(id int, order *models.Order) error {
	database.Database.Db.Find(&order, "id = ?", id)
	if order.ID == 0 {
		return errors.New("order does not exist")
	}
	return nil
}

func GetOrder(c *fiber.Ctx) error {
	var id, err = c.ParamsInt("id")

	if err != nil {
		return c.Status(400).JSON("please ensure that :id is an integer")
	}

	var order models.Order
	if err := FindOrder(id, &order); err != nil {
		return c.Status(404).JSON(err.Error())
	}

	var user models.User
	var product models.Product

	FindUser(order.UserRefer, &user)
	database.Database.Db.First(&product, order.ProductRefer)

	responseUser := CreateResponseUser(user)
	responseProduct := CreateResponseProduct(product)
	responseOrder := CreateResponseOrder(order, responseUser, responseProduct)

	return c.Status(200).JSON(responseOrder)
}
