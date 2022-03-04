// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

// Insight 0+0各包共同使用的辅助函数
package base

import (
	"regexp"
	"strconv"
	"strings"
)

// InputProcessor 输入处理器
type InputProcessor struct {
	reScript *regexp.Regexp
	reMark   *regexp.Regexp
	reEmail  *regexp.Regexp
	reUrl    *regexp.Regexp
}

// 创建输入处理
func NewInputProcessor() (ip *InputProcessor) {
	ip = new(InputProcessor)
	ip.reScript, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	ip.reMark, _ = regexp.Compile("^[A-Za-z0-9_-]+$")
	ip.reEmail, _ = regexp.Compile("^[A-Za-z0-9]+([_.-][A-Za-z0-9]+)*@[A-Za-z0-9]+([_.-][A-Za-z0-9]+)*.([A-Za-z]){2,5}$")
	ip.reUrl, _ = regexp.Compile("^[A-Za-z0-9]+://")
	return
}

// 替换危险字符
func (ip *InputProcessor) replaceText(text string) string {
	theReplaceMap := make(map[string]string)
	theReplaceMap["<"] = "&lt;"
	theReplaceMap[">"] = "&gt;"
	theReplaceMap["\""] = "&guot;"
	theReplaceMap["'"] = "&#039;"
	theReplaceMap["|"] = "&brvbar;"
	theReplaceMap["`"] = "&acute;"
	for index, value := range theReplaceMap {
		text = strings.Replace(text, index, value, -1)
	}
	return text
}

// 按要求判断字符长度是否符合要求
func (ip *InputProcessor) minOrMax(text string, min, max int64) (err int) {
	textLen := len(text)
	if int64(textLen) < min {
		err = 2
		return
	}
	if max != 0 && int64(textLen) > max {
		err = 3
		return
	}
	err = 0
	return
}

// 处理简单文本输入，替换可能引起注入的字符、判断长度等
func (ip *InputProcessor) Text(text string, cannull bool, min, max int64) (textc string, err int) {
	err = 0
	textc = ""
	text = strings.TrimSpace(text)
	textLen := len(text)
	if cannull == false && textLen == 0 {
		err = 1
		return
	}
	if textLen == 0 {
		return
	}
	text = ip.replaceText(text)
	err = ip.minOrMax(text, min, max)
	if err == 0 {
		textc = text
	}
	return
}

// 处理作为标记的的字符串，只能由字母数字和连字符组成
func (ip *InputProcessor) Mark(text string, cannull bool, min, max int64) (textc string, err int) {
	err = 0
	textc = ""
	text = strings.TrimSpace(text)
	textLen := len(text)
	if cannull == false && textLen == 0 {
		err = 1
		return
	}
	if textLen == 0 {
		return
	}
	err = ip.minOrMax(text, min, max)
	if err != 0 {
		return
	}
	ifmatch := ip.reMark.MatchString(text)
	if ifmatch == false {
		err = 4
		return
	}
	textc = text
	return
}

// 处理编辑器的输入，除了字数检测外，主要是过滤script标签
func (ip *InputProcessor) EditorIn(text string, cannull bool, min, max int64) (textc string, err int) {
	err = 0
	textc = ""
	text = strings.TrimSpace(text)
	textLen := len(text)
	if cannull == false && textLen == 0 {
		err = 1
		return
	}
	if textLen == 0 {
		return
	}

	text = ip.reScript.ReplaceAllString(text, "")
	text = strings.Replace(text, "`", "&acute;", -1)
	text = strings.Replace(text, "'", "''", -1)
	err = ip.minOrMax(text, min, max)
	if err != 0 {
		return
	}
	textc = text
	return
}

// 处理编辑器的重新编辑
func (ip *InputProcessor) EditorRe(text string) string {
	text = strings.Replace(text, "&acute;", "`", -1)
	text = strings.Replace(text, "&#039;", "''", -1)
	return text
}

// 处理文本域的输出，主要就是对换行符进行处理，将\r\n之类的转换成<p>或<br>，thetype为true则转成<p>，type为false则转为<br>
func (ip *InputProcessor) TextareaOut(text string, thetype bool) string {
	if thetype == true {
		theReplaceMap := make(map[string]string)
		theReplaceMap["\r\n"] = "</p><p>"
		theReplaceMap["\r"] = "</p><p>"
		theReplaceMap["\n"] = "</p><p>"
		for index, value := range theReplaceMap {
			text = strings.Replace(text, index, value, -1)
		}
		text = "<p>" + text + "</p>"
	} else {
		theReplaceMap := make(map[string]string)
		theReplaceMap["\r\n"] = "<br>"
		theReplaceMap["\r"] = "</br>"
		theReplaceMap["\n"] = "<br>"
		for index, value := range theReplaceMap {
			text = strings.Replace(text, index, value, -1)
		}
	}
	return text
}

// 输入是否为整数
func (ip *InputProcessor) Int(text string, cannull bool, min, max int64) (num int64, err int) {
	text = strings.TrimSpace(text)
	err = 0
	textLen := len(text)
	if cannull == false && textLen == 0 {
		num = 0
		err = 1
		return
	}
	if textLen == 0 {
		num = 0
		return
	}
	num, e := strconv.ParseInt(text, 10, 64)
	if e != nil {
		err = 2
		return
	}
	if num < min || num > max {
		err = 3
		return
	}
	return
}

// 输入是否为浮点
func (ip *InputProcessor) Float(text string, cannull bool, min, max float64) (num float64, err int) {
	text = strings.TrimSpace(text)
	err = 0
	textLen := len(text)
	if cannull == false && textLen == 0 {
		num = 0
		err = 1
		return
	}
	if textLen == 0 {
		num = 0
		return
	}
	num, e := strconv.ParseFloat(text, 64)
	if e != nil {
		err = 2
		return
	}
	if num < min || num > max {
		err = 3
		return
	}
	return
}

// 枚举，判断提供的字符串是否出现在字符串切片里
func (ip *InputProcessor) Enum(text string, cannull bool, enum []string) (textc string, err int) {
	err = 0
	textc = ""
	text = strings.TrimSpace(text)
	textLen := len(text)
	if cannull == false && textLen == 0 {
		err = 1
		return
	}
	if textLen == 0 {
		return
	}
	for _, v := range enum {
		if text == v {
			textc = text
			return
		}
	}
	err = 2
	return
}

// 查看格式是不是邮箱地址
func (ip *InputProcessor) Email(text string, cannull bool, min, max int64) (textc string, err int) {
	err = 0
	textc = ""
	text = strings.TrimSpace(text)
	textLen := len(text)
	if cannull == false && textLen == 0 {
		err = 1
		return
	}
	if textLen == 0 {
		return
	}
	err = ip.minOrMax(text, min, max)
	if err != 0 {
		return
	}
	ifmatch := ip.reEmail.MatchString(text)
	if ifmatch == false {
		err = 4
		return
	}
	textc = text
	return
}

// 处理连接，主要是如果没有协议的话，就加上http://
func (ip *InputProcessor) Url(text string, cannull bool, min, max int64) (textc string, err int) {
	err = 0
	textc = ""
	text = strings.TrimSpace(text)
	textLen := len(text)
	if cannull == false && textLen == 0 {
		err = 1
		return
	}
	if textLen == 0 {
		return
	}
	err = ip.minOrMax(text, min, max)
	if err != 0 {
		return
	}
	ifmatch := ip.reUrl.MatchString(text)
	if ifmatch == false {
		textc = "http://" + text
	} else {
		textc = text
	}
	return
}

// 正则判断，提供一个正则表达，如果匹配则返回字符串
func (ip *InputProcessor) Regular(text string, cannull bool, rg string) (textc string, err int) {
	err = 0
	textc = ""
	text = strings.TrimSpace(text)
	textLen := len(text)
	if cannull == false && textLen == 0 {
		err = 1
		return
	}
	if textLen == 0 {
		return
	}
	therp, ther := regexp.MatchString(rg, text)
	if ther != nil {
		err = 2
		return
	}
	if therp == false {
		err = 3
		return
	}
	textc = text
	return
}

// 处理密码
func (ip *InputProcessor) Password(text string, cannull bool) (textc string, err int) {
	err = 0
	textc = ""
	text = strings.TrimSpace(text)
	textLen := len(text)
	if cannull == false && textLen == 0 {
		err = 1
		return
	}
	if textLen == 0 {
		return
	}
	if textLen != 40 {
		err = 2
		return
	}
	textc = text
	return
}

// 查看两遍密码
func (ip *InputProcessor) PasswordTwo(text, text2 string, cannull bool) (textc string, err int) {
	err = 0
	textc = ""
	text = strings.TrimSpace(text)
	text2 = strings.TrimSpace(text2)
	textLen := len(text)
	if cannull == false && textLen == 0 {
		err = 1
		return
	}
	if textLen == 0 {
		return
	}
	if textLen != 40 {
		err = 2
		return
	}
	if text != text2 {
		err = 3
		return
	}
	textc = text
	return
}

// 字符串枚举
func (ip *InputProcessor) StringEnum(text string, canull bool, enum []string) (textc string, err int) {
	err = 1
	textc = ""
	text = strings.TrimSpace(text)
	textLen := len(text)
	if canull == false && textLen == 0 {
		err = 1
		return
	}
	if textLen == 0 {
		err = 0
		return
	}
	for _, one := range enum {
		if one == text {
			err = 0
			return
		}
	}
	return
}
