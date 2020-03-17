package net

import (
	"corn/iface"
	"corn/utils"
	"fmt"
	"net"
)

/*
连接模块
*/
type Connection struct {
	//当前连接的所持可人

	Conn *net.TCPConn

	//连接的ID
	ConnID uint32

	//当前的连接状态
	IsClosed bool



	//告知当前连接已经退出的 channel
	ExitChan chan bool

	//该链接处理的方法
	Router iface.IRouter
}

//初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, router iface.IRouter) *Connection {

	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		Router: router,
		IsClosed:  false,
		ExitChan:  make(chan bool, 1),
	}
	return c
}

//连接的读取业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader G0routine is running ...")
	defer fmt.Println("connID=", c.ConnID, "Reader  is exist,remote addr is : ", c.GetRemoteAddr().String())
	defer c.Stop()

	for {
		//读取客户端的数据到buf中

		buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("conn  read error:", err)
			continue
		}

		//得到当前连接的request对象数据

		req := Request{
			conn: c,
			data: buf,
		}

		//执行注册的路由方法
		go func(request iface.IRequest) {
			//调用路由，注册绑定的connection
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)

	}

}

//启动连接
func (c *Connection) Start() {
	fmt.Println("Conn Start()...ConnID=", c.ConnID)
	//启动从当前连接的读业务
	go c.StartReader()
	//TODO 启动从当前写数据的业务

}

//停止连接
func (c *Connection) Stop() {
	fmt.Println("Conn stop().. ConnID:", c.ConnID)

	//判断当前连接是否关闭
	if c.IsClosed == true {
		return
	}

	c.IsClosed = true
	//尝试关闭

	c.Conn.Close()
	//回收资源
	close(c.ExitChan)
}

//获取当前连接所绑定的connection
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取当前连接模块所绑定的连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端的TCP状态 IP Port
func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//发送数据
func (c *Connection) Send(data []byte) error {
	return nil
}
