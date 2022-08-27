package routes

import (
	"encoding/json"
	"fmt"
	"github.com/XxThunderBlastxX/chamting-api/database"
	"github.com/XxThunderBlastxX/chamting-api/repository"
	"github.com/XxThunderBlastxX/chamting-api/service"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http/httptest"
	"testing"
)

func TestInitialRouter(t *testing.T) {
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

	//Instance of authentication handler/service/repository
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

		fmt.Println("Test Case: " + fmt.Sprint(rune(testNo+1)) + " Passed")
	}
}
