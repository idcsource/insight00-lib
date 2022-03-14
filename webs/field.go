// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package webs

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/idcsource/insight00-lib/base"
	"github.com/idcsource/insight00-lib/jconf"
)

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

type FieldType uint

const (
	// 表单输入项内容
	FIELD_FORM_OTHER        FieldType = iota // 未知类型
	FIELD_FORM_TEXT                          // 文本：test
	FIELD_FORM_TEXTAREA                      // 文本域：testara
	FIELD_FORM_MARKDOWN                      // Markdown文本：markdown
	FIELD_FORM_EDITOR                        // 编辑器：editor
	FIELD_FORM_MARK                          // 也就是A-Za-z0-9_-：mark
	FIELD_FORM_PATH                          // 路径：path
	FIELD_FORM_LINK                          // 链接：也就是http://或https://或ftp://开头：link
	FIELD_FORM_EMAIL                         // 邮箱：email
	FIELD_FORM_FILE                          // 文件：file
	FIELD_FORM_TIME                          // 时间：time
	FIELD_FORM_DATE                          // 日期：date
	FIELD_FORM_INT                           // 数字:int
	FIELD_FORM_FLOAT                         // 浮点:float
	FIELD_FORM_MONEY                         // 钱:money
	FIELD_FORM_PASSWORD                      // 密码:password
	FIELD_FORM_PASSWORD_TWO                  // 密码两遍输入:passwordt
	FIELD_FORM_ENUM_STRING                   // 字符串类枚举:enums
	FIELD_FORM_ENUM_INT                      // 数字类枚举:enumi
	FIELD_FORM_ENUM_FLOAT                    // 浮点型枚举:enumf
)

type CheckStatus uint8

const (
	// 检查状态
	CHECK_STATUS_NO     CheckStatus = iota // 检查状态不对
	CHECK_STATUS_OK                        // 检查OK
	CHECK_STATUS_NULL                      // 为空
	CHECK_STATUS_ERROR                     // 错误
	CHECK_STATUS_ERROR2                    // 错误2
	CHECK_STATUS_ERROR3                    // 错误3
)

// 单个字段信息
type FieldConfig struct {
	Name          string    // 字段名
	DisName       string    // 显示名称
	DatabaseField string    // 对应数据库字段名
	RoleField     string    // 对应角色内字段名
	UseIt         bool      // 是否启用
	Type          FieldType // 类型
	CanNull       bool      // 是否可以为空
	Info          string    // 说明信息
	Min           int64     // 最小长度或最小数字
	Max           int64     // 最大长度或最大数字
	EnumString    []string  // 字符串枚举的限定值范围
	EnumInt       []int64   // 数字型枚举的限定值范围
	EnumFloat     []float64 // 浮点型枚举的限定值范围
	Other         string    // 其他状态
}

// 静态字段设置
type StaticFields map[string]*FieldConfig

// 运行时字段设置
type RuntimeFields map[string]*FieldConfig

type FormData struct {
	FieldsConfig RuntimeFields
	R            *http.Request
	Ip           *base.InputProcessor
}
type OneFormData struct {
	Int    int64
	Float  float64
	String string
	Bool   bool
	Null   bool
}

// 准备运行时字段设置
//
// fields是所有在运行时使用到的字段，bool为true则是可以为null。
// config是针对某个节点的特有字段设定，主要是哪些启用，以及名字什么的。
// sfc则是unit设定的所有字段信息。
// sfc会根据config调整运行时的配置。
func NewFormData(config *jconf.JsonConf, sfc StaticFields, r *http.Request, maxmemory int64) (fd *FormData, err error) {
	fd = &FormData{
		R:  r,
		Ip: base.NewInputProcessor(),
	}
	afc, err := fd.prepareRuntimeFields(config, sfc)
	if err != nil {
		return
	}
	fd.FieldsConfig = afc
	if fd.R.PostForm == nil {
		if maxmemory == 0 {
			fd.R.ParseMultipartForm(defaultMaxMemory)
		} else {
			fd.R.ParseMultipartForm(maxmemory << 20)
		}
	}
	return
}

// 准备运行时字段设置
func (fd *FormData) prepareRuntimeFields(config *jconf.JsonConf, sfc StaticFields) (rfc RuntimeFields, err error) {
	rfc = make(map[string]*FieldConfig)

	for k, sf := range sfc {
		cf, errs := config.GetEnum(k)
		if errs != nil {
			sf.UseIt = false
		} else {
			err = fd.getOneFieldConfig(sf, cf)
			if err != nil {
				return
			}
		}
		rfc[k] = sf
	}
	return
}

func (fd *FormData) getOneFieldConfig(field *FieldConfig, config []string) (err error) {
	//cfA := strings.Split(config, ",")
	for i, ocf := range config {
		switch i {
		case 0:
			//是否使用
			if ocf == "true" {
				field.UseIt = true
			} else {
				field.UseIt = false
				continue
			}
		case 1:
			//显示名
			field.DisName = ocf
		case 2:
			field.Info = ocf
		case 3:
			//最小值
			min, err := strconv.ParseInt(ocf, 10, 64)
			if err != nil {
				return err
			}
			if min > field.Min {
				field.Min = min
			}
		case 4:
			//最大值
			max, err := strconv.ParseInt(ocf, 10, 64)
			if err != nil {
				return err
			}
			if max > field.Max {
				field.Max = max
			}
			field.Max, _ = strconv.ParseInt(ocf, 10, 64)
		}
	}
	return
}

// 获取一个的值。
// 任何不符合要求的输入值，check都会返回CHECK_STATUS_ERROR或CHECK_STATUS_NO（包括输入为空但不允许为空的时候）。
// 任何符合要求的输入，check会返回CHECK_STATUS_OK。
// 如果允许为空，且实际也为空，则会返回CHECK_STATUS_NULL。
// 如果字段在运行时标注为不使用，则check不做检查直接返回CHECK_STATUS_OK,通过fdata可以取得各类型的默认值。
func (fd *FormData) Get(field string, cannull bool) (fdata OneFormData, check CheckStatus) {
	f, have := fd.FieldsConfig[field]
	if have == false {
		check = CHECK_STATUS_ERROR
		return
	}
	if f.UseIt == false {
		fdata = OneFormData{}
		check = CHECK_STATUS_OK
		return
	}
	formd, find := fd.R.PostForm[field]
	if find == false {
		check = CHECK_STATUS_ERROR
		return
	}
	fdata, check = fd.checkOne(formd, 0, f, cannull)
	return
}

// 获取一个的值的集合（post的同名情况）
// 任何不符合要求的输入值，check都会返回CHECK_STATUS_ERROR或CHECK_STATUS_NO（包括输入为空但不允许为空的时候）。
// 任何符合要求的输入，check会返回CHECK_STATUS_OK。
// 如果允许为空，且实际也为空，则会返回CHECK_STATUS_NULL。
// 如果字段在运行时标注为不使用，则check不做检查直接返回CHECK_STATUS_OK,通过data可以获得一项默认的数值。
func (fd *FormData) GetAll(field string, cannull bool) (data []OneFormData, check CheckStatus) {
	f, have := fd.FieldsConfig[field]
	if have == false {
		check = CHECK_STATUS_ERROR
		return
	}
	if f.UseIt == false {
		data = []OneFormData{OneFormData{}}
		check = CHECK_STATUS_OK
		return
	}
	formd, find := fd.R.PostForm[field]
	if find == false {
		check = CHECK_STATUS_ERROR
		return
	}
	data = make([]OneFormData, len(formd))
	for i, _ := range formd {
		fdata := OneFormData{}
		fdata, check = fd.checkOne(formd, i, f, cannull)
		if check != CHECK_STATUS_OK {
			return
		}
		data[i] = fdata
		i++
	}
	return
}

// 获取fields指定的全部字段的值。
// fields的map中，string为字段名，bool为是否可为空（cannull）。
// 内部使用Get方法而不是GetAll方法。
// 只要有一个输入不符合要求，check就会返回CHECK_STATUS_ERROR或CHECK_STATUS_NO。
func (fd *FormData) Fields(fields map[string]bool) (data map[string]OneFormData, check CheckStatus) {
	data = make(map[string]OneFormData)
	for name, cannull := range fields {
		data[name], check = fd.Get(name, cannull)
		if check != CHECK_STATUS_OK && check != CHECK_STATUS_NULL {
			return
		}
	}
	check = CHECK_STATUS_OK
	return
}

//func (fd *FormData) GetAll() (all map[string]OneFormData, check bool, checks string) {
//	all = make(map[string]OneFormData)
//	if fieldnames == nil {
//		for key, field := range fd.FieldsConfig {
//			one := OneFormData{Int: 0, Float: 0, String: "", Bool: false, Null: false}
//			if field.UseIt == false {
//				all[key] = one
//				continue
//			} else {
//				check, checks = fd.checkOne(field, &one, fd.R)
//				if check == false {
//					return
//				} else {
//					all[key] = one
//				}
//			}
//		}
//	} else {
//		for _, fieldname := range fieldnames {
//			field, find := fd.FieldsConfig[fieldname]
//			if find == false {
//				check = false
//				checks = "There is no field " + fieldname + " in config."
//				return
//			}
//			one := OneFormDataReturn{Int: 0, Float: 0, String: "", Bool: false}
//			if field.UseIt == false {
//				all[fieldname] = one
//				continue
//			} else {
//				check, checks = fd.checkOne(field, &one, fd.R)
//				if check == false {
//					return
//				} else {
//					all[fieldname] = one
//				}
//			}
//		}
//	}
//	return
//}

func (fd *FormData) checkOne(s []string, i int, fc *FieldConfig, cannull bool) (now OneFormData, check CheckStatus) {
	now = OneFormData{}
	s1 := strings.TrimSpace(s[i])
	if fc.UseIt == false {
		check = CHECK_STATUS_OK
		return
	}
	if len(s1) == 0 {
		if cannull == true {
			check = CHECK_STATUS_NULL
		} else {
			check = CHECK_STATUS_NO
		}
		return
	}
	var ec int
	switch fc.Type {
	case FIELD_FORM_TEXT:
		now.String, ec = fd.Ip.Text(s[i], cannull, fc.Min, fc.Max)
	case FIELD_FORM_TEXTAREA:
		now.String, ec = fd.Ip.Text(s[i], cannull, fc.Min, fc.Max)
	case FIELD_FORM_MARKDOWN:
		now.String, ec = fd.Ip.EditorIn(s[i], cannull, fc.Min, fc.Max)
	case FIELD_FORM_EDITOR:
		now.String, ec = fd.Ip.EditorIn(s[i], cannull, fc.Min, fc.Max)
	case FIELD_FORM_MARK:
		now.String, ec = fd.Ip.Mark(s[i], cannull, fc.Min, fc.Max)
	case FIELD_FORM_PATH:
		now.String, ec = fd.Ip.Text(s[i], cannull, fc.Min, fc.Max)
	case FIELD_FORM_LINK:
		now.String, ec = fd.Ip.Url(s[i], cannull, fc.Min, fc.Max)
	case FIELD_FORM_EMAIL:
		now.String, ec = fd.Ip.Email(s[i], cannull, fc.Min, fc.Max)
	case FIELD_FORM_FILE:
		now.String, ec = fd.Ip.Text(s[i], cannull, fc.Min, fc.Max)
	case FIELD_FORM_TIME:
		now.Int, ec = fd.Ip.Int(s[i], cannull, 0, 9999999999)
	case FIELD_FORM_DATE:
		now.Int, ec = fd.Ip.Int(s[i], cannull, 0, 9999999999)
	case FIELD_FORM_INT:
		now.Int, ec = fd.Ip.Int(s[i], cannull, fc.Min, fc.Max)
	case FIELD_FORM_FLOAT:
		now.Float, ec = fd.Ip.Float(s[i], cannull, float64(fc.Min), float64(fc.Max))
	case FIELD_FORM_MONEY:
		now.Float, ec = fd.Ip.Float(s[i], cannull, float64(fc.Min), float64(fc.Max))
	case FIELD_FORM_PASSWORD:
		now.String, ec = fd.Ip.Password(s[i], cannull)
	case FIELD_FORM_PASSWORD_TWO:
		now.String, ec = fd.Ip.PasswordTwo(s[i], s[i+1], cannull)
	case FIELD_FORM_ENUM_STRING:
		now.String, ec = fd.Ip.StringEnum(s[i], cannull, fc.EnumString)
	}
	if ec != 0 {
		check = CHECK_STATUS_NO
		return
	} else {
		check = CHECK_STATUS_OK
		return
	}
}

//func (fd *FormData) checkOne(fc *FieldConfig, now *OneFormDataReturn, r *http.Request) (check bool, checks string) {
//	check = false
//	var ec int
//	formd := r.PostForm[fc.DatabaseField]
//	if len(formd) > 0 {
//		switch fc.Type {
//		case FIELD_FORM_TEXT:
//			now.String, ec = fd.Ip.Text(formd[0], fc.CanNull, fc.Min, fc.Max)
//		case FIELD_FORM_TEXTAREA:
//			now.String, ec = fd.Ip.Text(formd[0], fc.CanNull, fc.Min, fc.Max)
//		case FIELD_FORM_MARKDOWN:
//			now.String, ec = fd.Ip.EditorIn(formd[0], fc.CanNull, fc.Min, fc.Max)
//		case FIELD_FORM_EDITOR:
//			now.String, ec = fd.Ip.EditorIn(formd[0], fc.CanNull, fc.Min, fc.Max)
//		case FIELD_FORM_MARK:
//			now.String, ec = fd.Ip.Mark(formd[0], fc.CanNull, fc.Min, fc.Max)
//		case FIELD_FORM_PATH:
//			now.String, ec = fd.Ip.Text(formd[0], fc.CanNull, fc.Min, fc.Max)
//		case FIELD_FORM_LINK:
//			now.String, ec = fd.Ip.Url(formd[0], fc.CanNull, fc.Min, fc.Max)
//		case FIELD_FORM_EMAIL:
//			now.String, ec = fd.Ip.Email(formd[0], fc.CanNull, fc.Min, fc.Max)
//		case FIELD_FORM_FILE:
//			now.String, ec = fd.Ip.Text(formd[0], fc.CanNull, fc.Min, fc.Max)
//		case FIELD_FORM_TIME:
//			now.Int, ec = fd.Ip.Int(formd[0], fc.CanNull, 0, 9999999999)
//		case FIELD_FORM_DATE:
//			now.Int, ec = fd.Ip.Int(formd[0], fc.CanNull, 0, 9999999999)
//		case FIELD_FORM_INT:
//			now.Int, ec = fd.Ip.Int(formd[0], fc.CanNull, fc.Min, fc.Max)
//		case FIELD_FORM_FLOAT:
//			now.Float, ec = fd.Ip.Float(formd[0], fc.CanNull, float64(fc.Min), float64(fc.Max))
//		case FIELD_FORM_MONEY:
//			now.Float, ec = fd.Ip.Float(formd[0], fc.CanNull, float64(fc.Min), float64(fc.Max))
//		case FIELD_FORM_PASSWORD:
//			now.String, ec = fd.Ip.Password(formd[0], fc.CanNull)
//		case FIELD_FORM_PASSWORD_TWO:
//			now.String, ec = fd.Ip.PasswordTwo(formd[0], formd[1], fc.CanNull)
//		}
//		if ec != 0 {
//			checks = fc.DatabaseField + " fields data not ok"
//			return
//		} else {
//			check = true
//			return
//		}
//	} else {
//		checks = "fields not from formdata: " + fc.DatabaseField
//		return
//	}
//}
