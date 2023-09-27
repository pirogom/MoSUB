package main

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	strip "github.com/grokify/html-strip-tags-go"
)

type smiSync struct {
	Sync string
	PTag string
	Text string
}

type smiParser struct {
	header string
	body   string
	footer string
	subs   []smiSync
}

/**
*	newSMIParser
**/
func newSMIParser(strBuf string) (*smiParser, error) {
	nd := smiParser{}

	bodyIdx := strings.Index(strBuf, "<BODY>")
	if bodyIdx == -1 {
		bodyIdx = strings.Index(strBuf, "<body>")
	}

	if bodyIdx == -1 {
		return nil, errors.New("<body> not found")
	}

	nd.header = strBuf[0 : bodyIdx+len("<BODY>")]
	//

	bodyEndIdx := strings.Index(strBuf, "</BODY>")
	if bodyEndIdx == -1 {
		bodyEndIdx = strings.Index(strBuf, "</body>")
	}

	if bodyEndIdx == -1 {
		return nil, errors.New("</body> not found")
	}

	nd.body = strBuf[bodyIdx+len("<BODY>") : bodyEndIdx]

	nd.footer = strBuf[bodyEndIdx:]
	//

	pbErr := nd.parsingBody()

	if pbErr != nil {
		return nil, pbErr
	}

	return &nd, nil
}

/**
*	parsingBody
**/
func (d *smiParser) parsingBody() error {

	ts := d.body

	for {
		ot := smiSync{}

		syncStart := strings.Index(ts, "<SYNC")

		if syncStart == -1 {
			syncStart = strings.Index(ts, "<sync")
			if syncStart == -1 {
				break
			}
		}

		ts = ts[syncStart:]

		syncEnd := strings.Index(ts, ">")
		if syncEnd == -1 {
			return errors.New("Bad Syntex sync closer error")
		}
		ot.Sync = ts[:syncEnd+1]

		// PTag
		ts = ts[syncEnd+1:]

		pStart := strings.Index(ts, "<P")

		if pStart == -1 {
			pStart = strings.Index(ts, "<p")
			if pStart == -1 {
				return errors.New("Bad Syntex P tag not found")
			}
		}

		ts = ts[pStart:]

		pEnd := strings.Index(ts, ">")

		if pEnd == -1 {
			return errors.New("Bad Syntex p tag close error")
		}

		ot.PTag = ts[:pEnd+1]

		// Text
		ts = ts[pEnd+1:]

		nextSync := strings.Index(ts, "<SYNC")
		if nextSync == -1 {
			nextSync = strings.Index(ts, "<sync")
			if nextSync == -1 {
				ot.Text = ts
				d.subs = append(d.subs, ot)
				break
			}
		}

		ot.Text = ts[:nextSync]

		ot.Text = strings.Replace(ot.Text, "<br>", "↵", -1)
		ot.Text = strip.StripTags(ot.Text)
		ot.Text = strings.Replace(ot.Text, "\n", "", -1)
		ot.Text = strings.Replace(ot.Text, "\r", "", -1)

		d.subs = append(d.subs, ot)
	}

	return nil
}

/**
*	getResult
**/
func (d *smiParser) getResult() string {

	smiStr := d.header

	for i := 0; i < len(d.subs); i++ {
		txt := strings.Replace(d.subs[i].Text, "↵", "<br>", -1)
		smiStr += d.subs[i].Sync + d.subs[i].PTag + txt + "\r\n"
	}

	smiStr += d.footer

	return smiStr
}

/**
*	SMI  업로드 파일 처리
**/
func smiSUBProc(smiString string, fname string, passportKey string) int {
	smip, smipErr := newSMIParser(smiString)

	txtBlockMap := make(map[int]string)
	var txtBlockCnt int
	var oneTxtBlock string
	var txtCnt int

	if smipErr == nil {

		for si := 0; si < len(smip.subs); si++ {
			if strings.Index(smip.subs[si].Text, "&nbsp") == -1 {
				txtCnt++

				if len(oneTxtBlock)+len(smip.subs[si].Text)+len("\r\n") < SpellCheckLimit {
					oneTxtBlock += smip.subs[si].Text + "\r\n"
				} else {
					txtBlockMap[txtBlockCnt] = oneTxtBlock
					txtBlockCnt++

					oneTxtBlock = smip.subs[si].Text + "\r\n"
				}
			}
		}

		txtBlockMap[txtBlockCnt] = oneTxtBlock

		if txtCnt > 0 {
			go spellCheckProc(txtBlockMap, fname, passportKey)
		}
	}

	return txtCnt
}

/**
*	saveSmiSUBData
**/
func saveSmiSUBData(data string, fname string, savefname string) map[string]string {
	ret := make(map[string]string)
	ret["result"] = "ERR"

	var lj saveJSON

	jerr := json.Unmarshal([]byte(data), &lj)

	if jerr != nil {
		ret["result"] = "ERR"
	} else {
		//
		oldSmiData, oldSmiDataErr := os.ReadFile(workFilepath(fname))

		if oldSmiDataErr != nil {
			ret["result"] = "ERR"
		} else {

			smip, smipErr := newSMIParser(string(oldSmiData))

			var txtCnt int

			if smipErr == nil {

				for si := 0; si < len(smip.subs); si++ {
					if strings.Index(smip.subs[si].Text, "&nbsp") == -1 {
						if txtCnt < len(lj) {
							smip.subs[si].Text = lj[txtCnt].Value
						}
						txtCnt++
					}
				}

				werr := os.WriteFile(workFilepath(savefname), []byte(smip.getResult()), 0644)

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
