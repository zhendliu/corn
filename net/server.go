package net

import (
	"corn/iface"
	"corn/utils"
	"fmt"
	"net"
)

//IServer的接口实现
type Server struct {
	//属性
	//服务名称
	Name string
	//绑定ip版本
	IPVersion string
	//绑定IP
	IP string
	//绑定端口
	Port int
	//当前server 的消息管理模块，用来绑定MsgId和对应的处理业务的API关系
	MsgHandle iface.IMsgHandle
}

//启动
func (s *Server) Start() {
	fmt.Printf("[Corn] Server Name: %s,listenner at IP:%s ,Port:%d, is starting \n", s.Name, s.IP, s.Port)

	//1 获取一个tcp的addr

	go func() {

		//0 开启消息队列及worker工作池

		s.MsgHandle.StartWorkerPool()

		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))

		if err != nil {
			fmt.Println("get  addr error :", err)
			return
		}
		//监听服务器地址
		listen, err := net.ListenTCP(s.IPVersion, addr)

		if err != nil {
			fmt.Println("listen  error :", err)
			return
		}

		fmt.Println("Success start server")
		var cid uint32 = 0
		for {
			conn, err := listen.AcceptTCP()
			if err != nil {
				fmt.Printf("Accept err:%s", err)
				continue
			}
			//已经建立了连接，测试回写echo
			//将处理新的连接业务方法 和 conn 进行绑定，得到我们的连接模块对象
			dealConn := NewConnection(conn, cid, s.MsgHandle)
			cid++
			//启动当前的连接业务处理
			go dealConn.Start()
		}

	}()

	//阻塞的等待客户端连接
}

//停止
func (s *Server) Stop() {
	//TODO  停止连接框架

}

func (s *Server) Serve() {
	s.Start()

	//阻塞状态
	select {}
}

//添加一个路由方法
func (s *Server) AddRouter(msgID uint32,router iface.IRouter) {

	s.MsgHandle.AddRouter(msgID,router)
	fmt.Println("Add Router Success!!")
}

//初始化server的模块

func NewServer(name string) iface.IServer {
	s := Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		MsgHandle: NewMsgHandle(),
	}

	return &s
}
