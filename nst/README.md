# Package nst
一个TCP的服务和客户端实现，包含发送和接收数据的一系列方法

## 如何使用

### 服务端的使用

准备一个logs包中的Logs，作为服务端运行时可以发送错误日志的地方。

准备一个符合nst.ConnExecer接口的执行器，这个执行器将接收*nst.ConnExec所传来的客户端数据。

使用nst.NewServer()方法建立服务器，方法将返回*Server。如果需要启用TSL加密，则可以继续用*Server.ToTLS()方法配置密钥文件。最后用*Server.Start()打开监听。

可使用ConnExec.GetData()接收数据，也可以使用ConnExec.SendData()发送数据，直到服务端或客户端使用了ConnExec.SendClose()。

最终使用*Server.Close()关闭服务器。

### 客户端的使用

使用NewClient()或NewClientL()创建一个客户端，都是*Client。

使用Client.OpenProgress()或Client.OpenConnect()拿出一个可以使用的连接，也就是*CConnect。

使用*CConnect.SendAndReturn()向服务端发送数据并接收服务器端的返回数据，直到服务器端关闭连接或客户端执行*CConnect.Close()。

最后使用*Client.Close()关闭整个客户端。

