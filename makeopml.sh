#!/bin/sh
cd "$(dirname "${0}")/build"

cwd="$(pwd)"
opml="${cwd}/denkmal.opml"

cat > "${opml}" <<FOO
<?xml version="1.0" encoding="utf8"?>
<?xml-stylesheet type="text/xsl" href="assets/opml2html.xslt"?>
<opml xmlns:a="http://www.w3.org/2005/Atom" xmlns:dct="http://purl.org/dc/terms/" version="2.0">
  <!--
    Lizenz: CC BY-SA 3.0 DE
  -->
  <!--
    <a:link rel='license'>http://creativecommons.org/licenses/by-sa/3.0/de/</a:link>
    <a:link rel='self'></a:link>
    <a:link rel='via'>http://www.ardmediathek.de/tv/sendungen-a-z?sendungsTyp=sendung</a:link>
    <a:link rel='hub'></a:link>
    validates against https://raw.githubusercontent.com/mro/opml-schema/hotfix/typo/schema.rng
  -->
  <head>
    <title>unverbindliche Denkmallisten des Bayerischen Landesamts für Denkmalpflege</title>
    <!-- <dateCreated/> see file timestamp -->
    <ownerId>http://purl.mro.name/denkmaeler</ownerId>
  </head>
  <body>
FOO

for bundesland in `find . -mindepth 1 -maxdepth 1 -type d -name '??' | cut -c3- | sort` ; do  
  cd "${cwd}/${bundesland}"
  echo "    <outline id='AGS-${bundesland}' text='$(head -n 1 README.txt)'>" >> "${opml}"
  for regbez in `find . -mindepth 1 -maxdepth 1 -type d -name '?' | cut -c3- | sort` ; do
    cd "${cwd}/${bundesland}/${regbez}"
    echo "      <outline id='AGS-${bundesland}-${regbez}' text='$(head -n 1 README.txt)'>" >> "${opml}"
    for landkreis in `find . -mindepth 1 -maxdepth 1 -type d -name '??' | cut -c3- | sort` ; do
      cd "${cwd}/${bundesland}/${regbez}/${landkreis}"
      echo "        <outline id='AGS-${bundesland}-${regbez}-${landkreis}' text='$(head -n 1 README.txt)'>" >> "${opml}"
      for gemeinde in `find . -mindepth 1 -maxdepth 1 -type d -name '???' | cut -c3- | sort` ; do
        cd "${cwd}/${bundesland}/${regbez}/${landkreis}/${gemeinde}"
        xmlUrl="${bundesland}/${regbez}/${landkreis}/${gemeinde}/denkmal.rdf"
        htmlUrl="${xmlUrl}"
        echo "          <outline id='AGS-${bundesland}-${regbez}-${landkreis}-${gemeinde}' text='$(head -n 1 README.txt)' language='de' type='rdf' version='rdf' xmlUrl='${xmlUrl}' htmlUrl='${htmlUrl}'/>" >> "${opml}"
      done
      echo "        </outline>" >> "${opml}"
    done
    echo "      </outline>" >> "${opml}"
  done
  echo "    </outline>" >> "${opml}"
done

cat >> "${opml}" <<FOO
  </body>
</opml>
FOO
