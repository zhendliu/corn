package net

import (
	"corn/iface"
	"corn/utils"
	"errors"
	"fmt"
	"io"
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

	//用于读写goroutine 之间的消息通信

	msgChan chan []byte

	//消息id对应的业务处理API关系
	MsgHandle iface.IMsgHandle
}

//初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, msgHandler iface.IMsgHandle) *Connection {

	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		MsgHandle: msgHandler,
		IsClosed:  false,
		ExitChan:  make(chan bool, 1),
		msgChan:   make(chan []byte),
	}
	return c
}

//连接的读取业务方法
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running]")
	defer fmt.Println("connID=", c.ConnID, "Reader  is exist,remote addr is : ", c.GetRemoteAddr().String())
	defer c.Stop()

	for {

		//创建一个拆包解包的对象
		dp := NewDataPack()

		//读取客户端的MsgHead（8个字节 二进制流）
		//读取head
		headData := make([]byte, dp.GetHeadLen())

		_, err := io.ReadFull(c.GetTCPConnection(), headData)

		if err != nil {
			fmt.Println("read msg head error:", err)
			break
		}

		//拆包，得到msgID，msgDataLen 放在msg 消息中
		msg, err := dp.UnPack(headData)

		if err != nil {
			fmt.Println(" unpack error:", err)
			break
		}
		//根据dataLen 读取Data，放在msg.data 属性中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())

			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg  data error:", err)
				break
			}
		}
		msg.SetMsgData(data)

		//得到当前连接的request对象数据
		req := Request{
			conn: c,
			msg:  msg,
		}


		if utils.GlobalObject.WorkerPoolSize >0{
			//开启了工作池机制
			c.MsgHandle.SendMsgToTaskQueue(&req)
		}else{
			//执行注册的路由方法
			go c.MsgHandle.DoMsgHandler(&req)
		}
	}

}

/*
	写消息的goroutine，专门给用户将消息发送给客户端
*/

func (c  *Connection)StartWriter(){
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.GetRemoteAddr().String(),"Connection Writer exist!")

	for{
		select {
		case data := <-c.msgChan:
			if _,err :=c.Conn.Write(data);err !=nil{
				fmt.Println("Send data error:",err)
				return
			}
		case  <-c.ExitChan:
			//reader 关闭，指示退出
			return
		}
	}

}

//启动连接
func (c *Connection) Start() {
	fmt.Println("Conn Start()...ConnID=", c.ConnID)
	//启动从当前连接的读业务
	go c.StartReader()
	//TODO 启动从当前写数据的业务
	go c.StartWriter()

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

	//告知Writer 关闭

	c.ExitChan <- true
	//回收资源
	close(c.ExitChan)
	close(c.msgChan)
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

//提供一个sendMsg方法，将我们要发送给客户端的数据

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.IsClosed == true {
		return errors.New("connection is close when send message")
	}

	//将data进行封包
	//创建一个封包对象

	dp := NewDataPack()
	//封装第一个msg

	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("pack msg  err  :", err)
		return err
	}
	//将数据发送到客户端

	c.msgChan  <-binaryMsg

	return nil
}
