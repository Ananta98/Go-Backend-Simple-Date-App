package controllers

import (
	"backend-dating-app/models"
	"backend-dating-app/utils"
	"backend-dating-app/utils/token"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginUserInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateUserPasswordInput struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type RegisterUserInput struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Gender   string `json:"gender" binding:"required"`
}

// Register godoc
// @Summary Register new user or create new user for dating app.
// @Description registering a user to get access dating app.
// @Tags Auth
// @Param Body body RegisterUserInput true "json body to register new user"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /auth/register [post]
func Register(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)
	var input RegisterUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !utils.IsValidEmail(input.Email) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Email (Please check make sure email format correct)"})
		return
	}
	newUser := models.User{
		Name:     input.Name,
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
		Gender:   input.Gender,
	}
	if _, err := newUser.SaveUser(db); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "registration success", "data": newUser})
}

// Login godoc
// @Summary Login user in dating app.
// @Description Logging in to get jwt token to access dating app.
// @Tags Auth
// @Param Body body LoginUserInput true "json body for login user"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /auth/login [post]
func Login(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)
	var input LoginUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var foundUser models.User
	if err := db.Where("username = ?", input.Username).First(&foundUser).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(input.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := token.GenerateToken(foundUser.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "login success", "token": token})
}

func UpdatePassword(ctx *gin.Context) {

}
