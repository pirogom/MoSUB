package main

import (
	"bytes"
	"errors"
	"io/ioutil"

	"golang.org/x/net/html/charset"
)

//   "utf-8":               utf8,
//   "utf8":                utf8,
//   "utf-16be":            utf16be,
//   "utf-16le":            utf16le,
//   "euc-kr":              euckr,
func convToUTF8(strBytes []byte, origEncoding string) []byte {
	byteReader := bytes.NewReader(strBytes)
	reader, _ := charset.NewReaderLabel(origEncoding, byteReader)
	strBytes, _ = ioutil.ReadAll(reader)
	return strBytes
}

/**
*   getUTF8buffer
**/
func getUTF8buffer(rs []byte) []byte {
	if len(string(rs)) < 2 {
		return nil
	}

	if rs[0] == 0xEF && rs[1] == 0xBB && rs[2] == 0xBF { // UTF8
		return rs // 변환 없음
	} else if rs[0] == 0xFE && rs[1] == 0xFF { // UTF-16 Big Endian
		return convToUTF8(rs, "utf-16be")
	} else if rs[0] == 0xFF && rs[1] == 0xFE { // UTF-16 Little Endian
		return convToUTF8(rs, "utf-16le")
	}
	return convToUTF8(rs, "euc-kr")
}

/**
*   getUTF8bufferFile
**/
func getUTF8bufferFile(fname string) ([]byte, error) {
	buf, err := ioutil.ReadFile(fname)

	if err != nil {
		return nil, err
	}
	ret := getUTF8buffer(buf)

	if ret == nil {
		return nil, errors.New("empty file")
	}

	return ret, nil
}
