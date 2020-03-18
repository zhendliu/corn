package iface

/*
	连接管理抽象层模块
*/

type IConnManager interface {
	//添加连接
	Add(conn IConnection)
	//删除连接
	Remove(conn IConnection)
	//根据connID获取链接
	Get(connID uint32)(IConnection,error)
	//得到当前连接总数
	Len()int
	//清除并终止所有连接
	ClearConn()

}
