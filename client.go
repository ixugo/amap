package amap

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	// BaseURL 高德地图API基础URL
	BaseURL = "https://restapi.amap.com"
	// API版本
	APIVersion = "v3"
)

// Client 高德地图API客户端
type Client struct {
	APIKey     string
	HTTPClient *http.Client
	BaseURL    string
	Cache      Cache // 可选缓存接口
}

// NewClient 创建新的高德地图API客户端
func NewClient(apiKey string) *Client {
	return &Client{
		APIKey:  apiKey,
		BaseURL: BaseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:    10,
				IdleConnTimeout: time.Minute,
			},
		},
	}
}

// NewClientWithCache 创建带缓存的高德地图API客户端
func NewClientWithCache(apiKey string, cache Cache) *Client {
	return &Client{
		APIKey:  apiKey,
		BaseURL: BaseURL,
		Cache:   cache,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetHTTPClient 设置 HTTP 客户端
func (c *Client) SetHTTPClient(cli *http.Client) *Client {
	c.HTTPClient = cli
	return c
}

// SetCache 设置缓存
func (c *Client) SetCache(cache Cache) {
	c.Cache = cache
}

// SetTimeout 设置HTTP请求超时时间
func (c *Client) SetTimeout(timeout time.Duration) {
	c.HTTPClient.Timeout = timeout
}

// doRequest 执行HTTP请求
func (c *Client) doRequest(endpoint string, params url.Values) ([]byte, error) {
	// 添加API Key
	params.Set("key", c.APIKey)

	// 构建完整URL
	fullURL := fmt.Sprintf("%s/%s/%s?%s", c.BaseURL, APIVersion, endpoint, params.Encode())

	// 发送GET请求
	resp, err := c.HTTPClient.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("http get err: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body err: %w", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP状态码错误: %d, 响应: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// doRequestWithCache 执行带缓存的HTTP请求
func (c *Client) doRequestWithCache(endpoint string, params url.Values, cacheKeyParams interface{}) ([]byte, error) {
	// 如果没有缓存，直接请求
	if c.Cache == nil {
		return c.doRequest(endpoint, params)
	}

	// 生成缓存键
	cacheKey := generateCacheKey(endpoint, cacheKeyParams)

	// 尝试从缓存获取
	if cachedData, found := c.Cache.Get(cacheKey); found {
		return cachedData, nil
	}

	// 缓存未命中，执行实际请求
	data, err := c.doRequest(endpoint, params)
	if err != nil {
		return nil, err
	}

	// 将结果存入缓存
	c.Cache.Set(cacheKey, data)

	return data, nil
}

// BaseResponse 基础响应结构
type BaseResponse struct {
	Status   string `json:"status"`
	Info     string `json:"info"`
	InfoCode string `json:"infocode"`
}

// IsSuccess 检查响应是否成功
func (r *BaseResponse) IsSuccess() bool {
	return r.Status == "1"
}

// GetError 获取错误信息
func (r *BaseResponse) GetError() error {
	if r.IsSuccess() {
		return nil
	}
	return fmt.Errorf("API错误: %s (代码: %s)", r.Info, r.InfoCode)
}
