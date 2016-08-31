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

// http://golang.org/pkg/testing/
// http://blog.stretchr.com/2014/03/05/test-driven-development-specifically-in-golang/
// https://xivilization.net/~marek/blog/2015/05/04/go-1-dot-4-2-for-raspberry-pi/

package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXmlEncoding(t *testing.T) {
	x := pdf2xml{
		Page: []page{page{
			Text: []text{text{
				Bold:  "Foo",
				Value: "Bar",
			}},
		}},
	}
	buf := new(bytes.Buffer)
	err := xml.NewEncoder(buf).Encode(x)
	assert.Nil(t, err, "jaja")
	assert.Equal(t, "<pdf2xml><page><text top=\"0\" left=\"0\" font=\"0\"><b>Foo</b>Bar</text></page></pdf2xml>", buf.String(), "echt?")
}

func TestDataFromXmlFileName(t *testing.T) {
	pdf, err := rawFromXmlFileName("")
	assert.Equal(t, "open : no such file or directory", err.Error(), "soso")

	pdf, err = rawFromXmlFileName("testdata/162000-liste.xml")
	assert.Nil(t, err, "soso")
	assert.Equal(t, 890, len(pdf.Page), "page count")
	assert.Equal(t, 50, len(pdf.Page[0].Text), "text count")
	{
		tx := pdf.Page[0].Text[0]
		assert.Equal(t, int32(51), tx.Top, "text top")
		assert.Equal(t, int32(602), tx.Left, "text left")
		assert.Equal(t, int8(1), tx.Font, "text font")
		assert.Equal(t, "", tx.Bold, "text bold")
		assert.Equal(t, "Regierungsbezirk Oberbayern", tx.Value, "text value")
	}
	{
		tx := pdf.Page[0].Text[3]
		assert.Equal(t, int32(202), tx.Top, "text top")
		assert.Equal(t, int32(38), tx.Left, "text left")
		assert.Equal(t, int8(2), tx.Font, "text font")
		assert.Equal(t, "München", tx.Bold, "text bold")
		assert.Equal(t, "", tx.Value, "text value")
	}
	{
		tx := pdf.Page[0].Text[5]
		assert.Equal(t, int32(260), tx.Top, "text top")
		assert.Equal(t, int32(38), tx.Left, "text left")
		assert.Equal(t, int8(3), tx.Font, "text font")
		assert.Equal(t, "E-1-62-000-30", tx.Bold, "text bold")
		assert.Equal(t, "", tx.Value, "text value")
	}
	{
		tx := pdf.Page[0].Text[6]
		assert.Equal(t, int32(260), tx.Top, "text top")
		assert.Equal(t, int32(204), tx.Left, "text left")
		assert.Equal(t, int8(3), tx.Font, "text font")
		assert.Equal(t, "Ensemble Wohnanlagen am Loehleplatz. ", tx.Bold, "text bold")
		assert.Equal(t, " Das Ensemble der Wohnanlage am", tx.Value, "text value")
	}
	{
		tx := pdf.Page[544].Text[12]
		assert.Equal(t, int32(372), tx.Top, "text top")
		assert.Equal(t, int32(204), tx.Left, "text left")
		assert.Equal(t, int8(3), tx.Font, "text font")
		assert.Equal(t, "Moltkestraße 11; Unertlstraße 1; Unertlstraße 2; Unertlstraße 3; Unertlstraße 4;", tx.Bold, "text bold")
		assert.Equal(t, "", tx.Value, "text value")
	}
	{
		tx := pdf.Page[544].Text[14]
		assert.Equal(t, int32(414), tx.Top, "text top")
		assert.Equal(t, int32(204), tx.Left, "text left")
		assert.Equal(t, int8(3), tx.Font, "text font")
		assert.Equal(t, "Viktoriastraße 32; Viktoriastraße 34. ", tx.Bold, "text bold")
		assert.Equal(t, " Wohnanlagen an der Unertlstraße, im Abschnitt", tx.Value, "text value")
	}

	assert.Equal(t, 16, len(pdf.Page[889].Text), "text count")

	for ip, vp := range pdf.Page {
		for it, vt := range vp.Text {
			if "nonono D-1-62-000-7825" == vt.Bold {
				fmt.Printf("page %d text %d\n", ip, it)
			}
		}
	}
}

func TestFineFromRawSmall(t *testing.T) {
	pdf, err := rawFromXmlFileName("testdata/189159-liste.xml")
	ds, err := fineFromRaw(pdf)
	assert.Nil(t, err, "soso")

	// assert.Equal(t, "Regierungsbezirk Oberbayern", l.gemeinde.regierungsbezirk, "soso")
	// assert.Equal(t, "Traunstein", l.gemeinde.landkreis, "soso")
	// assert.Equal(t, "Übersee", l.gemeinde.gemeinde, "soso")
	assert.Equal(t, 46, len(ds), "soso")
	{
		d := ds[0]
		assert.Equal(t, "D-1-89-159-30", d.aktennummer, "soso")
		assert.Equal(t, "09 1 89 159", d.gemeindeschlüssel, "soso")
		assert.Equal(t, "Baudenkmäler", d.typ, "soso")
		assert.Equal(t, 1, len(d.adresse), "soso")
		assert.Equal(t, "Achenzipf", d.adresse[0], "soso")
		assert.Equal(t, "Nikolauskapelle, sog. Achenzipfkapelle, bez. 1808 und 1952; mit Ausstattung; am Chiemseeufer.", d.beschreibung, "soso")
		assert.Equal(t, "nachqualifiziert", d.verfahrensstand, "soso")
	}

	{
		d := ds[41]
		assert.Equal(t, "D-1-8140-0209", d.aktennummer, "soso")
		assert.Equal(t, "09 1 89 159", d.gemeindeschlüssel, "soso")
		assert.Equal(t, "Bodendenkmäler", d.typ, "soso")
		assert.Equal(t, 1, len(d.adresse), "soso")
		assert.Equal(t, "", d.adresse[0], "soso")
		assert.Equal(t, "Untertägige spätmittelalterliche und frühneuzeitliche Befunde im Bereich der Kath. Pfarrkirche St. Nikolaus in Übersee und ihrer Vorgängerbauten.", d.beschreibung, "soso")
		assert.Equal(t, "nachqualifiziert", d.verfahrensstand, "soso")
	}
}

func TestFineFromRawLarge(t *testing.T) {
	pdf, err := rawFromXmlFileName("testdata/162000-liste.xml")
	ds, err := fineFromRaw(pdf)
	assert.Nil(t, err, "soso")

	// assert.Equal(t, "Regierungsbezirk Oberbayern", l.gemeinde.regierungsbezirk, "soso")
	// assert.Equal(t, "Traunstein", l.gemeinde.landkreis, "soso")
	// assert.Equal(t, "Übersee", l.gemeinde.gemeinde, "soso")
	assert.Equal(t, 7171, len(ds), "soso")
	{
		d := ds[0]
		assert.Equal(t, "E-1-62-000-30", d.aktennummer, "soso")
		assert.Equal(t, "09 1 62 000", d.gemeindeschlüssel, "soso")
		assert.Equal(t, "Baudenkmäler", d.typ, "soso")
		assert.Equal(t, 1, len(d.adresse), "soso")
		assert.Equal(t, "Ensemble Wohnanlagen am Loehleplatz", d.adresse[0], "soso")
		assert.Equal(t, "Das Ensemble der Wohnanlage am Loehleplatz, zwischen 1907 bis 1926 errichtet, stellt ein Beispiel des genossenschaftlichen Wohnungsbaus in München dar. Die Bebauung, vom Ersten Weltkrieg unterbrochen, erfolgte durch den „Verein für Verbesserung der Wohnungsverhältnisse in München“ unter der Führung von Johann Mund und unter Beteiligung von Richard Fuchs, Hans Wagner, Paul Liebergesell und Feodor Lehmann. Der 1899 gegründete Verein zählt zu den vielen, seit der Mitte des 19. Jahrhunderts in Deutschland entstehenden Wohnungsbaugenossenschaften, die als Antwort auf die drängende Wohnungsfrage insbesondere für die Bevölkerungsgruppen mit kleinem Einkommen ein gesundes Wohnumfeld schaffen wollten.\n\nDie Ausführung der Anlage am Loehleplatz entspricht dem Grundgedanken der kurz zuvor in Kraft getretenen Staffelbauordnung des Stadterweiterungsbüros von Theodor Fischer, welche Neubauprojekte in ein übergeordnetes, gesamt-städtebauliches Konzept einzubinden suchte. Dementsprechend ist die äußere Bebauung an der Rosenheimer Straße als Ausfallstraße viergeschossig und die inneren Bauten am Loehleplatz, an der Abenthum- und Wollanistraße von drei- zu zweigeschossigen Mehrfamilienhäusern herabgestaffelt. An der Weißkopfstraße sind schließlich eingeschossige Reihenhauszeilen zu finden. Durch die Ausgestaltung der Eckbauten an der Mündung der Maria-Lehner- Straße wird städtebaulich ein Zugang zu den Straßen- und Platzräumen im Innern der Anlage geschaffen. Die Nord-Süd-Achsen sind auf die Ramersdorfer Kirche ausgerichtet. Die Baukörper sind, besonders aus dem Anfang der Bautätigkeit noch vor dem Ersten Weltkrieg, mittels abwechslungsreicher Dachformen, Gauben, Zwerchhäusern, Erkerbauten, Loggien und Putzdekor reich gegliedert und dabei sowohl symmetrisch wie asymmetrisch zusammengeordnet. Die um einen Hof geschlossene Blockbebauung wird ebenso aufgelockert wie die Folgen von Reihenhäusern. Der Stilwandel zur Nachkriegsarchitektur wird, besonders bei den jüngeren Bauten an der Rosenheimer Straße, spürbar, bleibt jedoch im vorgegebenen Rahmen.", d.beschreibung, "soso")
		assert.Equal(t, "", d.verfahrensstand, "soso")
	}

	{
		d := ds[6812]
		assert.Equal(t, "D-1-7734-0101", d.aktennummer, "soso")
		assert.Equal(t, "09 1 62 000", d.gemeindeschlüssel, "soso")
		assert.Equal(t, "Bodendenkmäler", d.typ, "soso")
		assert.Equal(t, 1, len(d.adresse), "soso")
		assert.Equal(t, "", d.adresse[0], "soso")
		assert.Equal(t, "Grabhügel mit Bestattungen der Hallstattzeit sowie Siedlung vorgeschichtlicher Zeitstellung.", d.beschreibung, "soso")
		assert.Equal(t, "nachqualifiziert", d.verfahrensstand, "soso")
	}
}
