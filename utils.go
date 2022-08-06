package isuconechoutils

import (
	"math/rand"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Usage:
//   import isuconechoutils "github.com/hirosuzuki/isucon-echo-utils"
//   isuconechoutils.Setup(e)

// Nginx:
//   log_format with_time '$remote_addr $sent_http_x_request_id $sent_http_x_unique_id [$time_local] '
//     '"$request" $status $body_bytes_sent '
//     '"$http_referer" "$http_user_agent" $request_time';
//   access_log /var/log/nginx/access.log with_time;

func Setup(e *echo.Echo) {
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return "r:" + strconv.FormatUint(rand.Uint64()|(1<<63), 16)
		},
	}))
	e.Use(UniqueID)
	e.Use(Logger())
	e.GET("/_/routes", ShowRoutes)
	e.GET("/_/kataribe", ShowKataribe)
}
