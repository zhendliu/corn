package net

import (
	"corn/iface"
	"corn/utils"
	"fmt"
)

/*
	消息处理模块的实现
*/

type MsgHandle struct {
	//存放每个msgID 对应的方法
	Apis map[uint32]iface.IRouter

	//负责Worker读取任务的消息队列
	TaskQueue []chan iface.IRequest
	//业务工作池的worker数量
	WorkerPoolSize uint32
}

//提供一个创建MsgHandle的方法

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]iface.IRouter),
		//全局配置中获取
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan iface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

func (mh *MsgHandle) DoMsgHandler(request iface.IRequest) {
	//从request中找到msgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if ok == false {
		fmt.Printf("api msgID =%d  is not fond router ,need register !!!", request.GetMsgID())
		return
	}
	//根据msgID 调度对应的router业务即可
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (mh *MsgHandle) AddRouter(msgID uint32, router iface.IRouter) {
	//判断当前的msg绑定的处理方法是否已经存在
	if _, ok := mh.Apis[msgID]; ok == true {
		fmt.Println("exist msgID,add error")
	}
	fmt.Println("添加:", msgID, "对应路由")
	//添加msg与api的绑定关系
	mh.Apis[msgID] = router
	fmt.Printf("add api MsgId = %d, success ", msgID)
}

//启动一个Worker工作池
func (mh *MsgHandle) StartWorkerPool() {
	//根据workPoolSize 分别开启Worker ，分别开启Worker 用一个Goroutine 来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个Worker被启动
		//为每个worker分配管道，开辟空间
		mh.TaskQueue[i] = make(chan iface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//启动当前的Worker ，阻塞消息的到来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}

}

//启动一个Worker工作流程

func (mh *MsgHandle) StartOneWorker(workID int, taskQueue chan iface.IRequest) {
	fmt.Printf("Worker ID = %d is started ... \n", workID)
	//不断的阻塞对应的消息
	for {
		select {
		//如果有消息过来，制定当前的Request 所绑定的业务
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

//

func (mh *MsgHandle) SendMsgToTaskQueue(request iface.IRequest) {

	//将消息平均分配给不同的worker
	//基本的轮训法则
	workID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(), "Request ID = ",
		request.GetMsgID(), "To Worker ID:", workID)
	//将消息发送给对应的worker 的TaskQueue
	mh.TaskQueue[workID] <- request

}
