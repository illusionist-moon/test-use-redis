package api

import (
	"ChildrenMath/models"
	"ChildrenMath/pkg/e"
	"ChildrenMath/pkg/util"
	"ChildrenMath/pkg/validation"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Login(ctx *gin.Context) {
	username := strings.TrimSpace(ctx.PostForm("username"))
	password := ctx.PostForm("password")

	verify := &validation.UserLogin{
		UserName: username,
		Password: password,
	}
	if err := ctx.ShouldBind(verify); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.InvalidParams,
			"msg":  validation.GetValidMsg(err, verify),
		})
		return
	}

	userID, getPassword, exists := models.Exists(models.DB, username)
	if !exists {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.ErrorNotExistUser,
			"msg":  e.GetMsg(e.ErrorNotExistUser),
		})
		return
	}
	// 校验密码
	if !util.ComparePwd(getPassword, password) {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.ErrorIncorrectPwd,
			"msg":  e.GetMsg(e.ErrorIncorrectPwd),
		})
		return
	}

	token, err := util.GenerateToken(userID, username)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.ErrorAuthToken,
			"msg":  e.GetMsg(e.ErrorAuthToken),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":  e.Success,
		"msg":   e.GetMsg(e.Success),
		"token": token,
	})
}

func Register(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	rePassword := ctx.PostForm("re-password")

	verify := &validation.UserRegister{
		UserName:   username,
		Password:   password,
		RePassword: rePassword,
	}
	if err := ctx.ShouldBind(verify); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.InvalidParams,
			"msg":  validation.GetValidMsg(err, verify),
		})
		return
	}

	// 判断用户是否存在
	_, _, exists := models.Exists(models.DB, username)
	if exists {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.ErrorExistUser,
			"msg":  e.GetMsg(e.ErrorExistUser),
		})
		return
	}
	hash, err := util.GetBcryptPwd(password)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.InvalidParams,
			"msg":  "用户密码加密时失败",
		})
		return
	}
	err = models.CreateUser(models.DB, username, hash)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.Error,
			"msg":  "Create User Error: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": e.Success,
		"msg":  e.GetMsg(e.Success),
	})
}

func Logout(ctx *gin.Context) {
	code := e.Success
	ctx.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
	})
}

func GetPointsRank(ctx *gin.Context) {
	idVal, exist := ctx.Get("userid")
	// 下面这种情况理论是不存在，但还是需要写出处理
	if !exist {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.ErrorNotExistUser,
			"data": nil,
			"msg":  "用户获取出现问题",
		})
		return
	}
	userID := idVal.(int)

	var (
		rank      []models.Rank
		ownPoints int
		err       error
	)
	rank, err = models.GetPointsRankFromRedis()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.Error,
			"data": nil,
			"msg":  "积分排名拉取失败: " + err.Error(),
		})
		return
	}

	ownPoints, err = models.GetOwnPointsFromRedisWithSave(userID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.Error,
			"data": nil,
			"msg":  "个人积分拉取失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": e.Success,
		"data": map[string]any{
			"max_count":  models.RankSize,
			"rank":       rank,
			"own_points": ownPoints,
		},
		"msg": e.GetMsg(e.Success),
	})
}

func GetUserPoints(ctx *gin.Context) {
	idVal, exist := ctx.Get("userid")
	// 下面这种情况理论是不存在，但还是需要写出处理
	if !exist {
		ctx.JSON(http.StatusOK, gin.H{
			"code":   e.ErrorNotExistUser,
			"points": nil,
			"msg":    "用户获取出现问题",
		})
		return
	}
	userID := idVal.(int)

	ownPoints, err := models.GetOwnPointsFromRedisWithSave(userID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":   e.Error,
			"points": nil,
			"msg":    "个人积分拉取失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":   e.Success,
		"points": ownPoints,
		"msg":    e.GetMsg(e.Success),
	})
}
