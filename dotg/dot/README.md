# Package dot

一个最小数据存储单位的实现

## 版本1中dot的文件数据说明

### dot数据体

文件名： (dot id的40位sha1散列)_data

数据结构： uint8的应用版本(1bit)|dot id(255bit定长)|time类型的操作时间(15bit)|uint64的操作版本(8bit)|数据体

### dot上下文关系总索引

文件名： (dot id的40位sha1散列)_context

数据结构： uint8的应用版本(1bit)|time类型的操作时间(15bit)|uint64的操作版本(8bit)|[]string的上下文关系索引

### dot单个上下文关系

文件名： (dot id的40位sha1散列)_context_(上下文id的40位sha1散列)

数据结构： uint8的应用版本(1bit)|上下文id(255bit定长)|time类型的操作时间(15bit)|uint64的操作版本(8bit)|上下文关系