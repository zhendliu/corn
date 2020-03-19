package iface

//定义一个服务器接口

type IServer interface {
	//启动服务器
	Start()
	//停止服务器
	Stop()
	//运行服务器
	Serve()
	//添加的一个路由的方法
	//给当前的路由添加一个处理路由方法
	AddRouter(uint32, IRouter)

	GetConnMgr() IConnManager

	//注册OnConnStart 钩子函数的方法
	SetOnConnStart(func (conn IConnection))

	//注册 OnConnStop 钩子函数的方法
	SetOnConnStop(func (conn IConnection))

	//调用OnConnStart 钩子函数的方法
	CallOnConnStart(conn IConnection)

	//调用 OnConnStop 钩子函数的方法
	CallOnConnStop(conn IConnection)

}