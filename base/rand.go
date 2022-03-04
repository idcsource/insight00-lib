// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

// random包提供一套生成随机字符串的方法。
// 包括根据字符串生成sha1和sha256散列值，生成40位或成倍的Unid的方法。
package base

import (
	crand "crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	mrand "math/rand"
	"strconv"
	"time"
)

// 获取字符串的Sha1散列值。
func GetSha1Sum(text string) string {
	t := sha1.New()
	t.Write([]byte(text))
	return hex.EncodeToString(t.Sum(nil))
}

// 获取字节切片的Sha1散列值。
func GetSha1SumBytes(b []byte) string {
	t := sha1.New()
	t.Write(b)
	return hex.EncodeToString(t.Sum(nil))
}

func GetSha256Sum(text string) string {
	s256 := sha256.Sum256([]byte(text))
	write := make([]byte, 0, len(s256))
	for i := range s256 {
		write = append(write, s256[i])
	}
	return hex.EncodeToString(write)
}

// 获取字符串的Sha512散列值。
func GetSha512Sum(text string) string {
	s512 := sha512.Sum512([]byte(text))
	write := make([]byte, 0, len(s512))
	for i := range s512 {
		write = append(write, s512[i])
	}
	return hex.EncodeToString(write)
}

// 获取len长度随机字符串。
func GetRand(len int) string {
	rands := make([]byte, len)
	crand.Read(rands)
	return hex.EncodeToString(rands)
}

// 获取不大于max的随机正整数。
func GetRandNum(max int) int {
	if max <= 0 {
		return 0
	}
	r := mrand.New(mrand.NewSource(time.Now().UnixNano()))
	return r.Intn(max)
}

// 根据输入的text生成40字符长度的（最终是sha1）随机字符串。
func unid_one(text ...string) string {
	var rt string
	rt += strconv.Itoa(GetRandNum(1000))
	rt += GetRand(GetRandNum(100))
	for _, thet := range text {
		rt += GetSha512Sum(thet)
	}
	rt += strconv.FormatInt(time.Now().UnixNano(), 10)

	return GetSha1Sum(rt)
}

// 根据输入的text生成唯一ID。length最小需要为1，返回的字符串的长度为length*40个字符长，为length个sha1的值串联。
func Unid(length uint, text ...string) string {
	var allstring string
	for i := length; i >= 1; i-- {
		allstring += unid_one(text...)
	}
	return allstring
}
