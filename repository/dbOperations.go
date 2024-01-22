package repository

import (
	"blogpost/middleware"
	"blogpost/models"
	"blogpost/utilities"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type DbConnection struct {
	DB     *gorm.DB
	Logger *log.Logger
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
	GetPostStatistics(post *models.Post, postCount, commentCount *int64) error
	AddComments(mail string, comment *models.Comments) error
	UpdateCommentByID(mail string, commentID string, data map[string]interface{}) error
	DeleteCommentByID(mail string, commentID string, comment *models.Comments) error
	GetCommentsBasedOnUser(mail string, comment *[]models.Comments) error
	GetCommentsBasedOnPostID(postID string, comment *[]models.Comments) error
}

func NewDbConnection(db *gorm.DB, logger *log.Logger) *DbConnection {
	return &DbConnection{DB: db, Logger: logger}
}

func (db *DbConnection) AddUser(user *models.User) error {
	user.ID = uuid.New()

	err := utilities.ValidateStruct(user)
	if err != nil {
		db.Logger.Printf("Error validating the struct")
		return err
	}

	if err := db.DB.Create(&user).Error; err != nil {
		db.Logger.Printf("Error creating user: %v", err)
		return err
	}

	db.Logger.Println("Added user Successfully with ID", user.ID)
	return nil
}

// Login
func (db *DbConnection) Login(mail, password string) (string, error) {
	if mail == "" || password == "" {
		db.Logger.Printf("mail or password can't be empty")
		return "", fmt.Errorf("mail or password can't be empty")
	}

	var checkingUser models.User

	if err := db.DB.Debug().First(&checkingUser, "mail=?", mail).Error; err != nil {
		db.Logger.Printf("%v", err)
		return "", err
	}

	err := bcrypt.CompareHashAndPassword([]byte(checkingUser.Password), []byte(password))
	if err != nil {
		db.Logger.Printf("Error in validating the user: %v", err)
		return "", err
	}

	var token string

	if checkingUser.Role == "admin" {
		token, err = middleware.AdminToken(mail, checkingUser.ID)
		if err != nil {
			db.Logger.Printf("Error creating the token: %v", err)
			return "", nil
		}
	} else if checkingUser.Role == "user" {
		token, err = middleware.MemberToken(mail, checkingUser.ID)
		if err != nil {
			db.Logger.Printf("Error creating the token: %v", err)
			return "", nil
		}
	} else {
		db.Logger.Printf("invalid user role")
		return "", fmt.Errorf("invalid user role")
	}

	db.Logger.Printf("User with mailID %v logged in successfully", mail)
	return token, nil
}

// GetRoleID
func (db *DbConnection) GetRoleID(user *models.User) error {
	if err := db.DB.Debug().First(&user, "mail=?", user.Mail).Error; err != nil {
		db.Logger.Printf("invalid MailID!!! kindly Check it")
		return fmt.Errorf("invalid MailID!!! kindly Check it")
	}

	db.Logger.Println("user role id retrived")
	return nil
}

// ----------------------------------------------POST---------------------------------------------------------------------------
func (db *DbConnection) AddPost(post *models.Post) error {
	var checkingUser models.User

	err := utilities.ValidateStruct(post)
	if err != nil {
		db.Logger.Printf("Error validating the struct")
		return err
	}

	if err := db.DB.Debug().First(&checkingUser, "id=?", post.RoleID).Error; err != nil {
		db.Logger.Println("invalid RoleID! Kindly Check it")
		return fmt.Errorf("invalid RoleID! Kindly Check it")
	}
	post.ID = uuid.New()
	post.PostDate = time.Now()
	if err := db.DB.Create(&post).Error; err != nil {
		db.Logger.Printf("Error creating the post: %v", err)
		return err
	}

	db.Logger.Printf("Added post with ID: %v", post.ID)
	return nil
}

// Get post Id by title
func (db *DbConnection) GetPostID(post *models.Post) error {
	if err := db.DB.Debug().First(&post, "title=?", post.Title).Error; err != nil {
		db.Logger.Printf("invalid Title!!! kindly Check it")
		return fmt.Errorf("invalid Title!!! kindly Check it")
	}

	db.Logger.Printf("Retrive the post with ID: %v", post.ID)
	return nil
}

// Search all the posts
func (db *DbConnection) SearchAllPost(post *[]models.Post) error {
	if err := db.DB.Debug().Select("id", "role_id", "category", "title", "description", "post_date", "comment_count").Find(&post).Error; err != nil {
		db.Logger.Printf("%v", err)
		return err
	}

	db.Logger.Println("Retrived all the post")
	return nil
}

// Update the post content
func (db *DbConnection) UpdatePostByID(PostID string, mail string, data map[string]interface{}) (*models.Post, error) {
	post := models.Post{}
	user := models.User{}

	if mail == "" {
		db.Logger.Printf("mailID can not be empty")
		return nil, fmt.Errorf("mailID can not be empty")
	}

	if err := db.DB.Debug().Where("mail=?", mail).Where("role=?", "admin").First(&user).Error; err != nil {
		db.Logger.Printf("Error creating the token: %v", err)
		return nil, fmt.Errorf("unauthorized")
	}

	if PostID == "" {
		db.Logger.Printf("PostID can not be empty")
		return nil, fmt.Errorf("PostID can not be empty")
	}

	if err := db.DB.Debug().Where("role_id=?", user.ID).First(&post, "id=?", PostID).Error; err != nil {
		db.Logger.Printf("Error creating the token: %v", err)
		return nil, err
	}

	for d, value := range data {
		if err := db.DB.Debug().Model(&post).Where("id=?", PostID).Update(d, value).Error; err != nil {
			db.Logger.Printf("Error creating the token: %v", err)
			return nil, err
		}
	}

	db.Logger.Printf("Updated the post content with ID: %v", post.ID)
	return &post, nil
}

// Delete the post by Id dbOperation
func (db *DbConnection) DeletePostByID(mail string, postID string, post *models.Post) error {
	user := models.User{}

	if mail == "" {
		db.Logger.Printf("mailID can not be empty")
		return fmt.Errorf("mailID can not be empty")
	}

	if err := db.DB.Debug().Where("mail=?", mail).Where("role=?", "admin").First(&user).Error; err != nil {
		db.Logger.Printf("unauthorized")
		return fmt.Errorf("unauthorized")
	}

	if postID == "" {
		db.Logger.Printf("PostID can not be empty")
		return fmt.Errorf("PostID can not be empty")
	}

	if err := db.DB.Debug().Where("role_id=?", user.ID).First(&post, "id=?", postID).Error; err != nil {
		db.Logger.Printf("Error, %v Occured when searching the post with ID: %v", err, postID)
		return err
	}

	if err := db.DB.Debug().Delete(&post, "ID=?", postID).Error; err != nil {
		db.Logger.Printf("Error %v Occured when deleting the post with ID: %v", err, postID)
		return err
	}

	db.Logger.Printf("Deleted the post with ID: %v", post.ID)
	return nil
}

// to get the post based on role id db operation
func (db *DbConnection) GetPostBasedOnRoleID(mail string, post *[]models.Post) error {
	user := models.User{}

	if mail == "" {
		db.Logger.Printf("mailID can not be empty")
		return fmt.Errorf("mailID can not be empty")
	}

	if err := db.DB.Debug().Where("mail=?", mail).Where("role=?", "admin").First(&user).Error; err != nil {
		db.Logger.Printf("Error %v Occured when searching the user with mailID: %v", err, mail)
		return fmt.Errorf("unauthorized")
	}

	if err := db.DB.Debug().Where("role_id=?", user.ID).Find(&post).Error; err != nil {
		db.Logger.Printf("Error %v Occured when searching the post posted by the user with roleID: %v", err, user.ID)
		return err
	}

	db.Logger.Printf("Retrived the posts posted the user with ID:%v", user.ID)
	return nil
}

// to get the post based on category db operation
func (db *DbConnection) GetPostBasedOnCategory(category string, post *[]models.Post) error {
	if err := db.DB.Debug().Where("category=?", category).Find(&post).Error; err != nil {
		db.Logger.Printf("Error %v Occured when searching the post based on the category: %v", err, category)
		return err
	}

	db.Logger.Printf("Retrived the posts based ont the category:%v", category)
	return nil
}

// to get the post based on category db operation
func (db *DbConnection) GetAllCategory(post *[]models.Post) error {
	if err := db.DB.Debug().Select("category").Find(&post).Error; err != nil {
		db.Logger.Printf("Error, %v Occured when retriving all the categories", err)
		return err
	}

	db.Logger.Printf("Retrived all the categories of the post")
	return nil
}

// to get the post statistics db operation
func (db *DbConnection) GetPostStatistics(post *models.Post, postCount, commentCount *int64) error {
	if err := db.DB.Debug().Model(&post).Count(postCount).Error; err != nil {
		db.Logger.Printf("Error, %v Occured when counting the post", err)
		return err
	}

	if err := db.DB.Debug().Model(&models.Comments{}).Count(commentCount).Error; err != nil {
		db.Logger.Printf("Error %v Occured when counting the comment", err)
		return err
	}

	db.Logger.Printf("Post Satistics has been retrived")
	return nil
}

// ---------------------------------Comments---------------------------------------------------------------------------
// to add comments db operation
func (db *DbConnection) AddComments(mail string, comment *models.Comments) error {
	comment.ID = uuid.New()

	if mail == "" {
		db.Logger.Printf("mailID can not be empty")
		return fmt.Errorf("mailID can not be empty")
	}

	if err := db.DB.Debug().Where("mail=?", mail).Where("role=?", "user").First(&models.User{}).Error; err != nil {
		db.Logger.Printf("Error, %v Occured when searching the user based on the mail: %v", err, mail)
		return fmt.Errorf("unauthorized")
	}

	err := utilities.ValidateStruct(comment)
	if err != nil {
		db.Logger.Printf("Error validating the struct")
		return err
	}

	if err := db.DB.Create(&comment).Error; err != nil {
		db.Logger.Printf("Error, %v Occured when creating the comment", err)
		return err
	}

	var commentCount int64
	if err := db.DB.Debug().Model(&comment).Where("post_id = ?", comment.PostID).Count(&commentCount).Error; err != nil {
		db.Logger.Printf("Error, %v Occured when counting the comment count", err)
		return err
	}

	if err := db.DB.Model(&models.Post{}).Where("id", comment.PostID).Update("comment_count", commentCount).Error; err != nil {
		db.Logger.Printf("Error, %v Occured when updating the comment count", err)
		return err
	}

	db.Logger.Printf("Added new comment with ID:%v", comment.ID)
	return nil
}

// Update the comment added by the user
func (db *DbConnection) UpdateCommentByID(mail string, commentID string, data map[string]interface{}) error {
	var comment models.Comments
	user := models.User{}

	if mail == "" {
		db.Logger.Printf("mailID can not be empty")
		return fmt.Errorf("mailID can not be empty")
	}

	if err := db.DB.Debug().Where("mail=?", mail).Where("role=?", "user").First(&user).Error; err != nil {
		db.Logger.Printf("Error %v Occured when searching the user with mail: %v", err, mail)
		return fmt.Errorf("unauthorized")
	}

	if commentID == "" {
		db.Logger.Printf("commentID can not be empty")
		return fmt.Errorf("commentID can not be empty")
	}

	if err := db.DB.Debug().Where("id=?", commentID).Where("role_id=?", user.ID).First(&comment).Error; err != nil {
		db.Logger.Printf("Error %v Occured when searching the comment posted by the user with ID: %v", err, user.ID)
		return fmt.Errorf("unauthorized")
	}

	if err := db.DB.Debug().Where("id=?", commentID).First(&comment).Error; err != nil {
		db.Logger.Printf("Error %v Occured when searching the comment with ID: %v", err, commentID)
		return err
	}

	fmt.Println("Data in update comment in handler", data)

	for d, value := range data {
		if err := db.DB.Model(&models.Comments{}).Where("id=?", commentID).Update(d, value).Error; err != nil {
			db.Logger.Printf("Error %v Occured when updating the comment with ID: %v", err, commentID)
			return err
		}
	}

	db.Logger.Printf("Updated the comment with ID:%v", commentID)
	return nil
}

// Delete comment by ID  DBOperation
func (db *DbConnection) DeleteCommentByID(mail string, commentID string, comment *models.Comments) error {
	user := &models.User{}
	if mail == "" {
		db.Logger.Printf("mailID can not be empty")
		return fmt.Errorf("mailID can not be empty")
	}

	if err := db.DB.Debug().Where("mail=?", mail).Where("role=?", "user").First(&user).Error; err != nil {
		db.Logger.Printf("Error, %v Occured when searching the user with mail: %v", err, mail)
		return fmt.Errorf("unauthorized")
	}

	if commentID == "" {
		db.Logger.Printf("commentID can not be empty")
		return fmt.Errorf("commentID can not be empty")
	}

	if err := db.DB.Debug().Where("id=?", commentID).Where("role_id=?", user.ID).First(&comment).Error; err != nil {
		db.Logger.Printf("Error, %v Occured when searching the comments with ID: %v", err, commentID)
		return fmt.Errorf("unauthorized")
	}

	if err := db.DB.Debug().Where("id=?", commentID).Delete(&comment).Error; err != nil {
		db.Logger.Printf("Error, %v Occured when deleting the comment with ID: %v", err, commentID)
		return err
	}

	var commentCount int64
	if err := db.DB.Debug().Model(&comment).Where("post_id = ?", comment.PostID).Count(&commentCount).Error; err != nil {
		db.Logger.Printf("Error, %v Occured when calculating the comment count", err)
		return err
	}

	if err := db.DB.Model(&models.Post{}).Where("id", comment.PostID).Update("comment_count", commentCount).Error; err != nil {
		db.Logger.Printf("Error, %v Occured when updating the comment count=", err)
		return err
	}

	db.Logger.Printf("Deleted the comment with ID:%v", commentID)
	return nil
}

// Get all the comments
func (db *DbConnection) GetAllComments(comment *[]models.Comments) error {
	if err := db.DB.Debug().Find(&comment).Error; err != nil {
		db.Logger.Printf("Error %v Occured when searching the comments", err)
		return err
	}

	db.Logger.Printf("Retrived all the comments")
	return nil
}

// Get all the comments added by the user
func (db *DbConnection) GetCommentsBasedOnUser(mail string, comment *[]models.Comments) error {
	user := models.User{}
	if mail == "" {
		db.Logger.Printf("mailID can not be empty")
		return fmt.Errorf("mailID can not be empty")
	}

	if err := db.DB.Debug().Where("mail=?", mail).Where("role=?", "user").First(&user).Error; err != nil {
		db.Logger.Printf("Error, %v Occured when searching the user with mail: %v", err, mail)
		return fmt.Errorf("unauthorized")
	}

	if err := db.DB.Debug().Where("role_id=?", user.ID).Find(&comment).Error; err != nil {
		db.Logger.Printf("Error %v Occured when searching the comment posted by the user with ID: %v", err, user.ID)
		return err
	}

	db.Logger.Printf("Retrived all the comments posted by the user with ID:%v", user.ID)
	return nil
}

// Get all the comments based ont the postID
func (db *DbConnection) GetCommentsBasedOnPostID(postID string, comment *[]models.Comments) error {
	// if mail == "" {
	// 	return fmt.Errorf("mailID can not be empty")
	// }

	// if err := db.DB.Debug().Where("mail=?", mail).Where("role=?", "user").First(&models.User{}).Error; err != nil {
	// 	return fmt.Errorf("unauthorized")
	// }
	if postID == "" {
		db.Logger.Printf("postID can not be empty")
		return fmt.Errorf("postID can not be empty")
	}

	if err := db.DB.Debug().Where("post_id=?", postID).Find(&comment).Error; err != nil {
		db.Logger.Printf("Error %v Occured when searching the comment for the post with ID: %v", err, postID)
		return err
	}

	db.Logger.Printf("Retrived all the comments for the post with ID:%v", postID)
	return nil
}
