package tests

import (
	"GolangStore/handlers"
	"GolangStore/middleware"
	"GolangStore/models"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

/* =============================== MODELS TESTS =============================== */
// Test the validity of different combinations of username/password
func TestUserValidity(t *testing.T) {
	if !models.IsUserValid("user1", "pass1") {
		t.Fail()
	}

	if models.IsUserValid("user2", "pass1") {
		t.Fail()
	}

	if models.IsUserValid("user1", "") {
		t.Fail()
	}

	if models.IsUserValid("", "pass1") {
		t.Fail()
	}

	if models.IsUserValid("User1", "pass1") {
		t.Fail()
	}
}

// Test if a new user can be registered with valid username/password
func TestValidUserRegistration(t *testing.T) {
	saveLists()

	u, err := models.RegisterNewUser("newuser", "newpass")

	if err != nil || u.Username == "" {
		t.Fail()
	}

	restoreLists()
}

// Test that a new user cannot be registered with invalid username/password
func TestInvalidUserRegistration(t *testing.T) {
	saveLists()

	// Try to register a user with a used username
	u, err := models.RegisterNewUser("user1", "pass1")

	if err == nil || u != nil {
		t.Fail()
	}

	// Try to register with a blank password
	u, err = models.RegisterNewUser("newuser", "")

	if err == nil || u != nil {
		t.Fail()
	}

	restoreLists()
}

// Test the function that checks for username availability
func TestUsernameAvailability(t *testing.T) {
	saveLists()

	// This username should be available
	if !models.IsUsernameAvailable("newuser") {
		t.Fail()
	}

	// This username should not be available
	if models.IsUsernameAvailable("user1") {
		t.Fail()
	}

	// Register a new user
	models.RegisterNewUser("newuser", "newpass")

	// This newly registered username should not be available
	if models.IsUsernameAvailable("newuser") {
		t.Fail()
	}

	restoreLists()
}

/* =============================== HANDLERS TESTS =============================== */

// Test that a GET request to the login page returns
// an HTTP error with code 401 for an authenticated user
func TestShowLoginPageAuthenticated(t *testing.T) {
	// Create a response recorder
	w := httptest.NewRecorder()

	// Get a new router
	r := getRouter(true)

	// Set the token cookie to simulate an authenticated user
	http.SetCookie(w, &http.Cookie{Name: "token", Value: "123"})

	// Define the route similar to its definition in the routes file
	r.GET("/u/login", middleware.EnsureNotLoggedIn(), handlers.ShowLoginPage)

	// Create a request to send to the above route
	res := w.Result()
	defer res.Body.Close()
	req, _ := http.NewRequest("GET", "/u/login", nil)
	req.Header = http.Header{"Cookie": res.Header["Set-Cookie"]}

	// Create the service and process the above request.
	r.ServeHTTP(w, req)

	// Test that the http status code is 401
	if w.Code != http.StatusUnauthorized {
		t.Fail()
	}
}

// Test that a GET request to the login page returns the login page with
// the HTTP code 200 for an unauthenticated user
func TestShowLoginPageUnauthenticated(t *testing.T) {
	r := getRouter(true)

	// Define the route similar to its definition in the routes file
	r.GET("/u/login", middleware.EnsureNotLoggedIn(), handlers.ShowLoginPage)

	// Create a request to send to the above route
	req, _ := http.NewRequest("GET", "/u/login", nil)

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		// Test that the http status code is 200
		statusOK := w.Code == http.StatusOK

		// Test that the page title is "Login"
		p, err := ioutil.ReadAll(w.Body)
		pageOK := err == nil && strings.Index(string(p), "<title>Login</title>") > 0

		return statusOK && pageOK
	})
}

// Test that a POST request to the login route returns
// an HTTP error with code 401 for an authenticated user
func TestLoginAuthenticated(t *testing.T) {
	// Create a response recorder
	w := httptest.NewRecorder()

	// Get a new router
	r := getRouter(true)

	// Set the token cookie to simulate an authenticated user
	http.SetCookie(w, &http.Cookie{Name: "token", Value: "123"})

	// Define the route similar to its definition in the routes file
	r.POST("/u/login", middleware.EnsureNotLoggedIn(), handlers.PerformLogin)

	// Create a request to send to the above route
	res := w.Result()
	defer res.Body.Close()
	loginPayload := getLoginPOSTPayload()
	req, _ := http.NewRequest("POST", "/u/login", strings.NewReader(loginPayload))
	req.Header = http.Header{"Cookie": res.Header["Set-Cookie"]}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(loginPayload)))

	// Create the service and process the above request.
	r.ServeHTTP(w, req)

	// Test that the http status code is 401
	if w.Code != http.StatusUnauthorized {
		t.Fail()
	}
}

// Test that a POST request to login returns a success message for
// an unauthenticated user
func TestLoginUnauthenticated(t *testing.T) {
	// Create a response recorder
	w := httptest.NewRecorder()

	// Get a new router
	r := getRouter(true)

	// Define the route similar to its definition in the routes file
	r.POST("/u/login", middleware.EnsureNotLoggedIn(), handlers.PerformLogin)

	// Create a request to send to the above route
	loginPayload := getLoginPOSTPayload()
	req, _ := http.NewRequest("POST", "/u/login", strings.NewReader(loginPayload))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(loginPayload)))

	// Create the service and process the above request.
	r.ServeHTTP(w, req)

	// Test that the http status code is 200
	if w.Code != http.StatusOK {
		t.Fail()
	}

	// Test that the page title is "Successful Login"
	// You can carry out a lot more detailed tests using libraries that can
	// parse and process HTML pages
	p, err := ioutil.ReadAll(w.Body)
	if err != nil || strings.Contains(string(p), "<title>Successful Login</title>") != true {
		t.Fail()
	}
}

// Test that a POST request to login returns an error when using
// incorrect credentials
func TestLoginUnauthenticatedIncorrectCredentials(t *testing.T) {
	// Create a response recorder
	w := httptest.NewRecorder()

	// Get a new router
	r := getRouter(true)

	// Define the route similar to its definition in the routes file
	r.POST("/u/login", middleware.EnsureNotLoggedIn(), handlers.PerformLogin)

	// Create a request to send to the above route
	loginPayload := getRegistrationPOSTPayload()
	req, _ := http.NewRequest("POST", "/u/login", strings.NewReader(loginPayload))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(loginPayload)))

	// Create the service and process the above request.
	r.ServeHTTP(w, req)

	// Test that the http status code is 200
	if w.Code != http.StatusBadRequest {
		t.Fail()
	}
}

// Test that a GET request to the registration page returns
// an HTTP error with code 401 for an authenticated user
func TestShowRegistrationPageAuthenticated(t *testing.T) {
	// Create a response recorder
	w := httptest.NewRecorder()

	// Get a new router
	r := getRouter(true)

	// Set the token cookie to simulate an authenticated user
	http.SetCookie(w, &http.Cookie{Name: "token", Value: "123"})

	// Define the route similar to its definition in the routes file
	r.GET("/u/register", middleware.EnsureNotLoggedIn(), handlers.ShowRegistrationPage)

	// Create a request to send to the above route
	res := w.Result()
	defer res.Body.Close()
	req, _ := http.NewRequest("GET", "/u/register", nil)
	req.Header = http.Header{"Cookie": res.Header["Set-Cookie"]}

	// Create the service and process the above request.
	r.ServeHTTP(w, req)

	// Test that the http status code is 401
	if w.Code != http.StatusUnauthorized {
		t.Fail()
	}
}

// Test that a GET request to the registration page returns the registration
// page with the HTTP code 200 for an unauthenticated user
func TestShowRegistrationPageUnauthenticated(t *testing.T) {
	r := getRouter(true)

	// Define the route similar to its definition in the routes file
	r.GET("/u/register", middleware.EnsureNotLoggedIn(), handlers.ShowRegistrationPage)

	// Create a request to send to the above route
	req, _ := http.NewRequest("GET", "/u/register", nil)

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		// Test that the http status code is 200
		statusOK := w.Code == http.StatusOK

		// Test that the page title is "Login"
		p, err := ioutil.ReadAll(w.Body)
		pageOK := err == nil && strings.Index(string(p), "<title>Register</title>") > 0

		return statusOK && pageOK
	})
}

// Test that a POST request to the registration route returns
// an HTTP error with code 401 for an authenticated user
func TestRegisterAuthenticated(t *testing.T) {
	// Create a response recorder
	w := httptest.NewRecorder()

	// Get a new router
	r := getRouter(true)

	// Set the token cookie to simulate an authenticated user
	http.SetCookie(w, &http.Cookie{Name: "token", Value: "123"})

	// Define the route similar to its definition in the routes file
	r.POST("/u/register", middleware.EnsureNotLoggedIn(), handlers.Register)

	// Create a request to send to the above route
	registrationPayload := getRegistrationPOSTPayload()
	res := w.Result()
	defer res.Body.Close()
	req, _ := http.NewRequest("POST", "/u/register", strings.NewReader(registrationPayload))
	req.Header = http.Header{"Cookie": res.Header["Set-Cookie"]}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(registrationPayload)))

	// Create the service and process the above request.
	r.ServeHTTP(w, req)

	// Test that the http status code is 401
	if w.Code != http.StatusUnauthorized {
		t.Fail()
	}
}

// Test that a POST request to register returns a success message for
// an unauthenticated user
func TestRegisterUnauthenticated(t *testing.T) {
	// Create a response recorder
	w := httptest.NewRecorder()

	// Get a new router
	r := getRouter(true)

	// Define the route similar to its definition in the routes file
	r.POST("/u/register", middleware.EnsureNotLoggedIn(), handlers.Register)

	// Create a request to send to the above route
	registrationPayload := getRegistrationPOSTPayload()
	req, _ := http.NewRequest("POST", "/u/register", strings.NewReader(registrationPayload))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(registrationPayload)))

	// Create the service and process the above request.
	r.ServeHTTP(w, req)

	// Test that the http status code is 200
	if w.Code != http.StatusOK {
		t.Fail()
	}

	// Test that the page title is "Successful registration &amp; Login"
	// You can carry out a lot more detailed tests using libraries that can
	// parse and process HTML pages
	p, err := ioutil.ReadAll(w.Body)
	if err != nil || strings.Contains(string(p), "<title>Successful registration &amp; Login</title>") != true {
		t.Fail()
	}
}

// Test that a POST request to register returns a an error when
// the username is already in use
func TestRegisterUnauthenticatedUnavailableUsername(t *testing.T) {
	// Create a response recorder
	w := httptest.NewRecorder()

	// Get a new router
	r := getRouter(true)

	// Define the route similar to its definition in the routes file
	r.POST("/u/register", middleware.EnsureNotLoggedIn(), handlers.Register)

	// Create a request to send to the above route
	registrationPayload := getLoginPOSTPayload()
	req, _ := http.NewRequest("POST", "/u/register", strings.NewReader(registrationPayload))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(registrationPayload)))

	// Create the service and process the above request.
	r.ServeHTTP(w, req)

	// Test that the http status code is 400
	if w.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func getLoginPOSTPayload() string {
	params := url.Values{}
	params.Add("username", "user1")
	params.Add("password", "pass1")

	return params.Encode()
}

func getRegistrationPOSTPayload() string {
	params := url.Values{}
	params.Add("username", "u1")
	params.Add("password", "p1")

	return params.Encode()
}
