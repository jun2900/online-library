package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/jun2900/online-library/database"
	"github.com/jun2900/online-library/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Signup(c *fiber.Ctx) error {
	db := database.DBConn

	type inputUser struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
		Role     string `json:"role"`
	}
	input := new(inputUser)

	if err := c.BodyParser(input); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	inputErrors := HandlingInput(input)
	if inputErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "err": inputErrors})
	}

	if err := db.Where(&models.User{Email: input.Email}).First(&models.User{}).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": fmt.Sprintf("user with the email %s already exist", input.Email)})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't hash password", "data": err})
	}

	var role models.Role
	if err := db.Where("name = ?", input.Role).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "role not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "error on processing roles", "err": err})
	}

	user := &models.User{Email: input.Email, Password: string(hashedPassword), RoleId: role.ID}
	db.Session(&gorm.Session{FullSaveAssociations: true}).Create(&user)
	return c.JSON(fiber.Map{"status": "success", "message": "user created"})
}

func Login(c *fiber.Ctx) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := database.DBConn
	user := new(models.User)
	input := new(models.User)

	if err := c.BodyParser(input); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	inputErrors := HandlingInput(*input)
	if inputErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "err": inputErrors})
	}

	if err := db.Where(&models.User{Email: input.Email}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "User not found"})
		}
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Error on email", "data": err})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid password"})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["role_id"] = user.RoleId
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	t, err := token.SignedString([]byte(os.Getenv("ACCESS_TOKEN_SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Token successfully created", "token": t})
}

func Logout(c *fiber.Ctx) error {
	var ctx = context.Background()
	token := c.Get("x-access-token")

	rdb := database.Rdb
	if err := rdb.SAdd(ctx, "blacklist_token", token).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "err": err})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "you are logged out"})
}
