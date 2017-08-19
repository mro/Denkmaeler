// Copyright (c) 2016-2017 Marcus Rohrmoser, http://mro.name/~me
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
	"time"

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
	assert.Equal(t, "<pdf2xml><page><text top=\"0\" left=\"0\" width=\"0\" font=\"0\"><b>Foo</b>Bar</text></page></pdf2xml>", buf.String(), "echt?")
}

func TestDataFromXmlFileName(t *testing.T) {
	pdf, err := rawFromXmlFileName("")
	assert.Equal(t, "open : no such file or directory", err.Error(), "soso")

	pdf, err = rawFromXmlFileName("testdata/162000.xml")
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
	pdf, err := rawFromXmlFileName("testdata/189159.xml")
	ds, mod, name, err := fineFromRaw(pdf)
	assert.Nil(t, err, "soso")
	assert.Equal(t, "2016-08-13T00:00:00+02:00", mod.Format(time.RFC3339), "huhu")
	assert.Equal(t, "Übersee", name, "huhu")
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
		assert.Equal(t, 0, len(d.adresse), "soso")
		assert.Equal(t, "Untertägige spätmittelalterliche und frühneuzeitliche Befunde im Bereich der Kath. Pfarrkirche St. Nikolaus in Übersee und ihrer Vorgängerbauten.", d.beschreibung, "soso")
		assert.Equal(t, "nachqualifiziert", d.verfahrensstand, "soso")
	}
}

func TestFineFromRawLarge(t *testing.T) {
	pdf, err := rawFromXmlFileName("testdata/162000.xml")
	ds, mod, name, err := fineFromRaw(pdf)
	assert.Nil(t, err, "soso")
	assert.Equal(t, "2016-08-23T00:00:00+02:00", mod.Format(time.RFC3339), "huhu")
	assert.Equal(t, "München", name, "huhu")
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
		d := ds[1]
		assert.Equal(t, "E-1-62-000-70", d.aktennummer, "soso")
		assert.Equal(t, "09 1 62 000", d.gemeindeschlüssel, "soso")
		assert.Equal(t, "Baudenkmäler", d.typ, "soso")
		assert.Equal(t, 1, len(d.adresse), "soso")
		assert.Equal(t, "Ensemble Olympiapark", d.adresse[0], "soso")
		assert.Equal(t, "Das Ensemble Olympiapark umfasst die in dem künstlich gestalteten Landschaftspark zur Ausrichtung der XX. Olympischen Spiele der Neuzeit 1972 angelegten Sportstätten mit den sportlichen und funktionalen Nebeneinrichtungen, dem Olympiaturm, den Verkehrsanlagen sowie dem Olympischen Dorf.\nDer Olympiapark befindet sich auf der ausgedehnten Ebene des Oberwiesenfelds im Nordwesten Münchens. Die Fläche war seit dem späten 18. Jahrhundert Exerziergelände und von 1929 bis zur Eröffnung des Flughafens Riem 1939 der erste Münchner Verkehrsflughafen. Nach der Zerstörung Münchens im Zweiten Weltkrieg wurde das Areal für den Räumungsschutt genutzt. Südlich des durch den Pak verlaufenden Nymphenburg-Biedersteiner Kanal entstand bis 1958, neben den Endkippen in Neuhofen und im Luitpoldpark, der dritte und umfangreichste Schuttberg. Das 1965 als Erholungszone ausgewiesene Gebiet war inzwischen nur mit vereinzelten, öffentlichen Gebäuden – mit der Eissporthalle und dem Fernmeldehochturm der Bundespost – bebaut. In Planung waren zu diesem Zeitpunkt Teilflächen für eine Hochschulsportanlage, für eine Studentenwohnanlage und für den noch fehlenden nordwestlichen Abschnitt des im Entstehen begriffenen Mittleren Rings zu nutzen. Diese Pläne wurden, als München 1966 den Zuschlag als Austragungsort für die XX. Olympischen Spiele bekam, in das Gesamtkonzept integriert. Für die Gestaltung der olympischen Sportstätten schrieb man 1967 ein Architektenwettbewerb aus, den das Büro Behnisch und Partner gewann.\n\nIn der Gesamtgliederung des bis 1972 fertiggestellten Olympiaparks sind zwei Großkomplexe deutlich voneinander zu unterschieden, die durch das breite, ost-westlich verlaufende, das Gelände halbierende Verkehrsband des Mittleren Rings räumlich scharf getrennt werden. Im Süden bilden die Hauptsportstätten (Stadion, Sporthalle, Schwimmhalle) das Herzstück der Anlage und im Norden befindet sich das Olympische Dorf. Diesen Großkomplexen sind Nebeneinrichtungen beigegeben, die Werner-von- Linde Halle und das Radstadion in der südwestlichen Ecke des Geländes und die Hochschulsportanlage westlich des Olympischen Dorfs. Hinzu kommen noch eine Reihe ebenerdiger Anlagen, wie die verschiedenen Spiel-, Sport- und Trainingsplätze sowie der Parkplatz an der Westseite des Stadions.\n\nDas von Günther Behnisch für die Hauptsportstätten entwickelte, übergeordnete Gestaltungskonzept geht von der künstlichen Landschaftsform des Schuttbergs aus, welcher das Gelände im Süden weitgehend gegen die Stadt abschirmt. Seine zufällige Haldenform wird zum Leitbild für die Anlage. Der sog. Olympiaberg erfährt in variierender Wiederholung eine nach Norden abnehmende Staffelung. An dessen nördlichen Abhang wurde der Kanal zu einem kurvenreichen, die Bergfußlinie aufnehmenden See aufgestaut. In dessen größten Halbkreisform Bauchung liegt ein kleines Freilufttheater mit Seebühne. Jenseits des Sees ist eine weitere künstliche Aufschüttung geschaffen worden, an welche sich die großen Sportkampfstätten anlehnen. Stadion, Sport- und Schwimmhalle sind wiederum durch ein zusammenhängendes Zeltdach miteinander verknüpft, dessen bewegte Gestalt an die naturhaften Haufenformen der benachbarten Landschaft erinnert. Das charakteristische Zeltdach geht auf den Entwurf Frei Ottos und auf den statischen Berechnungen von Fritz Leonhardt und Wolfhard Andrä zurück. Auf mächtigen Pylonen gestützt, hält die vorgespannte Seilnetzkonstruktion eine Dachhaut aus durchsichtigen Acrylplatten. Das „Dach ohne Schatten“ beschirmt in regelmäßigen Schwüngen alle drei Sportstätten, überdeckt das gesamte Oval der Sporthalle, schafft eine Torsituation zwischen Sport- und Schwimmhalle und endet auf der Hans-Braun-Brücke in einem einzelnen Pylon. Ein hoher Stellenwert innerhalb der Gesamtkomposition des Olympiaparks kommt der gärtnerischen Gestaltung zu, die in Händen von Günther Grzimek lag. Ähnlich durchdacht, wie die künstlich geschaffenen Landschaftsformen des Olympiaparks sind seine Wegesysteme, seine Ruheplätze, seine Ausstattung mit Kleinarchitekturen und Sitzbänken. Dem entspricht auch eine ebenso kunstvoll eingesetzte Vegetation, bei der etwa Leitbäume die einzelnen Bereiche prägen. So ist der Schuttberg mit Bergkiefern besetzt worden, die Wege sind durch Linden markiert, entlang den Wasserläufen wachsen Silberweiden und dem Parkplatzbereich sind Spitzahornbäume zugeordnet. An herausgehobenen Stellen des Parks sind Plastiken aufgestellt.\nNeben den zentralen Sportbauten sind ebenso die Nebeneinrichtungen im Süden des Olympiaparks zu erwähnen. Sie nehmen gegenüber den Hauptstätten zwar eine bewusst zurückhaltende Gestaltung ein, sind aber dennoch für sich gesehen wichtige Bestandteile des Ensembles und für den Ablauf der Spiele 1972 unverzichtbar. Das Eissportstadion entstand 1966/67 nach dem Entwurf von Rolf Schütze. Zu den Olympischen Spielen konnte es als Boxhalle genutzt werden, da Schütze an eine mögliche Mehrzweckfunktion gedacht hatte. Neben dem Olympiastadion befindet sich die sog. Werner-von-Linde- Halle, die ehemalige Aufwärmhalle für die Athleten. Sie ist zu diesem Zweck unmittelbar mit dem Stadion durch einen unterirdischen Tunnel verbunden. Das Radsportstadion nach Entwurf Herbert Schürmann u. a. nimmt sich ebenfalls zurück. Es ragt nicht in die Höhe, sondern ist in die Landschaft eingebettet. In unmittelbarer Nähe, an der westlichen Stadiontribüne, befindet sich die sog. Parkharfe. Auch deren sichelförmiger Grundriss gehört zum bewussten Gestaltungskonzept des Parks. Die einzelnen Parkbereiche sind mit Hecken und Spitzahornbäumen eingeteilt. Ebenso gestalterisch bedeutsam ist das Kreuzungsbauwerk der Landshuter Allee mit dem Georg-Brauchle-Ring. Der rechtwinklige Sprung des Mittleren Rings von einer Straße auf die andere wird hier mittels weit geschwungener Überführungen bewerkstelligt, die in ihrem Verlauf auf die Kurvung der westlichen Stadiontribüne antworten. Die Bedeutung des Kreuzungsbauwerks ist auch durch die Art seiner Beleuchtung hervorgehoben: mit Hilfe der Beleuchtungskörper, hoher Masten, die bis zu ihrer Spitze mit Strahlern bestückt sind, kommt es zu einer Art Licht-\"Inszenierung\". Zur weiteren verkehrstechnischen Erschließung dienen drei durch radial geführte Fußwege mit den Hauptsportstätten verbundene Haltepunkte des öffentlichen Nahverkehrs: der U-Bahnhof der Olympialinie an der Lerchenauer Straße im Osten, der aus einem bereits bestehenden Industriegleis gewonnene S-Bahnhof im Westen und schließlich die Straßenbahnschleife an der Schwere-Reiter-Straße im Süden. Über allem thront hier in der Südhälfte des Olympiaparks der Fernsehturm. Ehemals von der Deutschen Bundespost zur besseren Sendeleistung des Fernmeldenetzes errichtet, entwickelte sich der Turm zum Wahrzeichen. Der von Sebastian Rosenthal zwischen 1965-67 gebaute Turm ist von überall aus sichtbar und eröffnet von seiner Plattform aus einen freien Blick über den Park, somit auch über den Ring in die Nordhälfte.\nDen Norden erschließen, genauso wie den Süden, auf Dämme geführte Wege, wobei drei Brücken über die trennende Schneise des Mittleren Rings hinwegführen. Die Hauptlinien der Dammwege bündeln sich auf der breit angelegten Hanns-Braun-Brücke. Der in gerader Fortsetzung der Brücke nach Norden ausgerichtete Zweig dieses Wegenetzes spaltet den nördlichen Teil des Olympia-Geländes in zwei Hälften, deren östliche das Olympische Dorf von Werner Wirsing, Günther Eckert, Erwin Heinle und Robert Wischer einnimmt. Die Gestalt des Olympischen Dorfs beruht auf dem Zusammenwirken verschiedener Konzepte. Die Trabantenstadt mit eigenem Zentrum ist hier antikonzentrisch in der Form eines Dreistrahls verwirklicht. Ihr Aufbau basiert auf der konsequenten vertikalen Trennung von Auto- und Fußgängerverkehr und ist vom Gedanken der Terrassenanlage bestimmt. Ihre Struktur lebt von der Verbindung groß dimensionierter Wohnblöcke mit kleineren Einheiten und kleinsten Reihenhauszeilen und der Durchsetzung des Gebauten mit ausgedehnten Grünzonen. Das Zentrum des Olympischen Dorfs ist durch eine Reihe von Hochhauszeilen markiert, die parallel zur Lerchenauer Straße stehen. Diese Hochhäuser bilden die zentrale Ladenstraße entlang des Helene-Mayer-Rings aus. Die Straßbergerstraße, Nadistraße und Connollystraße erschließen von hier aus als Verkehrswege das Wohngebiet. Die entlang dieser Straßen entwickelten Wohnarme strahlen in Form dreier hoher, in ihrem Verlauf mehrfach gebrochener Gebäudeäste nach Westen aus. Die nach Süden ausgerichteten Terrassenbauten umgreifen breite, muldenartige Höfe von parkartigem Charakter. Ihnen sind, ebenfalls terrassenförmig zu den Parkhöfen hin, kleinere Zeilen von Reihenhäusern vorgelagert. Der Anlage ist südlich das seinerzeitige Olympische Dorf der Frauen vorgelagert. Die niedrig gehaltene Kleinsthaussiedlung in Reihenanordnung wird jetzt als Studentendorf genutzt. Die Gebäudegruppen des Olympischen Dorfs sind in ihrer Formgebung gänzlich von ihrer Bauweise in Beton-Fertigteilen abhängig. In bewusstem Kontrast zu diesem betonsichtigen Baukastenprinzip sind die Fußgängerwege mit mehrfarbigen Ziegelsteinen ausgelegt. Mitentscheidend für das charakteristische Erscheinungsbild des Dorfes ist zudem die intensive Bepflanzung der Terrassen. Die damit ermöglichte Fassadenbegrünung ergänzt die unmittelbar angrenzenden, parkartigen Höfe und den sich nach Westen anschließenden Landschaftspark mit Kleinarenen, künstlichen Wasserläufen und Rundplätzen. Auf diese Weise wird die begrünte Architekturlandschaft mit der Parklandschaft verzahnt. Wie der gesamte Olympiapark – mit Beschriftungen, Wegweisern, Logos und Piktogrammen in codierter Farbigkeit – unterliegt auch das Dorf einem durchdachten Orientierungssystem. Das Wegeleitsystem des Designers Otl Aicher ist durch Farben und Symbole (Kreis, Quadrat, Dreieck) gekennzeichnet, wobei sich die Farbigkeit (gelb in der Straßberger-, grün in der Nadi- und blau in der Connollystraße) sowohl an den Decken und Seitenwänden des Fahrgeschosses als auch in den Fußgängerebenen und Wohnbereichen wiederfindet. Innerhalb der Straßenzüge wirkt es durch aufgeständerte, farbige Rohrbahnen, die sog. „Media Linien“ von Hans Hollein, sogar raumbestimmend. Diese spielerisch-dekorativ eingesetzten Elemente schaffen eine eigene Kommunikationsebene und erleichtern generell die Orientierung im Olympischen Dorf.\nGegenüber im Westen befindet sich die Zentrale Hochschulsportanlage. Sie wurde 1972 als Volleyball- und Gymnastikhalle mit Rundfunk- und Fernsehzentrum genutzt. Der Anlage von Erwin Heinle und Robert Wischer liegt eine strenge Rasterstruktur zugrunde. Ihre dementsprechend kubisch wirkenden Bauten leben vom Kontrast zwischen den rostbraunen Teilen des Stahlgerüsts und den hellen Ausfachungen. Über dem zentralen Atriumhof schwebt an einem Stahlrahmen der sog. Lichtsatellit von Otto Piene, ein Glaskörper in Form eines geschliffenen Diamanten. Um die Gebäudegruppe liegen ausgedehnte Sportkampf- und Spielbahnen.\n\nDer Olympiapark hat nachträgliche Eingriffe erfahren. Das vormalige Olympische Dorf der Frauen ist mit Ausnahme von 12 Bungalows vollständig abgebrochen und durch Neubauten ersetzt. Weitgehend hat man zudem die Hochschulsportanlage abgebrochen. Mit der BMW-Welt, dem Sea Life Centre, der sog. Kleine Olympiahalle und dem BFTS- Bau wurden – teils aufgrund ihrer Größe störende – Neubauten in die Gesamtanlage eingefügt. Doch trotz der erwähnten Eingriffe hat der Olympiapark nichts an seiner herausragenden Bedeutung als gebautes Zeugnis für die noch junge Bundesrepublik Deutschland vor 1972 verloren. Er war das wichtigste Großbauprojekt der Bundesrepublik in der Zeit um 1970 und genießt in dieser Hinsicht und in der beschriebenen besonderen Gestaltungsweise internationale Bedeutung und Beachtung.", d.beschreibung, "soso")
		assert.Equal(t, "", d.verfahrensstand, "soso")
	}

	{
		d := ds[6812]
		assert.Equal(t, "D-1-7734-0101", d.aktennummer, "soso")
		assert.Equal(t, "09 1 62 000", d.gemeindeschlüssel, "soso")
		assert.Equal(t, "Bodendenkmäler", d.typ, "soso")
		assert.Equal(t, 0, len(d.adresse), "soso")
		assert.Equal(t, "Grabhügel mit Bestattungen der Hallstattzeit sowie Siedlung vorgeschichtlicher Zeitstellung.", d.beschreibung, "soso")
		assert.Equal(t, "nachqualifiziert", d.verfahrensstand, "soso")
	}
}

func TestFineFromRawNichtNachqualifiziertIssue1(t *testing.T) {
	pdf, err := rawFromXmlFileName("testdata/572115.xml")
	ds, _, _, err := fineFromRaw(pdf)
	assert.Nil(t, err, "soso")

	{
		d := ds[13]
		assert.Equal(t, "D-5-72-115-51", d.aktennummer, "soso")
		assert.Equal(t, "09 5 72 115", d.gemeindeschlüssel, "soso")
		assert.Equal(t, "Baudenkmäler", d.typ, "soso")
		assert.Equal(t, 1, len(d.adresse), "soso")
		assert.Equal(t, "Grenzstein", d.adresse[0], "soso")
		assert.Equal(t, "Jagdgrenzstein, bez. 1565 und 1781; an der Straße nach Bubenreuth.", d.beschreibung, "soso")
		assert.Equal(t, "nicht nachqualifiziert, im Bayerischen Denkmal-Atlas nicht kartiert", d.verfahrensstand, "soso")
	}
}

func TestFineFromRawPageBreakIssue2(t *testing.T) {
	pdf, err := rawFromXmlFileName("testdata/172129.xml")
	assert.Equal(t, 10, len(pdf.Page), "page count")
	{
		assert.Equal(t, 50, len(pdf.Page[2].Text), "text count")
		tx := pdf.Page[2].Text[45]
		assert.Equal(t, int8(3), tx.Font, "text font")
		assert.Equal(t, int32(1116), tx.Top, "text top")
		assert.Equal(t, int32(204), tx.Left, "text left")
		assert.Equal(t, int32(569), tx.Width, "text width")
		assert.Equal(t, "Engertalm. ", tx.Bold, "text bold")
		assert.Equal(t, " Kaser der Engertalm, eingeschossiger überkämmter Blockbau auf", tx.Value, "text value")
	}
	{
		assert.Equal(t, 50, len(pdf.Page[3].Text), "text count")
		tx := pdf.Page[3].Text[3]
		assert.Equal(t, int8(4), tx.Font, "text font")
		assert.Equal(t, int32(185), tx.Top, "text top")
		assert.Equal(t, int32(204), tx.Left, "text left")
		assert.Equal(t, int32(300), tx.Width, "text width")
		assert.Equal(t, "", tx.Bold, "text bold")
		assert.Equal(t, "südöstlich unterm Gernhorn, 965m Höhe.", tx.Value, "text value")
	}

	ds, _, _, err := fineFromRaw(pdf)

	assert.Nil(t, err, "soso")
	{
		d := ds[22]
		assert.Equal(t, "D-1-72-129-64", d.aktennummer, "soso")
		assert.Equal(t, "09 1 72 129", d.gemeindeschlüssel, "soso")
		assert.Equal(t, "Baudenkmäler", d.typ, "soso")
		assert.Equal(t, 1, len(d.adresse), "soso")
		assert.Equal(t, "Engertalm", d.adresse[0], "soso")
		assert.Equal(t, "Kaser der Engertalm, eingeschossiger überkämmter Blockbau auf Bruchsteinsockel, Flachsatteldach mit Legschindeldeckung, Firstpfette bez. 1801; südöstlich unterm Gernhorn, 965m Höhe.", d.beschreibung, "soso")
		assert.Equal(t, "nachqualifiziert", d.verfahrensstand, "soso")
	}
}

func TestFineFromRawWagingIssue4(t *testing.T) {
	pdf, err := rawFromXmlFileName("testdata/189162.xml")
	ds, mod, name, err := fineFromRaw(pdf)
	assert.Nil(t, err, "soso")
	assert.Equal(t, "2017-08-19T00:00:00+02:00", mod.Format(time.RFC3339), "huhu")
	assert.Equal(t, "Waging a.See", name, "huhu")
	// Equal.assert(t, "Regierungsbezirk Oberbayern", l.gemeinde.regierungsbezirk, "soso")
	// assert.Equal(t, "Traunstein", l.gemeinde.landkreis, "soso")
	// assert.Equal(t, "Übersee", l.gemeinde.gemeinde, "soso")
	assert.Equal(t, 113, len(ds), "soso")
	{
		d := ds[0]
		assert.Equal(t, "E-1-89-162-1", d.aktennummer, "soso")
		assert.Equal(t, "09 1 89 162", d.gemeindeschlüssel, "soso")
		assert.Equal(t, "Baudenkmäler", d.typ, "soso")
		assert.Equal(t, 1, len(d.adresse), "soso")
		assert.Equal(t, "Ensemble Ortskern Markt Waging", d.adresse[0], "soso")
		assert.Equal(t, "Das Ensemble umfasst die vier, aus verschiedenen Richtungen an dem kleinen Marktplatz zusammentreffenden Gassen des Marktortes mit ihrer historischen Bebauung. - Waging, im Voralpenland nahe dem Westufer des Waginger Sees gelegen, war bereits in keltischer und römischer Zeit besiedelt, wird im 8. Jh. erstmals als Besitz des Salzburger Nonnbergklosters genannt, erhielt 1385 Marktrechte und gehörte bis 1803 zum Erzstift Salzburg. - Als Salzburger- und Bahnhofstraße durchzieht in gewundenem Lauf eine alte Durchgangsstraße, die sog. Untere Salzstraße, den Ort, in dessen Mitte sie sich zu einem kleinen Marktplatz ausweitet. Innerhalb des historischen Ortsbereichs, der ehemals durch hölzerne Gatter abgegrenzt war, ist dieser Straßenzug im Gegensatz zu der außerhalb folgenden offenen Bebauung des späteren 19. und 20. Jh. geschlossen bebaut; er weist zwei- und dreigeschossige Wohn-, Handwerker- und Gasthäuser auf, die meist dem späteren 18. und dem 19. Jh. entstammen, im Kern aber oft älter sind. Es handelt sich ausschließlich um Putzbauten, einige mit Putzgliederungen und Stuckdekor an den Fronten. Ein Teil der Häuser erinnert mit seinen weit vorstehenden Flachsatteldächern an den älteren hölzernen Haustyp, der nach den zahlreichen Ortsbränden vom 17. bis zum 19. Jh. mehr und mehr zurückgedrängt wurde. Ein anderer Teil ist dem Haustyp der Inn-Salzach-Städte mit hinter Blendgiebeln und Vorschussmauern versenkten Dächern verpflichtet. - In der Bahnhofstraße manifestiert sich in dem ehem. Salzburgischen Pfleggerichtsgebäude, jetzt Schwemmbräu, die erzstiftische Herrschaft über den Ort; zugleich lässt auch das große, vorkragende Krüppelwalmdach salzburgischen Einfluss erkennen. Ein ähnlicher Bau ist der große, den Marktplatz beherrschende ehem. Brauereigasthof. - In der Salzburger Straße dominiert die große 1878 entstandene, dreiteilige Front des Hotels Waginger Hof das Straßenbild; der Bau macht gleichzeitig die Anfänge der Entwicklung von Waging als Fremdenverkehrsort deutlich. - Die südöstliche der beiden Nebengassen, die Wilhelm- Scharnow-Straße, ist eine Handwerkergasse, die sich durch lebendige Vielfalt ihrer Häuserfronten und Dachformen sowie malerische Durchblicke auszeichnet, während das Straßenbild der nördlichen Gasse, der Seestraße, von der Pfarrkirche St. Martin und dem Martinihof, dem ehem. Pfarrhof, bestimmt wird. - Die hochgelegene, über der abschüssigen Gasse und einer Zeile gut erhaltener bürgerlicher Giebelhäuser des frühen 19. Jh. aufragende Kirche ist eine nach dem Brand von 1611 neu errichtete Wandpfeileranlage, die bis in das 19. Jh. weiter ausgebaut wurde. Sie ist vom ehem. Kirchhof umgeben, dessen hohe Stützmauern einen Teil der Seestraße einfassen. Der Pfarrhof ist ein strenger, schlossartiger Walmdachbau des frühen 18. Jh., der die Handwerkeranwesen der nördlichen Seestraße eindrucksvoll überragt. - Mariensäule und Brunnen, beide aus dem Jahr 1854, setzen am Übergang zwischen Marktplatz und Seestraße einen städtebaulich bedeutenden Akzent.", d.beschreibung, "soso")
		assert.Equal(t, "", d.verfahrensstand, "soso")
	}
	{
		d := ds[1]
		assert.Equal(t, "D-1-89-162-1", d.aktennummer, "soso")
		assert.Equal(t, "09 1 89 162", d.gemeindeschlüssel, "soso")
		assert.Equal(t, "Baudenkmäler", d.typ, "soso")
		assert.Equal(t, 1, len(d.adresse), "soso")
		assert.Equal(t, "Bahnhofstraße 17", d.adresse[0], "soso")
		assert.Equal(t, "nachqualifiziert", d.verfahrensstand, "soso")
	}
}
