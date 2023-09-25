//go:build darwin || !windows
// +build darwin !windows

package main

import (
	"fmt"
	"os"
	"os/exec"
)

var gVer = "1.4"
var gDiceUtil *dice

func main() {

	gDiceUtil = newDice()

	fixWorkingDirectory()

	if len(os.Args) > 1 {
		for ai := 1; ai < len(os.Args); ai++ {
			if os.Args[ai] == "--dev-web" {
				gIsDebug = true
			}
		}
	}

	initJQuery()

	webserv("127.0.0.1", gWebPort)
}

/**
* 웹페이지 열기
**/
func openWeb(url string) {
	err := exec.Command("open", url).Start()

	if err != nil {
		fmt.Println(err.Error())
	}
}
