package iface

/*
	将请求的消息封装到一个Message中，定义抽象接口
*/

type IMessage interface {
	//获取消息ID
	GetMsgId() uint32
	//获取消息长度
	GetMsgLen() uint32
	//获取消息的内容
	GetData() []byte
	//设置消息的ID
	SetMsgId(uint32)
	//设置消息的内容
	SetMsgData([]byte)
	//设置消息的长度
	SetDataLen(uint32)
}
