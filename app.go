package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx    context.Context
	config *Config
}

// Config represents the application configuration
type Config struct {
	Token     string `json:"token"`
	VehicleID string `json:"vehicleId"`
}

// VehicleData represents the vehicle information
type VehicleData struct {
	VehicleName    string        `json:"vehicleName"`
	BmsSoc         string        `json:"bmssoc"`
	HmiRidableMile string        `json:"hmiRidableMile"`
	VehiclePicUrl  string        `json:"vehiclePicUrl"`
	Location       *LocationData `json:"location"`
}

// LocationData represents location information
type LocationData struct {
	CoordinateSystem string  `json:"coordinateSystem"`
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	LocationTime     string  `json:"locationTime"`
	Altitude         float64 `json:"altitude"`
	Address          string  `json:"address,omitempty"`
}

// AmapResponse represents the response from Amap API
type AmapResponse struct {
	Status    string `json:"status"`
	Regeocode struct {
		FormattedAddress string `json:"formatted_address"`
	} `json:"regeocode"`
}

// APIResponse represents the API response structure
type APIResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    VehicleData `json:"data"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	app := &App{}
	app.loadConfig()
	return app
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// GetVehicleData fetches vehicle data from the API
func (a *App) GetVehicleData() (*VehicleData, error) {
	// 检查配置是否存在
	if a.config.Token == "" || a.config.VehicleID == "" {
		return nil, fmt.Errorf("请先配置Token和车架号")
	}

	url := fmt.Sprintf("https://tapi.zeehoev.com/v1.0/app/cfmotoserverapp/vehicle/widgets/%s", a.config.VehicleID)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 使用配置中的token
	req.Header.Set("Authorization", "Bearer "+a.config.Token)
	req.Header.Set("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", "acw_tc=0b32824217388280172008957ec68b4a84c95e1b5efd8b103d6c69b40480d9")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查是否返回HTML（错误页面）
	if len(body) > 0 && body[0] == '<' {
		return nil, fmt.Errorf("认证失败，请检查Token是否正确")
	}

	// 检查非200状态码
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API返回错误状态码 %d: %s", resp.StatusCode, string(body))
	}

	var apiResponse APIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	// 检查API返回的状态码
	if apiResponse.Code != "10000" {
		return nil, fmt.Errorf("API返回错误: %s", apiResponse.Message)
	}

	// 获取地址信息
	if apiResponse.Data.Location != nil {
		address, err := a.getAddressFromLocation(apiResponse.Data.Location.Longitude, apiResponse.Data.Location.Latitude)
		if err == nil {
			apiResponse.Data.Location.Address = address
		}
	}

	return &apiResponse.Data, nil
}

// GetMockVehicleData 返回模拟数据用于测试
func (a *App) GetMockVehicleData() *VehicleData {
	return &VehicleData{
		VehicleName:    "Zeeho AE8",
		BmsSoc:         "85",
		HmiRidableMile: "120",
		VehiclePicUrl:  "https://via.placeholder.com/150x150/333/fff?text=Zeeho",
		Location: &LocationData{
			CoordinateSystem: "2",
			Latitude:         22.584569557649274,
			Longitude:        114.11163989303444,
			LocationTime:     "01/15 14:30",
			Altitude:         54.2,
			Address:          "广东省深圳市南山区科技园",
		},
	}
}

// getAddressFromLocation 根据经纬度获取地址信息
func (a *App) getAddressFromLocation(longitude, latitude float64) (string, error) {
	amapURL := fmt.Sprintf("http://restapi.amap.com/v3/geocode/regeo?output=json&location=%f,%f&key=948b3368e904b58cca42531bcfc5b064&extensions=all", longitude, latitude)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(amapURL)
	if err != nil {
		return "", fmt.Errorf("请求地址信息失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取地址响应失败: %v", err)
	}

	var amapResponse AmapResponse
	if err := json.Unmarshal(body, &amapResponse); err != nil {
		return "", fmt.Errorf("解析地址JSON失败: %v", err)
	}

	if amapResponse.Status != "1" {
		return "", fmt.Errorf("地址解析失败")
	}

	return amapResponse.Regeocode.FormattedAddress, nil
}

// 配置文件路径
func (a *App) getConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".zeeho-config.json")
}

// 加载配置
func (a *App) loadConfig() {
	configPath := a.getConfigPath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		a.config = &Config{}
		return
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		a.config = &Config{}
		return
	}

	a.config = &config
}

// 保存配置
func (a *App) saveConfig() error {
	configPath := a.getConfigPath()
	data, err := json.MarshalIndent(a.config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// GetConfig 获取当前配置
func (a *App) GetConfig() *Config {
	return a.config
}

// ValidateAndSaveConfig 验证并保存配置
func (a *App) ValidateAndSaveConfig(token, vehicleId string) error {
	// 创建临时配置进行验证
	tempConfig := &Config{
		Token:     token,
		VehicleID: vehicleId,
	}

	// 验证配置是否有效
	if err := a.validateConfig(tempConfig); err != nil {
		return fmt.Errorf("配置验证失败: %v", err)
	}

	// 验证成功，保存配置
	a.config = tempConfig
	if err := a.saveConfig(); err != nil {
		return fmt.Errorf("保存配置失败: %v", err)
	}

	return nil
}

// 验证配置
func (a *App) validateConfig(config *Config) error {
	if config.Token == "" {
		return fmt.Errorf("Token不能为空")
	}

	if config.VehicleID == "" {
		return fmt.Errorf("车架号不能为空")
	}

	// 尝试调用API验证
	url := fmt.Sprintf("https://tapi.zeehoev.com/v1.0/app/cfmotoserverapp/vehicle/widgets/%s", config.VehicleID)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("创建验证请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Authorization", "Bearer "+config.Token)
	req.Header.Set("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", "acw_tc=0b32824217388280172008957ec68b4a84c95e1b5efd8b103d6c69b40480d9")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("验证请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取验证响应失败: %v", err)
	}

	// 检查响应状态
	if resp.StatusCode != 200 {
		return fmt.Errorf("API返回错误状态码 %d，请检查Token和车架号是否正确", resp.StatusCode)
	}

	// 检查是否返回HTML（认证失败）
	if len(body) > 0 && body[0] == '<' {
		return fmt.Errorf("认证失败，请检查Token是否正确")
	}

	// 尝试解析JSON
	var apiResponse APIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return fmt.Errorf("API响应格式错误: %v", err)
	}

	// 检查API返回的状态码
	if apiResponse.Code != "10000" {
		return fmt.Errorf("API返回错误: %s", apiResponse.Message)
	}

	return nil
}

// SetWindowPosition 设置窗口位置
func (a *App) SetWindowPosition(x, y int) {
	runtime.WindowSetPosition(a.ctx, x, y)
}

// MoveToCorner 移动窗口到指定角落
func (a *App) MoveToCorner(corner string) {
	// 获取屏幕尺寸
	screen, err := runtime.ScreenGetAll(a.ctx)
	if err != nil || len(screen) == 0 {
		return
	}

	screenWidth := screen[0].Size.Width
	screenHeight := screen[0].Size.Height

	// 窗口尺寸
	windowWidth := 440
	windowHeight := 300

	var x, y int

	switch corner {
	case "top-left":
		x, y = 20, 20
	case "top-right":
		x, y = screenWidth-windowWidth-20, 20
	case "bottom-left":
		x, y = 20, screenHeight-windowHeight-60
	case "bottom-right":
		x, y = screenWidth-windowWidth-20, screenHeight-windowHeight-60
	default:
		// 默认右上角
		x, y = screenWidth-windowWidth-20, 20
	}

	runtime.WindowSetPosition(a.ctx, x, y)
}

// MinimizeToTray 最小化到系统托盘
func (a *App) MinimizeToTray() {
	runtime.WindowHide(a.ctx)
}

// ShowWindow 显示窗口
func (a *App) ShowWindow() {
	runtime.WindowShow(a.ctx)
}

// StartDrag 开始拖动窗口 (在 Wails v2 中通过 CSS 实现)
func (a *App) StartDrag() {
	// Wails v2 中拖动功能通过前端 CSS 实现
	// 这个方法保留用于兼容性，实际拖动在前端处理
}
