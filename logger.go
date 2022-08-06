package isuconechoutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/valyala/fasttemplate"
)

// Usage:
//  import ieu "github.com/hirosuzuki/isucon-echo-utils"
// 	e.Use(ieu.Logger())
// 	e.Use(ieu.Logger())
// 	e.Use(ieu.Logger())

const (
	LOG_FORMAT = `${remote_ip} ${uid} - [${time_apache}] "${method} ${uri} ${protocol}" ${status} ${bytes_out} "${referer}" "${user_agent}" ${latency_sec}`
	// Tags to construct the logger format.
	//
	// - time_unix
	// - time_unix_milli
	// - time_unix_micro
	// - time_unix_nano
	// - time_rfc3339
	// - time_rfc3339_nano
	// - time_apache
	// - id (Request ID / X-Request-ID Header)
	// - uid (Unique ID / X-Unique-ID Response Header)
	// - remote_ip
	// - uri
	// - host
	// - method
	// - path
	// - protocol
	// - referer
	// - user_agent
	// - status
	// - error
	// - latency (In nanoseconds)
	// - latency_sec (In seconds)
	// - latency_msec (In miniseconds)
	// - latency_human (Human readable)
	// - bytes_in (Bytes received)
	// - bytes_out (Bytes sent)
	// - header:<NAME>
	// - query:<NAME>
	// - form:<NAME>
)

func Logger() echo.MiddlewareFunc {
	logfile := os.Getenv("APPLOG")
	if logfile == "" {
		logfile = "/tmp/app.log"
	}
	output, _ := os.Create(logfile)
	template := fasttemplate.New(LOG_FORMAT+"\n", "${", "}")
	pool := &sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 256))
		},
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if output == nil {
				return
			}

			req := c.Request()
			res := c.Response()
			start := time.Now()
			if err = next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()
			buf := pool.Get().(*bytes.Buffer)
			buf.Reset()
			defer pool.Put(buf)

			if _, err = template.ExecuteFunc(buf, func(w io.Writer, tag string) (int, error) {
				switch tag {
				case "time_unix":
					return buf.WriteString(strconv.FormatInt(time.Now().Unix(), 10))
				case "time_unix_milli":
					// go 1.17 or later, it supports time#UnixMilli()
					return buf.WriteString(strconv.FormatInt(time.Now().UnixNano()/1000000, 10))
				case "time_unix_micro":
					// go 1.17 or later, it supports time#UnixMicro()
					return buf.WriteString(strconv.FormatInt(time.Now().UnixNano()/1000, 10))
				case "time_unix_nano":
					return buf.WriteString(strconv.FormatInt(time.Now().UnixNano(), 10))
				case "time_rfc3339":
					return buf.WriteString(time.Now().Format(time.RFC3339))
				case "time_rfc3339_nano":
					return buf.WriteString(time.Now().Format(time.RFC3339Nano))
				case "time_apache":
					return buf.WriteString(time.Now().Format("02/Jan/2006:15:04:05 -0700"))
				case "id":
					id := req.Header.Get(echo.HeaderXRequestID)
					if id == "" {
						id = res.Header().Get(echo.HeaderXRequestID)
					}
					return buf.WriteString(id)
				case "uid":
					id := res.Header().Get("X-Unique-ID")
					if id == "" {
						id = "-"
					}
					return buf.WriteString(id)
				case "remote_ip":
					return buf.WriteString(c.RealIP())
				case "host":
					return buf.WriteString(req.Host)
				case "uri":
					return buf.WriteString(req.RequestURI)
				case "method":
					return buf.WriteString(req.Method)
				case "path":
					p := req.URL.Path
					if p == "" {
						p = "/"
					}
					return buf.WriteString(p)
				case "protocol":
					return buf.WriteString(req.Proto)
				case "referer":
					return buf.WriteString(req.Referer())
				case "user_agent":
					return buf.WriteString(req.UserAgent())
				case "status":
					return buf.WriteString(strconv.Itoa(res.Status))
				case "error":
					if err != nil {
						// Error may contain invalid JSON e.g. `"`
						b, _ := json.Marshal(err.Error())
						b = b[1 : len(b)-1]
						return buf.Write(b)
					}
				case "latency":
					l := stop.Sub(start)
					return buf.WriteString(strconv.FormatInt(int64(l), 10))
				case "latency_sec":
					l := stop.Sub(start)
					return buf.WriteString(fmt.Sprintf("%.3f", float64(l)/1000000000))
				case "latency_msec":
					l := stop.Sub(start)
					return buf.WriteString(fmt.Sprintf("%.3f", float64(l)/1000000))
				case "latency_human":
					return buf.WriteString(stop.Sub(start).String())
				case "bytes_in":
					cl := req.Header.Get(echo.HeaderContentLength)
					if cl == "" {
						cl = "0"
					}
					return buf.WriteString(cl)
				case "bytes_out":
					return buf.WriteString(strconv.FormatInt(res.Size, 10))
				default:
					switch {
					case strings.HasPrefix(tag, "header:"):
						return buf.Write([]byte(c.Request().Header.Get(tag[7:])))
					case strings.HasPrefix(tag, "query:"):
						return buf.Write([]byte(c.QueryParam(tag[6:])))
					case strings.HasPrefix(tag, "form:"):
						return buf.Write([]byte(c.FormValue(tag[5:])))
					case strings.HasPrefix(tag, "cookie:"):
						cookie, err := c.Cookie(tag[7:])
						if err == nil {
							return buf.Write([]byte(cookie.Value))
						}
					}
				}
				return 0, nil
			}); err != nil {
				return
			}
			_, err = output.Write(buf.Bytes())
			return
		}
	}
}
