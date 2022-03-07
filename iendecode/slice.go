// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package iendecode

import "bytes"

func SliceStringToBytes(s []string) (b []byte) {
	var b_buf bytes.Buffer
	thecount := int64(len(s))
	b_buf.Write(Int64ToBytes(thecount))
	for i := range s {
		str_b := []byte(s[i])
		str_b_len := int64(len(str_b))
		b_buf.Write(Int64ToBytes(str_b_len))
		b_buf.Write(str_b)
	}
	b = b_buf.Bytes()
	return
}

func BytesToSliceString(b []byte) (s []string) {
	b_buf := bytes.NewBuffer(b)
	thecount := BytesToInt64(b_buf.Next(8))
	s = make([]string, thecount)
	var i int64 = 0
	for {
		if i >= thecount {
			break
		}
		thelen := BytesToInt64(b_buf.Next(8))
		s[i] = string(b_buf.Next(int(thelen)))
		i++
	}
	return
}
