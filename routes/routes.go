package routes

import (
	"GolangStore/handlers"
	"GolangStore/middleware"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func GinSetup() {
	// Set Gin to production mode
	//gin.SetMode(gin.ReleaseMode)

	// Set the router as the default one provided by Gin
	router = gin.Default()

	// Process the templates at the start so that they don't have to be loaded
	// from the disk again. This makes serving HTML pages very fast.
	router.LoadHTMLGlob("templates/*")

	// Initialize the routes
	initializeRoutes()

	// Start serving the application
	router.Run()
}

func initializeRoutes() {

	// Use the setUserStatus middleware for every route to set a flag
	// indicating whether the request was from an authenticated user or not
	router.Use(middleware.SetUserStatus())

	// Handle the index route
	router.GET("/", handlers.ShowIndexPage)
	/* Grouping routes together allows you to apply middleware on
	all routes in a group instead of doing so separately for each route.
	In the above snippet, the first route will use the showRegistrationPage
	function to display the registration page at the /u/register path.
	The second route will handle all the POST requests to the same path,
	making use of the registerroute handler. */
	// Group user related routes together
	userRoutes := router.Group("/u")
	{
		// Handle the GET requests at /u/login
		// Show the login page
		// Ensure that the user is not logged in by using the middleware
		userRoutes.GET("/login", middleware.EnsureNotLoggedIn(), handlers.ShowLoginPage)

		// Handle POST requests at /u/login
		// Ensure that the user is not logged in by using the middleware
		userRoutes.POST("/login", middleware.EnsureNotLoggedIn(), handlers.PerformLogin)

		// Handle GET requests at /u/logout
		// Ensure that the user is logged in by using the middleware
		userRoutes.GET("/logout", middleware.EnsureLoggedIn(), handlers.Logout)

		// Handle the GET requests at /u/register
		// Show the registration page
		// Ensure that the user is not logged in by using the middleware
		userRoutes.GET("/register", middleware.EnsureNotLoggedIn(), handlers.ShowRegistrationPage)

		// Handle POST requests at /u/register
		// Ensure that the user is not logged in by using the middleware
		userRoutes.POST("/register", middleware.EnsureNotLoggedIn(), handlers.Register)
	}

	// Group article related routes together
	articleRoutes := router.Group("/article")
	{
		// Handle GET requests at /article/view/some_article_id
		articleRoutes.GET("/view/:article_id", handlers.GetArticle)

		/* Handle the GET requests at /article/create Show the article creation
		page. Ensure that the user is logged in by using the middleware*/
		articleRoutes.GET("/create", middleware.EnsureLoggedIn(), handlers.ShowArticleCreationPage)

		/* Handle POST requests at /article/create Ensure that the user is
		logged in by using the middleware */
		articleRoutes.POST("/create", middleware.EnsureLoggedIn(), handlers.CreateArticle)
	}

}
