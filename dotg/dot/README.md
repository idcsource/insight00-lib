# Package dot

一个最小数据存储单位的实现（含Dot-Block）

## Dot是什么

Dot是存储数据的最小单元，它包括一个数据存储点，以及多组上下文关系（Context）。每个Dot由多个文件组成。

## Dot-Block是什么

Dot-Block（点块）为一个路径结构，用来存储Dot。可以理解为SQL中的表或一个普通单机版数据库。

## Context是什么

上下文关系，每个Dot依靠Context与其他Dot发生联系。

每个上下文关系包含一个上游（Up）指向其他Dot的id，和多个下游（Down）指向其他Dot的id，每个下游还可以携带一个字符串信息。

每个Dot都可以携带多组上下文关系。

## 版本1中dot的文件数据说明【计划废弃，不做兼容】

### dot数据体

文件名： (dot id的40位sha1散列)_data

数据结构： uint8的应用版本(1bit) | dot id(255bit定长) | time类型的操作时间(15bit) | uint64的操作版本(8bit)|数据体

### dot上下文关系总索引

文件名： (dot id的40位sha1散列)_context

数据结构： uint8的应用版本(1bit) | time类型的操作时间(15bit) | uint64的操作版本(8bit) | []string的上下文关系索引

### dot单个上下文关系

文件名： (dot id的40位sha1散列)_context_(上下文id的40位sha1散列)

数据结构： uint8的应用版本(1bit) | 上下文id(255bit定长) | time类型的操作时间(15bit) | uint64的操作版本(8bit) | 上下文关系


## 版本2中dot的文件数据说明【TODO】

### dot数据体

文件名： (dot id的40位sha1散列)_data

数据结构： uint8的应用版本(1bit) | dot id(255bit定长) | uint64的操作版本(8bit)|数据体

### dot上下文关系总索引

文件名： (dot id的40位sha1散列)_context

数据结构： uint8的应用版本(1bit) | uint64的操作版本(8bit) | uint64的index位置编号(从0开始)(8bit) | uint8的状态位 | 上下文关系名(255bit定长) | uint64的index位置编号(从0开始)(8bit) | uint8的状态位 | 上下文关系id(255bit定长) | ……

说明：状态位见DOT_CONTENT_INDEX_*

### dot上下文关系总索引已经删除位置

文件名： (dot id的40位sha1散列)_context_del_index

数据结构： uint8的应用版本(1bit) | uint64的操作版本(8bit) | []uint64(上下文的index位置编号)

### dot单个上下文关系

文件名： (dot id的40位sha1散列)_context_(上下文关系id的40位sha1散列)

数据结构： uint8的应用版本(1bit) | uint64的操作版本(8bit) | UP关系的dot id(255bit定长) | uint8的UP关系配置数据状态 | UP关系配置数据(255bit定长) | uint64的index位置编号(从0开始)(8bit) | uint8的状态位 | DOWN关系名(255bit定长) | DOWN关系配置数据(255bit定长) | uint64的index位置编号(从0开始)(8bit) | uint8的状态位 | DOWN关系名(255bit定长) | DOWN关系配置数据(255bit定长) | ……

说明：状态位见DOT_CONTENT_UP_DOWN_INDEX_*

### dot单个上下文关系中DOWN的已经删除位置

文件名： (dot id的40位sha1散列)_context_(上下文关系id的40位sha1散列)_del_index

数据结构： uint8的应用版本(1bit) | uint64的操作版本(8bit) | []uint64(DOWN的index位置编号)

### dot单个上下文关系中UP超出255bit配置数据时

文件名： (dot id的40位sha1散列)_context_(上下文关系id的40位sha1散列)_UP

数据结构：完整数据

### dot单个上下文关系中DOWN超出255bit配置数据时

文件名： (dot id的40位sha1散列)_context_(上下文关系id的40位sha1散列)_DOWN_(DOWN的index位置编号)

数据结构：完整数据


