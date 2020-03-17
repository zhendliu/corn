package iface


/*
	封包和拆包的模块
*/
type IDatapack interface {
	//获取包头的长度方法
	GetHeadLen() uint32
	//封包方法
	Pack(IMessage)([]byte,error)
	//拆包方法
	UnPack([]byte)(IMessage,error)
}