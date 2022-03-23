package main

import (
	"fmt"
)

var jqCallbackHead string
var jqCallbackTail string
var jqCallbackTailNum int64

func initJQuery() {

	diceUtil := newDice()

	jqCallbackHead = fmt.Sprintf("%d%d%d%d%d", diceUtil.randNumber(10000, 19999), diceUtil.randNumber(10000, 19999), diceUtil.randNumber(10000, 19999), diceUtil.randNumber(10000, 19999), diceUtil.randNumber(100, 199))
	jqCallbackTailNum = diceUtil.randNumber(1000, 9999)
	jqCallbackTail = fmt.Sprintf("%d%d%d%d", diceUtil.randNumber(100, 199), diceUtil.randNumber(100, 199), diceUtil.randNumber(100, 199), jqCallbackTailNum)
}

func getJQCallback() string {
	return fmt.Sprintf("jQuery%s_%s", jqCallbackHead, jqCallbackTail)
}

func getJQReqDummy() int64 {
	ret := jqCallbackTailNum
	jqCallbackTailNum++
	return ret
}
