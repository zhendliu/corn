package net

type Message struct {
	Id      uint32 //消息ID
	DataLen uint32 //消息长度
	Data    []byte //消息的内容
}

//获取消息ID
func (m *Message)GetMsgId() uint32 {
	return m.Id
}

//获取消息长度
func (m *Message)GetMsgLen() uint32 {
	return m.DataLen
}

//获取消息的内容
func (m *Message)GetData() []byte {
	return m.Data
}

//设置消息的ID
func (m *Message)SetMsgId(id uint32) {
	m.Id = id
}

//设置消息的内容
func (m *Message)GetMsgData(data []byte) {
	m.Data = data
}

//设置消息的长度
func (m *Message)setDataLen(len uint32) {
	m.DataLen = len
}
