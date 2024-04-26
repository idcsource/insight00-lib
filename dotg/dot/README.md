# Package dot

一个最小数据存储单位的实现（含Dot-Block）

当前为版本2（版本1已经废弃且不进行兼容），目前这部分的基本功能已经实现，并进行了基本的单进程下的测试。

## Dot是什么

Dot是存储数据的最小单元，它包括一个数据存储点，以及多组上下文关系（这里称为Context）。每个Dot由多个文件组成。

## Dot-Block是什么

Dot-Block（点块）为一个路径结构，用来存储Dot。可以理解为SQL中的表或一个普通单机版数据库。

Block实现了运行状态标记和对Dot级别的读写锁，确保一个Block在同一时间只能由一个进程启动，且同一个dot只能同时由一个进程（或协程）进行更改。因此必须使用StartBlock()启动并使用StopBlock()关闭。

目前已经实现了内部读写锁和显式声明的外部锁，但没有对可能的事务下的情况进行单独的锁设计。

## Context是什么

上下文关系，每个Dot依靠Context与其他Dot发生联系。

每个上下文关系包含一个上游（Up）指向其他dot的id，和多个下游（Down）指向其他dot的id，每个上游或下游还可以携带一个配置数据体。

每个Dot都可以携带多组上下文关系。

但在Context关系中，并不验证必须记录的是已经存在的dot，也就是说，在Context关系中，dotid实际可以放任何你想放的字符串内容，只要长度允许。

## 版本2中dot的文件数据结构说明

### dot数据体

文件名： (dotid的40位sha1散列)_body

数据结构： uint8的应用版本(1bit) | dotid(200bit定长) | uint64的操作版本(8bit)|数据体

### 一个dot所有Context的总索引

文件名： (dotid的40位sha1散列)_context_index

数据结构： uint8的应用版本(1bit) | uint64的操作版本(8bit) | uint8的状态位 | 上下文关系名(200bit定长) | uint8的状态位 | 上下文关系id(255bit定长) | ……

说明：状态位见类型 _DotContextIndex_Status

### 一个dot所有Context的总索引中对删除的记录

文件名： (dotid的40位sha1散列)_context_del_index

数据结构： uint8的应用版本(1bit) | uint64的操作版本(8bit) | []uint64(上下文的index位置编号)

### 一个dot的单个Context

文件名： (dotid的40位sha1散列)_context_(上下文关系id的40位sha1散列)

数据结构： 

uint8的应用版本(1bit) | uint64的操作版本(8bit)| 上下文关系名(200bit定长) | 

UP关系的dotid(200bit定长) | uint8的UP状态位 | uint64的UP关系配置数据实际长度(8bit) | UP关系配置数据(500bit定长) | 

uint8的down状态位 | DOWN关系名(200bit定长) | uint64的DOWN关系配置数据实际长度(8bit) | DOWN关系配置数据(500bit定长) | 

uint8的down状态位 | DOWN关系名(200bit定长) | uint64的DOWN关系配置数据实际长度(8bit) | DOWN关系配置数据(500bit定长) | 

……

说明：状态位见类型 _DotContextUpDownIndex_Status，目前没有实现DOT_CONTEXT_UP_DOWN_INDEX_OUTDATA_NODEL

### dot的单个Context中DOWN的已经删除位置

文件名： (dotid的40位sha1散列)_context_(上下文关系id的40位sha1散列)_del_index

数据结构： uint8的应用版本(1bit) | uint64的操作版本(8bit) | []uint64(DOWN的index位置编号)

### dot的单个Context中UP超出500bit配置数据时

文件名： (dotid的40位sha1散列)_context_(上下文关系id的40位sha1散列)_UP

数据结构：完整数据

### dot的单个Context中DOWN超出500bit配置数据时

文件名： (dotid的40位sha1散列)_context_(上下文关系id的40位sha1散列)_DOWN_(DOWN的关系名的40位sha1散列)

数据结构：完整数据

### 备注

dotid的长度显示根据这个常数设定：DOT_ID_MAX_LENGTH_V2

Context的数据大小限制根据这个常数设定：DOT_CONTENT_MAX_IN_DATA_V2

## 应该还会去做的（TODO）

但现在还没有这个打算，先缓缓再说，看心情。

### dot单个上下文关系文件的分页处理

如果这个文件的大小超过一个数值后（比如超过1GB），则会新建文件继续存放更多的数据。同时在进行检索操作的时候，可以多进程（线程/携程，随便怎么叫）并行处理。

### 对已删除位置的资源回收

除了直接删除整个dot外，无论是删除Context还是删除Context中的DOWN关系，在当前的删除操作中，并不会彻底删除所有相关数据，而是标记为删除。这种情况出现在上面的“一个dot所有Context的总索引”和“一个dot的单个Context”这两个文件中，也正因此还为这两个文件配备了专门的删除索引文件。

目前，当新增Context的DOWN关系时，是可以复用被删除的位置，但新增Context的时候，并没有这样的设计。但无论如何，随着使用，不可避免会有一些被标记为删除的空位占用磁盘空间。所以需要一个进行资源回收的机制，可能会类似PostgreSQL的“VACUUM FULL”，也就是完全重建文件和索引。

### 对单一dot的分区存放

如果一个dot拥有特别多的Context，且其关系中的数据体超过DOT_CONTENT_MAX_IN_DATA_V2的也比较多，会造成这个dot的物理文件非常多，并集中在一个路径下，可能超过物理文件系统的限制（文件数或总容量）。当前版本对这一情况没有任何处理办法。未来可能会增加支持dot分区的概念，可以将属于同一个dot的文件分别存放在多个物理存储位置上。