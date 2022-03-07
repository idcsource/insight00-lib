// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package iendecode

import (
	"bytes"
	"encoding/json"
)

// 将结构体转成Json的字符串
func StructToJson(e interface{}) (str string, err error) {
	var j_buff bytes.Buffer          //建立缓冲
	j_en := json.NewEncoder(&j_buff) //json开始编码
	err = j_en.Encode(e)             //json编码
	if err != nil {
		return
	}
	//j_b := j_buff.Bytes() //bytes.buffer转[]byte

	str = j_buff.String()
	return
}

// 将json的字符串转成结构体
func JsonToStruct(f_b string, stur interface{}) (err error) {
	j_buf := bytes.NewBuffer([]byte(f_b)) //将[]byte放入bytes的buffer中
	j_go := json.NewDecoder(j_buf)        //将buffer放入json的decoder中
	err = j_go.Decode(stur)               //将json解码放入stur
	return
}
