package validation

import (
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"reflect"
)

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	trans    ut.Translator
)

func init() {
	//translator := zh.New()
	//uni = ut.New(translator, translator)
	//trans, _ = uni.GetTranslator("zh")

	var ok bool
	//注册自定义验证方法
	if validate, ok = binding.Validator.Engine().(*validator.Validate); ok {
		validate.RegisterValidation("noCRLF", noCRLF)
	}

	//err := zhTranslations.RegisterDefaultTranslations(validate, trans)
	//if err != nil {
	//	log.Println(err)
	//}
}

func GetValidMsg(err error, obj interface{}) string {
	getObj := reflect.TypeOf(obj)
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			if f, exist := getObj.Elem().FieldByName(e.Field()); exist {
				return f.Tag.Get("msg") //错误信息不需要全部返回，当找到第一个错误的信息时，就可以结束
			}
		}
	}
	return err.Error()
}
