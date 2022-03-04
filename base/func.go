// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

// Insight 0+0各包共同使用的辅助函数
package base

import (
	"sort"
)

// Is odd number
func IsOdd(num int) bool {
	if num%2 == 0 {
		return false
	}
	return true
}

// Slice里是否包含String
func StringInSlice(list []string, s string) bool {
	sort.Strings(list)
	llen := len(list)
	i := sort.SearchStrings(list, s)
	if i < llen && list[i] == s {
		return true
	} else {
		return false
	}
}
