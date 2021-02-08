package network

import (
	"os"
	"runtime"
)

func IsSynologyDSM() bool  {
	if runtime.GOOS != "linux" {
		return false
	}
	file, err := os.Stat("/usr/syno")
	if err == nil && file.IsDir() {
		return true
	}
	_, err = os.Stat("/proc/syno_platform")
	if err == nil {
		return true
	}
	_, err = os.Stat("usr/syno_cpu_arch")
	if err == nil {
		return true
	}
	_, err = os.Stat("usr/synobios")
	if err == nil {
		return true
	}			
	return false
}