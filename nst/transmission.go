// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package nst

import (
	"crypto/tls"
	"fmt"
	"net"

	"github.com/idcsource/insight00-lib/iendecode"
)

// 发送接收数据结构。发送数据的方法必须使用对应类型的接收方法进行接收
type Transmission struct {
	tls      bool
	conn     net.Conn
	conn_tls *tls.Conn
	buf      int
}

// 新建发送接收
func NewTransmission(conn net.Conn) *Transmission {
	return &Transmission{tls: false, conn: conn, buf: 1024}
}

func NewTransmissionTLS(conn *tls.Conn) *Transmission {
	return &Transmission{tls: true, conn_tls: conn, buf: 1024}
}

func (t *Transmission) RemoteAddr() (addr net.Addr) {
	if t.tls == false {
		addr = t.conn.RemoteAddr()
	} else {
		addr = t.conn_tls.RemoteAddr()
	}
	return
}

// 设置缓冲大小
func (t *Transmission) SetBuf(buf int) {
	t.buf = buf
}

// 查看缓冲大小
func (t *Transmission) GetBuf() int {
	return t.buf
}

// 发送一个结构体（会在方法内部转换成gob），只能用GetStruct()方法接受
func (t *Transmission) SendStruct(stru interface{}) (errs error) {
	bytes, errs := iendecode.StructGobBytes(stru)
	if errs != nil {
		return
	}
	len := len(bytes)
	lens := uint64(len)
	errs = t.SendLen(lens) //发送长度
	if errs != nil {
		return
	}
	errs = t.SendBytes(bytes)
	return
}

// 接收一个结构体（会在方法内部从gob转出），只能接受SendStruct()方法发送的信息
func (t *Transmission) GetStruct(stru interface{}) (errs error) {
	lens, errs := t.GetLen()
	if errs != nil {
		return
	}
	bytes, errs := t.GetBytes(lens)
	if errs != nil {
		return
	}
	errs = iendecode.BytesGobStruct(bytes, stru)
	return
}

// 发送一串数据流（包括字节流的长度），必须用GetData()方法接收
func (t *Transmission) SendData(bytes []byte) (errs error) {
	len := len(bytes)
	lens := uint64(len)
	errs = t.SendLen(lens) //发送长度
	if errs != nil {
		return
	}
	errs = t.SendBytes(bytes)
	return
}

// 接收一串数据流（包括字节流的长度），特定接收SendData()方法发送的数据
func (t *Transmission) GetData() (bytes []byte, errs error) {
	lens, errs := t.GetLen()
	if errs != nil {
		return nil, errs
	}
	bytes, errs = t.GetBytes(lens)
	return
}

// 发送一个长度属性，也就是发送uint64
func (t *Transmission) SendLen(len uint64) (errs error) {
	vb := Uint64ToBytes(len)
	errs = t.SendBytes(vb)
	return
}

// 接收一个长度属性，也就是接收uint64
func (t *Transmission) GetLen() (len uint64, errs error) {
	bytes, errs := t.GetBytes(8)
	if errs != nil {
		return 0, errs
	}
	len = BytesToUint64(bytes)
	return
}

// 发送一个流程版本号，也可以发送流程编号，就是发送一个uint32
func (t *Transmission) SendVer(version uint32) (errs error) {
	vb := Uint32ToBytes(version)
	errs = t.SendBytes(vb)
	return
}

// 接收一个流程版本号，也就是接收一个uint32
func (t *Transmission) GetVer() (version uint32, errs error) {
	bytes, errs := t.GetBytes(4)
	version = BytesToUint32(bytes)
	return
}

// 发送状态，也就是发送uint8
func (t *Transmission) SendStat(status uint8) (errs error) {
	vb := Uint8ToBytes(status)
	errs = t.SendBytes(vb)
	return
}

// 接收状态，也就是接收uint8
func (t *Transmission) GetStat() (status uint8, errs error) {
	bytes, errs := t.GetBytes(1)
	if errs != nil {
		return
	}
	status = BytesToUint8(bytes)
	return
}

// 发送字节切片（没有字节长度信息）
func (t *Transmission) SendBytes(bytes []byte) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("nst2[Transmission]SendBytes: %v", e)
		}
	}()
	if t.tls == false {
		_, err = t.conn.Write(bytes)
	} else {
		_, err = t.conn_tls.Write(bytes)
	}
	if err != nil {
		return
	}
	return
}

// 接收字节切片（没有字节长度信息）
func (t *Transmission) GetBytes(len uint64) (returnByte []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("nst2[Transmission]GetBytes: %v", e)
		}
	}()
	returnByte = make([]byte, 0, len)
	for {
		var err error
		tempdata := []byte{}
		if len < uint64(t.buf) {
			tempdata = make([]byte, len)
		} else {
			tempdata = make([]byte, t.buf)
		}
		var r int
		if t.tls == false {
			r, err = t.conn.Read(tempdata)
		} else {
			r, err = t.conn_tls.Read(tempdata)
		}
		if err != nil {
			return returnByte, err
		}
		returnByte = append(returnByte, tempdata[:r]...)

		len = len - uint64(r)

		if len == 0 {
			break
		}
	}
	return returnByte, err
}

// 设置长连接模式
//func (t *TCP) SetKeepAlive(keepalive bool) error {
//	return t.tcp.SetKeepAlive(keepalive)
//}

// 关闭连接
func (t *Transmission) Close() (err error) {
	if t.tls == true {
		err = t.conn_tls.Close()
	} else {
		err = t.conn.Close()
	}
	return
}
