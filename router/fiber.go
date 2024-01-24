package router

import (
	"blogpost/handler"
	"blogpost/middleware"
	"blogpost/repository"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func Routing(db *repository.DbConnection) {
	h := handler.Newhandler(db)

	app := fiber.New()
	logger := log.New(log.Writer(), "Blog-Post ", log.LstdFlags)

	// Middleware to log requests
	app.Use(func(c *fiber.Ctx) error {
		logger.Printf("Incoming Request: %s %s", c.Method(), c.Path())
		return c.Next()
	})

	routes := app.Group("/blogpost/v1/")
	routes.Post("/signup", h.AddUser)
	routes.Post("/login/*", h.Login)
	routes.Get("/get-role-id", h.GetRoleID)
	routes.Get("/search-all-posts", h.SearchAllPost)
	routes.Get("/get-all-category", h.GetAllCategory)
	routes.Get("/get-post-by-category", h.GetPostBasedOnCategory)
	routes.Get("/get-post-statistics", h.GetPostStatistics)
	routes.Get("/get-comment-based-on-post", h.GetCommentsBasedOnPostID)

	adminroutes := app.Group("/blogpost/v1/admin")
	adminroutes.Post("/add-post", middleware.AdminAuthorize([]byte("secret"), h.AddPost))
	adminroutes.Get("/get-posts-by-role-id", middleware.AdminAuthorize([]byte("secret"), h.GetPostBasedOnRoleID))
	adminroutes.Put("/update-post-by-id", middleware.AdminAuthorize([]byte("secret"), h.UpdatePostByID))
	adminroutes.Delete("/delete-post-by-id", middleware.AdminAuthorize([]byte("secret"), h.DeletePostByID))

	memberRoutes := app.Group("/blogpost/v1/member")
	memberRoutes.Get("/get-post-by-id", middleware.MemberAuthorize([]byte("secret"), h.GetPostBasedOnPostID))
	memberRoutes.Post("/add-comment", middleware.MemberAuthorize([]byte("secret"), h.AddComments))
	memberRoutes.Put("/update-comment", middleware.MemberAuthorize([]byte("secret"), h.UpdateCommentByID))
	memberRoutes.Delete("/delete-comment", middleware.MemberAuthorize([]byte("secret"), h.DeleteCommentByID))
	memberRoutes.Get("/get-comment-based-on-user", middleware.MemberAuthorize([]byte("secret"), h.GetCommentsBasedOnUser))

	logger.Println("Server Started")
	if err := app.Listen(":8000"); err != nil {
		logger.Println("Server Ended")
		fmt.Println("Ended:", err)
		return
	}
}
