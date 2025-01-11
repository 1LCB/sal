package sal

import (
	"encoding/json"
	"net/http"
	"sync"
)

var (
	ctxPool = sync.Pool{
		New: func() any {
			return &Ctx{}
		},
	}
)

type Ctx struct {
	w       http.ResponseWriter
	Request *http.Request
}

type SalHandlerFunc func(c *Ctx)

func (f SalHandlerFunc) ServeHTTP(c *Ctx) {
	f(c)
}

func (c *Ctx) Json(content any, status int) error {
	c.w.Header().Add("Content-Type", "application/json")
	c.w.WriteHeader(status)
	if err := json.NewEncoder(c.w).Encode(content); err != nil {
		return err
	}
	return nil
}

func (c *Ctx) Error(message string, status int) error {
	return c.Json(map[string]string{"error": message}, status)
}

func (c *Ctx) Text(content string, status int) error {
	c.w.Header().Add("Content-Type", "text/plain")
	c.w.WriteHeader(status)
	if _, err := c.w.Write([]byte(content)); err != nil {
		return err
	}
	return nil
}

func (c *Ctx) Redirect(url string, status int) {
	http.Redirect(c.w, c.Request, url, status)
}

func (c *Ctx) HTML(content string, status int) error {
	c.w.Header().Add("Content-Type", "text/html")
	c.w.WriteHeader(status)
	if _, err := c.w.Write([]byte(content)); err != nil {
		return err
	}
	return nil
}

func (c *Ctx) File(filePath string) error {
	http.ServeFile(c.w, c.Request, filePath)
	return nil
}

func (c *Ctx) NoContent(status int) {
	c.w.WriteHeader(status)
}

func (c *Ctx) Header(key, value string) {
	c.w.Header().Set(key, value)
}

func (c *Ctx) Binary(content []byte, contentType string, status int) error {
	c.w.Header().Add("Content-Type", contentType)
	c.w.WriteHeader(status)
	if _, err := c.w.Write(content); err != nil {
		return err
	}
	return nil
}
