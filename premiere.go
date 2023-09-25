package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"strings"
)

type premiereXML struct {
	XMLName  xml.Name `xml:"xmeml"`
	Text     string   `xml:",chardata"`
	Version  string   `xml:"version,attr"`
	Sequence struct {
		Text     string `xml:",chardata"`
		ID       string `xml:"id,attr"`
		Name     string `xml:"name"`
		Duration string `xml:"duration"`
		Rate     struct {
			Text     string `xml:",chardata"`
			Timebase string `xml:"timebase"`
			Ntsc     string `xml:"ntsc"`
		} `xml:"rate"`
		Media struct {
			Text  string `xml:",chardata"`
			Video struct {
				Text   string `xml:",chardata"`
				Format struct {
					Text                  string `xml:",chardata"`
					Samplecharacteristics struct {
						Text             string `xml:",chardata"`
						Width            string `xml:"width"`
						Height           string `xml:"height"`
						Anamorphic       string `xml:"anamorphic"`
						Pixelaspectratio string `xml:"pixelaspectratio"`
						Fielddominance   string `xml:"fielddominance"`
					} `xml:"samplecharacteristics"`
				} `xml:"format"`
				Track []struct {
					Text     string `xml:",chardata"`
					Clipitem []struct {
						Text    string `xml:",chardata"`
						Name    string `xml:"name"`
						Enabled string `xml:"enabled"`
						Rate    struct {
							Text     string `xml:",chardata"`
							Timebase string `xml:"timebase"`
							Ntsc     string `xml:"ntsc"`
						} `xml:"rate"`
						In    string `xml:"in"`
						Out   string `xml:"out"`
						Start string `xml:"start"`
						End   string `xml:"end"`
						File  struct {
							Text    string `xml:",chardata"`
							ID      string `xml:"id,attr"`
							Name    string `xml:"name"`
							Pathurl string `xml:"pathurl"`
							Media   struct {
								Text  string `xml:",chardata"`
								Video struct {
									Text                  string `xml:",chardata"`
									Samplecharacteristics struct {
										Text             string `xml:",chardata"`
										Width            string `xml:"width"`
										Height           string `xml:"height"`
										Anamorphic       string `xml:"anamorphic"`
										Pixelaspectratio string `xml:"pixelaspectratio"`
										Fielddominance   string `xml:"fielddominance"`
									} `xml:"samplecharacteristics"`
								} `xml:"video"`
								Audio struct {
									Text         string `xml:",chardata"`
									In           string `xml:"in"`
									Out          string `xml:"out"`
									Channelcount string `xml:"channelcount"`
									Duration     string `xml:"duration"`
								} `xml:"audio"`
							} `xml:"media"`
						} `xml:"file"`
						Link []struct {
							Text       string `xml:",chardata"`
							Mediatype  string `xml:"mediatype"`
							Trackindex string `xml:"trackindex"`
							Clipindex  string `xml:"clipindex"`
							Groupindex string `xml:"groupindex"`
						} `xml:"link"`
					} `xml:"clipitem"`
					Generatoritem []struct {
						Text     string `xml:",chardata"`
						ID       string `xml:"id,attr"`
						Name     string `xml:"name"`
						Duration string `xml:"duration"`
						Rate     struct {
							Text     string `xml:",chardata"`
							Timebase string `xml:"timebase"`
							Ntsc     string `xml:"ntsc"`
						} `xml:"rate"`
						Start        string `xml:"start"`
						End          string `xml:"end"`
						In           string `xml:"in"`
						Out          string `xml:"out"`
						Enabled      string `xml:"enabled"`
						Anamorphic   string `xml:"anamorphic"`
						Alphatype    string `xml:"alphatype"`
						Masterclipid string `xml:"masterclipid"`
						Effect       struct {
							Text           string `xml:",chardata"`
							Name           string `xml:"name"`
							Effectid       string `xml:"effectid"`
							Effectcategory string `xml:"effectcategory"`
							Effecttype     string `xml:"effecttype"`
							Mediatype      string `xml:"mediatype"`
							Parameter      []struct {
								Text        string `xml:",chardata"`
								Parameterid string `xml:"parameterid"`
								Name        string `xml:"name"`
								Value       struct {
									Text  string `xml:",chardata"`
									Horiz string `xml:"horiz"`
									Vert  string `xml:"vert"`
									Alpha string `xml:"alpha"`
									Red   string `xml:"red"`
									Green string `xml:"green"`
									Blue  string `xml:"blue"`
								} `xml:"value"`
								Valuemin  string `xml:"valuemin"`
								Valuemax  string `xml:"valuemax"`
								Valuelist struct {
									Text       string `xml:",chardata"`
									Valueentry []struct {
										Text  string `xml:",chardata"`
										Name  string `xml:"name"`
										Value string `xml:"value"`
									} `xml:"valueentry"`
								} `xml:"valuelist"`
							} `xml:"parameter"`
						} `xml:"effect"`
						Sourcetrack struct {
							Text      string `xml:",chardata"`
							Mediatype string `xml:"mediatype"`
						} `xml:"sourcetrack"`
					} `xml:"generatoritem"`
				} `xml:"track"`
			} `xml:"video"`
			Audio struct {
				Text         string `xml:",chardata"`
				In           string `xml:"in"`
				Out          string `xml:"out"`
				Channelcount string `xml:"channelcount"`
				Duration     string `xml:"duration"`
				Track        []struct {
					Text     string `xml:",chardata"`
					Clipitem []struct {
						Text    string `xml:",chardata"`
						Name    string `xml:"name"`
						Enabled string `xml:"enabled"`
						In      string `xml:"in"`
						Out     string `xml:"out"`
						Start   string `xml:"start"`
						End     string `xml:"end"`
						File    struct {
							Text string `xml:",chardata"`
							ID   string `xml:"id,attr"`
						} `xml:"file"`
						Sourcetrack struct {
							Text       string `xml:",chardata"`
							Mediatype  string `xml:"mediatype"`
							Trackindex string `xml:"trackindex"`
						} `xml:"sourcetrack"`
					} `xml:"clipitem"`
				} `xml:"track"`
			} `xml:"audio"`
		} `xml:"media"`
	} `xml:"sequence"`
}

/**
*	removePremiereXMLIndent
**/
func removePremiereXMLIndent(origXML string) string {

	var ret string

	scanner := bufio.NewScanner(strings.NewReader(origXML))
	for scanner.Scan() {

		txt := scanner.Text()

		if strings.Index(txt, "<value>") != -1 && strings.Index(txt, "&#13;") != -1 {
			//↵
			txt = strings.Replace(txt, "&#13;", "↵", -1)
		}

		ret += strings.TrimSpace(txt)
	}

	return ret

}

// Entity can be used to map non-standard entity names to string replacements.
// The parser behaves as if these standard mappings are present in the map,
// regardless of the actual map content:
//
//	"lt": "<",
//	"gt": ">",
//	"amp": "&",
//	"apos": "'",
//	"quot": `"`,
func replacePremiereXMLSpecialCharacter(origXML string) string {
	var ret string

	scanner := bufio.NewScanner(strings.NewReader(origXML))
	for scanner.Scan() {

		txt := scanner.Text()

		if strings.Index(txt, "<value>") != -1 {

			if strings.Index(txt, "↵") != -1 {
				txt = strings.Replace(txt, "↵", "&#13;", -1)
			}
			if strings.Index(txt, "&amp;") != -1 {
				txt = strings.Replace(txt, "&amp;", "", -1)
			}
			if strings.Index(txt, "&lt;") != -1 {
				txt = strings.Replace(txt, "&lt;", "", -1)
			}
			if strings.Index(txt, "&gt;") != -1 {
				txt = strings.Replace(txt, "&gt;", "", -1)
			}
			if strings.Index(txt, "&#39;") != -1 {
				txt = strings.Replace(txt, "&#39;", "'", -1)
			}
			if strings.Index(txt, "&apos;") != -1 {
				txt = strings.Replace(txt, "&apos;", "'", -1)
			}
			if strings.Index(txt, "&#34;") != -1 {
				txt = strings.Replace(txt, "&#34;", "\"", -1)
			}
			if strings.Index(txt, "&quot;") != -1 {
				txt = strings.Replace(txt, "&quot;", "\"", -1)
			}
		}

		ret += strings.TrimSpace(txt)
	}

	return ret
}

/**
*	프리미어 업로드 파일 처리
**/
func premiereXMLProc(xmlString string, fname string) int {
	xmlString = removePremiereXMLIndent(xmlString)

	var xmlData premiereXML

	xmlerr := xml.Unmarshal([]byte(xmlString), &xmlData)

	txtBlockMap := make(map[int]string)
	var txtBlockCnt int
	var oneTxtBlock string
	var txtCnt int

	if xmlerr == nil {

		for ti := 0; ti < len(xmlData.Sequence.Media.Video.Track); ti++ {

			for gi := 0; gi < len(xmlData.Sequence.Media.Video.Track[ti].Generatoritem); gi++ {

				gitem := xmlData.Sequence.Media.Video.Track[ti].Generatoritem[gi]

				for epi := 0; epi < len(gitem.Effect.Parameter); epi++ {

					epitem := gitem.Effect.Parameter[epi]

					if epitem.Parameterid == "str" && epitem.Name == "Text" {

						txtCnt++

						if len(oneTxtBlock)+len(epitem.Value.Text)+len("\r\n") < SpellCheckLimit {
							oneTxtBlock += epitem.Value.Text + "\r\n"
						} else {
							txtBlockMap[txtBlockCnt] = oneTxtBlock
							txtBlockCnt++

							oneTxtBlock = epitem.Value.Text + "\r\n"
						}
					}
				}
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
*	savePremiereData
**/
func savePremiereData(data string, fname string, savefname string) map[string]string {
	ret := make(map[string]string)
	ret["result"] = "ERR"

	var lj saveJSON

	jerr := json.Unmarshal([]byte(data), &lj)

	if jerr != nil {
		ret["result"] = "ERR"
	} else {
		//
		oldXMLFileData, oldXMLFileDataErr := ioutil.ReadFile(workFilepath(fname))

		if oldXMLFileDataErr != nil {
			ret["result"] = "ERR"
		} else {

			oldXMLFileData = []byte(removePremiereXMLIndent(string(oldXMLFileData)))

			var xmlData premiereXML

			xmlerr := xml.Unmarshal(oldXMLFileData, &xmlData)

			var txtCnt int

			if xmlerr == nil {

				for ti := 0; ti < len(xmlData.Sequence.Media.Video.Track); ti++ {

					for gi := 0; gi < len(xmlData.Sequence.Media.Video.Track[ti].Generatoritem); gi++ {

						gitem := &xmlData.Sequence.Media.Video.Track[ti].Generatoritem[gi]

						for epi := 0; epi < len(gitem.Effect.Parameter); epi++ {

							epitem := &gitem.Effect.Parameter[epi]

							if epitem.Parameterid == "str" && epitem.Name == "Text" {
								if txtCnt < len(lj) {
									epitem.Value.Text = lj[txtCnt].Value
								}
								txtCnt++
							}
						}
					}
				}
				xmlMarData, xmlMarErr := xml.MarshalIndent(&xmlData, "", "  ")

				if xmlMarErr != nil {
					ret["result"] = "ERR"
				} else {

					replaceXML := replacePremiereXMLSpecialCharacter(string(xmlMarData))

					werr := ioutil.WriteFile(workFilepath(savefname), []byte(replaceXML), 0644)

					if werr != nil {
						ret["result"] = "ERR"
					} else {
						ret["result"] = "OK"
						ret["filename"] = savefname
					}
				}
			}
			//
		}
	}
	return ret
}
