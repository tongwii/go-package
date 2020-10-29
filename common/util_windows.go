package common

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"
)

/**
重启服务
只支持windows服务，不支持直接exe启动
*/
func RestartService() error {
	d, err := os.Getwd() // 获取当前工作目录
	if err != nil {
		return err
	}
	// 服务方式启动默认工作目录
	if strings.ToUpper(d) == "C:\\WINDOWS\\SYSTEM32" {
		Logger.Info("RestartService Restart")
		exe := os.Args[0]
		s := []string{"cmd.exe", "/C", "start", "/B", exe, "-s", "restart"}
		cmd_instance := exec.Command(s[0], s[1:]...)
		cmd_instance.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // 隐藏窗口
		out, err := cmd_instance.Output()
		if err != nil {
			Logger.Errorf("RestartService Output: %v", err)
			return err
		} else {
			out, _ = ConvertToUTF8(out)
			println(out)
		}
	} else {
		Logger.Info("RestartService Exit")
		os.Exit(0) // 直接结束当前进程
	}
	return nil
}

/**
重启服务(改,docker和mass在用)
只支持windows服务，不支持直接exe启动
*/
func RestartServer() error {
	d, err := os.Getwd() // 获取当前工作目录
	if err != nil {
		return err
	}
	// 服务方式启动默认工作目录
	if strings.ToUpper(d) == "C:\\WINDOWS\\SYSTEM32" {
		Logger.Info("RestartServer Restart")
		exe := os.Args[0]
		cmd_instance := exec.Command("cmd", "/C", fmt.Sprintf("%s", exe), "-s", "restart")
		cmd_instance.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // 隐藏窗口
		Logger.Infof(fmt.Sprintf("Args: %v", cmd_instance.Args))
		out, err := cmd_instance.Output()
		if err != nil {
			Logger.Errorf("RestartServer Output Error: %v", err)
			Logger.Infof("RestartServer Output: %s", string(out))
			return err
		} else {
			out, _ = ConvertToUTF8(out)
			Logger.Infof("RestartServer Output: %s", string(out))
		}
	} else {
		Logger.Info("RestartServer Exit")
		os.Exit(0) // 直接结束当前进程
	}
	return nil
}

/*
	获取目录所在磁盘空间
*/
func GetDiskSpace(path string) (availableBytes, totalBytes, freeBytes int64, err error) {
	h := syscall.MustLoadDLL("kernel32.dll")
	c := h.MustFindProc("GetDiskFreeSpaceExW")

	dir, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		Logger.Errorf("GetDiskSpace UTF16PtrFromString Error: %v", err)
		return
	}
	_, _, err = c.Call(uintptr(unsafe.Pointer(dir)),
		uintptr(unsafe.Pointer(&freeBytes)),
		uintptr(unsafe.Pointer(&totalBytes)),
		uintptr(unsafe.Pointer(&availableBytes)))
	if err != nil {
		Logger.Debugf("GetDiskSpace Syscall6 Error: %v", err)
		//return
	}
	Logger.Debugf("GetDiskSpace Available: %dMB", availableBytes/1024/1024.0)
	Logger.Debugf("GetDiskSpace Total: %dMB", totalBytes/1024/1024.0)
	Logger.Debugf("GetDiskSpace Free: %dMB", freeBytes/1024/1024.0)
	return
}

//func GetDiskSpace(path string) (availableBytes, totalBytes, freeBytes int64, err error) {
//	kernel32, err := syscall.LoadLibrary("Kernel32.dll")
//	if err != nil {
//		Logger.Errorf("GetDiskSpace LoadLibrary Error: %v", err)
//		return
//	}
//	defer func() { _ = syscall.FreeLibrary(kernel32) }()
//	GetDiskFreeSpaceEx, err := syscall.GetProcAddress(syscall.Handle(kernel32), "GetDiskFreeSpaceExW")
//	if err != nil {
//		Logger.Errorf("GetDiskSpace GetProcAddress Error: %v", err)
//		return
//	}
//	dir, err := syscall.UTF16PtrFromString(path)
//	if err != nil {
//		Logger.Errorf("GetDiskSpace UTF16PtrFromString Error: %v", err)
//		return
//	}
//	_, _, err = syscall.Syscall6(uintptr(GetDiskFreeSpaceEx), 4,
//		uintptr(unsafe.Pointer(dir)),
//		uintptr(unsafe.Pointer(&availableBytes)),
//		uintptr(unsafe.Pointer(&totalBytes)),
//		uintptr(unsafe.Pointer(&freeBytes)), 0, 0)
//	if err != nil {
//		Logger.Errorf("GetDiskSpace Syscall6 Error: %v", err)
//		return
//	}
//	Logger.Debugf("GetDiskSpace Available: %dMB", availableBytes/1024/1024.0)
//	Logger.Debugf("GetDiskSpace Total: %dMB", totalBytes/1024/1024.0)
//	Logger.Debugf("GetDiskSpace Free: %dMB", freeBytes/1024/1024.0)
//	return
//}
