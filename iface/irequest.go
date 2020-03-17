package iface

/*
IRequest
把客户端的请求的连接信息和请求数据包装到了request中
*/


type IRequest interface {

	//得到当前连接
	GetConnection() IConnection
	//得到请求的消息数据
	GetData() []byte

}