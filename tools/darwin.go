//go:build darwin

package tools

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa

#import <Cocoa/Cocoa.h>

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
                    // 设置为面板类型窗口
                    [win setStyleMask:NSWindowStyleMaskNonactivatingPanel];
                    [win setLevel:kCGDesktopWindowLevel];
                    [win setCollectionBehavior:NSWindowCollectionBehaviorCanJoinAllSpaces];

                    // 启用鼠标事件
                    [win setAcceptsMouseMovedEvents:YES];
                    [win setIgnoresMouseEvents:NO];
                    [win setMovableByWindowBackground:YES];

                    // 设置窗口属性
                    [win setOpaque:NO];
                    [win setHasShadow:YES];
                    [win setAlphaValue:1.0];
                    [win setBackgroundColor:[NSColor clearColor]];

                    // 确保窗口可见
                    [win orderFront:nil];
                    [NSApp activateIgnoringOtherApps:YES];

                    NSLog(@"Window properties set completed");
                });
            } else {
                NSLog(@"Window is nil");
            }
        } else {
            NSLog(@"No windows found");
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
import "fmt"

func SetupDesktopChildWidget(windowTitle string) error {
	fmt.Println("SetupDesktopChildWidget called with title:", windowTitle)
	C.setWindowLevelToDesktop()
	fmt.Println("SetupDesktopChildWidget completed")
	return nil
}
