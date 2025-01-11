package sal

import (
	"strconv"
	"sync"
)

var poolResponses = sync.Pool{
	New: func() any {
		return &Response{}
	},
}

type Response struct {
	Body        any
	Status      string
	Description string
}

type H = map[string]string
type EmptyResponse = map[string]string

func NewResponse(body any, status int) Response {
	if body == nil {
		body = EmptyResponse{}
	}
	resp := poolResponses.Get().(*Response)
	defer poolResponses.Put(resp)

	resp.Body = body
	resp.Status = strconv.Itoa(status)
	resp.Description = StatusCodes[status]
	return *resp
}