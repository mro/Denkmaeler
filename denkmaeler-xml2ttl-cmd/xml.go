// Copyright (c) 2016-2015 Marcus Rohrmoser, https://github.com/mro/
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and
// associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute,
// sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all copies or
// substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT
// NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES
// OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
//
// MIT License http://opensource.org/licenses/MIT

package main

import (
	"encoding/xml"
	"io"
	// "log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type gemeinde struct {
	regierungsbezirk  string
	landkreis         string
	gemeinde          string
	gemeindeschlüssel int
}

type denkmal struct {
	aktennummer       string
	gemeindeschlüssel int      // https://de.wikipedia.org/wiki/Amtlicher_Gemeindeschl%C3%BCssel
	typ               string   // baudenkmäler / bodendenkmäler
	adresse           []string //
	beschreibung      string
	verfahrensstand   string

	isFinished bool
}

func (d *denkmal) finish(l *[]denkmal) {
	if 0 == len(d.aktennummer) || d.isFinished {
		return
	}
	d.adresse = strings.Split(strings.Join(d.adresse, " "), ";")
	for i, a := range d.adresse {
		d.adresse[i] = strings.Trim(a, " .")
	}
	d.beschreibung = strings.TrimSpace(d.beschreibung)
	*l = append(*l, *d)
	d.isFinished = true
}

func fineFromRaw(pdf pdf2xml) ([]denkmal, error) {
	var ret []denkmal
	var d denkmal
	var gemeindeschlüssel int
	var typ string
	var oldTop int32
	var g gemeinde
	aktenPat := regexp.MustCompile(`([DE]-(\d)-(\d{2})-(\d{3})-\d+)|([D]-\d-\d{4}-\d{4})`)

	for _, page := range pdf.Page {
		for _, text := range page.Text {
			// fmt.Printf("%d b='%s' v='%s'\n", text.Font, text.Bold, text.Value)
			switch text.Font {
			case 1:
				switch {
				case 45 < text.Top && text.Top < 75:
					if !strings.HasPrefix(text.Value, "Regierungsbezirk ") {
						panic("oha")
					}
					g.regierungsbezirk = text.Value
				case 75 < text.Top && text.Top < 110:
					g.landkreis = text.Value
				default:
					g.gemeinde = text.Value
				}
			case 2, 5:
				switch text.Bold {
				case "Baudenkmäler", "Bodendenkmäler":
					typ = text.Bold
				case g.gemeinde:
				default:
					panic("87: " + text.Bold)
				}
			case 3:
				m := aktenPat.FindStringSubmatch(text.Bold)
				switch {
				case nil != m:
					d.finish(&ret)
					// log.Printf("%d b='%s' v='%s'\n", text.Font, text.Bold, text.Value)
					if 0 == gemeindeschlüssel && "" != m[1] {
						// log.Printf("🔑 %s\n", m[1])
						gemeindeschlüssel, _ = strconv.Atoi(m[2] + m[3] + m[4])
						g.gemeindeschlüssel = gemeindeschlüssel
					}
					d = denkmal{
						aktennummer:       text.Bold,
						gemeindeschlüssel: gemeindeschlüssel,
						typ:               typ,
					}
				case "nachqualifiziert" == text.Bold:
					d.verfahrensstand = text.Bold
				case strings.HasPrefix(text.Bold, "Anzahl "): // NOOP
					// log.Printf("%s\n", text.Bold)
				default:
					d.adresse = append(d.adresse, text.Bold)
					d.beschreibung = text.Value
				}
			case 4:
				switch {
				case strings.HasPrefix(text.Value, "Stand "): // TODO
				case strings.HasPrefix(text.Value, "Seite "): // NOOP
				case strings.HasPrefix(text.Value, "© Bayerisches Landesamt für Denkmalpflege"): // NOOP
				default:
					sep := " "
					if text.Top-oldTop > 30 {
						sep = "\n\n"
					}
					d.beschreibung += sep + text.Value
				}
			}
			oldTop = text.Top
		}
	}
	d.finish(&ret)
	return ret, nil
}

func rawFromXmlReader(xmlFile io.Reader) (pdf2xml, error) {
	x := pdf2xml{}
	err := xml.NewDecoder(xmlFile).Decode(&x)
	return x, err
}

func rawFromXmlFileName(xmlFileName string) (pdf2xml, error) {
	r, err := os.Open(xmlFileName)
	if nil != err {
		return pdf2xml{}, err
	}
	defer r.Close()
	return rawFromXmlReader(r)
}

type pdf2xml struct {
	XMLName xml.Name `xml:"pdf2xml"`
	Page    []page   `xml:"page"`
}

type page struct {
	XMLName xml.Name `xml:"page"`
	Text    []text   `xml:"text"`
}

type text struct {
	XMLName xml.Name `xml:"text"`
	Top     int32    `xml:"top,attr"`
	Left    int32    `xml:"left,attr"`
	Font    int8     `xml:"font,attr"`
	Bold    string   `xml:"b"`
	Value   string   `xml:",chardata"` // http://stackoverflow.com/a/20600762
}
