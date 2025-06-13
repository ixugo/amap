package amap

import (
	"encoding/json"
	"net/url"
)

// IPRequest IP定位请求参数
type IPRequest struct {
	IP string // IP地址，可选。若不填写则取客户端HTTP请求的IP进行定位
}

// IPResponse IP定位响应
type IPResponse struct {
	BaseResponse
	Province  string `json:"province"`  // 省份名称
	City      string `json:"city"`      // 城市名称
	AdCode    string `json:"adcode"`    // 城市的adcode编码
	Rectangle string `json:"rectangle"` // 所在城市矩形区域范围
}

// IP IP定位 - 根据IP地址获取位置信息
// 仅支持 IPV4，不支持国外 IP 解析。
// https://lbs.amap.com/api/webservice/guide/api/georegeo
func (c *Client) IP(req *IPRequest) (*IPResponse, error) {
	params := url.Values{}

	if req.IP != "" {
		params.Set("ip", req.IP)
	}

	// 使用带缓存的请求
	body, err := c.doRequestWithCache("ip", params, req)
	if err != nil {
		return nil, err
	}

	var resp IPResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if err := resp.GetError(); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetCurrentIP 获取当前客户端IP的位置信息
func (c *Client) GetCurrentIP() (*IPResponse, error) {
	return c.IP(&IPRequest{})
}
