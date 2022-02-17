package tests

import (
	"encoding/json"
	"encoding/xml"
	"gin-web-framework/handlers"
	"gin-web-framework/middleware"
	"gin-web-framework/models"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

/* =============================== MODELS TESTS =============================== */
// Test the function that fetches all articles
func TestGetAllArticles(t *testing.T) {
	alist := models.GetAllArticles()

	// Check that the length of the list of articles returned is the
	// same as the length of the global variable holding the list
	if len(alist) != len(models.ArticleList) {
		t.Fail()
	}

	// Check that each member is identical
	for i, v := range alist {
		if v.Content != models.ArticleList[i].Content ||
			v.ID != models.ArticleList[i].ID ||
			v.Title != models.ArticleList[i].Title {

			t.Fail()
			break
		}
	}
}

// Test the function that fetche an Article by its ID
func TestGetArticleByID(t *testing.T) {
	a, err := models.GetArticleByID(1)

	if err != nil || a.ID != 1 || a.Title != "Article 1" || a.Content != "Article 1 body" {
		t.Fail()
	}
}

// Test the function that creates a new article
func TestCreateNewArticle(t *testing.T) {
	// get the original count of articles
	originalLength := len(models.GetAllArticles())

	// add another article
	a, err := models.CreateNewArticle("New test title", "New test content")

	// get the new count of articles
	allArticles := models.GetAllArticles()
	newLength := len(allArticles)

	if err != nil || newLength != originalLength+1 ||
		a.Title != "New test title" || a.Content != "New test content" {

		t.Fail()
	}
}

/* =============================== HANDLERS TESTS =============================== */
// Test that a GET request to the home page returns the home page with
// the HTTP code 200 for an unauthenticated user
func TestShowIndexPageUnauthenticated(t *testing.T) {
	r := getRouter(true)

	r.GET("/", handlers.ShowIndexPage)

	// Create a request to send to the above route
	req, _ := http.NewRequest("GET", "/", nil)

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		// Test that the http status code is 200
		statusOK := w.Code == http.StatusOK

		// Test that the page title is "Home Page"
		// You can carry out a lot more detailed tests using libraries that can
		// parse and process HTML pages
		p, err := ioutil.ReadAll(w.Body)
		pageOK := err == nil && strings.Index(string(p), "<title>Home Page</title>") > 0

		return statusOK && pageOK
	})
}

// Test that a GET request to the home page returns the home page with
// the HTTP code 200 for an authenticated user
func TestShowIndexPageAuthenticated(t *testing.T) {
	// Create a response recorder
	w := httptest.NewRecorder()

	// Get a new router
	r := getRouter(true)

	// Set the token cookie to simulate an authenticated user
	http.SetCookie(w, &http.Cookie{Name: "token", Value: "123"})

	// Define the route similar to its definition in the routes file
	r.GET("/", handlers.ShowIndexPage)

	// Create a request to send to the above route
	res := w.Result()
	defer res.Body.Close()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header = http.Header{"Cookie": res.Header["Set-Cookie"]}
	// Create the service and process the above request.
	r.ServeHTTP(w, req)

	// Test that the http status code is 200
	if w.Code != http.StatusOK {
		t.Fail()
	}

	// Test that the page title is "Home Page"
	// You can carry out a lot more detailed tests using libraries that can
	// parse and process HTML pages
	p, err := ioutil.ReadAll(w.Body)
	if err != nil || strings.Contains(string(p), "<title>Home Page</title>") != true {
		t.Fail()
	}
}

//if err != nil || strings.Contains(string(p), "<title>Home Page</title>") {
// Test that a GET request to an article page returns the article page with
// the HTTP code 200 for an unauthenticated user
func TestArticleUnauthenticated(t *testing.T) {
	r := getRouter(true)

	// Define the route similar to its definition in the routes file
	r.GET("/article/view/:article_id", handlers.GetArticle)

	// Create a request to send to the above route
	req, _ := http.NewRequest("GET", "/article/view/1", nil)

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		// Test that the http status code is 200
		statusOK := w.Code == http.StatusOK

		// Test that the page title is "Article 1"
		// You can carry out a lot more detailed tests using libraries that can
		// parse and process HTML pages
		p, err := ioutil.ReadAll(w.Body)
		pageOK := err == nil && strings.Index(string(p), "<title>Article 1</title>") > 0

		return statusOK && pageOK
	})
}

// Test that a GET request to an article page returns the article page with
// the HTTP code 200 for an authenticated user
func TestArticleAuthenticated(t *testing.T) {
	// Create a response recorder
	w := httptest.NewRecorder()

	// Get a new router
	r := getRouter(true)

	// Set the token cookie to simulate an authenticated user
	http.SetCookie(w, &http.Cookie{Name: "token", Value: "123"})

	// Define the route similar to its definition in the routes file
	r.GET("/article/view/:article_id", handlers.GetArticle)

	// Create a request to send to the above route
	res := w.Result()
	defer res.Body.Close()
	req, _ := http.NewRequest("GET", "/article/view/1", nil)
	req.Header = http.Header{"Cookie": res.Header["Set-Cookie"]}

	// Create the service and process the above request.
	r.ServeHTTP(w, req)

	// Test that the http status code is 200
	if w.Code != http.StatusOK {
		t.Fail()
	}

	// Test that the page title is "Article 1"
	// You can carry out a lot more detailed tests using libraries that can
	// parse and process HTML pages
	p, err := ioutil.ReadAll(w.Body)
	if err != nil || strings.Contains(string(p), "<title>Article 1</title>") != true {
		t.Fail()
	}

}

// Test that a GET request to the home page returns the list of articles
// in JSON format when the Accept header is set to application/json
func TestArticleListJSON(t *testing.T) {
	r := getRouter(true)

	// Define the route similar to its definition in the routes file
	r.GET("/", handlers.ShowIndexPage)

	// Create a request to send to the above route
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept", "application/json")

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		// Test that the http status code is 200
		statusOK := w.Code == http.StatusOK

		// Test that the response is JSON which can be converted to
		// an array of Article structs
		p, err := ioutil.ReadAll(w.Body)
		if err != nil {
			return false
		}
		var articles []models.Article
		err = json.Unmarshal(p, &articles)

		return err == nil && len(articles) >= 2 && statusOK
	})
}

// Test that a GET request to an article page returns the article in XML
// format when the Accept header is set to application/xml
func TestArticleXML(t *testing.T) {
	r := getRouter(true)

	// Define the route similar to its definition in the routes file
	r.GET("/article/view/:article_id", handlers.GetArticle)

	// Create a request to send to the above route
	req, _ := http.NewRequest("GET", "/article/view/1", nil)
	req.Header.Add("Accept", "application/xml")

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		// Test that the http status code is 200
		statusOK := w.Code == http.StatusOK

		// Test that the response is JSON which can be converted to
		// an array of Article structs
		p, err := ioutil.ReadAll(w.Body)
		if err != nil {
			return false
		}
		var a models.Article
		err = xml.Unmarshal(p, &a)

		return err == nil && a.ID == 1 && len(a.Title) >= 0 && statusOK
	})
}

// Test that a GET request to the article creation page returns the
// article creation page with the HTTP code 200 for an authenticated user
func TestArticleCreationPageAuthenticated(t *testing.T) {
	// Create a response recorder
	w := httptest.NewRecorder()

	// Get a new router
	r := getRouter(true)

	// Set the token cookie to simulate an authenticated user
	http.SetCookie(w, &http.Cookie{Name: "token", Value: "123"})

	// Define the route similar to its definition in the routes file
	r.GET("/article/create", middleware.EnsureLoggedIn(), handlers.ShowArticleCreationPage)

	// Create a request to send to the above route
	res := w.Result()
	defer res.Body.Close()
	req, _ := http.NewRequest("GET", "/article/create", nil)
	req.Header = http.Header{"Cookie": res.Header["Set-Cookie"]}

	// Create the service and process the above request.
	r.ServeHTTP(w, req)

	// Test that the http status code is 200
	if w.Code != http.StatusOK {
		t.Fail()
	}

	// Test that the page title is "Create New Article"
	// You can carry out a lot more detailed tests using libraries that can
	// parse and process HTML pages
	p, err := ioutil.ReadAll(w.Body)
	if err != nil || strings.Contains(string(p), "<title>Create New Article</title>") != true {
		t.Fail()
	}

}

// Test that a GET request to the article creation page returns
// an HTTP 401 error for an unauthorized user
func TestArticleCreationPageUnauthenticated(t *testing.T) {
	r := getRouter(true)

	// Define the route similar to its definition in the routes file
	r.GET("/article/create", middleware.EnsureLoggedIn(), handlers.ShowArticleCreationPage)

	// Create a request to send to the above route
	req, _ := http.NewRequest("GET", "/article/create", nil)

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		// Test that the http status code is 401
		return w.Code == http.StatusUnauthorized
	})
}

// Test that a POST request to create an article returns
// an HTTP 200 code along with a success message for an authenticated user
func TestArticleCreationAuthenticated(t *testing.T) {
	// Create a response recorder
	w := httptest.NewRecorder()

	// Get a new router
	r := getRouter(true)

	// Set the token cookie to simulate an authenticated user
	http.SetCookie(w, &http.Cookie{Name: "token", Value: "123"})

	// Define the route similar to its definition in the routes file
	r.POST("/article/create", middleware.EnsureLoggedIn(), handlers.CreateArticle)

	// Create a request to send to the above route
	articlePayload := getArticlePOSTPayload()
	res := w.Result()
	defer res.Body.Close()
	req, _ := http.NewRequest("POST", "/article/create", strings.NewReader(articlePayload))
	req.Header = http.Header{"Cookie": res.Header["Set-Cookie"]}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(articlePayload)))

	// Create the service and process the above request.
	r.ServeHTTP(w, req)

	// Test that the http status code is 200
	if w.Code != http.StatusOK {
		t.Fail()
	}

	// Test that the page title is "Submission Successful"
	// You can carry out a lot more detailed tests using libraries that can
	// parse and process HTML pages
	p, err := ioutil.ReadAll(w.Body)
	if err != nil || strings.Contains(string(p), "<title>Submission Successful</title>") != true {
		t.Fail()
	}

}

// Test that a POST request to create an article returns
// an HTTP 401 error for an unauthorized user
func TestArticleCreationUnauthenticated(t *testing.T) {
	r := getRouter(true)

	// Define the route similar to its definition in the routes file
	r.POST("/article/create", middleware.EnsureLoggedIn(), handlers.CreateArticle)

	// Create a request to send to the above route
	articlePayload := getArticlePOSTPayload()
	req, _ := http.NewRequest("POST", "/article/create", strings.NewReader(articlePayload))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(articlePayload)))

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		// Test that the http status code is 401
		return w.Code == http.StatusUnauthorized
	})
}

func getArticlePOSTPayload() string {
	params := url.Values{}
	params.Add("title", "Test Article Title")
	params.Add("content", "Test Article Content")

	return params.Encode()
}
