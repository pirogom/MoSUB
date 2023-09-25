package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strings"
)

type finalcutXMLObj struct {
	XMLName       xml.Name `xml:"fcpxml"`
	Text          string   `xml:",chardata"`
	Version       string   `xml:"version,attr"`
	ImportOptions struct {
		Text   string `xml:",chardata"`
		Option struct {
			Text  string `xml:",chardata"`
			Value string `xml:"value,attr"`
			Key   string `xml:"key,attr"`
		} `xml:"option"`
	} `xml:"import-options"`
	Resources struct {
		Text   string `xml:",chardata"`
		Format struct {
			Text          string `xml:",chardata"`
			ID            string `xml:"id,attr"`
			FrameDuration string `xml:"frameDuration,attr"`
			Width         string `xml:"width,attr"`
			Height        string `xml:"height,attr"`
		} `xml:"format"`
		Asset struct {
			Text     string `xml:",chardata"`
			ID       string `xml:"id,attr"`
			Src      string `xml:"src,attr"`
			Format   string `xml:"format,attr"`
			Duration string `xml:"duration,attr"`
			Name     string `xml:"name,attr"`
			HasAudio string `xml:"hasAudio,attr"`
			HasVideo string `xml:"hasVideo,attr"`
		} `xml:"asset"`
		Effect struct {
			Text string `xml:",chardata"`
			ID   string `xml:"id,attr"`
			Name string `xml:"name,attr"`
			Uid  string `xml:"uid,attr"`
		} `xml:"effect"`
	} `xml:"resources"`
	Library struct {
		Text  string `xml:",chardata"`
		Event struct {
			Text    string `xml:",chardata"`
			Project struct {
				Text     string `xml:",chardata"`
				Name     string `xml:"name,attr"`
				Sequence struct {
					Text     string `xml:",chardata"`
					Duration string `xml:"duration,attr"`
					Format   string `xml:"format,attr"`
					TcStart  string `xml:"tcStart,attr"`
					Spine    struct {
						Text      string `xml:",chardata"`
						AssetClip []struct {
							Text     string `xml:",chardata"`
							Name     string `xml:"name,attr"`
							Start    string `xml:"start,attr"`
							Offset   string `xml:"offset,attr"`
							Duration string `xml:"duration,attr"`
							Ref      string `xml:"ref,attr"`
							Title    struct {
								Chardata string `xml:",chardata"`
								Lane     string `xml:"lane,attr"`
								Name     string `xml:"name,attr"`
								Ref      string `xml:"ref,attr"`
								Offset   string `xml:"offset,attr"`
								Start    string `xml:"start,attr"`
								Duration string `xml:"duration,attr"`
								Text     struct {
									Text      string `xml:",chardata"`
									TextStyle struct {
										Text string `xml:",chardata"`
										Ref  string `xml:"ref,attr"`
									} `xml:"text-style"`
								} `xml:"text"`
								TextStyleDef struct {
									Text      string `xml:",chardata"`
									ID        string `xml:"id,attr"`
									TextStyle struct {
										Text        string `xml:",chardata"`
										Alignment   string `xml:"alignment,attr"`
										FontColor   string `xml:"fontColor,attr"`
										Font        string `xml:"font,attr"`
										FontSize    string `xml:"fontSize,attr"`
										LineSpacing string `xml:"lineSpacing,attr"`
										Baseline    string `xml:"baseline,attr"`
										StrokeColor string `xml:"strokeColor,attr"`
										StrokeWidth string `xml:"strokeWidth,attr"`
									} `xml:"text-style"`
								} `xml:"text-style-def"`
							} `xml:"title"`
						} `xml:"asset-clip"`
					} `xml:"spine"`
				} `xml:"sequence"`
			} `xml:"project"`
		} `xml:"event"`
	} `xml:"library"`
}

/**
*	replaceFinalcutXMLSpecialCharacter
**/
func replaceFinalcutXMLSpecialCharacter(origXML string) string {
	var ret string

	scanner := bufio.NewScanner(strings.NewReader(origXML))
	for scanner.Scan() {

		txt := scanner.Text()

		if strings.Index(txt, "<text-style") != -1 {

			if strings.Index(txt, "↵") != -1 {
				txt = strings.Replace(txt, "↵", "\n", -1)
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
*	finalcutXMLProc
**/
func finalcutXMLProc(xmlString string, fname string) int {
	xmlString = removeFinalcutXMLIndent(xmlString)

	var xmlData finalcutXMLObj

	xmlerr := xml.Unmarshal([]byte(xmlString), &xmlData)

	txtBlockMap := make(map[int]string)
	var txtBlockCnt int
	var oneTxtBlock string
	var txtCnt int

	if xmlerr == nil {

		assetClip := xmlData.Library.Event.Project.Sequence.Spine.AssetClip

		for ai := 0; ai < len(assetClip); ai++ {
			subText := assetClip[ai].Title.Text.TextStyle.Text
			if len(subText) > 0 {
				txtCnt++

				if len(oneTxtBlock)+len(subText)+len("\r\n") < SpellCheckLimit {
					oneTxtBlock += subText + "\r\n"
				} else {

					txtBlockMap[txtBlockCnt] = oneTxtBlock
					txtBlockCnt++

					oneTxtBlock = subText + "\r\n"
				}
			}
		}

		/*		mediaAssetClip := xmlData.Resources.Media.Sequence.Spine.AssetClip

				for mai := 0; mai < len(mediaAssetClip); mai++ {
					for ti := 0; ti < len(mediaAssetClip[mai].Title); ti++ {
						for tsi := 0; tsi < len(mediaAssetClip[mai].Title[ti].Text.TextStyle); tsi++ {
							subText := mediaAssetClip[mai].Title[ti].Text.TextStyle[tsi].Text
							if len(subText) > 0 {
								txtCnt++

								if len(oneTxtBlock)+len(subText)+len("\r\n") < SpellCheckLimit {
									oneTxtBlock += subText + "\r\n"
								} else {

									txtBlockMap[txtBlockCnt] = oneTxtBlock
									txtBlockCnt++

									oneTxtBlock = subText + "\r\n"
								}
							}
						}
					}
				}*/

		txtBlockMap[txtBlockCnt] = oneTxtBlock

		if txtCnt > 0 {
			go spellCheckProc(txtBlockMap, fname)
		}
	}

	return txtCnt
}

/**
*	saveFinalcutData
**/
func saveFinalcutData(data string, fname string, savefname string) map[string]string {
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

			oldXMLFileData = []byte(removeFinalcutXMLIndent(string(oldXMLFileData)))

			var xmlData finalcutXMLObj

			xmlerr := xml.Unmarshal(oldXMLFileData, &xmlData)

			var txtCnt int

			if xmlerr == nil {

				for ai := 0; ai < len(xmlData.Library.Event.Project.Sequence.Spine.AssetClip); ai++ {
					if len(xmlData.Library.Event.Project.Sequence.Spine.AssetClip[ai].Title.Text.TextStyle.Text) > 0 {

						xmlData.Library.Event.Project.Sequence.Spine.AssetClip[ai].Title.Text.TextStyle.Text = lj[txtCnt].Value

						txtCnt++
					}
				}

				/*				for mai := 0; mai < len(xmlData.Resources.Media.Sequence.Spine.AssetClip); mai++ {
								for ti := 0; ti < len(xmlData.Resources.Media.Sequence.Spine.AssetClip[mai].Title); ti++ {
									for tsi := 0; tsi < len(xmlData.Resources.Media.Sequence.Spine.AssetClip[mai].Title[ti].Text.TextStyle); tsi++ {
										subText := xmlData.Resources.Media.Sequence.Spine.AssetClip[mai].Title[ti].Text.TextStyle[tsi].Text
										if len(subText) > 0 {

											xmlData.Resources.Media.Sequence.Spine.AssetClip[mai].Title[ti].Text.TextStyle[tsi].Text = lj[txtCnt].Value

											txtCnt++
										}
									}
								}
							}*/

				xmlMarData, xmlMarErr := xml.MarshalIndent(&xmlData, "", "  ")

				fmt.Println(string(xmlMarData))

				if xmlMarErr != nil {
					ret["result"] = "ERR"
				} else {
					replaceXML := replaceFinalcutXMLSpecialCharacter(string(xmlMarData))

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
