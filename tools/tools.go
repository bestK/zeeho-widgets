package tools

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"
)

var (
	user32 = syscall.NewLazyDLL("user32.dll")

	procFindWindowW                = user32.NewProc("FindWindowW")
	procFindWindowExW              = user32.NewProc("FindWindowExW")
	procSetWindowLongW             = user32.NewProc("SetWindowLongW")
	procGetWindowLongW             = user32.NewProc("GetWindowLongW")
	procSetWindowPos               = user32.NewProc("SetWindowPos")
	procSetParent                  = user32.NewProc("SetParent")
	procSetLayeredWindowAttributes = user32.NewProc("SetLayeredWindowAttributes")
	procGetDesktopWindow           = user32.NewProc("GetDesktopWindow")
	procSendMessage                = user32.NewProc("SendMessageW")
	procEnumWindows                = user32.NewProc("EnumWindows")
	procGetClassNameW              = user32.NewProc("GetClassNameW")
)

// 设置为桌面子窗口 - 抵抗显示桌面，融入桌面环境
func SetupDesktopChildWidget(windowTitle string) error {
	// 1. 查找我们的窗口
	hwnd, _, _ := procFindWindowW.Call(
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(windowTitle))),
	)
	if hwnd == 0 {
		return fmt.Errorf("找不到窗口: %s", windowTitle)
	}

	// 2. 创建桌面 WorkerW 窗口（关键技术）
	// 首先找到 Progman 窗口
	progmanHwnd, _, _ := procFindWindowW.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("Progman"))),
		0,
	)
	if progmanHwnd == 0 {
		return fmt.Errorf("找不到 Progman 窗口")
	}

	// 发送特殊消息创建 WorkerW 窗口
	// 这个技巧会在桌面壁纸和桌面图标之间创建一个新的 WorkerW 窗口
	procSendMessage.Call(progmanHwnd, 0x052C, 0x0000000D, 0)
	procSendMessage.Call(progmanHwnd, 0x052C, 0x0000000D, 1)

	// 3. 查找新创建的 WorkerW 窗口
	var workerHwnd uintptr
	enumCallback := syscall.NewCallback(func(hwnd, lParam uintptr) uintptr {
		// 获取窗口类名
		classBuf := make([]uint16, 256)
		procGetClassNameW.Call(hwnd, uintptr(unsafe.Pointer(&classBuf[0])), 256)
		className := syscall.UTF16ToString(classBuf)

		// 查找 WorkerW 窗口，并检查它是否有 SHELLDLL_DefView 子窗口
		if className == "WorkerW" {
			shellView, _, _ := procFindWindowExW.Call(hwnd, 0,
				uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("SHELLDLL_DefView"))), 0)
			if shellView != 0 {
				// 找到了包含桌面图标的 WorkerW 窗口
				workerHwnd = hwnd
				return 0 // 停止枚举
			}
		}
		return 1 // 继续枚举
	})

	// 枚举所有顶级窗口
	procEnumWindows.Call(enumCallback, 0)

	// 4. 如果找到了 WorkerW 窗口，将我们的窗口设为其子窗口
	var parentHwnd uintptr
	if workerHwnd != 0 {
		parentHwnd = workerHwnd
		log.Printf("找到 WorkerW 窗口: 0x%x", workerHwnd)
	} else {
		// 备选方案：直接使用 Progman 作为父窗口
		parentHwnd = progmanHwnd
		log.Println("使用 Progman 作为父窗口")
	}

	// 5. 设置窗口样式（在设置父窗口之前）
	gwlExstyle := uintptr(^uint32(19)) // -20
	exStyle, _, _ := procGetWindowLongW.Call(hwnd, gwlExstyle)

	const (
		WS_EX_TOOLWINDOW = 0x00000080 // 不在 Alt+Tab 显示
		WS_EX_LAYERED    = 0x00080000 // 支持透明
		WS_EX_NOACTIVATE = 0x08000000 // 不激活窗口
		WS_EX_APPWINDOW  = 0x00040000 // 在任务栏显示（要移除）
	)

	// 设置桌面小部件样式
	newStyle := (exStyle | WS_EX_TOOLWINDOW | WS_EX_LAYERED | WS_EX_NOACTIVATE) &^ WS_EX_APPWINDOW
	procSetWindowLongW.Call(hwnd, gwlExstyle, newStyle)

	// 6. 设置为桌面的子窗口（关键步骤）
	result, _, _ := procSetParent.Call(hwnd, parentHwnd)
	if result == 0 {
		return fmt.Errorf("设置父窗口失败")
	}

	// 7. 设置透明度
	procSetLayeredWindowAttributes.Call(hwnd, 0, 230, 0x02) // 稍微透明

	// 8. 最终位置设置 - 不需要设置 TOPMOST，父子关系会处理层级
	const (
		SWP_NOMOVE       = 0x0002
		SWP_NOSIZE       = 0x0001
		SWP_NOZORDER     = 0x0004 // 保持 Z 顺序由父子关系决定
		SWP_NOACTIVATE   = 0x0010
		SWP_FRAMECHANGED = 0x0020
		SWP_SHOWWINDOW   = 0x0040
	)

	procSetWindowPos.Call(
		hwnd,
		0, // 不指定插入位置
		0, 0, 0, 0,
		uintptr(SWP_NOMOVE|SWP_NOSIZE|SWP_NOZORDER|SWP_NOACTIVATE|SWP_FRAMECHANGED|SWP_SHOWWINDOW),
	)

	return nil
}
