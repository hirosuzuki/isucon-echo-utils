package isuconechoutils

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
)

// Usage:
//  import ieu "github.com/hirosuzuki/isucon-echo-utils"
//  e.GET("/_/routes", ieu.ShowRoutes)
//  e.GET("/_/kataribe", ieu.ShowKataribe)

func getAfter(s string, sep string) string {
	vs := strings.Split(s, sep)
	return vs[len(vs)-1]
}

func ShowRoutes(c echo.Context) error {
	html := ""
	html += "<table>\n"
	for _, v := range c.Echo().Routes() {
		name := getAfter(getAfter(v.Name, "/"), ".")
		html += fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td></tr>\n", v.Method, v.Path, name)
	}
	html += "</table>"
	return c.HTML(http.StatusOK, html)
}

// [[bundle]]
// regexp = '^(POST) /api/organizer/competition/[0-9a-zA-Z]+/finish '
// name = 'POST /api/organizer/competition/*/finish'

func ShowKataribe(c echo.Context) error {
	re := regexp.MustCompile(`:[0-9A-Za-z_]+`)
	result := ""
	for _, v := range c.Echo().Routes() {
		pattern := re.ReplaceAllString(v.Path, "[0-9A-Za-z_]+")
		result += "[[bundle]]\n"
		result += fmt.Sprintf("regexp = '^%s %s[ \\?]'\n", v.Method, pattern)
		result += fmt.Sprintf("name = '%-4s %s'\n", v.Method, v.Path)
		result += "\n"
	}
	return c.String(http.StatusOK, result)
}
