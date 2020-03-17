package net

import "corn/iface"

/*
实现router时，先嵌入这个baserouter 基类，然后根据需要对这个基类的方法就行重写
*/
type BaseRouter struct{}



//在处理conn业务之前的方法hook
func (br *BaseRouter)PreHandle(request  iface.IRequest){}
//在处理conn业务的主方法hook
func (br *BaseRouter)Handle(request  iface.IRequest){}
//在处理conn业务之后的方法hook
func (br *BaseRouter) PostHandle(request  iface.IRequest){}