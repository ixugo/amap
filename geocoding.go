package amap

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
)

// GeocodeRequest 地理编码请求参数
type GeocodeRequest struct {
	Address string // 结构化地址信息，必填
	City    string // 指定查询的城市，可选
}

// GeocodeResponse 地理编码响应
type GeocodeResponse struct {
	BaseResponse
	Count    string    `json:"count"`
	Geocodes []Geocode `json:"geocodes"`
}

// Geocode 地理编码信息
type Geocode struct {
	Country  string `json:"country"`  // 国家
	Province string `json:"province"` // 省份
	City     string `json:"city"`     // 城市
	CityCode string `json:"citycode"` // 城市编码
	District string `json:"district"` // 区县
	Street   string `json:"street"`   // 街道
	Number   string `json:"number"`   // 门牌号
	AdCode   string `json:"adcode"`   // 区域编码
	Location string `json:"location"` // 坐标点 "经度,纬度"
	Level    string `json:"level"`    // 匹配级别
}

// GetLongitude 获取经度
func (g *Geocode) GetLongitude() float64 {
	coords := strings.Split(g.Location, ",")
	if len(coords) >= 2 {
		if lng, err := strconv.ParseFloat(coords[0], 64); err == nil {
			return lng
		}
	}
	return 0
}

// GetLatitude 获取纬度
func (g *Geocode) GetLatitude() float64 {
	coords := strings.Split(g.Location, ",")
	if len(coords) >= 2 {
		if lat, err := strconv.ParseFloat(coords[1], 64); err == nil {
			return lat
		}
	}
	return 0
}

// Geocode 地理编码 - 将地址转换为经纬度坐标
// https://lbs.amap.com/api/webservice/guide/api/georegeo
func (c *Client) Geocode(req *GeocodeRequest) (*GeocodeResponse, error) {
	params := url.Values{}
	params.Set("address", req.Address)

	if req.City != "" {
		params.Set("city", req.City)
	}

	// 使用带缓存的请求
	body, err := c.doRequestWithCache("geocode/geo", params, req)
	if err != nil {
		return nil, err
	}

	var resp GeocodeResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if err := resp.GetError(); err != nil {
		return nil, err
	}

	return &resp, nil
}

// RegeoRequest 逆地理编码请求参数
type RegeoRequest struct {
	Location   string   // 经纬度坐标，必填 "经度,纬度"
	POIType    []string // 返回附近POI类型，可选
	Radius     int      // 搜索半径，默认1000米
	Extensions string   // 返回结果控制：base(默认) 或 all
	RoadLevel  int      // 道路等级：0(所有道路) 或 1(主干道路)
	HomeOrCorp int      // POI返回顺序优化：0(不优化) 1(居家相关) 2(公司相关)
}

// RegeoResponse 逆地理编码响应
type RegeoResponse struct {
	BaseResponse
	Regeocode Regeocode `json:"regeocode"`
}

// Regeocode 逆地理编码信息
type Regeocode struct {
	FormattedAddress string           `json:"formatted_address"`       // 格式化地址
	AddressComponent AddressComponent `json:"addressComponent"`        // 地址元素
	Pois             []POI            `json:"pois,omitempty"`          // POI信息列表
	Roads            []Road           `json:"roads,omitempty"`         // 道路信息列表
	RoadInters       []RoadInter      `json:"roadinters,omitempty"`    // 道路交叉口列表
	BusinessAreas    []BusinessArea   `json:"businessAreas,omitempty"` // 商圈列表
	AOIs             []AOI            `json:"aois,omitempty"`          // AOI信息列表
}

// AddressComponent 地址元素
type AddressComponent struct {
	Country       string           `json:"country"`       // 国家
	Province      string           `json:"province"`      // 省份
	City          string           `json:"city"`          // 城市
	CityCode      string           `json:"citycode"`      // 城市编码
	District      string           `json:"district"`      // 区县
	AdCode        string           `json:"adcode"`        // 行政区编码
	Township      string           `json:"township"`      // 乡镇/街道
	TownCode      string           `json:"towncode"`      // 乡镇街道编码
	Neighborhood  Neighborhood     `json:"neighborhood"`  // 社区信息
	Building      Building         `json:"building"`      // 楼信息
	StreetNumber  StreetNumber     `json:"streetNumber"`  // 门牌信息
	SeaArea       string           `json:"seaArea"`       // 海域信息
	BusinessAreas [][]BusinessArea `json:"businessAreas"` // 商圈列表（二维数组）
}

// Neighborhood 社区信息
type Neighborhood struct {
	Name []string `json:"name"` // 社区名称列表
	Type []string `json:"type"` // POI类型列表
}

// Building 楼信息
type Building struct {
	Name []string `json:"name"` // 建筑名称列表
	Type []string `json:"type"` // 类型列表
}

// StreetNumber 门牌信息
type StreetNumber struct {
	Street    interface{} `json:"street"`    // 街道名称（可能是字符串或数组）
	Number    interface{} `json:"number"`    // 门牌号（可能是字符串或数组）
	Location  interface{} `json:"location"`  // 坐标点（可能是字符串或数组）
	Direction interface{} `json:"direction"` // 方向（可能是字符串或数组）
	Distance  interface{} `json:"distance"`  // 距离（可能是字符串或数组）
}

// GetStreet 获取街道名称（处理字符串或数组）
func (s *StreetNumber) GetStreet() string {
	switch v := s.Street.(type) {
	case string:
		return v
	case []interface{}:
		if len(v) > 0 {
			if str, ok := v[0].(string); ok {
				return str
			}
		}
	}
	return ""
}

// GetNumber 获取门牌号（处理字符串或数组）
func (s *StreetNumber) GetNumber() string {
	switch v := s.Number.(type) {
	case string:
		return v
	case []interface{}:
		if len(v) > 0 {
			if str, ok := v[0].(string); ok {
				return str
			}
		}
	}
	return ""
}

// GetLocation 获取坐标点（处理字符串或数组）
func (s *StreetNumber) GetLocation() string {
	switch v := s.Location.(type) {
	case string:
		return v
	case []interface{}:
		if len(v) > 0 {
			if str, ok := v[0].(string); ok {
				return str
			}
		}
	}
	return ""
}

// GetDirection 获取方向（处理字符串或数组）
func (s *StreetNumber) GetDirection() string {
	switch v := s.Direction.(type) {
	case string:
		return v
	case []interface{}:
		if len(v) > 0 {
			if str, ok := v[0].(string); ok {
				return str
			}
		}
	}
	return ""
}

// GetDistance 获取距离（处理字符串或数组）
func (s *StreetNumber) GetDistance() string {
	switch v := s.Distance.(type) {
	case string:
		return v
	case []interface{}:
		if len(v) > 0 {
			if str, ok := v[0].(string); ok {
				return str
			}
		}
	}
	return ""
}

// POI POI信息
type POI struct {
	ID           string `json:"id"`           // POI ID
	Name         string `json:"name"`         // POI名称
	Type         string `json:"type"`         // POI类型
	Tel          string `json:"tel"`          // 电话
	Distance     string `json:"distance"`     // 距离
	Direction    string `json:"direction"`    // 方向
	Address      string `json:"address"`      // 地址
	Location     string `json:"location"`     // 坐标点
	BusinessArea string `json:"businessarea"` // 商圈名称
}

// Road 道路信息
type Road struct {
	ID        string `json:"id"`        // 道路ID
	Name      string `json:"name"`      // 道路名称
	Distance  string `json:"distance"`  // 距离
	Direction string `json:"direction"` // 方向
	Location  string `json:"location"`  // 坐标点
}

// RoadInter 道路交叉口
type RoadInter struct {
	Distance   string `json:"distance"`    // 距离
	Direction  string `json:"direction"`   // 方向
	Location   string `json:"location"`    // 坐标点
	FirstID    string `json:"first_id"`    // 第一条道路ID
	FirstName  string `json:"first_name"`  // 第一条道路名称
	SecondID   string `json:"second_id"`   // 第二条道路ID
	SecondName string `json:"second_name"` // 第二条道路名称
}

// BusinessArea 商圈信息
type BusinessArea struct {
	Location string `json:"location"` // 商圈中心点
	Name     string `json:"name"`     // 商圈名称
	ID       string `json:"id"`       // 商圈ID
}

// AOI AOI信息
type AOI struct {
	ID       string `json:"id"`       // AOI ID
	Name     string `json:"name"`     // AOI名称
	AdCode   string `json:"adcode"`   // 区域编码
	Location string `json:"location"` // 中心点坐标
	Area     string `json:"area"`     // 面积
	Distance string `json:"distance"` // 距离
	Type     string `json:"type"`     // AOI类型
}

// Regeo 逆地理编码 - 将经纬度坐标转换为地址
// https://lbs.amap.com/api/webservice/guide/api/georegeo
func (c *Client) Regeo(req *RegeoRequest) (*RegeoResponse, error) {
	params := url.Values{}
	params.Set("location", req.Location)

	if len(req.POIType) > 0 {
		params.Set("poitype", strings.Join(req.POIType, "|"))
	}

	if req.Radius > 0 {
		params.Set("radius", strconv.Itoa(req.Radius))
	}

	if req.Extensions != "" {
		params.Set("extensions", req.Extensions)
	}

	if req.RoadLevel > 0 {
		params.Set("roadlevel", strconv.Itoa(req.RoadLevel))
	}

	if req.HomeOrCorp > 0 {
		params.Set("homeorcorp", strconv.Itoa(req.HomeOrCorp))
	}

	// 使用带缓存的请求
	body, err := c.doRequestWithCache("geocode/regeo", params, req)
	if err != nil {
		return nil, err
	}

	var resp RegeoResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if err := resp.GetError(); err != nil {
		return nil, err
	}

	return &resp, nil
}
