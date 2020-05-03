package controllertests

import (
	"api/youtube-clone-backend/api/models"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	samples := []struct {
		inputJSON    string
		statusCode   int
		username     string
		email        string
		errorMessage string
	}{
		{
			inputJSON:    `{"username":"pet", "email": "pet@gmail.com", "password": "password"}`,
			statusCode:   201,
			username:     "pet",
			email:        "pet@gmail.com",
			errorMessage: "",
		},
		{
			inputJSON:    `{"username":"Frank", "email": "pet@gmail.com", "password": "password"}`,
			statusCode:   500,
			errorMessage: "Email Already Taken",
		},
		{
			inputJSON:    `{"username":"pet", "email": "grand@gmail.com", "password": "password"}`,
			statusCode:   500,
			errorMessage: "Username Already Taken",
		},
		{
			inputJSON:    `{"username":"pet", "email": "kangmail.com", "password": "password"}`,
			statusCode:   422,
			errorMessage: "Invalid Email",
		},
		{
			inputJSON:    `{"username":"", "email": "kan@gmail.com", "password": "password"}`,
			statusCode:   422,
			errorMessage: "Required Username",
		},
		{
			inputJSON:    `{"username":"kan", "email": "", "password": "password"}`,
			statusCode:   422,
			errorMessage: "Required Email",
		},
		{
			inputJSON:    `{"username":"kan", "email": "kan@gmail.com", "password": ""}`,
			statusCode:   422,
			errorMessage: "Required Password",
		},
	}
	for _, v := range samples {
		req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateUser)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)

		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 201 {
			assert.Equal(t, responseMap["username"], v.username)
			assert.Equal(t, responseMap["email"], v.email)
		}

		if v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}

}

func TestGetUsers(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	_, err = seedUsers()
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetUsers)
	handler.ServeHTTP(rr, req)

	var users []models.User
	err = json.Unmarshal([]byte(rr.Body.String()), &users)

	if err != nil {
		log.Fatalf("Cannot convert to json: %v\n", err)
	}

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(users), 2)
}

func TestGetUserByID(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatal(err)
	}
	userSample := []struct {
		id           string
		statusCode   int
		username     string
		email        string
		errorMessage string
	}{
		{
			id:         strconv.Itoa(int(user.ID)),
			statusCode: 200,
			username:   user.Username,
			email:      user.Email,
		},
		{
			id:         "unknown",
			statusCode: 400,
		},
	}
	for _, v := range userSample {
		req, err := http.NewRequest("GET", "/users", nil)
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetUser)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, user.Username, responseMap["username"])
			assert.Equal(t, user.Email, responseMap["email"])
		}
	}
}

func TestUpdateUser(t *testing.T) {
	var AuthEmail, AuthPassword string
	var AuthID uint32

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	users, err := seedUsers() // need atleast two users to check the update
	if err != nil {
		log.Fatalf("Error seeding user: %v \n", err)
	}
	// get only first user
	for _, user := range users {
		if user.ID == 2 {
			continue
		}
		AuthID = user.ID
		AuthEmail = user.Email
		AuthPassword = "password" // password in the database is already encrypted
	}
	// Login the user and get authentication token
	token, err := server.SignIn(AuthEmail, AuthPassword)
	if err != nil {
		log.Fatalf("cannot login: %v \n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)
	samples := []struct {
		id             string
		updateJSON     string
		statusCode     int
		updateUsername string
		updateEmail    string
		tokenGiven     string
		errorMessage   string
	}{
		{
			// convert int32 to int first before converting to string
			id:             strconv.Itoa(int(AuthID)),
			updateJSON:     `{"username":"Grand", "email": "grand@gmail.com", "password":"password"}`,
			statusCode:     200,
			updateUsername: "Grand",
			updateEmail:    "grand@gmail.com",
			tokenGiven:     tokenString,
			errorMessage:   "",
		},
		{
			// when password field is empty
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"username":"Woman", "email": "woman@gmail.com", "password":""}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Password",
		},
		{
			// when no token is passed
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"username":"Man", "email": "man@gmail.com", "password":"password"}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// when incorrect token is passed
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"username":"Man", "email": "man@gmail.com", "password":"password"}`,
			statusCode:   401,
			tokenGiven:   "Some incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			// if already email is taken
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"username":"Frank", "email": "kenny@gmail.com", "password":"password"}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Email Already Taken",
		},
		{
			// if already Username is taken
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"username":"Kenny Morris", "email": "grand@gmail.com", "password":"password"}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Username Already Taken",
		},
		{
			// if already email is taken
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"username":"kan", "email": "kennygmail.com", "password":"password"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Invalid Email",
		},
		{
			// if already email is taken
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"username":"", "email": "kan@gmail.com", "password":"password"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Username",
		},
		{
			// if already email is taken
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"username":"kan", "email": "", "password":"password"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Email",
		},
		{
			id:         "unknown",
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			// when using invalid token
			id:           strconv.Itoa(int(2)),
			updateJSON:   `{"username":"Mike", "email": "mike@gmail.com", "password":"password"}`,
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range samples {
		req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(v.updateJSON))

		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdateUser)
		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)
		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["username"], v.updateUsername)
			assert.Equal(t, responseMap["email"], v.updateEmail)
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestDeleteUser(t *testing.T) {
	var AuthEmail, AuthPassword string
	var AuthID uint32
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	users, err := seedUsers() //need at least two users
	if err != nil {
		log.Fatalf("Error seeding user: %v\n", err)
	}
	// get first and log in
	for _, user := range users {
		if user.ID == 2 {
			continue
		}
		AuthID = user.ID
		AuthEmail = user.Email
		AuthPassword = "password" // password already hashed
	}
	// Login user and get authentication token
	token, err := server.SignIn(AuthEmail, AuthPassword)
	if err != nil {
		log.Fatalf("cannot login : %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	userSample := []struct {
		id           string
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			// convert int32 to int first before converting to string
			id:           strconv.Itoa(int(AuthID)),
			tokenGiven:   tokenString,
			statusCode:   204,
			errorMessage: "",
		},
		{
			// When no token is given
			id:           strconv.Itoa(int(AuthID)),
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is given
			id:           strconv.Itoa(int(AuthID)),
			tokenGiven:   "Incorrect token given here",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// Unknown ID
			id:         "unknown",
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			// Trying to use another users token
			id:           strconv.Itoa(int(2)),
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range userSample {
		req, err := http.NewRequest("GET", "/users", nil)
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteUser)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 401 && v.errorMessage != "" {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)

			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
