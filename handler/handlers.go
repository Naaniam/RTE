package handler

import (
	"blogpost/models"
	"blogpost/repository"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	Repo repository.Operations
}

func Newhandler(db *repository.DbConnection) *Handler {
	return &Handler{Repo: db}
}

// ------------------------------------------------------------USER---------------------------------------------------------------
func (h *Handler) AddUser(c *fiber.Ctx) error {
	user := models.User{}

	// parse requestbody, attach to Post struct
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user.Password = string(hash)

	if err := h.Repo.AddUser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Created User!!! Kindly store/remeber the user id!!!!", "userID": user.ID})
}

// Login
func (h *Handler) Login(c *fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	token, err := h.Repo.Login(email, password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	c.Cookie(&fiber.Cookie{
		Name:  "access_token",
		Value: token,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Loggged in Successfully!!!", "Token": token})
}

// GetRoleID Handler function
func (h *Handler) GetRoleID(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.Repo.GetRoleID(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"RoleID": user.ID})
}

// ----------------------------------------------POSTS-----------------------------------------------------------------------
// AddPost handler function
func (h *Handler) AddPost(c *fiber.Ctx) error {
	post := models.Post{}

	// parse requestbody, attach to Post struct
	if err := c.BodyParser(&post); err != nil {
		fmt.Println("ERROR", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fiber.Map{"error": err.Error()}})
	}

	fmt.Println("POST Data", post)

	if err := h.Repo.AddPost(&post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fiber.Map{"error": err.Error()}})
	}

	return c.Status(fiber.StatusCreated).JSON("Created post!! Kindly note the post id:" + fmt.Sprint(post.ID))
}

// Search all the post handler function
func (h *Handler) SearchAllPost(c *fiber.Ctx) error {
	posts := []models.Post{}

	if err := h.Repo.SearchAllPost(&posts); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(posts)
}

// GetPostID handler function
func (h *Handler) GetPostID(c *fiber.Ctx) error {
	post := models.Post{}

	if err := c.BodyParser(&post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.Repo.GetPostID(&post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"PostID": post.ID})
}

// Update post by ID habdler function
func (h *Handler) UpdatePostByID(c *fiber.Ctx) error {
	data := make(map[string]interface{})
	postID := c.Query("post_id")

	if err := c.BodyParser(&data); err != nil {
		fmt.Println("ERROR", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fiber.Map{"error": err.Error()}})
	}

	cookie := c.Cookies("access_token")
	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fiber.Map{"error": err.Error()}})
	}

	payload := token.Claims.(jwt.MapClaims)

	fmt.Println("Payload", payload)

	post, err := h.Repo.UpdatePostByID(postID, payload["email"].(string), data)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Updated Successfully", "PostID": post.ID, "RoleID": post.RoleID, "Category": post.Category, "Title": post.Title, "Description": post.Description, "PostDate": post.PostDate, "CommentCount": post.CommentCount})
}

func (h *Handler) DeletePostByID(c *fiber.Ctx) error {
	post := models.Post{}
	postID := c.Query("post_id")
	cookie := c.Cookies("access_token")

	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fiber.Map{"error": err.Error()}})
	}

	payload := token.Claims.(jwt.MapClaims)

	fmt.Println("Payload", payload)

	err = h.Repo.DeletePostByID(payload["email"].(string), postID, &post)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "deleted the post Successfully"})
}

// get all the posts based on the role Id
func (h *Handler) GetPostBasedOnRoleID(c *fiber.Ctx) error {
	post := []models.Post{}

	cookie := c.Cookies("access_token")

	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fiber.Map{"error": err.Error()}})
	}

	payload := token.Claims.(jwt.MapClaims)

	fmt.Println("Payload", payload)

	err = h.Repo.GetPostBasedOnRoleID(payload["email"].(string), &post)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Retrived the post Successfully", "post": post})
}

// Get post based Category handler function
func (h *Handler) GetPostBasedOnCategory(c *fiber.Ctx) error {
	posts := []models.Post{}

	category := c.FormValue("category")
	if err := h.Repo.GetPostBasedOnCategory(category, &posts); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Post": posts})
}

// Get post based Category handler function
func (h *Handler) GetAllCategory(c *fiber.Ctx) error {
	posts := []models.Post{}

	if err := h.Repo.GetAllCategory(&posts); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	categories := make([]string, 0)

	for _, post := range posts {
		categories = append(categories, post.Category)
	}

	fmt.Println("Categories", categories)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Categories": categories})
}

func (h *Handler) GetPostStatistics(c *fiber.Ctx) error {
	var postCount int64
	post := models.Post{}

	if err := h.Repo.GetPostStatistics(&post, &postCount); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	fmt.Println("Post count in handlers", postCount)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Total posts": postCount})
}

// -------------------------------------------------Comments--------------------------------------------------------------------
func (h *Handler) AddComments(c *fiber.Ctx) error {
	comment := models.Comments{}
	cookie := c.Cookies("access_token")

	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fiber.Map{"error": err.Error()}})
	}

	payload := token.Claims.(jwt.MapClaims)

	postID := c.Query("post_id")
	roleID := payload["id"].(string)

	comment.PostID = uuid.MustParse(postID)
	comment.RoleID = uuid.MustParse(roleID)

	// parse requestbody, attach to Post struct
	if err := c.BodyParser(&comment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.Repo.AddComments(&comment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Message": "Added comment successfully"})
}
