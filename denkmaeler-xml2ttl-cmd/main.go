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
)

func main() {

	raw, err := rawFromXmlReader(os.Stdin)
	if nil != err {
		log.Fatal("aua")
	}
	ds, _, err := fineFromRaw(raw)
	if nil != err {
		log.Fatal("aua")
	}

	fmt.Printf("@prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .\n")
	fmt.Printf("@prefix dct: <http://purl.org/dc/terms/> .\n")
	fmt.Printf("@prefix gn: <http://www.geonames.org/ontology#> .\n")
	fmt.Printf("@prefix c: <http://nwalsh.com/rdf/contacts#> .\n")
	fmt.Printf("\n")

	for _, d := range ds {
		fmt.Printf("<http://linkeddata.mro.name/open/country/DE/AGS/%s/denkmal.rdf#%s>\n", strings.Replace(d.gemeindeschlüssel, " ", "/", -1), d.aktennummer)
		fmt.Printf("  dct:identifier \"%s\" ;\n", d.aktennummer)
		fmt.Printf("  gn:admin4Code \"%s\" ;\n", d.gemeindeschlüssel) // http://gis.stackexchange.com/q/7688
		fmt.Printf("  dct:subject <http://www.geodaten.bayern.de/denkmaltyp#%s> ;\n", d.typ)
		for _, a := range d.adresse {
			fmt.Printf("  c:street \"\"\"%s\"\"\" ;\n", a)
		}
		fmt.Printf("  dct:description \"\"\"%s\"\"\" ;\n", strings.Replace(d.beschreibung, "\"", "\\\"", -1))
		fmt.Printf(".\n")
	}
}

func commandHelp() {
	program := os.Args[0]
	fmt.Printf("Usage: %s < foo.xml > foo.ttl\n", program)
	fmt.Printf("\n")
}
