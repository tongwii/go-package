package common

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

const (
	XAPIToken      = "tongwii"                                   // 头文件信息
	XPoweredBy     = "tongwii.tongcloud; email=info@tongwii.com" // 头文件信息
	TokenExpiresAt = 1440                                        // Token到期时间 分钟
)

// http默认返回值
type Content struct {
	Data    interface{} `json:"data" xml:"data"`
	Message string      `json:"message" xml:"message"`
	Status  int         `json:"status" xml:"status"`
	Total   int64       `json:"total" xml:"total"`
}

// 默认token对象
type JwtCustomClaims struct {
	ID     string `json:"id"`     // 保存用户ID
	Update string `json:"update"` // 保存用户更新时间
	jwt.StandardClaims
}

// echo配合输出日志
func echoLogEntry(c echo.Context) *log.Entry {
	if c == nil {
		return Logger.WithFields(log.Fields{
			//"at": time.Now().Format("2006-01-02 15:04:05"),
			"module": "http",
		})
	}

	return Logger.WithFields(log.Fields{
		//"at":     time.Now().Format("2006-01-02 15:04:05"),
		"module": "http",
		"method": c.Request().Method,
		"uri":    c.Request().URL.String(),
		"ip":     c.Request().RemoteAddr,
	})
}

// 日志中间件
func SetLogging(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		echoLogEntry(c).Info("incoming request")
		return next(c)
	}
}

// 错误处理
func EchoErrorHandler(err error, c echo.Context) {
	report, ok := err.(*echo.HTTPError)
	if ok {
		report.Message = fmt.Sprintf("http error %d - %v", report.Code, report.Message)
	} else {
		report = echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	echoLogEntry(c).Error(report.Message)
	_ = c.HTML(report.Code, report.Message.(string))
}

/**
设置跨域信息
*/
func SetCORS() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		AllowHeaders:  []string{echo.HeaderContentType, echo.HeaderOrigin, echo.HeaderAuthorization, echo.HeaderAcceptEncoding},
		ExposeHeaders: []string{"Token", "X-Powered-By"},
	})
}

// 当前请求信息
type RequestInfo struct {
	User    interface{}
	Path    string // 请求路径 已过滤
	Visitor string // 访问者，mass、docker、client
}

// 生成token
func CreateToken(id, update, key string) (string, error) {
	claims := &JwtCustomClaims{
		ID:     id,     // 用户ID
		Update: update, // 用户更新时间
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(TokenExpiresAt)).Unix(), // 有效时间24小时
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // 生成token对象
	signed, err := token.SignedString([]byte(key))             // 秘钥加密
	return signed, err
}

// 自动设置返回头长度
func SetContentLength(c *echo.Context, v interface{}) error {
	retBytes, err := json.Marshal(v)
	if err == nil {
		(*c).Response().Header().Set("Content-Length", strconv.Itoa(len(retBytes)+1))
	}
	return err
}
