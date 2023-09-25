package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"

	strip "github.com/grokify/html-strip-tags-go"
)

type srtTags struct {
	Idx   string
	Tline string
	Text  string
}

type srtParser struct {
	subs []srtTags
}

/**
*	newSRTParser
**/
func newSRTParser(strBuf string) (*srtParser, error) {
	nd := srtParser{}

	scanner := bufio.NewScanner(strings.NewReader(strBuf))

	for {
		ot := srtTags{}

		// 번호
		if !scanner.Scan() {
			break
		}

		ot.Idx = scanner.Text()

		// 타임 정보
		if !scanner.Scan() {
			return nil, errors.New("time info read error")
		}

		ot.Tline = scanner.Text()

		for scanner.Scan() {
			if len(scanner.Text()) == 0 {
				break
			}

			if len(ot.Text) > 0 {
				ot.Text += "<br>" + scanner.Text()
			} else {
				ot.Text += scanner.Text()
			}
		}

		ot.Text = strings.Replace(ot.Text, "<br>", "↵", -1)
		ot.Text = strip.StripTags(ot.Text)

		nd.subs = append(nd.subs, ot)
	}

	if len(nd.subs) == 0 {
		return nil, errors.New("data is empty")
	}

	return &nd, nil
}

/**
*	getResult
**/
func (d *srtParser) getResult() string {

	var srtStr string

	for i := 0; i < len(d.subs); i++ {
		txt := strings.Replace(d.subs[i].Text, "↵", "<br>", -1)
		srtStr += d.subs[i].Idx + "\r\n" + d.subs[i].Tline + "\r\n" + txt + "\r\n\r\n"
	}

	return srtStr
}

/**
*	SRT  업로드 파일 처리
**/
func srtSUBProc(srtString string, fname string) int {
	smip, smipErr := newSRTParser(srtString)

	txtBlockMap := make(map[int]string)
	var txtBlockCnt int
	var oneTxtBlock string
	var txtCnt int

	if smipErr == nil {

		for si := 0; si < len(smip.subs); si++ {
			txtCnt++

			if len(oneTxtBlock)+len(smip.subs[si].Text)+len("\r\n") < SpellCheckLimit {
				oneTxtBlock += smip.subs[si].Text + "\r\n"
			} else {
				txtBlockMap[txtBlockCnt] = oneTxtBlock
				txtBlockCnt++

				oneTxtBlock = smip.subs[si].Text + "\r\n"
			}
		}

		txtBlockMap[txtBlockCnt] = oneTxtBlock

		if txtCnt > 0 {
			go spellCheckProc(txtBlockMap, fname)
		}
	}

	return txtCnt
}

/**
*	saveSrtSUBData
**/
func saveSrtSUBData(data string, fname string, savefname string) map[string]string {
	ret := make(map[string]string)
	ret["result"] = "ERR"

	var lj saveJSON

	jerr := json.Unmarshal([]byte(data), &lj)

	if jerr != nil {
		ret["result"] = "ERR"
	} else {
		//
		oldSrtData, oldSrtDataErr := ioutil.ReadFile(workFilepath(fname))

		if oldSrtDataErr != nil {
			ret["result"] = "ERR"
		} else {

			smip, smipErr := newSRTParser(string(oldSrtData))

			var txtCnt int

			if smipErr == nil {

				for si := 0; si < len(smip.subs); si++ {
					if txtCnt < len(lj) {
						smip.subs[si].Text = lj[txtCnt].Value
					}
					txtCnt++
				}

				werr := ioutil.WriteFile(workFilepath(savefname), []byte(smip.getResult()), 0644)

				if werr != nil {
					ret["result"] = "ERR"
				} else {
					ret["result"] = "OK"
					ret["filename"] = savefname
				}
			}
			//
		}
	}
	return ret
}
