// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> stephenfmqin@gmail.com
// This source code is governed by GNU LGPL v3 license

// The configure use YAML

package yconf

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/goccy/go-yaml"

	"github.com/idcsource/insight00-lib/base"
)

// 创建配置文件对象
func NewYamlConf() (yamlconf *YamlConf) {
	yamlconf = &YamlConf{
		yaml: make(map[string]interface{}),
	}
	return
}

// 从JSON文件中读取配置
func (j *YamlConf) ReadFile(fname string) (err error) {
	fname = base.LocalFile(fname)
	yamlstream, err := ioutil.ReadFile(fname) //这里返回的已经是[]byte
	if err != nil {
		return fmt.Errorf("yconf: %v", err)
	}
	yamlmap, err := j.doMapStart([]byte(yamlstream))
	if err != nil {
		return fmt.Errorf("yconf: %v", err)
	}
	j.yaml = yamlmap
	return nil
}

// 从JSON字符串中读取配置
func (j *YamlConf) ReadString(yamlstream string) (err error) {
	yamlmap, err := j.doMapStart([]byte(yamlstream))
	if err != nil {
		return fmt.Errorf("yconf: %v", err)
	}
	j.yaml = yamlmap
	return nil
}

// 从某个节点捡出配置，并返回interface{}
func (j *YamlConf) GetValue(node string) (oneNodeVal interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("yconf: %v", e)
		}
	}()

	nodes, err := j.nodeOp(node)
	if err != nil {
		err = fmt.Errorf("yconf: %v", err)
		return
	}
	nodelen := len(nodes)
	var oneNodeNode map[string]interface{}
	for i, one := range nodes {
		if i == 0 && i != (nodelen-1) {
			_, ok := j.yaml[one]
			if ok != true {
				err = fmt.Errorf("yconf: There is not have key \"%v\"", one)
				return
			}
			oneNodeNode = j.yaml[one].(map[string]interface{})
		} else if i == 0 && i == (nodelen-1) {
			_, ok := j.yaml[one]
			if ok != true {
				err = fmt.Errorf("yconf: There is not have key \"%v\"", one)
				return
			}
			oneNodeVal = j.yaml[one]
		} else if i == (nodelen - 1) {
			_, ok := oneNodeNode[one]
			if ok != true {
				err = fmt.Errorf("yconf: There is not have key \"%v\"", one)
				return
			}
			oneNodeVal = oneNodeNode[one]
		} else {
			_, ok := oneNodeNode[one]
			if ok != true {
				err = fmt.Errorf("yconf: There is not have key \"%v\"", one)
				return
			}
			var ok2 bool
			oneNodeNode, ok2 = oneNodeNode[one].(map[string]interface{})
			if ok2 != true {
				err = fmt.Errorf("yconf: The \"%v\" is a value key not a node key", one)
				return
			}
		}
	}
	return
}

// 从某个节点捡出配置，并返回string
func (j *YamlConf) GetString(node string) (str string, err error) {
	val, err := j.GetValue(node)
	if err != nil {
		return
	}
	var ok bool
	str, ok = val.(string)
	if ok != true {
		err = fmt.Errorf("yconf: The \"%v\" value not a string", node)
		return
	}
	return
}

// 从某个节点捡出配置，并返回int64
func (j *YamlConf) GetInt64(node string) (i64 int64, err error) {
	val, err := j.GetValue(node)
	if err != nil {
		return
	}
	i32, ok := val.(int)
	if ok != true {
		err = fmt.Errorf("yconf: The \"%v\" value not a int", node)
		return
	}
	i64 = int64(i32)
	return
}

// 从某个节点捡出配置，并返回float64
func (j *YamlConf) GetFloat64(node string) (f64 float64, err error) {
	val, err := j.GetValue(node)
	if err != nil {
		return
	}
	var ok bool
	f64, ok = val.(float64)
	if ok != true {
		err = fmt.Errorf("yconf: The \"%v\" value not a float", node)
		return
	}
	return
}

// 从某个节点捡出配置，并返回bool
func (j *YamlConf) GetBool(node string) (b bool, err error) {
	val, err := j.GetValue(node)
	if err != nil {
		return
	}
	var ok bool
	b, ok = val.(bool)
	if ok != true {
		err = fmt.Errorf("yconf: The \"%v\" value not a bool", node)
		return
	}
	return
}

// 从某个节点捡出配置，并返回枚举，枚举的样式为“a,b,c,d,e”的用逗号分割的字符串，返回为[]string
func (j *YamlConf) GetEnum(node string) (em []string, err error) {
	val, err := j.GetValue(node)
	if err != nil {
		return
	}
	b, ok := val.(string)
	if ok != true {
		err = fmt.Errorf("yconf: The \"%v\" value not a emum", node)
		return
	}
	em = strings.Split(b, ",")
	for i, one := range em {
		em[i] = strings.TrimSpace(one)
	}
	return
}

// 从某个节点捡出配置，并返回数组，返回为[]interface{}
func (j *YamlConf) GetArray(node string) (ar []interface{}, err error) {
	val, err := j.GetValue(node)
	if err != nil {
		return
	}
	var ok bool
	ar, ok = val.([]interface{})
	if ok != true {
		err = fmt.Errorf("yconf: The \"%v\" value not a array", node)
		return
	}
	return
}

// 从某个节点捡出配置，并返回你想要的数据结构，使用方法如下：
//
//	type Blocks struct {
//		Name  string `yaml:"name"`
//		Path  string `yaml:"path"`
//		Token string `yaml:"token"`
//	}
//
//	var blocc []Blocks
//	err = conf.GetStruct("blocks", &blocc)
func (j *YamlConf) GetStruct(node string, v interface{}) (err error) {
	val, err := j.GetValue(node)
	if err != nil {
		return
	}

	data, err := yaml.Marshal(val)
	if err != nil {
		err = fmt.Errorf("yconf: %v", err)
		return
	}

	err = yaml.Unmarshal(data, v)
	return
}

// 从某个节点下捡出节点，成为新的YamlConf
func (j *YamlConf) GetNode(node string) (newjconf *YamlConf, err error) {
	val, err := j.GetValue(node)
	if err != nil {
		return
	}
	b, ok := val.(map[string]interface{})
	if ok != true {
		err = fmt.Errorf("yconf: The \"%v\" not a node", node)
		return
	}
	newjconf = &YamlConf{
		yaml: b,
	}
	return
}

// 修改某个节点的值，可接受单一的值，也可以接受一个map[string]interface{}
func (j *YamlConf) SetValue(node string, value interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("yconf: %v", e)
		}
	}()

	nodes, err := j.nodeOp(node)
	if err != nil {
		err = fmt.Errorf("yconf: %v", err)
		return
	}
	nodelen := len(nodes)
	var oneNodeNode map[string]interface{}
	for i, one := range nodes {
		if i == 0 && i != (nodelen-1) {
			_, ok := j.yaml[one]
			if ok != true {
				err = fmt.Errorf("yconf: There is not have key \"%v\"", one)
				return
			}
			oneNodeNode = j.yaml[one].(map[string]interface{})
		} else if i == 0 && i == (nodelen-1) {
			_, ok := j.yaml[one]
			if ok != true {
				err = fmt.Errorf("yconf: There is not have key \"%v\"", one)
				return
			}
			var strb []byte
			strb, err = yaml.Marshal(value)
			if err == nil {
				onemap := make(map[string]interface{})
				err = j.doMap(strb, onemap)
				if err != nil {
					onearray := make([]interface{}, 0)
					err = yaml.Unmarshal(strb, &onearray)
					if err != nil {
						j.yaml[one] = value
						err = nil
					} else {
						j.yaml[one] = onearray
					}
				} else {
					j.yaml[one] = onemap
				}
			} else {
				j.yaml[one] = value
				err = nil
			}
		} else if i == (nodelen - 1) {
			_, ok := oneNodeNode[one]
			if ok != true {
				err = fmt.Errorf("yconf: There is not have key \"%v\"", one)
				return
			}
			var strb []byte
			strb, err = yaml.Marshal(value)
			if err == nil {
				onemap := make(map[string]interface{})
				err = j.doMap(strb, onemap)
				if err != nil {
					onearray := make([]interface{}, 0)
					err = yaml.Unmarshal(strb, &onearray)
					if err != nil {
						oneNodeNode[one] = value
						err = nil
					} else {
						oneNodeNode[one] = onearray
					}
				} else {
					oneNodeNode[one] = onemap
				}
			} else {
				oneNodeNode[one] = value
				err = nil
			}

		} else {
			_, ok := oneNodeNode[one]
			if ok != true {
				err = fmt.Errorf("yconf: There is not have key \"%v\"", one)
				return
			}
			var ok2 bool
			oneNodeNode, ok2 = oneNodeNode[one].(map[string]interface{})
			if ok2 != true {
				err = fmt.Errorf("yconf: The \"%v\" is a value key not a node key", one)
				return
			}
		}
	}
	return
}

// 在根节点下添加一个节点，可接受单一的值，也可以接受一个map[string]interface{}
func (j *YamlConf) AddValueInRoot(name string, value interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("yconf: %v", e)
		}
	}()

	if _, ok := j.yaml[name]; ok == true {
		err = fmt.Errorf("yconf: The \"%v\" node name is already have", name)
		return
	}
	var strb []byte
	strb, err = yaml.Marshal(value)
	if err == nil {
		onemap := make(map[string]interface{})
		err = j.doMap(strb, onemap)
		if err != nil {
			onearray := make([]interface{}, 0)
			err = yaml.Unmarshal(strb, &onearray)
			if err != nil {
				j.yaml[name] = value
				err = nil
			} else {
				j.yaml[name] = onearray
			}
		} else {
			j.yaml[name] = onemap
		}
	} else {
		j.yaml[name] = value
		err = nil
	}
	return
}

// 添加一个节点，可接受单一的值，也可以接受一个map[string]interface{}
func (j *YamlConf) AddValue(node string, name string, value interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("yconf: %v", e)
		}
	}()

	nodes, err := j.nodeOp(node)
	if err != nil {
		err = fmt.Errorf("yconf: %v", err)
		return
	}
	nodelen := len(nodes)
	var oneNodeNode map[string]interface{}
	for i, one := range nodes {
		if i == 0 && i != (nodelen-1) {
			_, ok := j.yaml[one]
			if ok != true {
				err = fmt.Errorf("yconf: There is not have key \"%v\"", one)
				return
			}
			oneNodeNode = j.yaml[one].(map[string]interface{})
		} else if i == 0 && i == (nodelen-1) {
			// 刚开始就到头
			_, ok := j.yaml[one]
			if ok != true {
				err = fmt.Errorf("yconf: There is not have key \"%v\"", one)
				return
			}
			var ok2 bool
			oneNodeNode, ok2 = j.yaml[one].(map[string]interface{})
			if ok2 != true {
				err = fmt.Errorf("yconf: The \"%v\" is a value key not a node key", one)
				return
			}
			if _, ok3 := oneNodeNode[name]; ok3 == true {
				err = fmt.Errorf("yconf: The \"%v\" node name is already have", name)
				return
			}
			var strb []byte
			strb, err = yaml.Marshal(value)
			if err == nil {
				onemap := make(map[string]interface{})
				err = j.doMap(strb, onemap)
				if err != nil {
					onearray := make([]interface{}, 0)
					err = yaml.Unmarshal(strb, &onearray)
					if err != nil {
						oneNodeNode[name] = value
						err = nil
					} else {
						oneNodeNode[name] = onearray
					}
				} else {
					oneNodeNode[name] = onemap
				}
			} else {
				oneNodeNode[name] = value
				err = nil
			}
		} else if i == (nodelen - 1) {
			// 到头了
			_, ok := oneNodeNode[one]
			if ok != true {
				err = fmt.Errorf("yconf: There is not have key \"%v\"", one)
				return
			}
			var ok2 bool
			oneNodeNode, ok2 = oneNodeNode[one].(map[string]interface{})
			if ok2 != true {
				err = fmt.Errorf("yconf: The \"%v\" is a value key not a node key", one)
				return
			}
			if _, ok3 := oneNodeNode[name]; ok3 == true {
				err = fmt.Errorf("yconf: The \"%v\" node name is already have", name)
				return
			}
			var strb []byte
			strb, err = yaml.Marshal(value)
			if err == nil {
				onemap := make(map[string]interface{})
				err = j.doMap(strb, onemap)
				if err != nil {
					onearray := make([]interface{}, 0)
					err = yaml.Unmarshal(strb, &onearray)
					if err != nil {
						oneNodeNode[one] = value
						err = nil
					} else {
						oneNodeNode[one] = onearray
					}
				} else {
					oneNodeNode[one] = onemap
				}
			} else {
				oneNodeNode[one] = value
				err = nil
			}
		} else {
			// 开始了但没到头
			_, ok := oneNodeNode[one]
			if ok != true {
				err = fmt.Errorf("yconf: There is not have key \"%v\"", one)
				return
			}
			var ok2 bool
			oneNodeNode, ok2 = oneNodeNode[one].(map[string]interface{})
			if ok2 != true {
				err = fmt.Errorf("yconf: The \"%v\" is a value key not a node key", one)
				return
			}
		}
	}
	return
}

// 从某个节点里删除某个名字的值
func (j *YamlConf) DelValue(node string, name string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("yconf: %v", e)
		}
	}()

	nodes, err := j.nodeOp(node)
	if err != nil {
		err = fmt.Errorf("yconf: %v", err)
		return
	}
	nodelen := len(nodes)
	var oneNodeNode map[string]interface{}
	for i, one := range nodes {
		if i == 0 && i != (nodelen-1) {
			_, ok := j.yaml[one]
			if ok != true {
				err = fmt.Errorf("yconf: There is not have key \"%v\"", one)
				return
			}
			oneNodeNode = j.yaml[one].(map[string]interface{})
		} else if i == 0 && i == (nodelen-1) {
			// 刚开始就到头
			_, ok := j.yaml[one]
			if ok != true {
				err = fmt.Errorf("yconf: There is not have key \"%v\"", one)
				return
			}
			var ok2 bool
			oneNodeNode, ok2 = j.yaml[one].(map[string]interface{})
			if ok2 != true {
				err = fmt.Errorf("yconf: The \"%v\" is a value key not a node key", one)
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
				err = fmt.Errorf("yconf: There is not have key \"%v\"", one)
				return
			}
			var ok2 bool
			oneNodeNode, ok2 = oneNodeNode[one].(map[string]interface{})
			if ok2 != true {
				err = fmt.Errorf("yconf: The \"%v\" is a value key not a node key", one)
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
				err = fmt.Errorf("yconf: There is not have key \"%v\"", one)
				return
			}
			var ok2 bool
			oneNodeNode, ok2 = oneNodeNode[one].(map[string]interface{})
			if ok2 != true {
				err = fmt.Errorf("yconf: The \"%v\" is a value key not a node key", one)
				return
			}
		}
	}
	return
}

// 在根节点下删除一个节点
func (j *YamlConf) DelValueInRoot(name string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("yconf: %v", e)
		}
	}()

	if _, ok := j.yaml[name]; ok == true {
		delete(j.yaml, name)
		return
	}
	return
}

// 输入成为YAML，TODO
// func (j *YamlConf) OutputYaml() (str string, err error) {
// 	// TODO
// 	return
// }

// 处理节点的标记，形如“abcd>dbce>dddd>dee”
func (j *YamlConf) nodeOp(node string) (nodes []string, err error) {
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
func (j *YamlConf) doMap(stream []byte, mapResult map[string]interface{}) (err error) {
	if err := yaml.Unmarshal(stream, &mapResult); err != nil {
		return fmt.Errorf("The JSON format is wrong: %v", err)
	}
	return nil
}

func (j *YamlConf) doMapStart(stream []byte) (mapResult map[string]interface{}, err error) {
	if err = yaml.Unmarshal(stream, &mapResult); err != nil {
		err = fmt.Errorf("The JSON format is wrong: %v", err)
		return
	}
	return
}

func (j *YamlConf) Println() {
	fmt.Println(j.yaml)
}

func (j *YamlConf) MarshalBinary() (data []byte, err error) {
	data, err = yaml.Marshal(j.yaml)
	if err != nil {
		err = fmt.Errorf("yconf: %v", err)
		return
	}
	return
}

func (j *YamlConf) UnmarshalBinary(data []byte) (err error) {
	j.yaml = make(map[string]interface{})
	err = j.doMap(data, j.yaml)
	if err != nil {
		return fmt.Errorf("yconf: %v", err)
	}
	return nil
}
