package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HealthState struct {
    Ready bool
}

func NewHealthState() *HealthState {
    return &HealthState{Ready: false}
}

// мидлварь, которая вешает маршруты и использует HealthState
func HealthMiddleware(state *HealthState) gin.HandlerFunc {
    return func(c *gin.Context) {
        path := c.Request.URL.Path

        switch path {
        case "/livez":
            c.Status(http.StatusOK)
            c.Abort()
            return
        case "/readyz":
            // проверяем готовность зависимостей
            if state.Ready {
                c.Status(http.StatusOK)
            } else {
                c.Status(http.StatusServiceUnavailable)
            }
            c.Abort()
            return
        }

        c.Next()
    }
}
