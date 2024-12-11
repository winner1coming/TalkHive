package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

// RunRouteTests 运行路由测试
func RunRouteTests(router *gin.Engine) {
	tests := []struct {
		Name       string
		Method     string
		URL        string
		Body       interface{}
		StatusCode int
	}{
		{
			Name:       "测试用户注册",
			Method:     "POST",
			URL:        "/register",
			Body:       map[string]interface{}{"id": "testuser", "email": "test@example.com", "password": "123456"},
			StatusCode: http.StatusOK,
		},
		{
			Name:       "测试用户登录",
			Method:     "POST",
			URL:        "/login",
			Body:       map[string]interface{}{"id": "testuser", "password": "123456"},
			StatusCode: http.StatusOK,
		},
		{
			Name:       "测试获取用户资料",
			Method:     "GET",
			URL:        "/profile/testuser",
			Body:       nil,
			StatusCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		log.Printf("运行测试: %s", test.Name)

		// 创建请求
		var bodyBytes []byte
		if test.Body != nil {
			bodyBytes, _ = json.Marshal(test.Body)
		}
		req := httptest.NewRequest(test.Method, test.URL, bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// 创建响应记录器
		w := httptest.NewRecorder()

		// 运行路由
		router.ServeHTTP(w, req)

		// 验证结果
		if w.Code != test.StatusCode {
			log.Printf("测试失败: %s - 预期状态码 %d, 实际状态码 %d", test.Name, test.StatusCode, w.Code)
		} else {
			log.Printf("测试成功: %s", test.Name)
		}
	}
}
