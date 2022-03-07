// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package iendecode

import (
	"encoding/binary"
)

// Uint64转[]byte
func Uint64ToBytes(i uint64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, i)
	return buf
}

func Int64ToBytes(i int64) (b []byte) {
	b = make([]byte, 8)
	b[0] = byte(i >> 56)
	b[1] = byte(i >> 48)
	b[2] = byte(i >> 40)
	b[3] = byte(i >> 32)
	b[4] = byte(i >> 24)
	b[5] = byte(i >> 16)
	b[6] = byte(i >> 8)
	b[7] = byte(i)
	return
}

func IntToBytes(i int) (b []byte) {
	return Int64ToBytes(int64(i))
}

func UintToBytes(i uint) (b []byte) {
	return Uint64ToBytes(uint64(i))
}

// []byte转uint64
func BytesToUint64(buf []byte) uint64 {
	return binary.BigEndian.Uint64(buf)
}

func BytesToInt64(b []byte) (i int64) {
	_ = b[7]
	return int64(b[7]) | int64(b[6])<<8 | int64(b[5])<<16 | int64(b[4])<<24 | int64(b[3])<<32 | int64(b[2])<<40 | int64(b[1])<<48 | int64(b[0])<<56
}

func BytesToInt(b []byte) (i int) {
	return int(BytesToInt64(b))
}

func BytesToUint(b []byte) (i uint) {
	return uint(BytesToUint64(b))
}

// Uint32转[]byte
func Uint32ToBytes(i uint32) []byte {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, i)
	return buf
}

func Int32ToBytes(i int32) (b []byte) {
	b = make([]byte, 4)
	b[0] = byte(i >> 24)
	b[1] = byte(i >> 16)
	b[2] = byte(i >> 8)
	b[3] = byte(i)
	return
}

// []byte转uint32
func BytesToUint32(buf []byte) uint32 {
	return binary.BigEndian.Uint32(buf)
}

func BytesToInt32(b []byte) (i int32) {
	_ = b[3]
	return int32(b[3]) | int32(b[2])<<8 | int32(b[1])<<16 | int32(b[0])<<24
}

// Uint16转[]byte
func Uint16ToBytes(i uint16) []byte {
	var buf = make([]byte, 2)
	binary.BigEndian.PutUint16(buf, i)
	return buf
}

func Int16ToBytes(i int16) (b []byte) {
	b = make([]byte, 2)
	b[0] = byte(i >> 8)
	b[1] = byte(i)
	return
}

// []byte转uint16
func BytesToUint16(buf []byte) uint16 {
	return binary.BigEndian.Uint16(buf)
}

func BytesToInt16(b []byte) (i int16) {
	_ = b[1]
	return int16(b[1]) | int16(b[0])<<8
}

// Uint8转[]byte
func Uint8ToBytes(i uint8) []byte {
	var buf = []byte{i}
	return buf
}

func Int8ToBytes(i int8) (b []byte) {
	b = []byte{uint8(i)}
	return
}

// []byte转uint8
func BytesToUint8(buf []byte) uint8 {
	return uint8(buf[0])
}

func BytesToInt8(b []byte) (i int8) {
	return int8(b[0])
}

func BoolToBytes(bo bool) (b []byte) {
	b = make([]byte, 1)
	if bo == true {
		b[0] = 1
	} else {
		b[0] = 0
	}
	return
}

func BytesToBool(b []byte) (bo bool) {
	if b[0] == 1 {
		bo = true
	} else {
		bo = false
	}
	return
}
