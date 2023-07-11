package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"embed"
	_ "embed"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

//go:embed www/*
var wwwFS embed.FS

/**
* builtinTypesLower
**/
var builtinTypesLower = map[string]string{
	".css":  "text/css; charset=utf-8",
	".gif":  "image/gif",
	".htm":  "text/html; charset=utf-8",
	".html": "text/html; charset=utf-8",
	".jpeg": "image/jpeg",
	".jpg":  "image/jpeg",
	".js":   "text/javascript; charset=utf-8",
	".json": "application/json",
	".mjs":  "text/javascript; charset=utf-8",
	".pdf":  "application/pdf",
	".png":  "image/png",
	".svg":  "image/svg+xml",
	".wasm": "application/wasm",
	".webp": "image/webp",
	".xml":  "text/xml; charset=utf-8",
}

type oneSpellInfo struct {
	TrackID         int    `json:"TrackID"`
	GeneratorItemID int    `json:"GeneratorItemID"`
	EffectParamID   int    `json:"EffectParamID"`
	OrigStr         string `json:"origStr"`
	EditStr         string `json:"editStr"`
}

type spellCheckMsg struct {
	Message struct {
		Result struct {
			ErrataCount int    `json:"errata_count"`
			OriginHTML  string `json:"origin_html"`
			HTML        string `json:"html"`
			NotagHTML   string `json:"notag_html"`
		} `json:"result"`
	} `json:"message"`
}

type spellResult struct {
	OriginHTML string `json:"origin_html"`
	HTML       string `json:"html"`
}

type saveJSON []struct {
	Idx   int    `json:"idx"`
	Value string `json:"value"`
}

/**
*	forbiddenResponse
**/
func forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(403)
	w.Write([]byte(http.StatusText(403)))
}

/**
*	notFoundResponse
**/
func notFoundResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	w.Write([]byte(http.StatusText(404)))
}

/**
*	clientIsLocalHost
**/
func clientIsLocalHost(w http.ResponseWriter, r *http.Request) bool {
	sip := strings.Split(r.RemoteAddr, ":")

	if sip[0] != "127.0.0.1" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Bad Request"))
		return false
	}
	return true
}

/**
*	decodeGzip
**/
func decodeGzip(data []byte) ([]byte, error) {
	// Write gzipped data to the client
	gr, err := gzip.NewReader(bytes.NewBuffer(data))
	defer gr.Close()
	data, err = ioutil.ReadAll(gr)
	if err != nil {
		return nil, err
	}
	return data, nil
}

/**
*	getSpellCheck
**/
func getSpellCheck(origTxt string) ([]byte, error) {
	encodeTxt := url.QueryEscape(origTxt)
	reqURL := fmt.Sprintf("https://m.search.naver.com/p/csearch/ocontent/util/SpellerProxy?_callback=%s&q=%s&where=nexearch&color_blindness=0&_=%d", getJQCallback(), encodeTxt, getJQReqDummy())

	request, rerr := http.NewRequest("GET", reqURL, nil)

	if rerr != nil {
		return nil, rerr
	}

	request.Header.Set("accept", "application/json, text/plain, */*")
	request.Header.Set("accept-encoding", "gzip, deflate, br")
	request.Header.Set("accept-language", "en-US,en;q=0.9,ko;q=0.8")
	//	request.Header.Set("Cookie", getCookies())
	request.Header.Set("sec-fetch-mode", "cors")
	request.Header.Set("sec-fetch-site", "same-origin")
	request.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36")
	request.Header.Set("x-accept-language", "ko-KR")
	request.Header.Set("referer", "https://search.naver.com/search.naver?sm=top_hty&fbm=1&ie=utf8&query=%EB%84%A4%EC%9D%B4%EB%B2%84+%EB%A7%9E%EC%B6%A4%EB%B2%95+%EA%B2%80%EC%82%AC%EA%B8%B0")

	//

	client := &http.Client{}
	resp, respErr := client.Do(request)

	if respErr != nil {
		return nil, respErr
	}
	defer resp.Body.Close()

	contentEncoding := resp.Header.Get("content-encoding")

	htmlBuf, htmlErr := ioutil.ReadAll(resp.Body)

	if htmlErr == nil {
		if contentEncoding == "gzip" {
			decodeBuf, decodeErr := decodeGzip(htmlBuf)

			if decodeErr == nil {
				return decodeBuf, nil
			}
		} else {
			return htmlBuf, nil
		}
	}

	return nil, errors.New("response decode falied ")
}

/**
*	spellCheckProc
**/
func spellCheckProc(tm map[int]string, fname string) {

	origm := make(map[int]string)
	htmlm := make(map[int]string)

	for k, v := range tm {

		buf, bufErr := getSpellCheck(v)

		if bufErr == nil {
			rmjq := strings.Replace(string(buf), getJQCallback()+"(", "", -1)
			rmta := strings.Replace(rmjq, ");", "", -1)

			so := spellCheckMsg{}

			jerr := json.Unmarshal([]byte(rmta), &so)

			if jerr == nil {
				htmlm[k] = so.Message.Result.HTML
				origm[k] = so.Message.Result.OriginHTML
			}
		}

		time.Sleep(100 * time.Millisecond)
	}

	//	var tmpfbuf string
	var resJSON []spellResult

	for i := 0; i < len(origm); i++ {
		origsp := strings.Split(origm[i], "<br>")
		htmlsp := strings.Split(htmlm[i], "<br>")

		if len(origsp) != len(htmlsp) {
			fmt.Println("결과처리 오류1")
		} else {
			for ri := 0; ri < len(origsp); ri++ {
				var ospres spellResult

				if len(htmlsp[ri]) > 0 && len(origsp[ri]) > 0 {
					ospres.HTML = htmlsp[ri]
					ospres.OriginHTML = origsp[ri]

					resJSON = append(resJSON, ospres)
				}
			}
		}
	}

	resJSONBuf, resJSONBufErr := json.Marshal(&resJSON)

	if resJSONBufErr == nil {
		resfname := workFilepath(fname) + ".tmp"
		os.Remove(resfname)

		ioutil.WriteFile(resfname, resJSONBuf, 0644)
	}

	debug.FreeOSMemory()
}

/**
*	removePremiereXMLIndent
**/
func removeFinalcutXMLIndent(origXML string) string {

	var ret string

	scanner := bufio.NewScanner(strings.NewReader(origXML))
	for scanner.Scan() {

		txt := scanner.Text()

		len := len(txt)

		if len > 0 {
			if txt[len-1] != '>' {
				txt += "↵"
			}
		}
		ret += strings.TrimSpace(txt)
	}

	return ret

}

/**
*	uploadXML
**/
func uploadXML(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(4096)

	defer func() {
		r.Body.Close()
	}()

	if clientIsLocalHost(w, r) == false {
		return
	}

	var txtCnt int

	workType := r.PostFormValue("workType")
	upfile, upfileHeader, upfileErr := r.FormFile("uploadFile")

	if upfileErr == nil {

		ufr, ufrErr := ioutil.ReadAll(upfile)

		if ufrErr == nil {
			fname := upfileHeader.Filename
			os.Remove(workFilepath(fname))
			ioutil.WriteFile(workFilepath(fname), ufr, 0644)

			xmlBuf, xmlBufErr := getUTF8bufferFile(workFilepath(fname))

			if xmlBufErr == nil {
				switch workType {
				case "app": // 프리미어 프로
					txtCnt = premiereXMLProc(string(xmlBuf), fname)
					break
				case "fcp": // 파이널 컷
					txtCnt = finalcutXMLProc(string(xmlBuf), fname)
					break
				case "smi": // smi
					txtCnt = smiSUBProc(string(xmlBuf), fname)
					break
				case "srt": // srt
					txtCnt = srtSUBProc(string(xmlBuf), fname)
					break
				}
			}

		}

	}

	rd := make(map[string]string)
	rd["TXTCNT"] = fmt.Sprintf("%d", txtCnt)

	jsonData, err := json.Marshal(&rd)

	if err != nil {
		return
	}

	w.Write(jsonData)

	debug.FreeOSMemory()
}

/**
* getResult
**/
func getResult(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	defer func() {
		r.Body.Close()
	}()

	if clientIsLocalHost(w, r) == false {
		return
	}

	resfname := r.Form.Get("resfile")

	readBuf, readErr := ioutil.ReadFile(workFilepath(resfname))

	if readErr != nil {
		return
	}

	w.Write(readBuf)
}

/**
*	checkAlive
**/
func checkAlive(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	defer func() {
		r.Body.Close()
	}()

}

/**
*	saveData
**/
func saveData(w http.ResponseWriter, r *http.Request) {
	if clientIsLocalHost(w, r) == false {
		return
	}

	r.ParseForm()

	defer func() {
		r.Body.Close()
	}()

	fname := r.Form.Get("fname")
	data := r.Form.Get("saveData")
	workType := r.Form.Get("workType")

	fnameBase, _, fnameErr := splitFilename(fname)

	if fnameErr != nil {
		fnameBase = fname
	}

	var ret map[string]string

	switch workType {
	case "app":
		savefname := fmt.Sprintf("%s_checked.xml", fnameBase)
		ret = savePremiereData(data, fname, savefname)
		break
	case "fcp":
		savefname := fmt.Sprintf("%s_checked.fcpxml", fnameBase)
		ret = saveFinalcutData(data, fname, savefname)
		break
	case "smi":
		savefname := fmt.Sprintf("%s_checked.smi", fnameBase)
		ret = saveSmiSUBData(data, fname, savefname)
		break
	case "srt":
		savefname := fmt.Sprintf("%s_checked.srt", fnameBase)
		ret = saveSrtSUBData(data, fname, savefname)
		break
	}

	jsonData, err := json.Marshal(&ret)

	if err != nil {
		return
	}

	w.Write(jsonData)
	debug.FreeOSMemory()
}

/**
*	getResFile
**/
func getResFile(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	defer func() {
		r.Body.Close()
	}()

	if clientIsLocalHost(w, r) == false {
		return
	}

	///getResFile/test.xml
	ruri, _ := url.QueryUnescape(r.RequestURI)

	urisp := strings.Split(ruri, "/")

	if len(urisp) != 3 {
		http.Error(w, "bad request", 404)
	} else {
		http.ServeFile(w, r, workFilepath(urisp[2]))
	}
}

/**
*	getVersion
**/
func getVersion(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	defer func() {
		r.Body.Close()
	}()

	ret := make(map[string]string)
	ret["ver"] = gVer

	jsonData, err := json.Marshal(&ret)

	if err != nil {
		return
	}

	w.Write(jsonData)
}

/**
*	shutdownServ
**/
func shutdownServ(w http.ResponseWriter, r *http.Request) {
	defer func() {
		r.Body.Close()
	}()

	if clientIsLocalHost(w, r) == false {
		return
	}

	go func() {
		openWeb("https://modu-print.com/category/%ec%97%85%eb%8d%b0%ec%9d%b4%ed%8a%b8/%eb%aa%a8%eb%91%90%ec%9d%98%ec%9e%90%eb%a7%89/")
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()
}

/**
*	webPageProc
**/
func webPageProc(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Cache-control", "no-cache")

	// Parsing From Data
	r.ParseMultipartForm(4096)
	defer func() {
		r.Body.Close()
	}()
	//

	// 불손한 시도 ㅋㅋㅋ
	if strings.Index(r.URL.Path, "..\\") != -1 ||
		strings.Index(r.URL.Path, "../") != -1 ||
		strings.Index(r.URL.Path, "<script") != -1 {
		w.WriteHeader(404)
		w.Write([]byte(http.StatusText(404)))
		return
	}

	if strings.Index(r.URL.Path, "/getResFile/") == 0 {
		getResFile(w, r)
		return
	}

	switch r.URL.Path {
	case "/uploadXML":
		uploadXML(w, r)
	case "/getResult":
		getResult(w, r)
	case "/checkAlive":
		checkAlive(w, r)
	case "/saveData":
		saveData(w, r)
	case "/getVer":
		getVersion(w, r)
	case "/shutdown":
		shutdownServ(w, r)
	default:
		if !gIsDebug {
			fpath := "www" + r.URL.Path

			if fpath == "www/" {
				fpath += "index.html"
			}

			data, dataErr := wwwFS.ReadFile(fpath)

			if dataErr != nil {
				notFoundResponse(w, r)
				return
			}

			ext := path.Ext(fpath)
			w.Header().Set("Content-Type", getMimeType(ext))
			w.Write(data)
		} else { // deubg mode
			//			revURL := strings.Replace(r.URL.Path, "/", pathSlash, -1)
			//			http.ServeFile(w, r, ".\\www"+revURL)
			http.ServeFile(w, r, "./www"+r.URL.Path)
		}
	}
}

/**
*	webserv
**/
func webserv(hostname string, servPort int) {

	http.HandleFunc("/", webPageProc)

	servAddr := fmt.Sprintf("%s:%d", hostname, servPort)

	go func() {

		time.Sleep(1 * time.Second)

		openWeb("http://127.0.0.1:" + strconv.Itoa(gWebPort))
	}()

	http.ListenAndServe(servAddr, nil)
}

/**
*	xmlEscape
**/
func xmlEscape(value string) string {
	escaped := &bytes.Buffer{}

	if err := xml.EscapeText(escaped, []byte(value)); err != nil {
		panic(err)
	}

	return escaped.String()
}

/**
*	getMimeType
**/
func getMimeType(fext string) string {
	ext := strings.ToLower(fext)

	if builtinTypesLower[ext] == "" {
		return "application/octet-stream"
	}
	return builtinTypesLower[ext]
}
