package repository

import (
	"blogpost/middleware"
	"blogpost/models"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type DbConnection struct {
	DB *gorm.DB
}

type Operations interface {
	AddUser(user *models.User) error
	Login(email, password string) (string, error)
	GetRoleID(user *models.User) error
	AddPost(*models.Post) error
	GetPostID(post *models.Post) error
	SearchAllPost(post *[]models.Post) error
	UpdatePostByID(ID string, mail string, data map[string]interface{}) (*models.Post, error)
	DeletePostByID(mail string, PostID string, post *models.Post) error
	GetPostBasedOnRoleID(mail string, post *[]models.Post) error
	GetPostBasedOnCategory(category string, post *[]models.Post) error
	GetAllCategory(Post *[]models.Post) error
	GetPostStatistics(post *models.Post, postCount *int64) error
	AddComments(mail string, comment *models.Comments) error
	UpdateCommentByID(mail string, commentID string, data map[string]interface{}) error
}

func NewDbConnection(db *gorm.DB) *DbConnection {
	return &DbConnection{DB: db}
}

func (db *DbConnection) AddUser(user *models.User) error {
	user.ID = uuid.New()

	if err := db.DB.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (db *DbConnection) Login(mail, password string) (string, error) {
	if mail == "" || password == "" {
		return "", fmt.Errorf("mail or password can't be empty")
	}

	var checkingUser models.User

	if err := db.DB.Debug().First(&checkingUser, "mail=?", mail).Error; err != nil {
		return "", err
	}

	err := bcrypt.CompareHashAndPassword([]byte(checkingUser.Password), []byte(password))
	if err != nil {
		return "", err
	}

	var token string

	if checkingUser.Role == "admin" {
		token, err = middleware.AdminToken(mail, checkingUser.ID)
		if err != nil {
			return "", nil
		}
	} else if checkingUser.Role == "user" {
		token, err = middleware.MemberToken(mail, checkingUser.ID)
		if err != nil {
			return "", nil
		}
	} else {
		return "", fmt.Errorf("invalid user role")
	}

	return token, nil
}

func (db *DbConnection) GetRoleID(user *models.User) error {
	if err := db.DB.Debug().First(&user, "mail=?", user.Mail).Error; err != nil {
		return fmt.Errorf("invalid MailID!!! kindly Check it")
	}
	return nil
}

// ----------------------------------------------POST---------------------------------------------------------------------------
func (db *DbConnection) AddPost(post *models.Post) error {
	var checkingUser models.User

	if err := db.DB.Debug().First(&checkingUser, "id=?", post.RoleID).Error; err != nil {
		return fmt.Errorf("invalid RoleID!!! kindly Check it")
	}
	post.ID = uuid.New()
	post.PostDate = time.Now()
	if err := db.DB.Create(&post).Error; err != nil {
		return err
	}
	return nil
}

// Get post Id by title
func (db *DbConnection) GetPostID(post *models.Post) error {
	if err := db.DB.Debug().First(&post, "title=?", post.Title).Error; err != nil {
		return fmt.Errorf("invalid Title!!! kindly Check it")
	}
	return nil
}

// Search all the posts
func (db *DbConnection) SearchAllPost(post *[]models.Post) error {
	if err := db.DB.Debug().Select("id", "role_id", "category", "title", "description", "post_date", "comment_count").Find(&post).Error; err != nil {
		return err
	}

	return nil
}

// Update the post content
func (db *DbConnection) UpdatePostByID(PostID string, mail string, data map[string]interface{}) (*models.Post, error) {
	post := models.Post{}
	user := models.User{}

	if mail == "" {
		return nil, fmt.Errorf("mailID can not be empty")
	}

	if err := db.DB.Debug().Where("mail=?", mail).Where("role=?", "admin").First(&user).Error; err != nil {
		return nil, fmt.Errorf("unauthorized")
	}

	if PostID == "" {
		return nil, fmt.Errorf("PostID can not be empty")
	}

	if err := db.DB.Debug().Where("role_id=?", user.ID).First(&post, "id=?", PostID).Error; err != nil {
		return nil, err
	}

	for d, value := range data {
		if err := db.DB.Debug().Model(&post).Where("id=?", PostID).Update(d, value).Error; err != nil {
			return nil, err
		}
	}

	return &post, nil
}

// Delete the post by Id dbOperation
func (db *DbConnection) DeletePostByID(mail string, PostID string, post *models.Post) error {
	user := models.User{}

	if mail == "" {
		return fmt.Errorf("mailID can not be empty")
	}

	fmt.Println("Mail", mail)

	if err := db.DB.Debug().Where("mail=?", mail).Where("role=?", "admin").First(&user).Error; err != nil {
		return fmt.Errorf("unauthorized")
	}

	if PostID == "" {
		return fmt.Errorf("PostID can not be empty")
	}

	if err := db.DB.Debug().Where("role_id=?", user.ID).First(&post, "id=?", PostID).Error; err != nil {
		return err
	}

	if err := db.DB.Debug().Delete(&post, "ID=?", PostID).Error; err != nil {
		return err
	}

	return nil
}

// to get the post based on role id db operation
func (db *DbConnection) GetPostBasedOnRoleID(mail string, post *[]models.Post) error {
	user := models.User{}
	if mail == "" {
		return fmt.Errorf("mailID can not be empty")
	}

	fmt.Println("Mail", mail)

	if err := db.DB.Debug().Where("mail=?", mail).Where("role=?", "admin").First(&user).Error; err != nil {
		return fmt.Errorf("unauthorized")
	}

	if err := db.DB.Debug().Where("role_id=?", user.ID).Find(&post).Error; err != nil {
		return err
	}
	return nil
}

// to get the post based on category db operation
func (db *DbConnection) GetPostBasedOnCategory(category string, post *[]models.Post) error {
	if err := db.DB.Debug().Where("category=?", category).Find(&post).Error; err != nil {
		return err
	}
	return nil
}

// to get the post based on category db operation
func (db *DbConnection) GetAllCategory(post *[]models.Post) error {
	if err := db.DB.Debug().Select("category").Find(&post).Error; err != nil {
		return err
	}
	return nil
}

// to get the post statistics db operation
func (db *DbConnection) GetPostStatistics(post *models.Post, postCount *int64) error {
	if err := db.DB.Debug().Model(&post).Count(postCount).Error; err != nil {
		return err
	}
	return nil
}

// ---------------------------------Comments---------------------------------------------------------------------------
// to add comments db operation
func (db *DbConnection) AddComments(mail string, comment *models.Comments) error {
	comment.ID = uuid.New()

	if mail == "" {
		return fmt.Errorf("mailID can not be empty")
	}

	if err := db.DB.Debug().Where("mail=?", mail).Where("role=?", "user").First(&models.User{}).Error; err != nil {
		return fmt.Errorf("unauthorized")
	}

	if err := db.DB.Create(&comment).Error; err != nil {
		return err
	}

	var commentCount int64
	if err := db.DB.Debug().Model(&comment).Where("post_id = ?", comment.PostID).Count(&commentCount).Error; err != nil {
		return err
	}

	if err := db.DB.Model(&models.Post{}).Where("id", comment.PostID).Update("comment_count", commentCount).Error; err != nil {
		return err
	}
	return nil
}

func (db *DbConnection) UpdateCommentByID(mail string, commentID string, data map[string]interface{}) error {
	var comment models.Comments
	if mail == "" {
		return fmt.Errorf("mailID can not be empty")
	}

	if err := db.DB.Debug().Where("mail=?", mail).Where("role=?", "user").First(&models.User{}).Error; err != nil {
		return fmt.Errorf("unauthorized")
	}

	if err := db.DB.Debug().Where("id=?", commentID).First(&comment).Error; err != nil {
		return err
	}

	fmt.Println("Data in update comment in handler", data)

	for d, value := range data {
		if err := db.DB.Model(&models.Comments{}).Where("id=?", commentID).Update(d, value).Error; err != nil {
			return err
		}
	}
	return nil
}
