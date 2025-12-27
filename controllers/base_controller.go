package controllers

import (
	"fmt"
	"log"
	"project-manager/pkg"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/validator/v10"

	ut "github.com/go-playground/universal-translator"
	zht "github.com/go-playground/validator/v10/translations/zh"
)

var (
	Role    = &RoleController{}
	User    = &UserController{}
	Team    = &TeamController{}
	Project = &ProjectController{}

	validate = validator.New()
	trans    ut.Translator
)

func init() {
	uni := ut.New(zh.New())
	trans, _ = uni.GetTranslator("zh")
	_ = zht.RegisterDefaultTranslations(validate, trans)
	//注册正则验证
	_ = validate.RegisterValidation("password", pkg.PasswordValidate)
}

func Handle(c *gin.Context, req any, fn func() (any, any)) {
	var err error

	if err := c.Bind(req); err != nil {
		pkg.Err(c, pkg.NewValidatorError(err), nil)
		return
	}

	// 设置默认值（如果请求对象有 SetDefaults 方法）
	if defaultSetter, ok := req.(interface{ SetDefaults() }); ok {
		defaultSetter.SetDefaults()
	}

	//打印请求参数
	log.Printf("请求参数: %+v\n", req)

	//校验
	if req != nil {
		// 检查是否为切片类型
		reqValue := reflect.ValueOf(req)

		// 如果是指针，获取指针指向的值
		if reqValue.Kind() == reflect.Ptr {
			reqValue = reqValue.Elem()
		}

		if reqValue.Kind() == reflect.Slice {
			log.Printf("切片校验，长度: %d\n", reqValue.Len())
			// 如果是切片，遍历验证每个元素
			for i := 0; i < reqValue.Len(); i++ {
				item := reqValue.Index(i).Interface()
				log.Printf("校验第 %d 个元素: %+v\n", i, item)
				err = validate.Struct(item)
				if err != nil {
					for _, err := range err.(validator.ValidationErrors) {
						pkg.Err(c, pkg.NewValidatorError(fmt.Errorf("%s", err.Translate(trans))), nil)
						return
					}
				}
			}
		} else {
			log.Printf("单个结构体校验\n")
			// 如果是单个结构体，按原来的逻辑处理
			err = validate.Struct(req)
			if err != nil {
				for _, err := range err.(validator.ValidationErrors) {
					pkg.Err(c, pkg.NewValidatorError(fmt.Errorf("%s", err.Translate(trans))), nil)
					return
				}
			}
		}
	}
	data, err1 := fn()
	if err1 != nil {
		pkg.Err(c, pkg.ReloadErr(err1), data)
		return
	}
	pkg.Success(c, data)
}
