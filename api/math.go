package api

import (
	"ChildrenMath/models"
	"ChildrenMath/pkg/e"
	"ChildrenMath/service/question"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"strconv"
	"strings"
)

func GetQuestions(ctx *gin.Context) {
	op, ok := ctx.GetQuery("op")
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.InvalidParams,
			"msg":  "未传入运算符",
			"data": nil,
		})
		return
	}
	switch op {
	case "plus":
		op = "+"
	case "minus":
		op = "-"
	case "multi":
		op = "*"
	case "div":
		op = "/"
	default:
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.InvalidParams,
			"msg":  "非法运算符",
			"data": nil,
		})
		return
	}
	data := question.GenerateQuestions(op)
	ctx.JSON(http.StatusOK, gin.H{
		"code": e.Success,
		"msg":  e.GetMsg(e.Success),
		"data": map[string]any{
			"count":     question.Count,
			"op":        op,
			"questions": data,
		},
	})
}

func JudgeQuestion(ctx *gin.Context) {
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

	val, exist := ctx.Get("username")
	// 下面这种情况理论是不存在，但还是需要写出处理
	if !exist {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.ErrorNotExistUser,
			"data": nil,
			"msg":  "用户获取出现问题",
		})
		return
	}
	username := val.(string)

	op, ok := ctx.GetPostForm("op")
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.InvalidParams,
			"data": nil,
			"msg":  "miss operator",
		})
		return
	}
	switch op {
	case "plus":
		op = "+"
	case "minus":
		op = "-"
	case "multi":
		op = "*"
	case "div":
		op = "/"
	default:
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.InvalidParams,
			"msg":  "非法运算符",
			"data": nil,
		})
		return
	}

	answers, ansOK := ctx.GetPostForm("answers")
	if !ansOK {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.InvalidParams,
			"data": nil,
			"msg":  "miss answers",
		})
		return
	}
	var ansSlice [][]int
	unmarshalErr := json.Unmarshal([]byte(answers), &ansSlice)
	if unmarshalErr != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.InvalidParams,
			"data": nil,
			"msg":  "answers反序列化失败",
		})
		return
	}
	if len(ansSlice) != question.Count {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.InvalidParams,
			"data": nil,
			"msg":  "invalid answer count, count should be " + strconv.Itoa(question.Count),
		})
		return
	}

	// 开启一个事务，保证错题库和积分的一致性
	tx := models.DB.Begin()
	var (
		addPoints int
		err       error
	)

	for i := 0; i < question.Count; i++ {
		correct := question.Judge(ansSlice[i], op)
		if correct {
			addPoints++
		} else {
			err = models.AddProblem(tx, username, op, ansSlice[i][0], ansSlice[i][1], ansSlice[i][2])
			if err != nil {
				tx.Rollback()
				ctx.JSON(http.StatusOK, gin.H{
					"code": e.Error,
					"data": nil,
					"msg":  "add incorrect answer failed",
				})
				return
			}
		}
	}
	err = models.IncreaseOwnPointsInRedis(userID, addPoints, username)
	if err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.Error,
			"data": nil,
			"msg":  "add point failed in Redis: " + err.Error(),
		})
		return
	}
	//err = models.UpdatePoints(tx, userID, addPoints)
	//if err != nil {
	//	tx.Rollback()
	//	ctx.JSON(http.StatusOK, gin.H{
	//		"code": e.Error,
	//		"data": nil,
	//		"msg":  "add point failed",
	//	})
	//	return
	//}
	ctx.JSON(http.StatusOK, gin.H{
		"code": e.Success,
		"data": addPoints,
		"msg":  "success",
	})
	tx.Commit()
}

func GetWrongList(ctx *gin.Context) {
	val, exist := ctx.Get("username")
	// 下面这种情况理论是不存在，但还是需要写出处理
	if !exist {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.ErrorNotExistUser,
			"data": nil,
			"msg":  "用户获取出现问题",
		})
		return
	}
	username := val.(string)

	pageStr, ok := ctx.GetQuery("page")
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.InvalidParams,
			"data": nil,
			"msg":  "miss page",
		})
		return
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.InvalidParams,
			"data": nil,
			"msg":  "invalid page",
		})
		return
	}

	wrongItems, total, getErr := models.GetWrongList(models.DB, username, page)
	if getErr != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.Error,
			"data": nil,
			"msg":  "拉取错题列表失败",
		})
		return
	}
	totalPages := int(math.Ceil(float64(total) / models.WrongListOffset))
	if page > totalPages {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.InvalidParams,
			"data": nil,
			"msg":  "页号溢出，最大为" + strconv.Itoa(totalPages),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": e.Success,
		"data": map[string]any{
			"record":       wrongItems,
			"record_count": len(wrongItems),
			"total_page":   totalPages,
		},
		"msg": e.GetMsg(e.Success),
	})
	return
}

func GetRedoProblem(ctx *gin.Context) {
	val, exist := ctx.Get("username")
	// 下面这种情况理论是不存在，但还是需要写出处理
	if !exist {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.ErrorNotExistUser,
			"data": nil,
			"msg":  "用户获取出现问题",
		})
		return
	}
	username := val.(string)

	res, err := models.GetRedoProblem(models.DB, username)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.Error,
			"data": nil,
			"msg":  "拉取错误题目失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": e.Success,
		"data": map[string]any{
			"count":     len(res),
			"questions": res,
		},
		"msg": e.GetMsg(e.Success),
	})
}

func JudgeRedoProblem(ctx *gin.Context) {
	answers, ansOK := ctx.GetPostFormArray("answer[]")
	if !ansOK {
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.InvalidParams,
			"data": nil,
			"msg":  "miss answers",
		})
		return
	}

	// 开启一个事务，保证错题库和积分的一致性
	tx := models.DB.Begin()
	var (
		correct     int
		id          int
		nums        = make([]int, 3) // 依次为 num1, num2, res
		op          string
		err         error
		deleteIDSet []int
	)
	for _, ans := range answers {
		data := strings.Split(ans, ",")
		if len(data) != 5 {
			tx.Rollback()
			ctx.JSON(http.StatusOK, gin.H{
				"code": e.InvalidParams,
				"data": nil,
				"msg":  "invalid data count in an answer, count must be 5: id, num1, num2, ans, op",
			})
			return
		}

		id, err = strconv.Atoi(data[0])
		if err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusOK, gin.H{
				"code": e.InvalidParams,
				"data": nil,
				"msg":  err.Error(),
			})
			return
		}
		for i := 1; i <= 3; i++ {
			nums[i-1], err = strconv.Atoi(data[i])
			if err != nil {
				tx.Rollback()
				ctx.JSON(http.StatusOK, gin.H{
					"code": e.InvalidParams,
					"data": nil,
					"msg":  "number convert into int failed",
				})
				return
			}
		}
		op = data[4]
		if op != "+" && op != "-" && op != "*" && op != "/" {
			ctx.JSON(http.StatusOK, gin.H{
				"code": e.InvalidParams,
				"data": nil,
				"msg":  "invalid op",
			})
			return
		}
		if question.Judge(nums, op) {
			correct++
			deleteIDSet = append(deleteIDSet, id)
		}
	}
	var deleteCount int64
	deleteCount, err = models.DeleteRedoProblem(tx, deleteIDSet)
	if err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.Error,
			"data": nil,
			"msg":  err.Error(),
		})
		return
	}
	if int(deleteCount) != len(deleteIDSet) {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{
			"code": e.Error,
			"data": nil,
			"msg":  "delete wrong problems failed",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    e.Success,
		"correct": correct,
		"msg":     "success",
	})
	tx.Commit()
}
