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
	/* Grouping routes together allows you to apply middleware on all
	   routes in a group instead of doing so separately for each route. */
	/*
		EnsureNotLoggedIn -> Ensures that the user is not logged in by using the middleware function
		EnsureLoggedIn    -> Ensure that the user is logged in by using the middleware
	*/
	// Group user related routes together
	userRoutes := router.Group("/u")
	{
		// Handle the GET requests at /u/login and show the login page
		userRoutes.GET("/login", middleware.EnsureNotLoggedIn(), handlers.ShowLoginPage)
		// Handle POST requests at /u/login
		userRoutes.POST("/login", middleware.EnsureNotLoggedIn(), handlers.PerformLogin)
		// Handle GET requests at /u/logout
		userRoutes.GET("/logout", middleware.EnsureLoggedIn(), handlers.Logout)
		// Handle the GET requests at /u/register and show the registration page
		userRoutes.GET("/register", middleware.EnsureNotLoggedIn(), handlers.ShowRegistrationPage)
		// Handle POST requests at /u/register
		userRoutes.POST("/register", middleware.EnsureNotLoggedIn(), handlers.Register)
	}

	// Group article related routes together
	articleRoutes := router.Group("/article")
	{
		// Handle GET requests at /article/view/some_article_id
		articleRoutes.GET("/view/:article_id", handlers.GetArticle)
		// Handle the GET requests at /article/create Show the article creation
		articleRoutes.GET("/create", middleware.EnsureLoggedIn(), handlers.ShowArticleCreationPage)
		// Handle POST requests at /article/create
		articleRoutes.POST("/create", middleware.EnsureLoggedIn(), handlers.CreateArticle)
		// articleRoutes.GET("/delete/:article_id", middleware.EnsureLoggedIn(), handlers.DeleteArticlePage)
		// articleRoutes.GET("/delete", middleware.EnsureLoggedIn(), handlers.ShowArticleDeletePage)
	}

	router.GET("/products", handlers.IndexPage)

}
