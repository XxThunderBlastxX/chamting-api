package auth

import (
	"github.com/XxThunderBlastxX/chamting-api/config"
	"github.com/XxThunderBlastxX/chamting-api/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"time"
)

//Get the user credential by email
func getUserByEmail(email string, c *fiber.Ctx) (*model.User, error) {
	//MongoDB instance
	userDb := config.ConnectDb().Database("chamting-app").Collection("user")

	//Instance of model.User
	var user *model.User

	//Finding the received email
	err := userDb.FindOne(c.Context(), fiber.Map{"email": email}).Decode(&user)
	if err != nil {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user by email not found"})
	}

	//Disconnect from MongoDB server
	defer func() {
		if err = config.ConnectDb().Disconnect(c.Context()); err != nil {
			c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "error": err})
		}
	}()
	return user, nil
}

//Get the user credential by username
func getUserByUsername(uname string, c *fiber.Ctx) (*model.User, error) {
	//MongoDB instance
	userDb := config.ConnectDb().Database("chamting-app").Collection("user")

	//Instance of model.User
	var user *model.User

	//Finding the received username
	err := userDb.FindOne(c.Context(), fiber.Map{"username": uname}).Decode(&user)
	if err != nil {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user by username not found"})
	}

	//Disconnect from MongoDB server
	defer func() {
		if err = config.ConnectDb().Disconnect(c.Context()); err != nil {
			c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "error": err})
		}
	}()
	return user, err
}

//CheckPassword checks the given password matches with the encrypted password in DB
func CheckPassword(pass string, email string, c *fiber.Ctx) bool {

	//Initialized model.User object
	var user *model.User

	//Database Instance of chamting-app of collection-user
	userColl := config.ConnectDb().Database("chamting-app").Collection("user")

	//Getting the user data of email provided
	err := userColl.FindOne(c.Context(), fiber.Map{"email": email}).Decode(&user)

	//Disconnect from MongoDB server
	defer func() {
		if err = config.ConnectDb().Disconnect(c.Context()); err != nil {
			c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "error": err})
		}
	}()

	//Getting hashed password from server
	hash := user.Password

	//Comparing Hashed and received password is correct
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	return err == nil
}

// Login handler
func Login(c *fiber.Ctx) error {
	//Input Structure
	type LoginInput struct {
		Identity string `json:"identity"`
		Password string `json:"password"`
	}

	// Structure of user same as Mongo
	type UserData struct {
		ID       primitive.ObjectID `json:"id"`
		Username string             `json:"username"`
		Email    string             `json:"email"`
		Password string             `json:"password"`
	}
	var input LoginInput
	var ud UserData

	//Getting InputData
	err := c.BodyParser(&input)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "something went wrong on login", "data": err})
	}

	identity := input.Identity
	pass := input.Password

	//Checking if the given identity matches with any of the email
	email, _ := getUserByEmail(identity, c)

	//Checking if the given identity matches with any of the username
	uname, _ := getUserByUsername(identity, c)

	//Checking if the identity entered is false
	if email != nil && uname != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "cant find any email or username with matching identity"})
	}

	if email == nil {
		ud = UserData{
			ID:       uname.Id,
			Username: uname.UserName,
			Email:    uname.Email,
			Password: uname.Password,
		}
	} else {
		ud = UserData{
			ID:       email.Id,
			Username: email.UserName,
			Email:    email.Email,
			Password: email.Password,
		}
	}

	//Checking if the entered password matches with the DB
	if !CheckPassword(pass, ud.Email, c) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "password does not match"})
	}

	//Generating JWT token
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["username"] = ud.Username
	claims["user_id"] = ud.ID
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()

	//Getting encoded JWT token
	t, err := token.SignedString([]byte(config.Env("SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Success login", "data": t})
}
