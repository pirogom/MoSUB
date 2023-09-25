//go:build !darwin || windows
// +build !darwin windows

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/gonutz/w32"
	"github.com/rodolfoag/gow32"
	"github.com/tadvi/systray"
)

var gDiceUtil *dice

const appIconResID = 8

// 콘솔 윈도우 숨기기
func hideConsole() {
	consoleWnd := w32.GetConsoleWindow()
	if consoleWnd == 0 {
		return // no console attached
	}

	// If this application is the process that created the console window, then
	// this program was not compiled with the -H=windowsgui flag and on start-up
	// it created a console along with the main application window. In this case
	// hide the console window.
	// See
	// http://stackoverflow.com/questions/9009333/how-to-check-if-the-program-is-run-from-a-console
	_, consoleProcID := w32.GetWindowThreadProcessId(consoleWnd)
	if w32.GetCurrentProcessId() == consoleProcID {
		w32.ShowWindowAsync(consoleWnd, w32.SW_HIDE)
	}
}

/**
*	checkMutex
**/
func checkMutex() bool {

	_, err := gow32.CreateMutex("mtx_running_mosub")

	if err != nil {

		if int(err.(syscall.Errno)) == gow32.ERROR_ALREADY_EXISTS {
			return false
		}
	}

	return true
}

func main() {
	gDiceUtil = newDice()
	fixWorkingDirectory()

	// 중복실행 방지
	if checkMutex() == false {
		return
	}

	// 콘솔 가리기
	hideConsole()

	if len(os.Args) > 1 {
		for ai := 1; ai < len(os.Args); ai++ {
			if os.Args[ai] == "--dev-web" {
				gIsDebug = true
			}
		}
	}

	initJQuery()

	go webserv("127.0.0.1", gWebPort)

	//
	tray, terr := systray.New()

	if terr != nil {
		panic(terr)
	}

	terr = tray.Show(appIconResID, "모두의 자막")
	if terr != nil {
		panic(terr)
	}

	tray.AppendMenu("피로곰's 모두의 자막", func() {
		openWeb("https://modu-print.tistory.com")
	})
	tray.AppendSeparator()
	tray.AppendMenu("작업페이지 열기", func() {
		openWeb("http://127.0.0.1:" + strconv.Itoa(gWebPort))
	})
	tray.AppendMenu("종료", func() {
		openWeb("https://modu-print.com/category/%ec%97%85%eb%8d%b0%ec%9d%b4%ed%8a%b8/%eb%aa%a8%eb%91%90%ec%9d%98%ec%9e%90%eb%a7%89/")
		tray.Stop()
		os.Exit(0)
	})

	err := tray.Run()

	if err != nil {
		fmt.Println(err.Error())
	}
	//

}

/**
* 웹페이지 열기
**/
func openWeb(url string) {
	r := strings.NewReplacer("&", "^&")
	stripURL := r.Replace(url)

	err := exec.Command("cmd.exe", "/C", "start", stripURL).Start()

	if err != nil {
		fmt.Println(err.Error())
	}
}
