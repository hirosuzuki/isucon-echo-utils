package middleware

import (
	"math/rand"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

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
