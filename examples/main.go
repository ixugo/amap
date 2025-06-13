package main

import (
	"fmt"
	"log"
	"time"

	amap "github.com/ixugo/amap"
)

func main() {
	// 创建高德地图API客户端
	// 请替换为你的API Key
	client := amap.NewClient("YOUR_API_KEY")

	// 可选：启用缓存以提高性能
	cache := amap.NewTTLMapCache(4 * time.Hour)
	client.SetCache(cache)

	// 示例1: 地理编码 - 地址转坐标
	fmt.Println("=== 地理编码示例 ===")
	geocodeReq := &amap.GeocodeRequest{
		Address: "北京市朝阳区阜通东大街6号",
		City:    "北京",
	}

	geocodeResp, err := client.Geocode(geocodeReq)
	if err != nil {
		log.Printf("地理编码失败: %v", err)
	} else if len(geocodeResp.Geocodes) > 0 {
		geocode := geocodeResp.Geocodes[0]
		fmt.Printf("地址: %s\n", geocodeReq.Address)
		fmt.Printf("坐标: %s\n", geocode.Location)
		fmt.Printf("经度: %.6f\n", geocode.GetLongitude())
		fmt.Printf("纬度: %.6f\n", geocode.GetLatitude())
		fmt.Printf("省份: %s\n", geocode.Province)
		fmt.Printf("城市: %s\n", geocode.City)
		fmt.Printf("区县: %s\n", geocode.District)
		fmt.Printf("匹配级别: %s\n", geocode.Level)
	}

	fmt.Println()

	// 示例2: 逆地理编码 - 坐标转地址
	fmt.Println("=== 逆地理编码示例 ===")
	regeoReq := &amap.RegeoRequest{
		Location:   "116.310003,39.991957",
		Radius:     1000,
		Extensions: "all", // 返回详细信息
	}

	regeoResp, err := client.Regeo(regeoReq)
	if err != nil {
		log.Printf("逆地理编码失败: %v", err)
	} else {
		regeo := regeoResp.Regeocode
		fmt.Printf("坐标: %s\n", regeoReq.Location)
		fmt.Printf("格式化地址: %s\n", regeo.FormattedAddress)
		fmt.Printf("省份: %s\n", regeo.AddressComponent.Province)
		fmt.Printf("城市: %s\n", regeo.AddressComponent.City)
		fmt.Printf("区县: %s\n", regeo.AddressComponent.District)
		fmt.Printf("乡镇/街道: %s\n", regeo.AddressComponent.Township)

		// 显示附近POI信息
		if len(regeo.Pois) > 0 {
			fmt.Printf("附近POI (%d个):\n", len(regeo.Pois))
			for i, poi := range regeo.Pois {
				if i >= 3 { // 只显示前3个
					break
				}
				fmt.Printf("  - %s (%s) 距离:%sm\n", poi.Name, poi.Type, poi.Distance)
			}
		}

		// 显示附近道路信息
		if len(regeo.Roads) > 0 {
			fmt.Printf("附近道路 (%d条):\n", len(regeo.Roads))
			for i, road := range regeo.Roads {
				if i >= 3 { // 只显示前3条
					break
				}
				fmt.Printf("  - %s 距离:%sm\n", road.Name, road.Distance)
			}
		}
	}

	fmt.Println()

	// 示例3: IP定位
	fmt.Println("=== IP定位示例 ===")

	// 定位指定IP
	ipReq := &amap.IPRequest{
		IP: "114.247.50.2", // 示例IP
	}

	ipResp, err := client.IP(ipReq)
	if err != nil {
		log.Printf("IP定位失败: %v", err)
	} else {
		fmt.Printf("IP地址: %s\n", ipReq.IP)
		fmt.Printf("省份: %s\n", ipResp.Province)
		fmt.Printf("城市: %s\n", ipResp.City)
		fmt.Printf("区域编码: %s\n", ipResp.AdCode)
		fmt.Printf("城市范围: %s\n", ipResp.Rectangle)
	}

	fmt.Println()

	// 示例4: 获取当前客户端IP位置
	fmt.Println("=== 当前IP定位示例 ===")
	currentIPResp, err := client.GetCurrentIP()
	if err != nil {
		log.Printf("当前IP定位失败: %v", err)
	} else {
		fmt.Printf("当前IP所在省份: %s\n", currentIPResp.Province)
		fmt.Printf("当前IP所在城市: %s\n", currentIPResp.City)
		fmt.Printf("当前IP区域编码: %s\n", currentIPResp.AdCode)
	}

	fmt.Println()

	// 示例5: 缓存演示
	fmt.Println("=== 缓存演示 ===")
	fmt.Println("第一次请求（会调用API并缓存）:")
	start := time.Now()
	_, err = client.Geocode(&amap.GeocodeRequest{
		Address: "上海市黄浦区南京东路",
		City:    "上海",
	})
	duration1 := time.Since(start)
	if err != nil {
		fmt.Printf("   请求失败: %v\n", err)
	} else {
		fmt.Printf("   请求耗时: %v\n", duration1)
	}

	fmt.Println("第二次相同请求（从缓存获取）:")
	start = time.Now()
	_, err = client.Geocode(&amap.GeocodeRequest{
		Address: "上海市黄浦区南京东路",
		City:    "上海",
	})
	duration2 := time.Since(start)
	if err != nil {
		fmt.Printf("   请求失败: %v\n", err)
	} else {
		fmt.Printf("   请求耗时: %v\n", duration2)
		if duration2 < duration1 {
			fmt.Println("   ✅ 缓存生效，第二次请求更快")
		}
	}

	fmt.Println()
	fmt.Println("=== 缓存功能说明 ===")
	fmt.Println("- 默认缓存8小时过期")
	fmt.Println("- 相同请求参数会复用缓存结果")
	fmt.Println("- 可以使用Redis等实现Cache接口")
	fmt.Println("- 缓存是可选的，不设置则每次都调用API")
}
