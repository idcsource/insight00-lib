# Package nst
一个TCP的服务和客户端实现，包含发送和接收数据的一系列方法

## 如何使用

### 服务端的使用

准备一个logs包中的Logs，作为服务端运行时可以发送错误日志的地方。

准备一个符合nst.ConnExecer接口的执行器，这个执行器将接收*nst.ConnExec所传来的客户端数据。

可使用ConnExec.GetData()接收数据，也可以使用ConnExec.SendData()发送数据，直到服务端或客户端使用了ConnExec.SendClose()。

nst.ConnExecer接口的执行器返回给服务端nst2.SendStat状态和错误状态，让服务端可以获取到。

### 客户端的使用

使用NewClient()或NewClientL()创建一个客户端。

使用Client.OpenProgress()或Client.OpenConnect()拿出一个可以使用的连接。

使用CConnect.SendAndReturn()向服务端发送数据并接收服务器端的返回数据，直到服务器端关闭连接或客户端执行CConnect.Close()。

