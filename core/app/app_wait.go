package app

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

//signalsListen
func (app *Application) signalsListen(sigs chan os.Signal) {

EXIT:
	switch runtime.GOOS {
	case "windows":
		for {
			signal.Notify(sigs, syscall.SIGQUIT,
				syscall.SIGTERM,
				syscall.SIGINT)
			switch <-sigs {
			case syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT:
				// syscall.SIGQUIT 关闭服务
				break EXIT
			}
		}
	default:
		for {
			signal.Notify(sigs, syscall.SIGQUIT,
				syscall.SIGTERM,
				syscall.SIGINT,
				syscall.SIGUSR1,
				syscall.SIGUSR2)

			switch <-sigs {
			case syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT:
				// SIGQUIT 关闭服务
				fmt.Println(" sig:quit")
				break EXIT
			case syscall.SIGUSR1:
				fmt.Println(" sig:SIGUSR1")
			case syscall.SIGUSR2:
				// 关闭服务
				break EXIT
			}
		}
	}

}
