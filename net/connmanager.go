package net

import (
	"corn/iface"
	"errors"
	"fmt"
	"sync"
)

/*
	当前的连接管理模块
*/


type ConnManager struct{
	connections map[uint32] iface.IConnection
	connLook sync.RWMutex
}

func NewConnManager()*ConnManager{

	return  &ConnManager{
		connections:make(map[uint32] iface.IConnection ),
	}
}


//添加连接
func (cm *ConnManager)Add(conn iface.IConnection){

	cm.connLook.Lock()
	defer cm.connLook.Unlock()
	//添加
	cm.connections[conn.GetConnID()] = conn
	fmt.Println("conn add to ConnManager success,connID:",conn.GetConnID(),"sum:",cm.Len())
}
//删除连接
func (cm *ConnManager)Remove(conn iface.IConnection){
	cm.connLook.Lock()
	defer cm.connLook.Unlock()
	//删除
	delete(cm.connections,conn.GetConnID())
	fmt.Println("conn remove   ConnManager success,connID:",conn.GetConnID(),"sum:",cm.Len())
}
//根据connID获取链接
func (cm *ConnManager)Get(connID uint32)(iface.IConnection,error){
	cm.connLook.RLock()
	defer cm.connLook.RUnlock()
	if conn ,ok :=cm.connections[connID];ok {
		return  conn,nil
	}
	return  nil,errors.New("not found conn")
}
//得到当前连接总数
func (cm *ConnManager)Len()int{
	return  len(cm.connections)
}
//清除并终止所有连接
func (cm *ConnManager)ClearConn(){
	cm.connLook.Lock()
	defer cm.connLook.Unlock()
	//删除conn，并停止conn的工作

	for connID,conn:= range cm.connections{
		//停止
		conn.Stop()

		//删除
		delete(cm.connections,connID)
	}
	fmt.Println("Clear All Connections Success,sum:",cm.Len())
}