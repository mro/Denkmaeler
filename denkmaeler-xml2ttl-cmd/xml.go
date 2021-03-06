// Copyright (c) 2016-2016 Marcus Rohrmoser, http://mro.name/~me
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
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

const bayern = "09" // https://de.wikipedia.org/wiki/Amtlicher_Gemeindeschl%C3%BCssel

type gemeinde struct {
	regierungsbezirk  string
	landkreis         string
	gemeinde          string
	gemeindeschlüssel string
}

type denkmal struct {
	aktennummer       string
	gemeindeschlüssel string   // https://de.wikipedia.org/wiki/Amtlicher_Gemeindeschl%C3%BCssel
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
	if 1 == len(d.adresse) && "" == d.adresse[0] {
		d.adresse = []string{}
	}
	d.beschreibung = strings.TrimSpace(d.beschreibung)
	*l = append(*l, *d)
	d.isFinished = true
}

func mustAtoi(s string) int {
	r, e := strconv.Atoi(s)
	if nil != e {
		panic(s)
	}
	return r
}

func fineFromRaw(pdf pdf2xml) ([]denkmal, time.Time, string, error) {
	tz, err := time.LoadLocation("Europe/Berlin")
	if nil != err {
		panic(err)
	}
	modifiedPat := regexp.MustCompile(`^Stand (\d+)\.(\d+)\.(\d{4})`)

	var modified time.Time
	var ret []denkmal
	var d denkmal
	var gemeindeschlüssel string
	var typ string
	var oldText text
	var g gemeinde
	aktenPat := regexp.MustCompile(`([DE]-(\d)-(\d{2})-(\d{3})-\d+)|([D]-\d-\d{4}-\d{4})`)

	for _, page := range pdf.Page {
		sort.Sort(page.Text)
		for _, text := range page.Text {
			// fmt.Printf("(%d,%d) f=%d b='%s' v='%s'\n", text.Top, text.Left, text.Font, text.Bold, text.Value)
			switch text.Font {
			case 1:
				if g.regierungsbezirk == "" {
					if !strings.HasPrefix(text.Value, "Regierungsbezirk ") {
						panic("oha")
					}
					g.regierungsbezirk = text.Value
				} else {
					if g.landkreis == "" {
						g.landkreis = text.Value
					} else {
						if g.gemeinde == "" {
							g.gemeinde = text.Value
						}
					}
				}
			case 2, 5:
				switch text.Bold {
				case "Baudenkmäler", "Bodendenkmäler":
					typ = text.Bold
				}
			case 3:
				m := aktenPat.FindStringSubmatch(text.Bold)
				switch {
				case nil != m:
					d.finish(&ret)
					// log.Printf("%d b='%s' v='%s'\n", text.Font, text.Bold, text.Value)
					if "" == gemeindeschlüssel && "" != m[1] {
						// log.Printf("🔑 %s\n", m[1])
						gemeindeschlüssel = bayern + " " + m[2] + " " + m[3] + " " + m[4]
						g.gemeindeschlüssel = gemeindeschlüssel
					}
					d = denkmal{
						aktennummer:       text.Bold,
						gemeindeschlüssel: gemeindeschlüssel,
						typ:               typ,
					}
					oldText = text
				case strings.HasPrefix(text.Bold, "nicht nachqualifiziert"), "nachqualifiziert" == text.Bold:
					d.verfahrensstand = text.Bold
					oldText = text
				case strings.HasPrefix(text.Bold, "Anzahl "): // NOOP
					// log.Printf("%s\n", text.Bold)
				default:
					d.adresse = append(d.adresse, text.Bold)
					d.beschreibung = text.Value
					oldText = text
				}
			case 4:
				m := modifiedPat.FindStringSubmatch(text.Value)
				switch {
				case nil != m:
					modified = time.Date(mustAtoi(m[3]), time.Month(mustAtoi(m[2])), mustAtoi(m[1]), 0, 0, 0, 0, tz)
				case strings.HasPrefix(text.Value, "Seite "): // NOOP
				case strings.HasPrefix(text.Value, "© Bayerisches Landesamt für"): // NOOP
				default:
					sep := " "
					first, _ := utf8.DecodeRuneInString(text.Value[0:])
					switch {
					case text.Top-oldText.Top > 30:
						sep = "\n\n"
					case strings.HasSuffix(oldText.Value, ".") && !unicode.IsLower(first) && oldText.Width < 550:
						// when is it a forced linefeed?
						// - previous line ends with "."
						// - current line does not start with lowercase
						// - previous width < 550
						sep = "\n"
					}
					d.beschreibung += sep + text.Value
					oldText = text
				}
			}
		}
	}
	d.finish(&ret)
	return ret, modified, g.gemeinde, nil
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

type textSlice []text

type page struct {
	XMLName xml.Name  `xml:"page"`
	Text    textSlice `xml:"text"`
}

type text struct {
	XMLName xml.Name `xml:"text"`
	Top     int32    `xml:"top,attr"`
	Left    int32    `xml:"left,attr"`
	Width   int32    `xml:"width,attr"`
	Font    int8     `xml:"font,attr"`
	Bold    string   `xml:"b"`
	Value   string   `xml:",chardata"` // http://stackoverflow.com/a/20600762
}

func (s textSlice) Len() int {
	return len(s)
}
func (s textSlice) Less(i, j int) bool {
	if s[i].Top < s[j].Top {
		return true
	}
	if s[i].Top > s[j].Top {
		return false
	}
	return s[i].Left < s[j].Left
}
func (s textSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
