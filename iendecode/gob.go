// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package iendecode

import (
	"bytes"
	"encoding/gob"
	"reflect"
)

// 将结构体数据转成Gob再转成[]Byte
func StructGobBytes(e interface{}) ([]byte, error) {
	var gob_buff bytes.Buffer           //建立缓冲
	gob_en := gob.NewEncoder(&gob_buff) //gob开始编码
	err := gob_en.Encode(e)             //gob编码
	if err != nil {
		return nil, err
	}
	gob_b := gob_buff.Bytes() //bytes.buffer转[]byte
	return gob_b, nil
}

// 将[]byte转成Gob再转成结构体
func BytesGobStruct(f_b []byte, stur interface{}) error {
	b_buf := bytes.NewBuffer(f_b) //将[]byte放入bytes的buffer中
	b_go := gob.NewDecoder(b_buf) //将buffer放入gob的decoder中
	err := b_go.Decode(stur)      //将gob解码放入stur
	if err != nil {
		return err
	} else {
		return nil
	}
}

// 将[]byte转成提供的反射Value
func BytesGobReflect(f_b []byte, v reflect.Value) error {
	b_buf := bytes.NewBuffer(f_b) //将[]byte放入bytes的buffer中
	b_go := gob.NewDecoder(b_buf) //将buffer放入gob的decoder中
	err := b_go.DecodeValue(v)    //将gob解码放入
	if err != nil {
		return err
	} else {
		return nil
	}
}
