package routers

type RouterGroup struct {
	Order OrderRouter
}

var OrderServiceRouterGroup = new(RouterGroup)
