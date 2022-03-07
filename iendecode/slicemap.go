// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package iendecode

import (
	"bytes"
	"fmt"
	"time"
)

// Support bool, byte, all int, all uint, all float, all complex, string, []byte, time.Time.
func SingleToBytes(thetype string, data interface{}) (b []byte, err error) {
	switch thetype {
	case "bool":
		b = BoolToBytes(data.(bool))
	case "byte":
		b = make([]byte, 1)
		b[0] = data.(byte)
	case "int8":
		b = Int8ToBytes(data.(int8))
	case "int16":
		b = Int16ToBytes(data.(int16))
	case "int32":
		b = Int32ToBytes(data.(int32))
	case "int64":
		b = Int64ToBytes(data.(int64))
	case "uint8":
		b = Uint8ToBytes(data.(uint8))
	case "uint16":
		b = Uint16ToBytes(data.(uint16))
	case "uint32":
		b = Uint32ToBytes(data.(uint32))
	case "uint64":
		b = Uint64ToBytes(data.(uint64))
	case "int":
		b = Int64ToBytes(data.(int64))
	case "uint":
		b = Uint64ToBytes(data.(uint64))
	case "float32":
		b, err = ToBinary(data)
		if err != nil {
			return
		}
	case "float64":
		b, err = ToBinary(data)
		if err != nil {
			return
		}
	case "complex64":
		b, err = ToBinary(data)
		if err != nil {
			return
		}
	case "complex128":
		b, err = ToBinary(data)
		if err != nil {
			return
		}
	case "string":
		b = []byte(data.(string))
	case "[]byte":
		b = data.([]byte)
	case "time.Time":
		b, err = data.(time.Time).MarshalBinary()
	default:
		err = fmt.Errorf("Type not support")
	}
	return
}

// Support bool, byte, all int, all uint, all float, all complex,
// which implements the encoding.BinaryMarshaler interface, which support gob encode.
func BytesToSingle(thetype string, b []byte) (data interface{}, err error) {
	switch thetype {
	case "bool":
		data = BytesToBool(b)
	case "byte":
		data = b[0]
	case "int8":
		data = BytesToInt8(b)
	case "int16":
		data = BytesToInt16(b)
	case "int32":
		data = BytesToInt32(b)
	case "int64":
		data = BytesToInt64(b)
	case "uint8":
		data = BytesToUint8(b)
	case "uint16":
		data = BytesToUint16(b)
	case "uint32":
		data = BytesToUint32(b)
	case "uint64":
		data = BytesToUint64(b)
	case "int":
		data = int(BytesToInt64(b))
	case "uint":
		data = uint(BytesToUint64(b))
	case "float32":
		err = FromBinary(b, data)
		if err != nil {
			return
		}
	case "float64":
		err = FromBinary(b, data)
		if err != nil {
			return
		}
	case "complex64":
		err = FromBinary(b, data)
		if err != nil {
			return
		}
	case "complex128":
		err = FromBinary(b, data)
		if err != nil {
			return
		}
	case "string":
		data = string(b)
	case "[]byte":
		data = b
	case "time.Time":
		var t time.Time
		err = t.UnmarshalBinary(b)
		if err != nil {
			return
		}
		data = t
	default:
		err = fmt.Errorf("Type not support")
	}
	return
}

// Support bool, int, uint, int8, uint8, int64, uint64, float64, complex128, string, time.Time slice.
func BytesToSlice(thetype string, b []byte) (slice interface{}, err error) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	b_buf := bytes.NewBuffer(b)
	switch thetype {
	case "[]bool":
		thecount := BytesToInt64(b_buf.Next(8))
		slicet := make([]bool, thecount)
		var i int64 = 0
		for {
			if i >= thecount {
				break
			}
			slicet[i] = BytesToBool(b_buf.Next(1))
			i++
		}
		slice = slicet
	case "[]int":
		thecount := BytesToInt64(b_buf.Next(8))
		slicet := make([]int, thecount)
		var i int64 = 0
		for {
			if i >= thecount {
				break
			}
			slicet[i] = int(BytesToInt64(b_buf.Next(8)))
			i++
		}
		slice = slicet
	case "[]uint":
		thecount := BytesToInt64(b_buf.Next(8))
		slicet := make([]uint, thecount)
		var i int64 = 0
		for {
			if i >= thecount {
				break
			}
			slicet[i] = uint(BytesToUint64(b_buf.Next(8)))
			i++
		}
		slice = slicet
	case "[]int8":
		thecount := BytesToInt64(b_buf.Next(8))
		slicet := make([]int8, thecount)
		var i int64 = 0
		for {
			if i >= thecount {
				break
			}
			slicet[i] = BytesToInt8(b_buf.Next(1))
			i++
		}
		slice = slicet
	case "[]uint8":
		thecount := BytesToInt64(b_buf.Next(8))
		slicet := make([]uint8, thecount)
		var i int64 = 0
		for {
			if i >= thecount {
				break
			}
			slicet[i] = BytesToUint8(b_buf.Next(1))
			i++
		}
		slice = slicet
	case "[]int64":
		thecount := BytesToInt64(b_buf.Next(8))
		slicet := make([]int64, thecount)
		var i int64 = 0
		for {
			if i >= thecount {
				break
			}
			slicet[i] = BytesToInt64(b_buf.Next(8))
			i++
		}
		slice = slicet
	case "[]uint64":
		thecount := BytesToInt64(b_buf.Next(8))
		slicet := make([]uint64, thecount)
		var i int64 = 0
		for {
			if i >= thecount {
				break
			}
			slicet[i] = BytesToUint64(b_buf.Next(8))
			i++
		}
		slice = slicet
	case "[]float64":
		thecount := BytesToInt64(b_buf.Next(8))
		slicet := make([]float64, thecount)
		var i int64 = 0
		for {
			if i >= thecount {
				break
			}
			var f float64
			FromBinary(b_buf.Next(8), &f)
			slicet[i] = f
			i++
		}
		slice = slicet
	case "[]complex128":
		thecount := BytesToInt64(b_buf.Next(8))
		slicet := make([]complex128, thecount)
		var i int64 = 0
		for {
			if i >= thecount {
				break
			}
			var f complex128
			FromBinary(b_buf.Next(16), &f)
			slicet[i] = f
			i++
		}
		slice = slicet
	case "[]string":
		thecount := BytesToInt64(b_buf.Next(8))
		slicet := make([]string, thecount)
		var i int64 = 0
		for {
			if i >= thecount {
				break
			}
			thelen := BytesToInt64(b_buf.Next(8))
			slicet[i] = string(b_buf.Next(int(thelen)))
			i++
		}
		slice = slicet
	case "[]time.Time":
		thecount := BytesToInt64(b_buf.Next(8))
		slicet := make([]time.Time, thecount)
		var i int64 = 0
		for {
			if i >= thecount {
				break
			}
			var t time.Time
			err = t.UnmarshalBinary(b_buf.Next(15))
			if err != nil {
				return
			}
			slicet[i] = t
			i++
		}
		slice = slicet
	default:
		err = fmt.Errorf("Type not support")
	}
	return
}

// Support bool, int, uint, int8, uint8, int64, uint64, float64, complex128, string, time.Time slice.
func SliceToBytes(thetype string, slice interface{}) (b []byte, err error) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	var b_buf bytes.Buffer
	switch thetype {
	case "[]bool":
		slicet, ok := slice.([]bool)
		if ok == false {
			err = fmt.Errorf("Type Error.")
			return
		}
		thecount := int64(len(slicet))
		b_buf.Write(Int64ToBytes(thecount))
		for i := range slicet {
			var d_b []byte
			d_b = BoolToBytes(slicet[i])
			b_buf.Write(d_b)
		}
	case "[]int":
		slicet, ok := slice.([]int)
		if ok == false {
			err = fmt.Errorf("Type Error.")
			return
		}
		thecount := int64(len(slicet))
		b_buf.Write(Int64ToBytes(thecount))
		for i := range slicet {
			var d_b []byte
			d_b = Int64ToBytes(int64(slicet[i]))
			b_buf.Write(d_b)
		}
	case "[]uint":
		slicet, ok := slice.([]uint)
		if ok == false {
			err = fmt.Errorf("Type Error.")
			return
		}
		thecount := int64(len(slicet))
		b_buf.Write(Int64ToBytes(thecount))
		for i := range slicet {
			var d_b []byte
			d_b = Uint64ToBytes(uint64(slicet[i]))
			b_buf.Write(d_b)
		}
	case "[]int8":
		slicet, ok := slice.([]int8)
		if ok == false {
			err = fmt.Errorf("Type Error.")
			return
		}
		thecount := int64(len(slicet))
		b_buf.Write(Int64ToBytes(thecount))
		for i := range slicet {
			var d_b []byte
			d_b = Int8ToBytes(slicet[i])
			b_buf.Write(d_b)
		}
	case "[]uint8":
		slicet, ok := slice.([]uint8)
		if ok == false {
			err = fmt.Errorf("Type Error.")
			return
		}
		thecount := int64(len(slicet))
		b_buf.Write(Int64ToBytes(thecount))
		for i := range slicet {
			var d_b []byte
			d_b = Uint8ToBytes(slicet[i])
			b_buf.Write(d_b)
		}
	case "[]int64":
		slicet, ok := slice.([]int64)
		if ok == false {
			err = fmt.Errorf("Type Error.")
			return
		}
		thecount := int64(len(slicet))
		b_buf.Write(Int64ToBytes(thecount))
		for i := range slicet {
			var d_b []byte
			d_b = Int64ToBytes(slicet[i])
			b_buf.Write(d_b)
		}
	case "[]uint64":
		slicet, ok := slice.([]uint64)
		if ok == false {
			err = fmt.Errorf("Type Error.")
			return
		}
		thecount := int64(len(slicet))
		b_buf.Write(Int64ToBytes(thecount))
		for i := range slicet {
			var d_b []byte
			d_b = Uint64ToBytes(slicet[i])
			b_buf.Write(d_b)
		}
	case "[]float64":
		slicet, ok := slice.([]float64)
		if ok == false {
			err = fmt.Errorf("Type Error.")
			return
		}
		thecount := int64(len(slicet))
		b_buf.Write(Int64ToBytes(thecount))
		for i := range slicet {
			var d_b []byte
			d_b, err = ToBinary(slicet[i])
			if err != nil {
				return
			}
			b_buf.Write(d_b)
		}
	case "[]complex128":
		slicet, ok := slice.([]complex128)
		if ok == false {
			err = fmt.Errorf("Type Error.")
			return
		}
		thecount := int64(len(slicet))
		b_buf.Write(Int64ToBytes(thecount))
		for i := range slicet {
			var d_b []byte
			d_b, err = ToBinary(slicet[i])
			if err != nil {
				return
			}
			b_buf.Write(d_b)
		}
	case "[]string":
		slicet, ok := slice.([]string)
		if ok == false {
			err = fmt.Errorf("Type Error.")
			return
		}
		thecount := int64(len(slicet))
		b_buf.Write(Int64ToBytes(thecount))
		for i := range slicet {
			str_b := []byte(slicet[i])
			str_b_len := int64(len(str_b))
			b_buf.Write(Int64ToBytes(str_b_len))
			b_buf.Write(str_b)
		}
	case "[]time.Time":
		slicet, ok := slice.([]time.Time)
		if ok == false {
			err = fmt.Errorf("Type Error.")
			return
		}
		thecount := int64(len(slicet))
		b_buf.Write(Int64ToBytes(thecount))
		for i := range slicet {
			var d_b []byte
			d_b, err = slicet[i].MarshalBinary()
			if err != nil {
				return
			}
			b_buf.Write(d_b)
		}
	default:
		err = fmt.Errorf("Type not support")

	}
	b = b_buf.Bytes()
	return

}

// Support map[string]string, map[string]time.Time, map[string]int64, map[string]uint64,
// map[string]float64, map[string]complex128.
func MapToBytes(thetype string, themaps interface{}) (b []byte, err error) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	var b_buf bytes.Buffer
	switch thetype {
	case "map[string]string":
		themap := themaps.(map[string]string)
		thecount := int64(len(themap))
		b_buf.Write(Int64ToBytes(thecount))
		for key, _ := range themap {
			key_b := []byte(key)
			key_b_len := int64(len(key_b))
			b_buf.Write(Int64ToBytes(key_b_len))
			b_buf.Write(key_b)
			var bs []byte
			bs = []byte(themap[key])
			b_buf.Write(Int64ToBytes(int64(len(bs))))
			b_buf.Write(bs)
		}
	case "map[string]time.Time":
		themap := themaps.(map[string]time.Time)
		thecount := int64(len(themap))
		b_buf.Write(Int64ToBytes(thecount))
		for key, _ := range themap {
			key_b := []byte(key)
			key_b_len := int64(len(key_b))
			b_buf.Write(Int64ToBytes(key_b_len))
			b_buf.Write(key_b)
			var bs []byte
			bs, err = themap[key].MarshalBinary()
			if err != nil {
				return
			}
			b_buf.Write(bs)
		}
	case "map[string]int64":
		themap := themaps.(map[string]int64)
		thecount := int64(len(themap))
		b_buf.Write(Int64ToBytes(thecount))
		for key, _ := range themap {
			key_b := []byte(key)
			key_b_len := int64(len(key_b))
			b_buf.Write(Int64ToBytes(key_b_len))
			b_buf.Write(key_b)
			var bs []byte
			bs = Int64ToBytes(themap[key])
			b_buf.Write(bs)
		}
	case "map[string]uint64":
		themap := themaps.(map[string]uint64)
		thecount := int64(len(themap))
		b_buf.Write(Int64ToBytes(thecount))
		for key, _ := range themap {
			key_b := []byte(key)
			key_b_len := int64(len(key_b))
			b_buf.Write(Int64ToBytes(key_b_len))
			b_buf.Write(key_b)
			var bs []byte
			bs = Uint64ToBytes(themap[key])
			b_buf.Write(bs)
		}
	case "map[string]float64":
		themap := themaps.(map[string]float64)
		thecount := int64(len(themap))
		b_buf.Write(Int64ToBytes(thecount))
		for key, _ := range themap {
			key_b := []byte(key)
			key_b_len := int64(len(key_b))
			b_buf.Write(Int64ToBytes(key_b_len))
			b_buf.Write(key_b)
			var bs []byte
			bs, err = ToBinary(themap[key])
			if err != nil {
				return
			}
			b_buf.Write(bs)
		}
	case "map[string]complex128":
		themap := themaps.(map[string]complex128)
		thecount := int64(len(themap))
		b_buf.Write(Int64ToBytes(thecount))
		for key, _ := range themap {
			key_b := []byte(key)
			key_b_len := int64(len(key_b))
			b_buf.Write(Int64ToBytes(key_b_len))
			b_buf.Write(key_b)
			var bs []byte
			bs, err = ToBinary(themap[key])
			if err != nil {
				return
			}
			b_buf.Write(bs)
		}
	default:
		err = fmt.Errorf("Type not support")
	}

	b = b_buf.Bytes()
	return
}

// Support map[string]string, map[string]time.Time, map[string]int64, map[string]uint64,
// map[string]float64, map[string]complex128.
func BytesToMap(thetype string, b []byte) (themaps interface{}, err error) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	b_buf := bytes.NewBuffer(b)
	switch thetype {
	case "map[string]string":
		thecount := BytesToInt64(b_buf.Next(8))
		themap := make(map[string]string)
		var i int64 = 0
		for {
			if i >= thecount {
				break
			}
			key_b_len := BytesToInt64(b_buf.Next(8))
			key := string(b_buf.Next(int(key_b_len)))
			v_b_len := BytesToInt64(b_buf.Next(8))
			v := string(b_buf.Next(int(v_b_len)))
			themap[key] = v
			i++
		}
		themaps = themap
	case "map[string]time.Time":
		thecount := BytesToInt64(b_buf.Next(8))
		themap := make(map[string]time.Time)
		var i int64 = 0
		for {
			if i >= thecount {
				break
			}
			key_b_len := BytesToInt64(b_buf.Next(8))
			key := string(b_buf.Next(int(key_b_len)))
			var v time.Time
			err = v.UnmarshalBinary(b_buf.Next(15))
			if err != nil {
				return
			}
			themap[key] = v
			i++
		}
		themaps = themap
	case "map[string]int64":
		thecount := BytesToInt64(b_buf.Next(8))
		themap := make(map[string]int64)
		var i int64 = 0
		for {
			if i >= thecount {
				break
			}
			key_b_len := BytesToInt64(b_buf.Next(8))
			key := string(b_buf.Next(int(key_b_len)))
			themap[key] = BytesToInt64(b_buf.Next(8))
			i++
		}
		themaps = themap
	case "map[string]uint64":
		thecount := BytesToInt64(b_buf.Next(8))
		themap := make(map[string]uint64)
		var i int64 = 0
		for {
			if i >= thecount {
				break
			}
			key_b_len := BytesToInt64(b_buf.Next(8))
			key := string(b_buf.Next(int(key_b_len)))
			themap[key] = BytesToUint64(b_buf.Next(8))
			i++
		}
		themaps = themap
	case "map[string]float64":
		thecount := BytesToInt64(b_buf.Next(8))
		themap := make(map[string]float64)
		var i int64 = 0
		for {
			if i >= thecount {
				break
			}
			key_b_len := BytesToInt64(b_buf.Next(8))
			key := string(b_buf.Next(int(key_b_len)))
			var f float64
			err = FromBinary(b_buf.Next(8), &f)
			themap[key] = f
			i++
		}
		themaps = themap
	case "map[string]complex128":
		thecount := BytesToInt64(b_buf.Next(8))
		themap := make(map[string]complex128)
		var i int64 = 0
		for {
			if i >= thecount {
				break
			}
			key_b_len := BytesToInt64(b_buf.Next(8))
			key := string(b_buf.Next(int(key_b_len)))
			var f complex128
			err = FromBinary(b_buf.Next(16), &f)
			themap[key] = f
			i++
		}
		themaps = themap
	default:
		err = fmt.Errorf("Type not support")
	}

	b = b_buf.Bytes()
	return
}
