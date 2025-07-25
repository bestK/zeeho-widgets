//go:build windows

package main

import (
	"context"
	"embed"
	"syscall"
	"time"
	"unsafe"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

// Windows API 声明
var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	procGetWindowLongW      = user32.NewProc("GetWindowLongW")
	procSetWindowLongW      = user32.NewProc("SetWindowLongW")
	procFindWindowW         = user32.NewProc("FindWindowW")
	procSetWindowPos        = user32.NewProc("SetWindowPos")
	procGetForegroundWindow = user32.NewProc("GetForegroundWindow")
	procGetWindowTextW      = user32.NewProc("GetWindowTextW")
	procIsWindowVisible     = user32.NewProc("IsWindowVisible")
)

var (
	widgetHwnd         uintptr   // 保存小部件窗口句柄
	isTopmost          bool      // 当前是否置顶状态
	lastForegroundHwnd uintptr   // 上次的前台窗口
	stateChangeTime    time.Time // 上次状态改变时间
)

const (
	GWL_EXSTYLE      = -20
	WS_EX_TOOLWINDOW = 0x00000080
	WS_EX_APPWINDOW  = 0x00040000
	WS_EX_TOPMOST    = 0x00000008
	WS_EX_NOACTIVATE = 0x08000000

	// SetWindowPos flags
	SWP_NOSIZE       = 0x0001
	SWP_NOMOVE       = 0x0002
	SWP_NOZORDER     = 0x0004
	SWP_FRAMECHANGED = 0x0020
	SWP_NOACTIVATE   = 0x0010
	SWP_NOREDRAW     = 0x0008

	// Special HWND values
	HWND_TOPMOST = ^uintptr(0) // -1
)

// 设置桌面小部件样式，不在 Alt+Tab 中显示，但不覆盖其他窗口
func setupDesktopWidget(windowTitle string) {
	// 查找窗口句柄
	hwnd, _, _ := procFindWindowW.Call(
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(windowTitle))),
	)
	if hwnd == 0 {
		println("找不到窗口:", windowTitle)
		return
	}

	// 获取当前扩展样式
	gwlExstyle := uintptr(^uint32(19)) // -20 的无符号表示
	exStyle, _, _ := procGetWindowLongW.Call(hwnd, gwlExstyle)

	// 设置桌面小部件样式：
	// - WS_EX_TOOLWINDOW: 不在 Alt+Tab 中显示
	// - WS_EX_NOACTIVATE: 不激活窗口，不抢夺焦点
	// - 移除 WS_EX_APPWINDOW: 不在任务栏显示
	// - 不设置 WS_EX_TOPMOST: 允许其他窗口覆盖
	newStyle := (exStyle | WS_EX_TOOLWINDOW | WS_EX_NOACTIVATE) &^ WS_EX_APPWINDOW

	// 应用新样式
	procSetWindowLongW.Call(hwnd, gwlExstyle, newStyle)

	// 通知系统样式已更改，但不设置为最顶层
	procSetWindowPos.Call(
		hwnd,
		0, // 不设置为 HWND_TOPMOST
		0, 0, 0, 0,
		uintptr(SWP_FRAMECHANGED|SWP_NOMOVE|SWP_NOSIZE|SWP_NOZORDER|SWP_NOACTIVATE),
	)

	// 保存窗口句柄用于后续监控
	widgetHwnd = hwnd

}

// 检查窗口是否为系统窗口或桌面
func isSystemWindow(hwnd uintptr) bool {
	if hwnd == 0 {
		return true
	}

	// 检查窗口是否可见
	visible, _, _ := procIsWindowVisible.Call(hwnd)
	if visible == 0 {
		return true
	}

	// 获取窗口标题
	titleBuf := make([]uint16, 256)
	procGetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&titleBuf[0])), 256)
	title := syscall.UTF16ToString(titleBuf)

	// 系统窗口或空标题窗口
	systemWindows := []string{"", "Program Manager", "Desktop Window Manager", "Windows Shell Experience Host"}
	for _, sysWin := range systemWindows {
		if title == sysWin {
			return true
		}
	}

	return false
}

// 设置窗口层级（优化版，减少闪烁）
func setWindowLevel(topmost bool) {
	if widgetHwnd == 0 || isTopmost == topmost {
		return
	}

	var hwndInsertAfter uintptr
	if topmost {
		hwndInsertAfter = HWND_TOPMOST
	} else {
		hwndInsertAfter = 1 // HWND_TOP
	}

	// 使用平滑的标志组合，减少视觉闪烁
	procSetWindowPos.Call(
		widgetHwnd,
		hwndInsertAfter,
		0, 0, 0, 0,
		uintptr(SWP_NOMOVE|SWP_NOSIZE|SWP_NOACTIVATE),
	)

	isTopmost = topmost

}

// 监控前台窗口变化（优化版，减少闪烁）
func monitorForegroundWindow() {
	const debounceDelay = 1 * time.Second // 防抖延迟

	for {
		time.Sleep(500 * time.Millisecond) // 降低检查频率到500ms

		if widgetHwnd == 0 {
			continue
		}

		// 获取当前前台窗口
		foregroundHwnd, _, _ := procGetForegroundWindow.Call()

		// 如果前台窗口没有变化，跳过
		if foregroundHwnd == lastForegroundHwnd {
			continue
		}

		// 如果前台窗口是小部件本身，跳过
		if foregroundHwnd == widgetHwnd {
			lastForegroundHwnd = foregroundHwnd
			continue
		}

		// 防抖：如果状态刚刚改变，等待一段时间再处理
		if time.Since(stateChangeTime) < debounceDelay {
			continue
		}

		lastForegroundHwnd = foregroundHwnd

		// 检查是否为系统窗口或桌面
		shouldBeTopmost := isSystemWindow(foregroundHwnd)

		// 只有当需要的状态与当前状态不同时才切换
		if shouldBeTopmost != isTopmost {
			setWindowLevel(shouldBeTopmost)
			stateChangeTime = time.Now()
		}
	}
}

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "Zeeho Widget",
		Width:  440,
		Height: 300,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 255},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)

			// 延迟执行，等待窗口完全创建
			go func() {
				time.Sleep(500 * time.Millisecond)
				setupDesktopWidget("Zeeho Widget")

				// 启动前台窗口监控
				go monitorForegroundWindow()
			}()
		},
		Bind: []interface{}{
			app,
		},
		Frameless:         true,
		AlwaysOnTop:       false,
		DisableResize:     false,
		StartHidden:       false,
		HideWindowOnClose: false,
		WindowStartState:  options.Normal,
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
