package router

import (
	"blogpost/handler"
	"blogpost/middleware"
	"blogpost/repository"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Routing(db *repository.DbConnection) {
	h := handler.Newhandler(db)
	app := fiber.New()

	routes := app.Group("/blogpost/v1")

	routes.Post("/signup", h.AddUser)
	routes.Post("/login", h.Login)
	//routes.Get("/get-role-id", h.GetRoleID)

	routes.Post("/add-post", middleware.AdminAuthorize([]byte("secret"), h.AddPost))
	routes.Get("/search-all-posts", h.SearchAllPost)
	routes.Get("/get-all-category", h.GetAllCategory)
	routes.Get("/get-posts-by-role-id", middleware.AdminAuthorize([]byte("secret"), h.GetPostBasedOnRoleID))
	routes.Put("/update-post-by-id", middleware.AdminAuthorize([]byte("secret"), h.UpdatePostByID))
	//routes.Get("/get-post-id", middleware.AdminAuthorize([]byte("secret"), h.GetPostID))
	routes.Delete("/delete-post-by-id", middleware.AdminAuthorize([]byte("secret"), h.DeletePostByID))
	routes.Get("/get-post-by-category", h.GetPostBasedOnCategory)
	routes.Get("/get-post-statistics", h.GetPostStatistics)

	routes.Post("/add-comment", middleware.MemberAuthorize([]byte("secret"), h.AddComments))
	routes.Put("/update-comment", middleware.MemberAuthorize([]byte("secret"), h.UpdateCommentByID))

	if err := app.Listen(":8000"); err != nil {
		fmt.Println("Ended:", err)
		return
	}
}
