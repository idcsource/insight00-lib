// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ]
// Stephen Fire Meditation Qin [ 火志溟 ] -> stephenfmqin@gmail.com
// This source code is governed by GNU LGPL v3 license

package logs

// 接口实现
type Logser interface {
	WriteLog(log string) (err error)      // 写一条log
	ReadLogs() (logs []string, err error) // 读全部log
	ReadLast() (l string, err error)      // 读最后一条
	Close()                               //关闭
}
