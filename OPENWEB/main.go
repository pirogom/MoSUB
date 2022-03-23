package main

import (
	"os/exec"
	"strconv"
)

var gWebPort = 30999

func main() {
	openWeb("http://127.0.0.1:" + strconv.Itoa(gWebPort))
}

/**
* 웹페이지 열기
**/
func openWeb(url string) {
	exec.Command("open", url).Start()
}
