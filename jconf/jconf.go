// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// This source code is governed by GNU LGPL v3 license

// The configure use JSON
package jconf

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/idcsource/insight00-lib/base"
)

func NewJsonConf() (jsonconf *JsonConf) {
	jsonconf = &JsonConf{
		json: make(map[string]interface{}),
	}
	return
}

// 从文件中读取配置
func (j *JsonConf) ReadFile(fname string) (err error) {
	fname = base.LocalFile(fname)
	jsonstream, err := ioutil.ReadFile(fname) //这里返回的已经是[]byte
	if err != nil {
		return fmt.Errorf("jconf: %v", err)
	}
	err = j.doMap([]byte(jsonstream), j.json)
	if err != nil {
		return fmt.Errorf("jconf: %v", err)
	}
	return nil
}

// 从字符串中读取配置
func (j *JsonConf) ReadString(jsonstream string) (err error) {
	err = j.doMap([]byte(jsonstream), j.json)
	if err != nil {
		return fmt.Errorf("jconf: %v", err)
	}
	return nil
}

// 从某个节点捡出配置，并返回interface{}
func (j *JsonConf) GetValue(node string) (oneNodeVal interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("jconf: %v", e)
		}
	}()

	nodes, err := j.nodeOp(node)
	if err != nil {
		err = fmt.Errorf("jconf: %v", err)
		return
	}
	nodelen := len(nodes)
	var oneNodeNode map[string]interface{}
	for i, one := range nodes {
		if i == 0 && i != (nodelen-1) {
			_, ok := j.json[one]
			if ok != true {
				err = fmt.Errorf("jconf: There is not have key \"%v\"", one)
				return
			}
			oneNodeNode = j.json[one].(map[string]interface{})
		} else if i == 0 && i == (nodelen-1) {
			_, ok := j.json[one]
			if ok != true {
				err = fmt.Errorf("jconf: There is not have key \"%v\"", one)
				return
			}
			oneNodeVal = j.json[one]
		} else if i == (nodelen - 1) {
			_, ok := oneNodeNode[one]
			if ok != true {
				err = fmt.Errorf("jconf: There is not have key \"%v\"", one)
				return
			}
			oneNodeVal = oneNodeNode[one]
		} else {
			_, ok := oneNodeNode[one]
			if ok != true {
				err = fmt.Errorf("jconf: There is not have key \"%v\"", one)
				return
			}
			var ok2 bool
			oneNodeNode, ok2 = oneNodeNode[one].(map[string]interface{})
			if ok2 != true {
				err = fmt.Errorf("jconf: The \"%v\" is a value key not a node key", one)
				return
			}
		}
	}
	return
}

// 从某个节点捡出配置，并返回string
func (j *JsonConf) GetString(node string) (str string, err error) {
	val, err := j.GetValue(node)
	if err != nil {
		return
	}
	var ok bool
	str, ok = val.(string)
	if ok != true {
		err = fmt.Errorf("jconf: The \"%v\" value not a string", node)
		return
	}
	return
}

// 从某个节点捡出配置，并返回int64
func (j *JsonConf) GetInt64(node string) (i64 int64, err error) {
	val, err := j.GetValue(node)
	if err != nil {
		return
	}
	f64, ok := val.(float64)
	if ok != true {
		err = fmt.Errorf("jconf: The \"%v\" value not a int", node)
		return
	}
	i64 = int64(f64)
	return
}

// 从某个节点捡出配置，并返回float64
func (j *JsonConf) GetFloat64(node string) (f64 float64, err error) {
	val, err := j.GetValue(node)
	if err != nil {
		return
	}
	var ok bool
	f64, ok = val.(float64)
	if ok != true {
		err = fmt.Errorf("jconf: The \"%v\" value not a float", node)
		return
	}
	return
}

// 从某个节点捡出配置，并返回bool
func (j *JsonConf) GetBool(node string) (b bool, err error) {
	val, err := j.GetValue(node)
	if err != nil {
		return
	}
	var ok bool
	b, ok = val.(bool)
	if ok != true {
		err = fmt.Errorf("jconf: The \"%v\" value not a bool", node)
		return
	}
	return
}

// 从某个节点捡出配置，并返回枚举，枚举的样式为“a,b,c,d,e”的字符串，返回为[]string
func (j *JsonConf) GetEnum(node string) (em []string, err error) {
	val, err := j.GetValue(node)
	if err != nil {
		return
	}
	b, ok := val.(string)
	if ok != true {
		err = fmt.Errorf("jconf: The \"%v\" value not a emum", node)
		return
	}
	em = strings.Split(b, ",")
	for i, one := range em {
		em[i] = strings.TrimSpace(one)
	}
	return
}

// 从某个节点下捡出节点，成为新的JsonConf
func (j *JsonConf) GetNode(node string) (newjconf *JsonConf, err error) {
	val, err := j.GetValue(node)
	if err != nil {
		return
	}
	b, ok := val.(map[string]interface{})
	if ok != true {
		err = fmt.Errorf("jconf: The \"%v\" not a node", node)
		return
	}
	newjconf = &JsonConf{
		json: b,
	}
	return
}

// 修改某个节点的值，可接受单一的值，也可以接受一个map[string]interface{}
func (j *JsonConf) SetValue(node string, value interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("jconf: %v", e)
		}
	}()

	nodes, err := j.nodeOp(node)
	if err != nil {
		err = fmt.Errorf("jconf: %v", err)
		return
	}
	nodelen := len(nodes)
	var oneNodeNode map[string]interface{}
	for i, one := range nodes {
		if i == 0 && i != (nodelen-1) {
			_, ok := j.json[one]
			if ok != true {
				err = fmt.Errorf("jconf: There is not have key \"%v\"", one)
				return
			}
			oneNodeNode = j.json[one].(map[string]interface{})
		} else if i == 0 && i == (nodelen-1) {
			_, ok := j.json[one]
			if ok != true {
				err = fmt.Errorf("jconf: There is not have key \"%v\"", one)
				return
			}
			j.json[one] = value
		} else if i == (nodelen - 1) {
			_, ok := oneNodeNode[one]
			if ok != true {
				err = fmt.Errorf("jconf: There is not have key \"%v\"", one)
				return
			}
			oneNodeNode[one] = value
		} else {
			_, ok := oneNodeNode[one]
			if ok != true {
				err = fmt.Errorf("jconf: There is not have key \"%v\"", one)
				return
			}
			var ok2 bool
			oneNodeNode, ok2 = oneNodeNode[one].(map[string]interface{})
			if ok2 != true {
				err = fmt.Errorf("jconf: The \"%v\" is a value key not a node key", one)
				return
			}
		}
	}
	return
}

// 在根节点下添加一个节点，可接受单一的值，也可以接受一个map[string]interface{}
func (j *JsonConf) AddValueInRoot(name string, value interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("jconf: %v", e)
		}
	}()

	if _, ok := j.json[name]; ok == true {
		err = fmt.Errorf("jconf: The \"%v\" node name is already have", name)
		return
	}
	j.json[name] = value
	return
}

// 添加一个节点，可接受单一的值，也可以接受一个map[string]interface{}
func (j *JsonConf) AddValue(node string, name string, value interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("jconf: %v", e)
		}
	}()

	nodes, err := j.nodeOp(node)
	if err != nil {
		err = fmt.Errorf("jconf: %v", err)
		return
	}
	nodelen := len(nodes)
	var oneNodeNode map[string]interface{}
	for i, one := range nodes {
		if i == 0 && i != (nodelen-1) {
			_, ok := j.json[one]
			if ok != true {
				err = fmt.Errorf("jconf: There is not have key \"%v\"", one)
				return
			}
			oneNodeNode = j.json[one].(map[string]interface{})
		} else if i == 0 && i == (nodelen-1) {
			// 刚开始就到头
			_, ok := j.json[one]
			if ok != true {
				err = fmt.Errorf("jconf: There is not have key \"%v\"", one)
				return
			}
			var ok2 bool
			oneNodeNode, ok2 = j.json[one].(map[string]interface{})
			if ok2 != true {
				err = fmt.Errorf("jconf: The \"%v\" is a value key not a node key", one)
				return
			}
			if _, ok3 := oneNodeNode[name]; ok3 == true {
				err = fmt.Errorf("jconf: The \"%v\" node name is already have", name)
				return
			}
			oneNodeNode[name] = value
		} else if i == (nodelen - 1) {
			// 到头了
			_, ok := oneNodeNode[one]
			if ok != true {
				err = fmt.Errorf("jconf: There is not have key \"%v\"", one)
				return
			}
			var ok2 bool
			oneNodeNode, ok2 = oneNodeNode[one].(map[string]interface{})
			if ok2 != true {
				err = fmt.Errorf("jconf: The \"%v\" is a value key not a node key", one)
				return
			}
			if _, ok3 := oneNodeNode[name]; ok3 == true {
				err = fmt.Errorf("jconf: The \"%v\" node name is already have", name)
				return
			}
			oneNodeNode[name] = value
		} else {
			// 开始了但没到头
			_, ok := oneNodeNode[one]
			if ok != true {
				err = fmt.Errorf("jconf: There is not have key \"%v\"", one)
				return
			}
			var ok2 bool
			oneNodeNode, ok2 = oneNodeNode[one].(map[string]interface{})
			if ok2 != true {
				err = fmt.Errorf("jconf: The \"%v\" is a value key not a node key", one)
				return
			}
		}
	}
	return
}

// 从某个节点里删除某个名字的值
func (j *JsonConf) DelValue(node string, name string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("jconf: %v", e)
		}
	}()

	nodes, err := j.nodeOp(node)
	if err != nil {
		err = fmt.Errorf("jconf: %v", err)
		return
	}
	nodelen := len(nodes)
	var oneNodeNode map[string]interface{}
	for i, one := range nodes {
		if i == 0 && i != (nodelen-1) {
			_, ok := j.json[one]
			if ok != true {
				err = fmt.Errorf("jconf: There is not have key \"%v\"", one)
				return
			}
			oneNodeNode = j.json[one].(map[string]interface{})
		} else if i == 0 && i == (nodelen-1) {
			// 刚开始就到头
			_, ok := j.json[one]
			if ok != true {
				err = fmt.Errorf("jconf: There is not have key \"%v\"", one)
				return
			}
			var ok2 bool
			oneNodeNode, ok2 = j.json[one].(map[string]interface{})
			if ok2 != true {
				err = fmt.Errorf("jconf: The \"%v\" is a value key not a node key", one)
				return
			}
			if _, ok3 := oneNodeNode[name]; ok3 == true {
				delete(oneNodeNode, name)
				return
			}
		} else if i == (nodelen - 1) {
			// 到头了
			_, ok := oneNodeNode[one]
			if ok != true {
				err = fmt.Errorf("jconf: There is not have key \"%v\"", one)
				return
			}
			var ok2 bool
			oneNodeNode, ok2 = oneNodeNode[one].(map[string]interface{})
			if ok2 != true {
				err = fmt.Errorf("jconf: The \"%v\" is a value key not a node key", one)
				return
			}
			if _, ok3 := oneNodeNode[name]; ok3 == true {
				delete(oneNodeNode, name)
				return
			}
		} else {
			// 开始了但没到头
			_, ok := oneNodeNode[one]
			if ok != true {
				err = fmt.Errorf("jconf: There is not have key \"%v\"", one)
				return
			}
			var ok2 bool
			oneNodeNode, ok2 = oneNodeNode[one].(map[string]interface{})
			if ok2 != true {
				err = fmt.Errorf("jconf: The \"%v\" is a value key not a node key", one)
				return
			}
		}
	}
	return
}

// 在根节点下删除一个节点
func (j *JsonConf) DelValueInRoot(name string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("jconf: %v", e)
		}
	}()

	if _, ok := j.json[name]; ok == true {
		delete(j.json, name)
		return
	}
	return
}

// 输入成为JSON
func (j *JsonConf) OutputJson() (str string, err error) {
	strb, err := json.Marshal(j.json)
	if err != nil {
		err = fmt.Errorf("jconf: %v", err)
		return
	}
	str = string(strb)
	return
}

// 处理节点的标记，形如“abcd>dbce>dddd>dee”
func (j *JsonConf) nodeOp(node string) (nodes []string, err error) {
	nodes = strings.Split(node, ">")
	for i, one := range nodes {
		one = strings.TrimSpace(one)
		if len(one) == 0 {
			err = errors.New("The Node format is wrong.")
			return
		}
		nodes[i] = one
	}
	return
}

// 转换JSON到map
func (j *JsonConf) doMap(stream []byte, mapResult map[string]interface{}) (err error) {
	if err := json.Unmarshal(stream, &mapResult); err != nil {
		return fmt.Errorf("The JSON format is wrong: %v", err)
	}
	return nil
}

func (j *JsonConf) Println() {
	fmt.Println(j.json)
}

func (j *JsonConf) MarshalBinary() (data []byte, err error) {
	data, err = json.Marshal(j.json)
	if err != nil {
		err = fmt.Errorf("jconf: %v", err)
		return
	}
	return
}

func (j *JsonConf) UnmarshalBinary(data []byte) (err error) {
	j.json = make(map[string]interface{})
	err = j.doMap(data, j.json)
	if err != nil {
		return fmt.Errorf("jconf: %v", err)
	}
	return nil
}
