package isuconechoutils

import (
	"math/rand"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Usage:
//  import ieu "github.com/hirosuzuki/isucon-echo-utils"
// 	e.Use(ieu.UniqueID)

// nginx:
// 	log_format with_time '$remote_addr $sent_http_x_unique_id $sent_http_x_request_id $cookie_user [$time_local] '
//    '"$request" $status $body_bytes_sent '
//    '"$http_referer" "$http_user_agent" $request_time';
//  access_log /var/log/nginx/access.log with_time;

func UniqueID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var uniqueID string = ""
		uidCookie, err := c.Cookie("uid")
		if err == nil {
			uniqueID = uidCookie.Value
		}
		if uniqueID == "" {
			uniqueID = "u:" + strconv.FormatUint(rand.Uint64()|(1<<63), 16)
			cookie := http.Cookie{
				Name:  "uid",
				Value: uniqueID,
			}
			c.SetCookie(&cookie)
		}
		res := c.Response()
		res.Header().Set("X-Unique-ID", uniqueID)
		return next(c)
	}
}
