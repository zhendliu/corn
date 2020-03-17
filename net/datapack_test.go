package net

import (
	"fmt"
	"io"
	"net"
	"testing"
)

//只是负责测试 拆包，和封包的单元测试
func TestDataPack(t *testing.T) {
	/*
		模拟一个服务器
	*/

	listenner, err := net.Listen("tcp", "0.0.0.0:8000")

	if err != nil {
		fmt.Println("server  lesten err :", err)
		return
	}

	//创建一个go 承载业务

	go func() {
		for {
			conn, err := listenner.Accept()
			if err != nil {
				fmt.Println("server accept error:", err)
				continue
			}
			go func(conn net.Conn) {
				//处理客户端的请求
				//定义一个拆包的对象
				dp := NewDataPack()
				for {
					//拆包的过程
					headData := make([]byte, dp.GetHeadLen())
					//读取head
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("readFull error:", err)
						break
					}
					msgHead, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println("dp.UnPack error:", err)
						break
					}
					if msgHead.GetMsgLen() > 0 {
						//msg中存在数据
						//根据head中的msgLen读取data

						msg := msgHead.(*Message)

						msg.Data = make([]byte, msgHead.GetMsgLen())

						if _, err := io.ReadFull(conn, msg.Data); err != nil {

							fmt.Println("ReadFull error :", err)
							break
						}
						//完整的一个消息已经读完
						fmt.Println("----> Recv data :", string(msg.Data), "dataLen :", msg.DataLen,"msgId:",msg.Id)
					}
				}
			}(conn)
		}
	}()

	//模拟客户端

	conn ,err  :=net.Dial("tcp","0.0.0.0:8000")

	if err  !=nil{
		fmt.Println("client dial err:",err)
		return
	}
	//创建一个封包对象

	dp :=NewDataPack()
	//封装第一个msg
	msg1 :=&Message{
		Id:1,
		DataLen:4,
		Data:[]byte{'c','o','r','n'},
	}
	sendData1,err :=dp.Pack(msg1)
	if err !=nil{
		fmt.Println("pack error:",err)
		return
	}
	//封装第二个msg
	msg2 :=&Message{
		Id:2,
		DataLen:7,
		Data:[]byte{'n','i','h','a','o','!','!'},
	}
	sendData2,err :=dp.Pack(msg2)
	if err !=nil{
		fmt.Println("pack error:",err)
		return
	}
	sendData1 =append(sendData1,sendData2...)
	//一次性发送
	conn.Write(sendData1)

	//阻塞
	select{}
}
