package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/bestk/zeeho-widgets/backend"
	"github.com/go-co-op/gocron"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx       context.Context
	config    *Config
	scheduler *gocron.Scheduler
}

// Config represents the application configuration
type Config struct {
	Token          string `json:"token"`
	VehicleID      string `json:"vehicleId"`
	UpdateInterval int    `json:"updateInterval"`
}

// VehicleData represents the vehicle information
type VehicleData struct {
	VinNo                      string        `json:"vinNo"`
	DeviceName                 string        `json:"deviceName"`
	BindStartTime              string        `json:"bindStartTime"`
	Fate                       string        `json:"fate"`
	VehicleName                string        `json:"vehicleName"`
	VehiclePicUrl              string        `json:"vehiclePicUrl"`
	VehicleBackPicUrl          string        `json:"vehicleBackPicUrl"`
	ShareName                  *string       `json:"shareName"`
	IsShared                   string        `json:"isShared"`
	RideMileageMonth           string        `json:"rideMileageMonth"`
	RidingTimeMonth            string        `json:"ridingTimeMonth"`
	RidingTimeMonthUnitMinute  string        `json:"ridingTimeMonthUnitMinute"`
	AvgVelocityMonth           string        `json:"avgVelocityMonth"`
	BmsSoc                     string        `json:"bmssoc"`
	HmiRidableMile             string        `json:"hmiRidableMile"`
	GsmRxLev                   string        `json:"gsmRxLev"`
	GsmRxLevValue              string        `json:"gsmRxLevValue"`
	BluetoothAddress           string        `json:"bluetoothAddress"`
	HmiBluetoothAddress        string        `json:"hmiBluetoothAddress"`
	ChargeState                string        `json:"chargeState"`
	FullChargeTime             string        `json:"fullChargeTime"`
	Pressure                   string        `json:"pressure"`
	PressureValue              string        `json:"pressureValue"`
	HeadLockState              string        `json:"headLockState"`
	RideState                  string        `json:"rideState"`
	GreenContribution          string        `json:"greenContribution"`
	OtaVersion                 string        `json:"otaVersion"`
	ShareUserId                *string       `json:"shareUserId"`
	BindingUserId              int64         `json:"bindingUserId"`
	CarMaster                  string        `json:"carMaster"`
	VehicleType                string        `json:"vehicleType"`
	VehicleTypeName            string        `json:"vehicleTypeName"`
	EncryptInfo                EncryptInfo   `json:"encryptInfo"`
	RedPoint                   int           `json:"redPoint"`
	HmiRidableMileAbnormalShow *string       `json:"hmiRidableMileAbnormalShow"`
	ExpectFullTimeDescribe     *string       `json:"expectFullTimeDescribe"`
	Location                   Location      `json:"location"`
	NavigationType             int           `json:"navigationType"`
	Navigation                 string        `json:"navigation"`
	Projection                 string        `json:"projection"`
	MotoPlay                   int           `json:"motoPlay"`
	WifiAddress                string        `json:"wifiAddress"`
	BluetoothSearch            bool          `json:"bluetoothSearch"`
	VehicleTypeDetailId        int           `json:"vehicleTypeDetailId"`
	ShareEndTime               *string       `json:"shareEndTime"`
	ResidualSeconds            *string       `json:"residualSeconds"`
	SupportNetworkUnlock       int           `json:"supportNetworkUnlock"`
	IntelligentType            *string       `json:"intelligentType"`
	TotalRideMile              string        `json:"totalRideMile"`
	MaxMileage                 string        `json:"maxMileage"`
	DeviceType                 int           `json:"deviceType"`
	BroadcastType              string        `json:"broadcastType"`
	SupportUnlock              int           `json:"supportUnlock"`
	BindDate                   *string       `json:"bindDate"`
	CyclingEventStatisticFlag  bool          `json:"cyclingEventStatisticFlag"`
	WhetherChargeState         bool          `json:"whetherChargeState"`
	FirstBindDate              string        `json:"firstBindDate"`
	MaxRangeMileage            string        `json:"maxRangeMileage"`
	OnlineStatus               string        `json:"onlineStatus"`
	ActivationDate             string        `json:"activationDate"`
	LastUseDate                int           `json:"lastUseDate"`
	RechargeEndDate            string        `json:"rechargeEndDate"`
	OpenCushionFlag            bool          `json:"openCushionFlag"`
	OpenStorageBoxFlag         bool          `json:"openStorageBoxFlag"`
	LoudlySearchCar            int           `json:"loudlySearchCar"`
	MmiUuid                    string        `json:"mmiUuid"`
	ProductKey                 string        `json:"productKey"`
	IotInstanceId              string        `json:"iotInstanceId"`
	IotProperties              []IotProperty `json:"iotProperties"`
	ServiceRechargeStatus      string        `json:"serviceRechargeStatus"`
	RefreshTime                string        `json:"refreshTime"`
	VehicleScalePicUrl         string        `json:"vehicleScalePicUrl"`
	GaodeLincenseVinNo         string        `json:"gaodeLincenseVinNo"`
	GaodeLincenseId            string        `json:"gaodeLincenseId"`
}

type EncryptInfo struct {
	Key          string `json:"key"`
	Iv           string `json:"iv"`
	EncryptValue string `json:"encryptValue"`
}

type Location struct {
	Longitude        float64 `json:"longitude"`
	Latitude         float64 `json:"latitude"`
	Altitude         float64 `json:"altitude"`
	CoordinateSystem string  `json:"coordinateSystem"`
	LocationTime     string  `json:"locationTime"`
	Address          string  `json:"address,omitempty"`
}

type IotProperty struct {
	Name         string  `json:"name"`
	Identify     string  `json:"identify"`
	Value        string  `json:"value"`
	Time         string  `json:"time"`
	DbUpdateTime *string `json:"dbUpdateTime"`
	Describe     *string `json:"describe"`
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
	Data    interface{} `json:"data"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	app := &App{
		scheduler: gocron.NewScheduler(time.UTC),
	}
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

	// 将 interface{} 转换为 map
	dataMap, ok := apiResponse.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("数据格式错误")
	}

	// 重新序列化为 JSON
	dataJSON, err := json.Marshal(dataMap)
	if err != nil {
		return nil, fmt.Errorf("数据转换失败: %v", err)
	}

	// 解析为 VehicleData
	var data VehicleData
	if err := json.Unmarshal(dataJSON, &data); err != nil {
		return nil, fmt.Errorf("数据解析失败: %v", err)
	}

	// 获取地址信息
	if data.Location.Longitude != 0 && data.Location.Latitude != 0 {
		address, err := a.getAddressFromLocation(data.Location.Longitude, data.Location.Latitude)
		if err == nil {
			data.Location.Address = address
		}
	}

	return &data, nil
}

// VehicleHomePage 获取车辆首页数据
func (a *App) VehicleHomePage() (*[]VehicleData, error) {
	url := "https://tapi.zeehoev.com/v1.0/app/cfmotoserverapp/vehicleHomePage"

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

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

	var apiResponse APIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	// 将 interface{} 转换为 []interface{}
	dataSlice, ok := apiResponse.Data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("数据格式错误")
	}

	// 重新序列化为 JSON
	dataJSON, err := json.Marshal(dataSlice)
	if err != nil {
		return nil, fmt.Errorf("数据转换失败: %v", err)
	}

	// 解析为 []VehicleData
	var data []VehicleData
	if err := json.Unmarshal(dataJSON, &data); err != nil {
		return nil, fmt.Errorf("数据解析失败: %v", err)
	}

	// 获取每个车辆的地址信息
	for i := range data {
		if data[i].Location.Longitude != 0 && data[i].Location.Latitude != 0 {
			address, err := a.getAddressFromLocation(data[i].Location.Longitude, data[i].Location.Latitude)
			if err == nil {
				data[i].Location.Address = address
			}
		}
	}

	return &data, nil
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

	runtime.EventsEmit(a.ctx, "configUpdate", string(data))

	return os.WriteFile(configPath, data, 0644)
}

// GetConfig 获取当前配置
func (a *App) GetConfig() *Config {
	return a.config
}

// ValidateAndSaveConfig 验证并保存配置
func (a *App) ValidateAndSaveConfig(token, vehicleId string, updateInterval int) error {
	// 创建临时配置进行验证
	tempConfig := &Config{
		Token:          token,
		VehicleID:      vehicleId,
		UpdateInterval: updateInterval,
	}

	if vehicleId != "" {
		// 验证配置是否有效
		if err := a.validateConfig(tempConfig); err != nil {
			return fmt.Errorf("配置验证失败: %v", err)
		}
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

// Widget 小部件
func (a *App) StartWidget() {
	backend.SetupDesktopChildWidget()
}

// Quit 退出应用程序
func (a *App) Quit() {
	runtime.Quit(a.ctx)
}

func (a *App) ScheduleRefresh() {
	// Stop any existing scheduled tasks
	a.scheduler.Clear()

	if a.config.UpdateInterval < 1 {
		log.Println("UpdateInterval must > 0")
	}

	log.Println("Schedule refresh based on config interval")
	// Schedule refresh based on config interval
	a.scheduler.Every(a.config.UpdateInterval).Minutes().Do(func() {
		// Refresh vehicle data
		data, err := a.VehicleHomePage()
		if err != nil {
			// Handle error - could emit event to frontend
			runtime.EventsEmit(a.ctx, "refreshError", err.Error())
		} else {
			// Emit success event to frontend
			runtime.EventsEmit(a.ctx, "dataRefreshed", data)
		}
	})

	// Start the scheduler
	a.scheduler.StartAsync()
}
