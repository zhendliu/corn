package utils

import (
	"corn/iface"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

/*
定义corn全局框架的参数，供框架使用
一些参数通过corn.json 文件读取
*/
type GlobalObj struct {

	/*
		Server
	*/

	TcpServer iface.IServer

	Host string //主机监听的IP

	TcpPort int //主机监听的端口

	Name string //当前服务器的名称

	/*
		Corn框架的配置
	*/

	Version string //当前框架的版本号

	MaxConn int //框架的最大连接数

	MaxPackageSize uint32 //单词读包最大字节数

	WorkerPoolSize uint32  //当前业务工作池的Goroutine数量
	MaxWorkerTaskLen uint32  //允许最多开辟的worker

}

/*
	定义一个全局的对外Globalobj
*/

var GlobalObject *GlobalObj

/*
	从corn.json加载自定义的参数，应用到服务中
*/
func (g *GlobalObj) Reload(){
	data ,err :=ioutil.ReadFile("G:\\GO_WORKSPACE\\src\\corn\\conf\\corn.json")
	if err !=nil{
		fmt.Println("Load corn.json error:",err)
		panic(err)
	}
	err = json.Unmarshal(data,&GlobalObject)
	if err !=nil{
		fmt.Println("json.Unmarshal  error:",err)
		panic(err)
	}
}
func init() {
	//首先加载默认参数
	GlobalObject = &GlobalObj{
		Name:    "Corn Server",
		Version: "V0.4",
		TcpPort: 8000,
		Host:    "0.0.0.0",
		MaxConn:1000,
		MaxPackageSize:4096,
		WorkerPoolSize:8,
		MaxWorkerTaskLen:1024,
	}

	//尝试从corn.json 加载配置文件

	GlobalObject.Reload()
}
