package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"PanshiCMS/internal/database"
	"PanshiCMS/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// setupTestServer 初始化一个用于测试的 Gin 服务器
func setupTestServer() *gin.Engine {
	// 设置测试配置
	viper.Set("database.type", "sqlite")
	viper.Set("database.dsn", "file::memory:?cache=shared")
	viper.Set("jwt.secret", "testsecret")

	// 初始化数据库（内存数据库）
	database.InitDB()

	// 设置 Gin 引擎
	r := gin.Default()
	router.SetupRouter(r)
	return r
}

// helper function to parse JSON response
func parseJSONBody(body []byte, v interface{}) error {
	return json.Unmarshal(body, v)
}

// TestAdminInitAndLogin 测试初始化管理员和登录功能
func TestAdminInitAndLogin(t *testing.T) {
	r := setupTestServer()

	// 1. 初始化管理员
	initBody := map[string]string{"username": "admin", "password": "password"}
	initData, _ := json.Marshal(initBody)
	req := httptest.NewRequest(http.MethodPost, "/init", bytes.NewBuffer(initData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("init admin expected status 200, got %d", w.Code)
	}

	// 2. 登录获取 token
	loginBody := map[string]string{"username": "admin", "password": "password"}
	loginData, _ := json.Marshal(loginBody)
	req = httptest.NewRequest(http.MethodPost, "/admin/login", bytes.NewBuffer(loginData))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("login expected status 200, got %d", w.Code)
	}
	// 解析 token
	var loginResp struct {
		Token string `json:"token"`
	}
	if err := parseJSONBody(w.Body.Bytes(), &loginResp); err != nil {
		t.Fatalf("parse login response failed: %v", err)
	}
	if loginResp.Token == "" {
		t.Fatalf("token should not be empty")
	}

	// 3. 创建一篇新闻
	newsBody := map[string]interface{}{
		"title":       "Test Title",
		"summary":     "Test summary",
		"content":     "<p>Test content</p>",
		"slug":        "test-title",
		"publishDate": "2024-01-01T00:00:00Z",
		"metaTitle":   "Test Meta Title",
	}
	newsData, _ := json.Marshal(newsBody)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/news", bytes.NewBuffer(newsData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+loginResp.Token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("create news expected status 200, got %d", w.Code)
	}
	// 读取返回的数据，获取 ID
	var createResp struct {
		Data struct {
			ID uint `json:"ID"`
		} `json:"data"`
	}
	if err := parseJSONBody(w.Body.Bytes(), &createResp); err != nil {
		t.Fatalf("parse create news response failed: %v", err)
	}
	if createResp.Data.ID == 0 {
		t.Fatalf("create news should return valid ID")
	}

	// 4. 获取新闻列表，验证搜索功能
	req = httptest.NewRequest(http.MethodGet, "/api/v1/news?page=1&pageSize=10&q=Test", nil)
	req.Header.Set("Authorization", "Bearer "+loginResp.Token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("list news expected status 200, got %d", w.Code)
	}
	var listResp struct {
		Data struct {
			List  []map[string]interface{} `json:"list"`
			Total int                      `json:"total"`
		} `json:"data"`
	}
	if err := parseJSONBody(w.Body.Bytes(), &listResp); err != nil {
		t.Fatalf("parse list response failed: %v", err)
	}
	if listResp.Data.Total == 0 {
		t.Fatalf("expected at least one news item in list")
	}

	// 5. 更新新闻
	updateBody := map[string]interface{}{
		"title":       "Updated Title",
		"summary":     "Updated summary",
		"content":     "Updated content",
		"slug":        "updated-title",
		"publishDate": "2024-02-02T00:00:00Z",
	}
	updateData, _ := json.Marshal(updateBody)
	// 更新请求
	req = httptest.NewRequest(http.MethodPut, "/api/v1/news/"+strconv.FormatUint(uint64(createResp.Data.ID), 10), bytes.NewBuffer(updateData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+loginResp.Token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		body, _ := ioutil.ReadAll(w.Body)
		t.Fatalf("update news expected status 200, got %d, body: %s", w.Code, string(body))
	}

	// 6. 删除新闻
	req = httptest.NewRequest(http.MethodDelete, "/api/v1/news/"+strconv.FormatUint(uint64(createResp.Data.ID), 10), nil)
	req.Header.Set("Authorization", "Bearer "+loginResp.Token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("delete news expected status 200, got %d", w.Code)
	}
}
