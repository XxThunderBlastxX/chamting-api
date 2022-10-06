package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/XxThunderBlastxX/chamting-api/database"
	"github.com/XxThunderBlastxX/chamting-api/repository"
	"github.com/XxThunderBlastxX/chamting-api/service"
	"github.com/XxThunderBlastxX/chamting-api/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http/httptest"
	"testing"
)

type UserDataModel struct {
	Id       string `json:"id" bson:"_id"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	UserName string `json:"username" bson:"username"`
	Name     string `json:"name" bson:"name"`
}

type Response struct {
	Success bool          `json:"success" bson:"success"`
	Data    UserDataModel `json:"data" bson:"data"`
	Error   string        `json:"error" bson:"error"`
	Token   string        `json:"token,omitempty" bson:"token,omitempty"`
}

type Request struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	UserName string `json:"username" bson:"username"`
	Name     string `json:"name" bson:"name"`
}

// TestInitialRoute is for testing initial route of the app - /
func TestInitialRoute(t *testing.T) {
	// structure of the body
	type bodyMap struct {
		CreatedBy string `json:"Created By"`
		Name      string `json:"Name"`
		Status    string `json:"Status"`
		Version   string `json:"Version"`
	}

	// all the test cases
	tests := []struct {
		description string

		// Test input
		route string

		// Expected output
		expectedError bool
		expectedCode  int
		expectedBody  bodyMap
	}{
		{
			route:         "/",
			description:   "Testing initial route",
			expectedBody:  bodyMap{CreatedBy: "Koustav Mondal <ThunderBlast>", Status: "Running", Version: "0.0.1", Name: "Chamting - API"},
			expectedCode:  200,
			expectedError: false,
		},
	}

	//Loads variables from .env
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}

	app := fiber.New()

	//Connect to mongo-database
	db, cancel, _ := database.DBConnect()
	defer cancel()

	//Instance of authentication handler/service_mock/repository
	authCollection := db.Collection("auth")
	authRepo := repository.NewAuthRepo(authCollection)
	authService := service.NewAuthService(authRepo)

	// Router instance
	Router(app, authService)

	// loop through all the test cases and test each case
	for testNo, test := range tests {
		// Create a new http request with the route from the test case
		req := httptest.NewRequest("GET", test.route, nil)

		// Perform the request plain with the app.
		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		// verify that no error occurred, since that is not expected
		assert.Equalf(t, test.expectedError, err != nil, test.description)

		// As expected errors lead to broken responses, the next test case needs to be processed
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)

		// Read the response body and map to bodyMap
		body, readErr := io.ReadAll(res.Body)

		// Reading the response body should work everytime, such that the readErr variable should be nil
		assert.Nilf(t, readErr, test.description)

		// bodyMap instance
		resBody := bodyMap{}

		_ = json.Unmarshal(body, &resBody)

		// Verify, that the response body equals the expected body
		assert.Equalf(t, test.expectedBody, resBody, test.description)

		fmt.Println("Test Case: " + fmt.Sprint(rune(testNo+1)) + " Passed ✅ ")
	}
}

// TestAuthSignup is for testing signup route - /auth/signup
func TestAuthSignup(t *testing.T) {
	// All the test cases
	tests := []struct {
		description string

		// Test route input
		route string

		// Expected output
		expectedError bool
		expectedCode  int
		reqBody       Request
		resBody       Response
	}{
		{
			route:       "/auth/signup",
			description: "Testing by adding new user to DB",
			reqBody: Request{
				Email:    "username@gmail.com",
				Password: "password@123",
				UserName: "username",
				Name:     "User Name",
			},
			resBody: Response{
				Success: true,
				Data: UserDataModel{
					Email:    "username@gmail.com",
					Password: "password@123",
					UserName: "username",
					Name:     "User Name",
				},
				Error: "",
				Token: "",
			},
			expectedCode:  200,
			expectedError: false,
		},
		{
			route:       "/auth/signup",
			description: "Trying to add duplicate user with same email",
			reqBody: Request{
				Email:    "username@gmail.com",
				Password: "password@123",
				UserName: "username",
				Name:     "User Name",
			},
			resBody: Response{
				Success: false,
				Data:    UserDataModel{},
				Error:   "email already exist",
			},
			expectedCode:  409,
			expectedError: true,
		},
	}

	// Loads variables from .env
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}

	app := fiber.New()

	// Connect to mongo-database
	db, cancel, _ := database.DBConnect()
	defer cancel()

	// Instance of authentication handler/service_mock/repository
	authCollection := db.Collection("auth_test")
	authRepo := repository.NewAuthRepo(authCollection)
	authService := service.NewAuthService(authRepo)

	// Router instance
	Router(app, authService)

	// loop through all the test cases and test each case
	for testNo, test := range tests {
		// Encoding json body
		jsonBody, jsonErr := json.Marshal(test.reqBody)
		if jsonErr != nil {
			continue
		}

		// Create a new http request with the route from the test case
		req := httptest.NewRequest(fiber.MethodPost, test.route, bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		// Perform the request plain with the app.
		// The -1 disables request latency.
		res, err := app.Test(req, 10*1000)

		//// As expected errors lead to broken responses, the next test case needs to be processed
		//if test.expectedError {
		//	continue
		//}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)

		// Read the response body and map to bodyMap
		body, readErr := io.ReadAll(res.Body)

		// Reading the response body should work everytime, such that the readErr variable should be nil
		assert.Nilf(t, readErr, test.description)

		// bodyMap instance
		resBody := Response{}

		_ = json.Unmarshal(body, &resBody)

		switch res.StatusCode {
		case 200:
			// verify that no error occurred, since that is not expected
			assert.Equalf(t, test.expectedError, err != nil, test.description)

			// Verify Object Id
			assert.Truef(t, primitive.IsValidObjectID(resBody.Data.Id), test.description)

			// Verify name
			assert.Equalf(t, test.resBody.Data.Name, resBody.Data.Name, test.description)

			// Verify email
			assert.Equalf(t, test.resBody.Data.Email, resBody.Data.Email, test.description)

			// Verify username
			assert.Equalf(t, test.resBody.Data.UserName, resBody.Data.UserName, test.description)

			// Verify password
			assert.Truef(t, utils.VerifyPassword(test.reqBody.Password, resBody.Data.Password) == nil, test.description)

			break

		case 409:
			// verify that no error occurred, since that is expected
			assert.Equalf(t, test.expectedError, err == nil, test.description)

			// Verify the error response body matches with expected error
			assert.Equalf(t, test.resBody.Error, resBody.Error, test.description)

			break
		}

		// Printing the case no. which is passed
		fmt.Println("Test Case: " + fmt.Sprint(rune(testNo+1)) + " Passed ✅ ")
	}
}

// TestAuthSignin is for testing signin route - /auth/signin
func TestAuthSignin(t *testing.T) {
	// All the test cases
	tests := []struct {
		description string

		// Test route input
		route string

		// Expected output
		expectedError bool
		expectedCode  int
		reqBody       Request
		resBody       Response
	}{
		{
			route:       "/auth/signin",
			description: "Testing by adding new user to DB",
			reqBody: Request{
				Email:    "username@gmail.com",
				Password: "password@123",
			},
			resBody: Response{
				Success: true,
				Data: UserDataModel{
					Email:    "username@gmail.com",
					Password: "password@123",
					UserName: "username",
					Name:     "User Name",
				},
				Error: "",
				Token: "",
			},
			expectedCode:  200,
			expectedError: false,
		},
		{
			route:       "/auth/signin",
			description: "Testing by adding new user to DB",
			reqBody: Request{
				Email:    "username456@gmail.com",
				Password: "pass@123",
			},
			resBody: Response{
				Success: false,
				Data:    UserDataModel{},
				Error:   "mongo: no documents in result",
				Token:   "",
			},
			expectedCode:  404,
			expectedError: true,
		},
		{
			route:       "/auth/signin",
			description: "Testing by adding new user to DB",
			reqBody: Request{
				Email:    "username@gmail.com",
				Password: "pass@123",
			},
			resBody: Response{
				Success: false,
				Data:    UserDataModel{},
				Error:   "crypto/bcrypt: hashedPassword is not the hash of the given password",
				Token:   "",
			},
			expectedCode:  401,
			expectedError: true,
		},
	}

	// Loads variables from .env
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}

	app := fiber.New()

	// Connect to mongo-database
	db, cancel, _ := database.DBConnect()
	defer cancel()

	// Instance of authentication handler/service_mock/repository
	authCollection := db.Collection("auth_test")
	authRepo := repository.NewAuthRepo(authCollection)
	authService := service.NewAuthService(authRepo)

	// Router instance
	Router(app, authService)

	// loop through all the test cases and test each case
	for testNo, test := range tests {
		// Encoding json body
		jsonBody, jsonErr := json.Marshal(test.reqBody)
		if jsonErr != nil {
			continue
		}

		// Create a new http request with the route from the test case
		req := httptest.NewRequest(fiber.MethodPost, test.route, bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		// Perform the request plain with the app.
		// The -1 disables request latency.
		res, err := app.Test(req, 10*1000)

		//// As expected errors lead to broken responses, the next test case needs to be processed
		//if test.expectedError {
		//	continue
		//}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)

		// Read the response body and map to bodyMap
		body, readErr := io.ReadAll(res.Body)

		// Reading the response body should work everytime, such that the readErr variable should be nil
		assert.Nilf(t, readErr, test.description)

		// bodyMap instance
		resBody := Response{}

		_ = json.Unmarshal(body, &resBody)

		switch res.StatusCode {
		case 200:
			// verify that no error occurred, since that is not expected
			assert.Equalf(t, test.expectedError, err != nil, test.description)

			// Verify Object Id
			assert.Truef(t, primitive.IsValidObjectID(resBody.Data.Id), test.description)

			// Verify name
			assert.Equalf(t, test.resBody.Data.Name, resBody.Data.Name, test.description)

			// Verify email
			assert.Equalf(t, test.resBody.Data.Email, resBody.Data.Email, test.description)

			// Verify username
			assert.Equalf(t, test.resBody.Data.UserName, resBody.Data.UserName, test.description)

			// Verify password
			assert.Truef(t, utils.VerifyPassword(test.reqBody.Password, resBody.Data.Password) == nil, test.description)

			break

		case 404:
			// verify that no error occurred, since that is expected
			assert.Equalf(t, test.expectedError, err == nil, test.description)

			// Verify the error response body matches with expected error
			assert.Equalf(t, test.resBody.Error, resBody.Error, test.description)

			break

		case 401:
			// verify that no error occurred, since that is expected
			assert.Equalf(t, test.expectedError, err == nil, test.description)

			// Verify the error response body matches with expected error
			assert.Equalf(t, test.resBody.Error, resBody.Error, test.description)

			break
		}

		// Printing the case no. which is passed
		fmt.Println("Test Case: " + fmt.Sprint(rune(testNo+1)) + " Passed ✅ ")
	}
}
