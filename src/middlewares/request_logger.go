package middlewares

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"io"
	"bytes"
	"fmt"
)

func RequestLogger() gin.HandlerFunc {
    return func(c *gin.Context) {
		log := slog.Default()
		body, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		headers := c.Request.Header

		query := c.Request.URL.Query()

		pathParams := make(map[string]string)
		for _, p := range c.Params {
			pathParams[p.Key] = p.Value
		}

		log.Info(fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path),
			"headers", headers,
			"query", query,
			"params", pathParams,
			"path", c.Request.URL.Path,
			"body", string(body),
		)

		c.Next()
	}
}