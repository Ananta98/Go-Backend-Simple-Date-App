package controllers

import (
	"backend-dating-app/models"
	"backend-dating-app/utils"
	"backend-dating-app/utils/token"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SwipeUserHistoryInput struct {
	UserID      uint      `json:"user_id" gorm:"not null"`
	UserMatchID uint      `json:"user_match_id" gorm:"not null"`
	Like        bool      `json:"like" gorm:"default:false;not null"`
	SwapDate    time.Time `json:"swap_date" gorm:"not null"`
}

// GetListMatching godoc
// @Summary Get list matching date with other user
// @Description Get list matching date with other user
// @Tags Post
// @Produce json
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Success 200 {object} map[string]interface{}
// @Router /post [get]
func ListMatching(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)
	userId, err := token.ExtractTokenID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	current_page := ctx.Query("current_page")
	page_size := ctx.Query("page_size")
	limit, offset, err := utils.GetPagination(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var usersMatching []models.User
	if err := db.Where("id != ?", userId).Offset(offset).Limit(limit).Find(&usersMatching).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var count int64
	if err := db.Model(&models.User{}).Where("id != ?", userId).Count(&count).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": usersMatching, "current_page": current_page, "page_size": page_size, "total_size": count})
}

func SwipeUser(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)
	userId, err := token.ExtractTokenID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// first check if premium user, user can swap without limits. If not one day just only 10 swap
	foundUser := models.User{}
	if err := db.Where("user_id = ?", userId).First(&foundUser).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if foundUser.PackageType == "standard" {
		var count int64
		currentDate := time.Now().Format("2006-01-02")
		if err := db.Model(&models.SwipeUserHistory{}).Where("user_id = ? and swap_date >= ?", userId, currentDate).Count(&count).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if count > 10 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "For standard user only 10 swap maximum. Please upgrade to premium version to get unlimeted access."})
			return
		}
	}

	// when swipe user add into history
	var input SwipeUserHistoryInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Model(&models.SwipeUserHistory{}).Save(&input).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": input})

}

func ProfileUser(ctx *gin.Context) {
	id := ctx.Param("id")
	db := ctx.MustGet("db").(*gorm.DB)
	var userProfile models.User
	if err := db.Where("id = ?", id).Take(&userProfile).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": userProfile})
}

func UpgradePremiumTier(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)
	userId, err := token.ExtractTokenID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var userFound models.User
	if err := db.Where("id = ?", userId).Take(&userFound).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Model(&userFound).Where("id = ?", userId).Update("package_type", "premium").Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Successfull upgrade to premium"})
}
