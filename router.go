package framework

type RoutePath string

type IRouterGroup interface {
	GetRouterGroup() RoutePath
	GET(route string, request IRouter) *RouterGroup
	POST(route string, request IRouter) *RouterGroup
	PUT(route string, request IRouter) *RouterGroup
	DELETE(route string, request IRouter) *RouterGroup
	GetRoutes() []RouterEndpointGroup
}

type RouterEndpointGroup struct {
	Method     string
	Route      string
	Controller IRouter
}

type RouterGroup struct {
	group  RoutePath
	routes []RouterEndpointGroup
}

func (r *RouterGroup) GetRoutes() []RouterEndpointGroup {
	return r.routes
}

func (r *RouterGroup) addRouter(method, route string, controller IRouter) {
	r.routes = append(r.routes, RouterEndpointGroup{
		Method:     method,
		Route:      route,
		Controller: controller,
	})
}
func (r *RouterGroup) GET(route string, request IRouter) *RouterGroup {
	//TODO implement me
	r.addRouter("GET", route, request)
	return r
}

func (r *RouterGroup) POST(route string, request IRouter) *RouterGroup {
	//TODO implement me
	r.addRouter("POST", route, request)
	return r
}

func (r *RouterGroup) PUT(route string, request IRouter) *RouterGroup {
	r.addRouter("PUT", route, request)
	return r
}

func (r *RouterGroup) DELETE(route string, request IRouter) *RouterGroup {
	r.addRouter("DELETE", route, request)
	return r
}

func (r *RouterGroup) GetRouterGroup() RoutePath {
	//TODO implement me
	return r.group
}

func NewRouterMapper(routerGroup RoutePath) *RouterGroup {
	return &RouterGroup{group: routerGroup}
}
