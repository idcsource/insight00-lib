// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package nst

type ConnExec struct {
	Transmission *Transmission
	long_shot    bool // if is long connection, the bool is true
}

func NewConnExec(trans *Transmission) (connExec *ConnExec) {
	return &ConnExec{Transmission: trans}
}

func (c *ConnExec) SetLong() {
	c.long_shot = true
}

func (c *ConnExec) SetShort() {
	c.long_shot = false
}

func (c *ConnExec) SendData(data []byte) (err error) {
	return c.Transmission.SendData(data)
}

func (c *ConnExec) GetData() (data []byte, err error) {
	return c.Transmission.GetData()
}

func (c *ConnExec) SendClose() (err error) {
	return c.Transmission.Close()
}
