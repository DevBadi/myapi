package controller

import (
	"net/http"

	"github.com/pilinux/gorest/database"
	"github.com/pilinux/gorest/lib/middleware"
	"github.com/pilinux/gorest/lib/renderer"

	"github.com/gin-gonic/gin"

	"myapi/model"
)

// GetPosts - GET /posts
func GetPosts(c *gin.Context) {
	db := database.GetDB()
	posts := []model.Post{}

	if err := db.Find(&posts).Error; err != nil {
		renderer.Render(c, gin.H{"msg": "not found"}, http.StatusNotFound)
	} else {
		renderer.Render(c, posts, http.StatusOK)
	}
}

// GetPost - GET /posts/:id
func GetPost(c *gin.Context) {
	db := database.GetDB()
	post := model.Post{}
	id := c.Params.ByName("id")

	if err := db.Where("post_id = ? ", id).First(&post).Error; err != nil {
		renderer.Render(c, gin.H{"msg": "not found"}, http.StatusNotFound)
	} else {
		renderer.Render(c, post, http.StatusOK)
	}
}

// CreatePost - POST /posts
func CreatePost(c *gin.Context) {
	db := database.GetDB()
	user := model.User{}
	post := model.Post{}
	postFinal := model.Post{}

	userIDAuth := middleware.AuthID

	// does the user have an existing profile
	if err := db.Where("id_auth = ?", userIDAuth).First(&user).Error; err != nil {
		renderer.Render(c, gin.H{"msg": "no user profile found"}, http.StatusForbidden)
		return
	}

	// bind JSON
	if err := c.ShouldBindJSON(&post); err != nil {
		renderer.Render(c, gin.H{"msg": "bad request"}, http.StatusBadRequest)
		return
	}

	// user must not be able to manipulate all fields
	postFinal.Title = post.Title
	postFinal.Body = post.Body
	postFinal.IDUser = user.UserID

	tx := db.Begin()
	if err := tx.Create(&postFinal).Error; err != nil {
		tx.Rollback()
		renderer.Render(c, gin.H{"msg": "internal server error"}, http.StatusInternalServerError)
	} else {
		tx.Commit()
		renderer.Render(c, postFinal, http.StatusCreated)
	}
}
