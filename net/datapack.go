package net

import (
	"bytes"
	"corn/iface"
	"corn/utils"
	"encoding/binary"
	"errors"
)

/*
	封包和拆包的方法实现
*/

type DataPack struct{}

//拆包封包实例初始化的方法

func NewDataPack() *DataPack {
	return &DataPack{}
}

//获取包头的长度方法
func (d *DataPack) GetHeadLen() uint32 {
	//Datalen uint32(4字节数据长度)+ ID（uint32协议号）
	return 8
}

//封包方法
func (d *DataPack) Pack(msg iface.IMessage) ([]byte, error) {
	//创建一个存放bytes 字节缓冲
	dataBuff := bytes.NewBuffer([]byte{})
	//将dataLen 写进DataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}

	//将MsgId 写进DataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	//将data数据 写进DataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

//拆包方法|----4字节数据长度---|---4字节协议号----|----DATA---|
func (d *DataPack) UnPack(binaryData []byte) (iface.IMessage, error) {
	dataBuff := bytes.NewReader(binaryData)
	//解压head信息，得到dataLen 和MsgId 信息
	msg := &Message{}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	//读取MsgId
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	//判断datalen是否已经超出了我们允许的最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && utils.GlobalObject.MaxPackageSize < msg.DataLen {
		return  nil,errors.New("too large msg data recv")
	}


	return msg,nil
}
