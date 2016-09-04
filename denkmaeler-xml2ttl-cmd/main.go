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
	"fmt"
	"log" // log.Fatal
	"os"
	"strings"
	"time"
)

func main() {

	raw, err := rawFromXmlReader(os.Stdin)
	if nil != err {
		log.Fatal("aua")
	}
	ds, date, err := fineFromRaw(raw)
	if nil != err {
		log.Fatal("aua")
	}

	fmt.Printf("@prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .\n")
	fmt.Printf("@prefix dct: <http://purl.org/dc/terms/> .\n")
	fmt.Printf("@prefix gn: <http://www.geonames.org/ontology#> .\n")
	fmt.Printf("@prefix c: <http://nwalsh.com/rdf/contacts#> .\n")
	fmt.Printf("\n")

	url := ""
	// data per entry
	for idx, d := range ds {
		if "" == url {
			url = "http://linkeddata.mro.name/open/country/DE/AGS/" + strings.Replace(d.gemeindeschlüssel, " ", "/", -1) + "/denkmal.rdf"
		}
		fmt.Printf("<%s#%s>\n", url, d.aktennummer)
		fmt.Printf("  # %d\n", idx)
		fmt.Printf("  dct:identifier \"%s\" ;\n", d.aktennummer)
		fmt.Printf("  gn:admin4Code \"%s\" ;\n", d.gemeindeschlüssel) // http://gis.stackexchange.com/q/7688
		fmt.Printf("  dct:subject <http://www.geodaten.bayern.de/denkmaltyp#%s> ;\n", d.typ)
		for _, a := range d.adresse {
			fmt.Printf("  c:street \"\"\"%s\"\"\" ;\n", a)
		}
		fmt.Printf("  dct:description \"\"\"%s\"\"\" .\n", strings.Replace(d.beschreibung, "\"", "\\\"", -1))
	}
	// Lists, Order
	for _, typ := range []string{"Baudenkmäler", "Bodendenkmäler"} {
		fmt.Printf("<%s#%s>\n", url, typ)
		fmt.Printf("  dct:title \"\"\"%s\"\"\" ; \n", typ)
		fmt.Printf("  dct:hasPart [ \n")
		for idx, d := range ds {
			if typ != d.typ {
				continue
			}
			fmt.Printf("    rdf:_%04d <%s#%s> ; \n", 1+idx, url, d.aktennummer)
		}
		fmt.Printf("  a rdf:Seq ] . \n")
	}
	// Modified date
	fmt.Printf("<%s>\n", url)
	fmt.Printf("  dct:modified \"%s\" .\n", date.Format(time.RFC3339))
}

func commandHelp() {
	program := os.Args[0]
	fmt.Printf("Usage: %s < foo.xml > foo.ttl\n", program)
	fmt.Printf("\n")
}
