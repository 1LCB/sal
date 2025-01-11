package sal

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/go-swagno/swagno/components/endpoint"
	"github.com/go-swagno/swagno/components/http/response"
	"github.com/go-swagno/swagno/components/mime"
	"github.com/go-swagno/swagno/components/parameter"
)


func getHeaders(headers []Header) []*parameter.Parameter {
	endpointsConfig := []*parameter.Parameter{}
	for _, v := range headers {
		if v.Required {
			endpointsConfig = append(
				endpointsConfig,
				parameter.StrParam(v.Name, parameter.Header, parameter.WithRequired()),
			)
			continue
		}

		endpointsConfig = append(
			endpointsConfig,
			parameter.StrParam(v.Name, parameter.Header, parameter.WithDefault(nil)),
		)
	}
	return endpointsConfig
}

func getPathParameters(route string) []*parameter.Parameter {
	r := regexp.MustCompile(`{(\w+)}`)
	matches := r.FindAllString(route, -1)

	endpointsConfig := []*parameter.Parameter{}

	for _, match := range matches {
		match = strings.Trim(match, "{}")

		endpointsConfig = append(endpointsConfig, parameter.StrParam(
			match, parameter.Path, parameter.WithRequired(),
		))
	}

	return endpointsConfig
}

func makeEndpoint(route, method, tag string, resp Response, body any, headers []Header) *endpoint.EndPoint {
	allParameters := append(getPathParameters(route), getHeaders(headers)...)

	endpoint := endpoint.New(
		endpoint.MethodType(method),
		route,
		endpoint.WithParams(
			allParameters...,
		),
		endpoint.WithConsume([]mime.MIME{mime.JSON}),
		endpoint.WithProduce([]mime.MIME{mime.JSON}),
		endpoint.WithBody(body),
		endpoint.WithTags(tag),
		endpoint.WithSuccessfulReturns([]response.Response{
			response.New(resp.Body, resp.Status, resp.Description),
		}),
	)

	return endpoint
}

func (r *router) addRoute(method, endpoint string, salHandler SalHandlerFunc) {
	if r.middlewares != nil {
		for _, middleware := range r.middlewares {
			salHandler = middleware(salHandler)
		}
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		c := ctxPool.Get().(*Ctx)
		defer ctxPool.Put(c)

		c.w = w
		c.Request = r

		salHandler(c)
	}

	pattern := method + " " + r.prefix + endpoint
	http.HandleFunc(pattern, LoggerMiddleware(handler))
}
