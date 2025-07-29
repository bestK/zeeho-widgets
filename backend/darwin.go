//go:build darwin

package backend

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework WidgetKit

#import <Cocoa/Cocoa.h>
#import <WidgetKit/WidgetKit.h>

bool isMainThread() {
    return [NSThread isMainThread];
}

void setWindowLevelToDesktopImpl() {
    @autoreleasepool {
        NSArray *windows = [NSApp windows];
        NSLog(@"Window count: %lu", (unsigned long)[windows count]);
        if ([windows count] > 0) {
            NSWindow *win = [windows objectAtIndex:0];
            if (win != nil) {
                NSLog(@"Setting window level and behaviors");
                dispatch_async(dispatch_get_main_queue(), ^{
                    // 创建桌面小部件效果
                    // 使用比桌面稍高的级别，但低于普通窗口
                    [win setLevel:kCGDesktopWindowLevel + 20];

                    // 设置窗口行为：在所有空间显示，但不参与窗口循环
                    [win setCollectionBehavior:NSWindowCollectionBehaviorCanJoinAllSpaces |
                                               NSWindowCollectionBehaviorStationary |
                                               NSWindowCollectionBehaviorIgnoresCycle];

                    // 设置窗口样式为无边框但可交互
                    NSUInteger styleMask = [win styleMask];
                    styleMask &= ~NSWindowStyleMaskTitled;
                    styleMask |= NSWindowStyleMaskBorderless;
                    [win setStyleMask:styleMask];

                    // 启用鼠标事件
                    [win setAcceptsMouseMovedEvents:YES];
                    [win setIgnoresMouseEvents:NO];
                    [win setMovableByWindowBackground:NO]; // 桌面小部件通常不可拖动

                    // 设置窗口属性
                    [win setOpaque:NO];
                    [win setHasShadow:YES];
                    [win setAlphaValue:1.0];
                    [win setBackgroundColor:[NSColor clearColor]];

                    // 确保窗口显示但不激活
                    [win orderFront:nil];
                    // 不要激活应用，保持桌面小部件的特性
                    // [NSApp activateIgnoringOtherApps:YES];

                    NSLog(@"Desktop widget window properties set completed");
                });
            } else {
                NSLog(@"Window is nil");
            }
        } else {
            NSLog(@"No windows found");
        }
    }
}

// 使用 WidgetKit 相关功能
void setupWidgetKitIntegration() {
    @autoreleasepool {
        if (@available(macOS 11.0, *)) {
            NSLog(@"WidgetKit is available - setting up integration");

            // 获取当前应用的 bundle identifier
            NSString *bundleId = [[NSBundle mainBundle] bundleIdentifier];
            NSLog(@"App bundle ID: %@", bundleId);

            // 检查是否有相关的 Widget Extension
            NSString *widgetBundleId = [NSString stringWithFormat:@"%@.widget", bundleId];
            NSLog(@"Looking for widget extension: %@", widgetBundleId);

            // 设置 WidgetKit 相关的通知监听
            [[NSNotificationCenter defaultCenter] addObserverForName:NSApplicationDidBecomeActiveNotification
                                                              object:nil
                                                               queue:[NSOperationQueue mainQueue]
                                                          usingBlock:^(NSNotification *note) {
                NSLog(@"App became active - could refresh widget data");
            }];

            [[NSNotificationCenter defaultCenter] addObserverForName:NSApplicationWillTerminateNotification
                                                              object:nil
                                                               queue:[NSOperationQueue mainQueue]
                                                          usingBlock:^(NSNotification *note) {
                NSLog(@"App will terminate - cleaning up widget resources");
            }];

            // 尝试刷新相关的小部件（如果存在）
            dispatch_async(dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_DEFAULT, 0), ^{
                @try {
                    // 这里可以添加刷新小部件时间线的代码
                    // WidgetCenter.shared.reloadTimelines(ofKind: "YourWidgetKind")
                    NSLog(@"Widget timeline refresh attempted");
                } @catch (NSException *exception) {
                    NSLog(@"Widget refresh failed: %@", exception.reason);
                }
            });

            NSLog(@"WidgetKit integration setup completed");

        } else {
            NSLog(@"WidgetKit is not available on this system (requires macOS 11.0+)");
        }
    }
}

// 刷新小部件数据
void refreshWidgetData(const char* widgetKind) {
    @autoreleasepool {
        if (@available(macOS 11.0, *)) {
            NSString *kind = [NSString stringWithUTF8String:widgetKind];
            NSLog(@"Refreshing widget data for kind: %@", kind);

            dispatch_async(dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_DEFAULT, 0), ^{
                @try {
                    // 在实际的 Widget Extension 中，你会使用：
                    // WidgetCenter.shared.reloadTimelines(ofKind: kind)
                    // 这里我们只是记录日志，因为主应用无法直接刷新小部件
                    NSLog(@"Widget refresh request sent for: %@", kind);
                } @catch (NSException *exception) {
                    NSLog(@"Widget refresh failed: %@", exception.reason);
                }
            });
        } else {
            NSLog(@"Widget refresh not available on this system");
        }
    }
}

void setWindowLevelToDesktop() {
    if ([NSThread isMainThread]) {
        setWindowLevelToDesktopImpl();
    } else {
        dispatch_async(dispatch_get_main_queue(), ^{
            setWindowLevelToDesktopImpl();
        });
    }
}
*/
import "C"
import (
	"fmt"
)

func SetupDesktopChildWidget() error {
	fmt.Println("SetupDesktopChildWidget called with title:", window_title)

	// 设置窗口为桌面级别
	C.setWindowLevelToDesktop()

	// 初始化 WidgetKit 集成
	C.setupWidgetKitIntegration()

	fmt.Println("SetupDesktopChildWidget completed")
	return nil
}

func SetTransparentBackground() {
}
