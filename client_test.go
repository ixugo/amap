package amap

import (
	"encoding/json"
	"os"
	"testing"
)

const testAPIKey = "YOUR_TEST_API_KEY" // 请替换为你的测试API Key

func getTestClient() *Client {
	apiKey := os.Getenv("AMAP_API_KEY")
	if apiKey == "" {
		apiKey = testAPIKey
	}
	return NewClient(apiKey)
}

func TestGeocoding(t *testing.T) {
	client := getTestClient()

	req := &GeocodeRequest{
		Address: "北京市朝阳区阜通东大街6号",
		City:    "北京",
	}

	resp, err := client.Geocode(req)
	if err != nil {
		t.Fatalf("地理编码失败: %v", err)
	}

	if len(resp.Geocodes) == 0 {
		t.Fatal("地理编码结果为空")
	}

	geocode := resp.Geocodes[0]
	if geocode.Location == "" {
		t.Error("坐标信息为空")
	}

	if geocode.GetLongitude() == 0 || geocode.GetLatitude() == 0 {
		t.Error("经纬度解析失败")
	}

	t.Logf("地理编码成功: %s -> %s", req.Address, geocode.Location)
}

func TestRegeo(t *testing.T) {
	client := getTestClient()

	req := &RegeoRequest{
		Location: "123.1238223,123.1239871",
		Radius:   100,
	}

	resp, err := client.Regeo(req)
	if err != nil {
		t.Fatalf("逆地理编码失败: %v", err)
	}

	if resp.Regeocode.FormattedAddress == "" {
		t.Error("格式化地址为空")
	}

	t.Logf("逆地理编码成功: %s -> %s", req.Location, resp.Regeocode.FormattedAddress)
}

func TestRegeoWithAll(t *testing.T) {
	client := getTestClient()

	req := &RegeoRequest{
		Location:   "116.310003,39.991957",
		Radius:     1000,
		Extensions: "all",
	}

	resp, err := client.Regeo(req)
	if err != nil {
		t.Fatalf("逆地理编码(详细)失败: %v", err)
	}

	regeo := resp.Regeocode
	if regeo.FormattedAddress == "" {
		t.Error("格式化地址为空")
	}

	if regeo.AddressComponent.Province == "" {
		t.Error("省份信息为空")
	}

	t.Logf("逆地理编码(详细)成功: %s -> %s", req.Location, regeo.FormattedAddress)
	t.Logf("省份: %s, 城市: %s, 区县: %s",
		regeo.AddressComponent.Province,
		regeo.AddressComponent.City,
		regeo.AddressComponent.District)

	if len(regeo.Pois) > 0 {
		t.Logf("附近POI数量: %d", len(regeo.Pois))
	}
}

func TestIP(t *testing.T) {
	client := getTestClient()

	req := &IPRequest{
		IP: "114.247.50.2", // 测试IP
	}

	resp, err := client.IP(req)
	if err != nil {
		t.Fatalf("IP定位失败: %v", err)
	}

	if resp.Province == "" {
		t.Error("省份信息为空")
	}

	if resp.City == "" {
		t.Error("城市信息为空")
	}

	t.Logf("IP定位成功: %s -> %s %s", req.IP, resp.Province, resp.City)
}

func TestGetCurrentIP(t *testing.T) {
	client := getTestClient()

	resp, err := client.GetCurrentIP()
	if err != nil {
		t.Fatalf("当前IP定位失败: %v", err)
	}

	// 注意：在某些测试环境中，可能返回"局域网"
	t.Logf("当前IP定位结果: %s %s", resp.Province, resp.City)
}

func TestClientTimeout(t *testing.T) {
	client := getTestClient()
	client.SetTimeout(1) // 设置1毫秒超时，应该会失败

	req := &GeocodeRequest{
		Address: "北京市朝阳区阜通东大街6号",
	}

	_, err := client.Geocode(req)
	if err == nil {
		t.Error("期望超时错误，但请求成功了")
	}

	t.Logf("超时测试通过: %v", err)
}

func TestRegeoJSONParsing(t *testing.T) {
	// 使用实际的API返回数据进行测试
	jsonData := `{"status":"1","regeocode":{"addressComponent":{"city":"合肥市","province":"安徽省","adcode":"340104","district":"蜀山区","towncode":"340104401000","streetNumber":{"number":"900号","location":"117.100235,31.832138","direction":"东南","distance":"601.566","street":"望江西路"},"country":"中国","township":"高新技术产业开发区","businessAreas":[[]],"building":{"name":[],"type":[]},"neighborhood":{"name":[],"type":[]},"citycode":"0551"},"formatted_address":"安徽省合肥市蜀山区高新技术产业开发区中安创谷科技园二期(北门)"},"info":"OK","infocode":"10000"}`

	var resp RegeoResponse
	err := json.Unmarshal([]byte(jsonData), &resp)
	if err != nil {
		t.Fatalf("JSON解析失败: %v", err)
	}

	// 验证基本字段
	if resp.Status != "1" {
		t.Errorf("期望状态为1，实际为: %s", resp.Status)
	}

	if resp.Info != "OK" {
		t.Errorf("期望信息为OK，实际为: %s", resp.Info)
	}

	if resp.InfoCode != "10000" {
		t.Errorf("期望信息代码为10000，实际为: %s", resp.InfoCode)
	}

	// 验证地址信息
	regeo := resp.Regeocode
	expectedAddress := "安徽省合肥市蜀山区高新技术产业开发区中安创谷科技园二期(北门)"
	if regeo.FormattedAddress != expectedAddress {
		t.Errorf("期望格式化地址为: %s，实际为: %s", expectedAddress, regeo.FormattedAddress)
	}

	// 验证地址组件
	addr := regeo.AddressComponent
	if addr.Province != "安徽省" {
		t.Errorf("期望省份为安徽省，实际为: %s", addr.Province)
	}

	if addr.City != "合肥市" {
		t.Errorf("期望城市为合肥市，实际为: %s", addr.City)
	}

	if addr.District != "蜀山区" {
		t.Errorf("期望区县为蜀山区，实际为: %s", addr.District)
	}

	// 验证数组字段不为nil
	if addr.Building.Name == nil {
		t.Error("Building.Name 不应该为 nil")
	}

	if addr.Building.Type == nil {
		t.Error("Building.Type 不应该为 nil")
	}

	if addr.Neighborhood.Name == nil {
		t.Error("Neighborhood.Name 不应该为 nil")
	}

	if addr.Neighborhood.Type == nil {
		t.Error("Neighborhood.Type 不应该为 nil")
	}

	if addr.BusinessAreas == nil {
		t.Error("BusinessAreas 不应该为 nil")
	}

	// 验证门牌信息
	street := addr.StreetNumber
	if street.Street != "望江西路" {
		t.Errorf("期望街道为望江西路，实际为: %s", street.Street)
	}

	if street.Number != "900号" {
		t.Errorf("期望门牌号为900号，实际为: %s", street.Number)
	}

	t.Log("JSON解析测试通过")
}
