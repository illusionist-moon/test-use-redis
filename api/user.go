package api

import (
	"ChildrenMath/models"
	"ChildrenMath/pkg/e"
	"ChildrenMath/pkg/util"
	"ChildrenMath/pkg/validation"
	"ChildrenMath/service/emailverify"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
)

func SendRegisterVCode(ctx *gin.Context) {
	email := ctx.PostForm("email")

	emailVerify := &validation.EmailValidator{Email: email}
	if err := ctx.ShouldBind(emailVerify); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.InvalidParams,
			"msg":  validation.GetValidMsg(err, emailVerify),
		})
		return
	}

	// 判断用户是否存在
	exists := models.Exists(models.DB, email)
	if exists {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.ErrorExistUser,
			"msg":  e.GetMsg(e.ErrorExistUser),
		})
		return
	}

	//产生六位数验证码
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vCode := fmt.Sprintf("%06v", rnd.Int31n(1000000))

	err := models.SaveVCodeInRedis(email, vCode)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.Error,
			"msg":  emailverify.ErrSend.Error(),
		})
		return
	}

	err = emailverify.SendRegisterEmail(email, vCode)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.Error,
			"msg":  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": e.Success,
		"msg":  e.GetMsg(e.Success),
	})
	return
}

func Register(ctx *gin.Context) {
	email := ctx.PostForm("email")
	emailVerify := &validation.EmailValidator{Email: email}
	if err := ctx.ShouldBind(emailVerify); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.InvalidParams,
			"msg":  validation.GetValidMsg(err, emailVerify),
		})
		return
	}

	vCode := ctx.PostForm("vcode")
	vCodeInRedis, err := models.GetVCodeFromRedis(email)
	if err != nil || vCode != vCodeInRedis {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.InvalidParams,
			"msg":  "验证码有误",
		})
		return
	}

	userName := ctx.PostForm("username")
	password := ctx.PostForm("password")
	rePassword := ctx.PostForm("re-password")

	verify := &validation.UserRegisterValidator{
		UserName:          userName,
		PasswordValidator: validation.PasswordValidator{Password: password},
		RePassword:        rePassword,
	}
	if err := ctx.ShouldBind(verify); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.InvalidParams,
			"msg":  validation.GetValidMsg(err, verify),
		})
		return
	}

	hash, err := util.GetBcryptPwd(password)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.Error,
			"msg":  "用户密码加密时失败",
		})
		return
	}
	var userID int
	userID, err = models.CreateUser(models.DB, userName, email, hash)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.Error,
			"msg":  "Create User Error: " + err.Error(),
		})
		return
	}

	_, err = models.InitPointsKeysInRedis(userID, userName)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.Error,
			"msg":  "Init Points Error: " + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": e.Success,
		"msg":  e.GetMsg(e.Success),
	})
}

// ------------------------------------------------------------

// Login 登录
func Login(ctx *gin.Context) {
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")

	verify := &validation.UserLoginValidator{
		EmailValidator:    validation.EmailValidator{Email: email},
		PasswordValidator: validation.PasswordValidator{Password: password},
	}
	if err := ctx.ShouldBind(verify); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.InvalidParams,
			"msg":  validation.GetValidMsg(err, verify),
		})
		return
	}

	userID, userName, getPassword, err := models.GetUserInfoByEmail(models.DB, email)
	if err != nil {
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

	token, err := util.GenerateToken(userID, userName, email)
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

func Logout(ctx *gin.Context) {
	code := e.Success
	ctx.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
	})
}

// ------------------------------------------------------------

// SendForgetPasswordVCode 用于发送忘记密码时的验证码
func SendForgetPasswordVCode(ctx *gin.Context) {
	email := ctx.PostForm("email")

	emailVerify := &validation.EmailValidator{Email: email}
	if err := ctx.ShouldBind(emailVerify); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.InvalidParams,
			"msg":  validation.GetValidMsg(err, emailVerify),
		})
		return
	}

	// 判断用户是否存在，如果不存在则返回假
	exists := models.Exists(models.DB, email)
	if !exists {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.ErrorNotExistUser,
			"msg":  e.GetMsg(e.ErrorNotExistUser),
		})
		return
	}

	//产生六位数验证码
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vCode := fmt.Sprintf("%06v", rnd.Int31n(1000000))

	err := models.SaveVCodeInRedis(email, vCode)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.Error,
			"msg":  emailverify.ErrSend.Error(),
		})
		return
	}

	err = emailverify.SendForgetPasswordEmail(email, vCode)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.Error,
			"msg":  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": e.Success,
		"msg":  e.GetMsg(e.Success),
	})
}

func UpdateForgetPassword(ctx *gin.Context) {
	email := ctx.PostForm("email")
	emailVerify := &validation.EmailValidator{Email: email}
	if err := ctx.ShouldBind(emailVerify); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.InvalidParams,
			"msg":  validation.GetValidMsg(err, emailVerify),
		})
		return
	}

	vCode := ctx.PostForm("vcode")
	vCodeInRedis, err := models.GetVCodeFromRedis(email)
	if err != nil || vCode != vCodeInRedis {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.Error,
			"msg":  emailverify.ErrEqual.Error(),
		})
		return
	}

	newPassword := ctx.PostForm("password")
	rePassword := ctx.PostForm("re-password")

	verify := &validation.UpdatePasswordValidator{
		PasswordValidator: validation.PasswordValidator{Password: newPassword},
		RePassword:        rePassword,
	}
	if err := ctx.ShouldBind(verify); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.InvalidParams,
			"msg":  validation.GetValidMsg(err, verify),
		})
		return
	}

	_, _, oldPassword, err := models.GetUserInfoByEmail(models.DB, email)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.Error,
			"msg":  "获取用户旧密码失败",
		})
		return
	}

	if util.ComparePwd(oldPassword, newPassword) {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.Error,
			"msg":  "新密码不能与旧密码相同",
		})
		return
	}

	hash, err := util.GetBcryptPwd(newPassword)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.Error,
			"msg":  "新密码加密时失败",
		})
		return
	}

	err = models.UpdatePassword(models.DB, email, hash)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.Error,
			"msg":  "更新新密码失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": e.Success,
		"msg":  e.GetMsg(e.Success),
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

	nameVal, exist := ctx.Get("username")
	// 下面这种情况理论是不存在，但还是需要写出处理
	if !exist {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.ErrorNotExistUser,
			"data": nil,
			"msg":  "用户获取出现问题",
		})
		return
	}
	userName := nameVal.(string)

	var (
		rank      []models.Rank
		ownPoints int
		err       error
	)

	fmt.Println(userID, userName)
	ownPoints, err = models.GetOwnPointsFromRedisWithSave(userID, userName)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.Error,
			"data": nil,
			"msg":  "个人积分拉取失败: " + err.Error(),
		})
		return
	}

	rank, err = models.GetPointsRankByIDFromRedis()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.Error,
			"data": nil,
			"msg":  "积分排名拉取失败: " + err.Error(),
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
			"code": e.ErrorNotExistUser,
			"data": nil,
			"msg":  "用户获取出现问题",
		})
		return
	}
	userID := idVal.(int)

	nameVal, exist := ctx.Get("username")
	// 下面这种情况理论是不存在，但还是需要写出处理
	if !exist {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.ErrorNotExistUser,
			"data": nil,
			"msg":  "用户获取出现问题",
		})
		return
	}
	userName := nameVal.(string)

	ownPoints, err := models.GetOwnPointsFromRedisWithSave(userID, userName)
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
