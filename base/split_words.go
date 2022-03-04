// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

// Insight 0+0各包共同使用的辅助函数
package base

import (
	"fmt"
	"regexp"
	"strings"
)

// Split the string, if trim is true, the split with out space.
func CommandSplit(command string, trim bool) (split []string, err error) {

	command = strings.TrimSpace(command)

	regexpm := make(map[string]*regexp.Regexp)
	regexpm["space"], _ = regexp.Compile(`[^ ]+`)
	regexpm["'b"], _ = regexp.Compile(`^'`)
	regexpm["b'"], _ = regexp.Compile(`'$`)
	regexpm["''"], _ = regexp.Compile(`^\'`)

	crune := []rune(command)
	//crunelen := len(crune)
	split = make([]string, 0)
	temprune := make([]rune, 0)
	inquot := false // 是否在引号里
	for i, onerune := range crune {
		if inquot == true {
			// 在引号里怎么处理
			if onerune == '"' && crune[i-1] == '\\' {
				// 转义字符 \"
				temprune[len(temprune)-1] = '"'
			} else if onerune == '"' {
				// 引号结尾 "
				split = append(split, string(temprune))
				temprune = make([]rune, 0)
				inquot = false
			} else {
				// 正常字符
				temprune = append(temprune, onerune)
			}
		} else {
			// 不在引号里
			if onerune == '"' {
				// 碰到引号怎么办
				inquot = true
			} else if onerune == ' ' && trim == false {
				// 如果是空格
				if len(temprune) == 0 {
					split = append(split, string(temprune))
					temprune = make([]rune, 0)
				} else {
					split = append(split, string(temprune))
					temprune = make([]rune, 0)
					split = append(split, string(temprune))
				}
			} else if onerune == ' ' {
				if len(temprune) != 0 {
					split = append(split, string(temprune))
					temprune = make([]rune, 0)
				}
			} else {
				// 正常字符
				temprune = append(temprune, onerune)
			}
		}

	}
	if inquot == true {
		err = fmt.Errorf("Command syntax error.")
		return
	}
	if len(temprune) != 0 {
		split = append(split, string(temprune))
	}

	return
}

// 将提供的字符串进行拆分词语处理
func SplitWords(str string) (normal [][]string) {
	strn := []rune(str)
	var tmpstring string
	tempslice := make([]string, 0)
	for _, one := range strn {
		// 碰到段落就新建一个切片
		if subsection(one) == true {
			if len(tmpstring) > 0 {
				tmpstring = optDot(tmpstring)
				tempslice = append(tempslice, tmpstring)
				tmpstring = ""
			}
			if len(tempslice) > 0 {
				normal = append(normal, tempslice)
				tempslice = make([]string, 0)
			}
			continue
		}
		// 如果碰到的是单字节的字，并且不是空格，就加入临时字符串
		if len([]byte(string(one))) == 1 && string(one) != " " {
			tmpstring += string(one)
			continue
		}
		// 如果碰到空格，如果临时字符串里有东西，则加入临时切片
		if string(one) == " " || string(one) == "\n" || string(one) == "\r" || string(one) == "\t" {
			if len(tmpstring) > 0 {
				tmpstring = optDot(tmpstring)
				tempslice = append(tempslice, tmpstring)
				tmpstring = ""
			}
			continue
		}
		// 普通的字符，就直接加入临时切片
		tempslice = append(tempslice, string(one))
	}
	if len(tmpstring) > 0 {
		tempslice = append(tempslice, tmpstring)
	}
	if len(tempslice) > 0 {
		normal = append(normal, tempslice)
	}
	return
}

// 处理最后的标点
func optDot(str string) string {
	str = strings.Trim(str, ".")
	str = strings.Trim(str, ",")
	str = strings.Trim(str, "!")
	return str
}

// 划分段落
func subsection(str rune) bool {
	pun := []string{
		"。",
		"　",
		"·",
		"，",
		"！",
		"；",
		";",
		"？",
		"：",
		"、",
		"“",
		"”",
		"\"",
		"'",
		"<",
		">",
		"《",
		"》",
		"(",
		")",
		"（",
		"）",
		"…",
		"}",
		"{",
		"\n",
		"\r",
		"\t",
	}
	strs := string(str)
	for _, one := range pun {
		if one == strs {
			return true
		}
	}
	return false
}

func strpos(str, substr string) int {
	// 子串在字符串的字节位置
	result := strings.Index(str, substr)
	if result >= 0 {
		// 获得子串之前的字符串并转换成[]byte
		prefix := []byte(str)[0:result]
		// 将子串之前的字符串转换成[]rune
		rs := []rune(string(prefix))
		// 获得子串之前的字符串的长度，便是子串在字符串的字符位置
		result = len(rs)
	}

	return result
}

// EncodeSalt 加盐
func EncodeSalt(str, stream, salt string) string {
	salt = GetSha1Sum(salt)
	tmpStream := ""
	lockLen := len(stream)
	j := 0
	k := 0
	streamb := []byte(stream)
	for i := 0; i < len(str); i++ {
		if k == len(salt) {
			k = 0
		}
		strb := []byte(str)
		stri := strb[i]
		saltb := []byte(salt)
		saltk := saltb[k]
		j = (strpos(stream, string(stri)) + int(saltk)) % (lockLen)

		streamj := streamb[j]
		tmpStream += string(streamj)
		k++
	}
	return tmpStream
}

// DecodeSalt 解盐
func DecodeSalt(str, stream, salt string) string {
	salt = GetSha1Sum(salt)
	tmpStream := ""
	lockLen := len(stream)
	j := 0
	k := 0
	streamb := []byte(stream)
	for i := 0; i < len(str); i++ {
		if k == len(salt) {
			k = 0
		}
		strb := []byte(str)
		stri := strb[i]
		saltb := []byte(salt)
		saltk := saltb[k]
		j = strpos(stream, string(stri)) - int(saltk)
		for j < 0 {
			j = j + lockLen
		}

		streamj := streamb[j]
		tmpStream += string(streamj)
		k++
	}
	return tmpStream
}
