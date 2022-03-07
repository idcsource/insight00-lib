// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

package nst

type SendStat uint8

const (
	SEND_STAT_NO          SendStat = iota // the stat have nothing
	SEND_STAT_OK                          // it's all ok
	SEND_STAT_NOT_OK                      // it's not ok
	SEND_STAT_CONN_LONG                   // this is long connect
	SEND_STAT_CONN_SHORT                  // this is shot connect
	SEND_STAT_CHECK_DATA                  // check the server or connect
	SEND_STAT_NORMAL_DATA                 // normal data
	SEND_STAT_DATA_GOON                   // data goon transmission
	SEND_STAT_DATA_CLOSE                  // data transmission close
	SEND_STAT_CONN_CLOSE                  // connect close

)

// Server's connect execution interface
type ConnExecer interface {
	NSTexec(ce *ConnExec) (SendStat, error)
}
