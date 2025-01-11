package sal

type Map = map[string]string
type Middleware func(next SalHandlerFunc) SalHandlerFunc

type router struct {
	prefix      string
	tag         string
	headers     []Header
	middlewares []Middleware
}

type Header struct {
	Name     string
	Required bool
}

func NewRouter(prefix, tag string) *router {
	return &router{
		prefix:      prefix,
		tag:         tag,
		headers:     nil,
		middlewares: nil,
	}
}

// UseMiddleware adds a middleware to the router. The middleware will be applied to all routes.
func (r *router) UseMiddleware(m Middleware) {
	if m == nil{
		return
	}
	
	r.middlewares = append(r.middlewares, m)
}

// UseHeader adds a required or optional header for the routes in the router.
func (r *router) UseHeader(name string, required bool) {
	r.headers = append(r.headers, Header{Name: name, Required: required})
}

func (r *router) POST(route string, body any, resp Response, handler SalHandlerFunc) {
	endpoint := makeEndpoint(r.prefix+route, "POST", r.tag, resp, body, r.headers)
	Swag.AddEndpoint(endpoint)

	r.addRoute("POST", route, handler)
}

func (r *router) GET(route string, resp Response, handler SalHandlerFunc) {
	endpoint := makeEndpoint(r.prefix+route, "GET", r.tag, resp, nil, r.headers)
	Swag.AddEndpoint(endpoint)

	r.addRoute("GET", route, handler)
}

func (r *router) PATCH(route string, body any, resp Response, handler SalHandlerFunc) {
	endpoint := makeEndpoint(r.prefix+route, "PATCH", r.tag, resp, body, r.headers)
	Swag.AddEndpoint(endpoint)

	r.addRoute("PATCH", route, handler)
}

func (r *router) PUT(route string, body any, resp Response, handler SalHandlerFunc) {
	endpoint := makeEndpoint(r.prefix+route, "PUT", r.tag, resp, body, r.headers)
	Swag.AddEndpoint(endpoint)

	r.addRoute("PUT", route, handler)
}

func (r *router) DELETE(route string, resp Response, handler SalHandlerFunc) {
	endpoint := makeEndpoint(r.prefix+route, "DELETE", r.tag, resp, nil, r.headers)
	Swag.AddEndpoint(endpoint)

	r.addRoute("DELETE", route, handler)
}
