# amap

高德地图 Web 服务 API 的 Go 语言 SDK，提供地理编码、逆地理编码和 IP 定位等功能。

## 功能特性

- ✅ **地理编码**: 将结构化地址转换为经纬度坐标
- ✅ **逆地理编码**: 将经纬度坐标转换为详细地址信息
- ✅ **IP定位**: 根据IP地址获取地理位置信息
- ✅ **智能缓存**: 可选的缓存系统，支持TTL Map和自定义实现
- ✅ **完整的数据结构**: 支持POI、道路、商圈等详细信息
- ✅ **错误处理**: 完善的错误处理机制

## 安装

```bash
go get github.com/ixugo/amap
```

## 快速开始

### 1. 获取API Key

首先需要在[高德开放平台](https://lbs.amap.com/dev/)注册账号并创建应用，获取Web服务API类型的Key。

### 2. 基本使用

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/ixugo/amap_go"
)

func main() {
    // 创建客户端
    client := amap.NewClient("YOUR_API_KEY")

    // 可选：启用缓存以提高性能
    cache := amap.NewTTLMapCache()
    client.SetCache(cache)

    // 地理编码示例
    geocodeResp, err := client.Geocode(&amap.GeocodeRequest{
        Address: "北京市朝阳区阜通东大街6号",
        City:    "北京",
    })
    if err != nil {
        log.Fatal(err)
    }

    if len(geocodeResp.Geocodes) > 0 {
        geocode := geocodeResp.Geocodes[0]
        fmt.Printf("坐标: %s\n", geocode.Location)
        fmt.Printf("经度: %.6f, 纬度: %.6f\n",
            geocode.GetLongitude(), geocode.GetLatitude())
    }
}
```

## API 文档

### 地理编码

将结构化地址转换为经纬度坐标。

```go
req := &amap.GeocodeRequest{
    Address: "北京市朝阳区阜通东大街6号", // 必填：结构化地址
    City:    "北京",                    // 可选：指定城市
}

resp, err := client.Geocode(req)
if err != nil {
    log.Fatal(err)
}

// 获取结果
for _, geocode := range resp.Geocodes {
    fmt.Printf("地址: %s\n", geocode.Location)
    fmt.Printf("省份: %s\n", geocode.Province)
    fmt.Printf("城市: %s\n", geocode.City)
    fmt.Printf("区县: %s\n", geocode.District)
    fmt.Printf("匹配级别: %s\n", geocode.Level)

    // 便捷方法获取经纬度
    lng := geocode.GetLongitude()
    lat := geocode.GetLatitude()
}
```

### 逆地理编码

将经纬度坐标转换为详细地址信息。

```go
req := &amap.RegeoRequest{
    Location:   "116.310003,39.991957", // 必填：经纬度坐标
    Radius:     1000,                   // 可选：搜索半径(米)
    Extensions: "all",                  // 可选：返回详细信息
    POIType:    []string{"050000"},     // 可选：POI类型过滤
    RoadLevel:  1,                      // 可选：道路等级
    HomeOrCorp: 1,                      // 可选：POI排序优化
}

resp, err := client.Regeo(req)
if err != nil {
    log.Fatal(err)
}

regeo := resp.Regeocode
fmt.Printf("格式化地址: %s\n", regeo.FormattedAddress)
fmt.Printf("省份: %s\n", regeo.AddressComponent.Province)
fmt.Printf("城市: %s\n", regeo.AddressComponent.City)

// 附近POI信息
for _, poi := range regeo.Pois {
    fmt.Printf("POI: %s (%s) 距离:%sm\n",
        poi.Name, poi.Type, poi.Distance)
}

// 附近道路信息
for _, road := range regeo.Roads {
    fmt.Printf("道路: %s 距离:%sm\n", road.Name, road.Distance)
}
```

### IP定位

根据IP地址获取地理位置信息。

```go
// 定位指定IP
req := &amap.IPRequest{
    IP: "114.247.50.2", // 可选：IP地址，不填则定位当前客户端IP
}

resp, err := client.IP(req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("省份: %s\n", resp.Province)
fmt.Printf("城市: %s\n", resp.City)
fmt.Printf("区域编码: %s\n", resp.AdCode)
fmt.Printf("城市范围: %s\n", resp.Rectangle)

// 或者直接获取当前IP位置
currentResp, err := client.GetCurrentIP()
```

## 配置选项

### 缓存配置

```go
// 方式1: 创建带缓存的客户端
cache := amap.NewTTLMapCache(4*time.Hour)
client := amap.NewClientWithCache("YOUR_API_KEY", cache)

// 方式2: 为现有客户端设置缓存
client := amap.NewClient("YOUR_API_KEY")
client.SetCache(amap.NewTTLMapCache(4*time.Hour))

```

### 自定义缓存实现

你可以实现`Cache`接口来使用Redis等外部缓存：

```go
type RedisCache struct {
    client *redis.Client
}

func (r *RedisCache) Get(key string) ([]byte, bool) {
    val, err := r.client.Get(context.Background(), key).Bytes()
    if err != nil {
        return nil, false
    }
    return val, true
}

func (r *RedisCache) Set(key string, value []byte, ) {
    r.client.Set(context.Background(), key, value, 4*time.Hour)
}


// 使用Redis缓存
redisCache := &RedisCache{client: redisClient}
client := amap.NewClientWithCache("YOUR_API_KEY", redisCache)
```

### 自定义HTTP客户端

```go
client := amap.NewClient("YOUR_API_KEY")
client.HTTPClient = &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        // 自定义传输配置
    },
}
```


## 运行测试

```bash
# 设置API Key环境变量
export AMAP_API_KEY="your_api_key_here"

# 运行测试
go test -v

# 运行示例
cd examples
go run main.go
```

## 数据结构说明

### 地理编码匹配级别

| 级别 | 说明 | 示例 |
|------|------|------|
| 国家 | 国家级别 | 中国 |
| 省 | 省级别 | 河北省、北京市 |
| 市 | 市级别 | 宁波市 |
| 区县 | 区县级别 | 北京市朝阳区 |
| 乡镇 | 乡镇级别 | 回龙观镇 |
| 村庄 | 村庄级别 | 三元村 |
| 热点商圈 | 商圈级别 | 上海市黄浦区老西门 |
| 道路 | 道路级别 | 北京市朝阳区阜通东大街 |
| 门牌号 | 门牌级别 | 朝阳区阜通东大街6号 |

### POI类型代码

常用POI类型代码：

- `010000`: 汽车服务
- `020000`: 汽车销售
- `030000`: 汽车维修
- `040000`: 摩托车服务
- `050000`: 餐饮服务
- `060000`: 购物服务
- `070000`: 生活服务
- `080000`: 体育休闲服务
- `090000`: 医疗保健服务
- `100000`: 住宿服务

完整的POI分类码表可以从[高德开放平台](https://lbs.amap.com/api/webservice/download)下载。

## 许可证

本项目采用 MIT 许可证，详见 [LICENSE](LICENSE) 文件。

## 贡献

欢迎提交 Issue 和 Pull Request！

## 相关链接

- [高德开放平台](https://lbs.amap.com/)
- [Web服务API文档](https://lbs.amap.com/api/webservice/summary)
- [地理/逆地理编码API](https://lbs.amap.com/api/webservice/guide/api/georegeo)
- [IP定位API](https://lbs.amap.com/api/webservice/guide/api/ipconfig)
