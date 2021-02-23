package httpproxy

import (
	"github.com/gogf/gf/text/gstr"
	"net/http"
)

// 从HTTP请求中获取用户ID。
// 注意：该方法与业务相关，按照业务既定规则获取数据。SessionName必须包含"__"字符串分隔UserID。
func parseUserId(r *http.Request) string {
	var (
		userId      string
		sessionName string
	)
	// Cookie
	if v, err := r.Cookie("sess"); err == nil {
		sessionName = v.Value
	}
	// Query
	if sessionName == "" || !gstr.Contains(sessionName, "__") {
		values := r.URL.Query()
		sessionName = values.Get("sess")
		if sessionName == "" || !gstr.Contains(sessionName, "__") {
			sessionName = values.Get("sessName")
		}
	}
	// Header
	if sessionName == "" || !gstr.Contains(sessionName, "__") {
		sessionName = r.Header.Get("X-Sess-Name")
	}

	if p := gstr.Pos(sessionName, "__"); p > 0 {
		userId = sessionName[:p]
	}
	return userId
}
