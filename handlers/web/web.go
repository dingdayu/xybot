package web

import (
	"io"
	"net/http"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	// 将字符串通过回写指针返回给浏览器
	io.WriteString(w, string("Hello Word!"))
}
