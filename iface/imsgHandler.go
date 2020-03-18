package iface

/*
	消息管理抽象层
*/
type IMsgHandle interface {

	//调度执行路由器
	DoMsgHandler(request IRequest)
	//为消息添加具体的路由
	AddRouter(msgID uint32, router IRouter)
	//启动worker工作池
	StartWorkerPool()

	SendMsgToTaskQueue(request IRequest)
}
