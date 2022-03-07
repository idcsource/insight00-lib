// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package nst

import (
	"crypto/tls"
	"fmt"
	"net"
	"sync"
)

type Client struct {
	addr         string        // ip address an port
	iftls        bool          // if use tls
	max_count    int           // the max connect count
	now_count    int           // now the connect count
	lconnection  bool          // if use long connection
	connect_pool []*CConnect   // the connect pool
	closechan    chan bool     // the connect closed channel
	closed       bool          // the Client closed
	lock         *sync.RWMutex // the lock for assign the connection
}

// the Client connection
type CConnect struct {
	conn_exec   *ConnExec // the tcp data transmission
	addr        string    // address
	iftls       bool      // if use tls
	lconnection bool      // if use long connection
	paused      bool      // if it paused
	closed      bool      // if is closed
	closechan   chan bool // the conncet close channel from Client
	lock        bool      // the lock
}

// 新建一个客户端
func NewClient(addr string, max_count int, iftls bool) (c *Client, err error) {
	c = &Client{
		iftls:        iftls,
		addr:         addr,
		max_count:    max_count,
		now_count:    0,
		lconnection:  false,
		connect_pool: make([]*CConnect, 0),
		closechan:    make(chan bool),
		closed:       false,
		lock:         new(sync.RWMutex),
	}
	go c.checkCount()
	return
}

// 新建一个长连接客户端
func NewClientL(addr string, max_count int, iftls bool) (c *Client, err error) {
	c, err = NewClient(addr, max_count, iftls)
	c.lconnection = true
	return
}

func (c *Client) checkCount() {
	for {
		if c.closed == true {
			return
		}
		<-c.closechan
		c.now_count--
	}
}

// 关闭客户端
func (c *Client) Close() (err error) {
	c.closed = true
	if c.lconnection == true {
		for _, one := range c.connect_pool {
			one.closed = true
			err = one.conn_exec.Transmission.SendStat(uint8(SEND_STAT_CONN_CLOSE))
			if err == nil {
				one.conn_exec.Transmission.GetStat()
			}
			one.conn_exec.Transmission.Close()
		}
	}
	return
}

// 开始处理，也就是真正去做这个连接，与OpenConnect一样
func (c *Client) OpenProgress() (cc *CConnect, err error) {
	return c.OpenConnect()
}

// 开始处理，也就是真正去做这个连接
func (c *Client) OpenConnect() (cc *CConnect, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	// check if is the long connect
	if c.lconnection == true {
		// if long connect
		cc, err = c.selectFromPool()
	} else {
		// if short connect
		// if connections to much
		if c.now_count > c.max_count {
			err = fmt.Errorf("The number of connections is exceeded.")
			return
		}
		// create the connect.
		cc, err = c.dial()
		if err != nil {
			return
		}
		c.now_count++
	}
	return
}

func (c *Client) selectFromPool() (cc *CConnect, err error) {
	cnum := 0
	clang := len(c.connect_pool)
	selected := false
	failure := make([]int, 0)
	for {
		if len(c.connect_pool) == 0 {
			break
		}
		//breakfor := false
		//select {
		// case c.connect_pool[cnum].lock <- true:
		if c.connect_pool[cnum].lock == false {
			cc = c.connect_pool[cnum]
			// send the SEND_STAT_DATA_GOON
			err = cc.conn_exec.Transmission.SendStat(uint8(SEND_STAT_DATA_GOON))
			if err != nil {
				err = nil
				failure = append(failure, cnum)
				if cnum < clang-1 {
					cnum++
				} else {
					//breakfor = true
					break
				}
			}
			// get the SEND_STAT_OK
			_, err = cc.conn_exec.Transmission.GetStat()
			if err != nil {
				err = nil
				failure = append(failure, cnum)
				if cnum < clang-1 {
					cnum++
				} else {
					//breakfor = true
					break
				}
			} else {
				selected = true
				cc.lock = true
				//breakfor = true
				break
			}
			//default:
		} else {
			if cnum < clang-1 {
				cnum++
			} else {
				//breakfor = true
				break
			}
		}
		//		if breakfor == true {
		//			break
		//		}
	}
	// clean up the failure link
	//fmt.Println("failuer", failure)
	//fmt.Println("now_count", c.now_count)
	//fmt.Println("t_count", len(c.connect_pool))
	if len(failure) != 0 {
		c.cleanFailure(failure)
	}
	//fmt.Println("now_count2", c.now_count)
	//fmt.Println("t_count2", len(c.connect_pool))
	if selected == false {
		if c.now_count > c.max_count {
			err = fmt.Errorf("The number of connections is exceeded.")
			return
		}
		cc, err = c.dial()
		if err != nil {
			return
		}
		//cc.lock <- true
		cc.lock = true
		c.connect_pool = append(c.connect_pool, cc)
		// send the SEND_STAT_DATA_GOON
		err = cc.conn_exec.Transmission.SendStat(uint8(SEND_STAT_DATA_GOON))
		if err != nil {
			return
		}
		// get the SEND_STAT_OK
		_, err = cc.conn_exec.Transmission.GetStat()
		if err != nil {
			return
		}
		c.now_count++
	}
	//fmt.Println("now_count3", c.now_count)
	//fmt.Println("t_count3", len(c.connect_pool))
	return
}

func (c *Client) cleanFailure(failure []int) {
	newpool := make([]*CConnect, 0)
	for i, one := range c.connect_pool {
		todel := false
		for _, j := range failure {
			if i == j {
				todel = true
			}
		}
		if todel == false {
			newpool = append(newpool, one)
		}
	}
	c.connect_pool = newpool
	c.now_count = len(c.connect_pool)
}

func (c *Client) dial() (cc *CConnect, err error) {
	var connecter *net.TCPConn
	var ipadrr *net.TCPAddr
	ipadrr, err = net.ResolveTCPAddr("tcp", c.addr)
	if err != nil {
		return
	}
	connecter, err = net.DialTCP("tcp", nil, ipadrr)
	if err != nil {
		return
	}
	var transmission *Transmission
	if c.iftls == true {
		conf := &tls.Config{
			InsecureSkipVerify: true,
		}
		transmission = NewTransmissionTLS(tls.Client(connecter, conf))
	} else {
		transmission = NewTransmission(connecter)
	}
	ce := NewConnExec(transmission)
	// send SEND_STAT_CONN_LONG or SEND_STAT_CONN_SHORT
	if c.lconnection == true {
		ce.SetLong()
		err = ce.Transmission.SendStat(uint8(SEND_STAT_CONN_LONG))
		if err != nil {
			return
		}
		// get the SEND_STAT_OK
		_, err = ce.Transmission.GetStat()
		if err != nil {
			return
		}
	} else {
		ce.SetShort()
		err = ce.Transmission.SendStat(uint8(SEND_STAT_CONN_SHORT))
		if err != nil {
			return
		}
		// get the SEND_STAT_OK
		_, err = ce.Transmission.GetStat()
		if err != nil {
			return
		}
	}

	cc = &CConnect{
		conn_exec:   ce,
		addr:        c.addr,
		iftls:       c.iftls,
		lconnection: c.lconnection,
		closechan:   c.closechan,
		closed:      false,
		lock:        false,
		//lock:        make(chan bool, 1),
	}
	return
}

// 关闭已经发起的连接
func (cc *CConnect) Close() (err error) {
	if cc.lconnection == false {
		err = cc.conn_exec.SendClose()
		cc.closed = true
		cc.closechan <- true
		cc = nil
	} else {
		//<-cc.lock
		cc.lock = false
	}
	return
}

// 发送一个数据体，并接受返回数据体
func (cc *CConnect) SendAndReturn(data []byte) (returndata []byte, err error) {
	err = cc.conn_exec.SendData(data)
	if err != nil {
		err = fmt.Errorf("nst[CConnect]SendAndReturn: %v", err)
		return
	}
	returndata, err = cc.conn_exec.GetData()
	if err != nil {
		err = fmt.Errorf("nst[CConnect]SendAndReturn: %v", err)
		return
	}
	return
}
