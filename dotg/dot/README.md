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

## 版本1中dot的文件数据说明

### dot数据体

文件名： (dot id的40位sha1散列)_data

数据结构： uint8的应用版本(1bit) | dot id(255bit定长) | time类型的操作时间(15bit) | uint64的操作版本(8bit)|数据体

### dot上下文关系总索引

文件名： (dot id的40位sha1散列)_context

数据结构： uint8的应用版本(1bit) | time类型的操作时间(15bit) | uint64的操作版本(8bit) | []string的上下文关系索引

### dot单个上下文关系

文件名： (dot id的40位sha1散列)_context_(上下文id的40位sha1散列)

数据结构： uint8的应用版本(1bit) | 上下文id(255bit定长) | time类型的操作时间(15bit) | uint64的操作版本(8bit) | 上下文关系
