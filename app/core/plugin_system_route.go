package core

// IRouteList -
type IRouteList interface {
	Routes() []IRoute
}

// IRoute -
type IRoute interface {
	Method() string
	Path() string
}
