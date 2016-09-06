#!/bin/sh
cd "$(dirname "$0")"
#  Copyright (c) 2016-2016 Marcus Rohrmoser, http://mro.name/~me

pdftohtml 2>/dev/null
[ 1 -eq $? ] || { echo "Please install pdftohtml." 1>&2 && exit 1 ; }
rapper --version >/dev/null || { echo "Please install raptor." 1>&2 && exit 1 ; }

dst="build"
deploy_base_url="http://linkeddata.mro.name/open/country/DE/AGS/"

######################################################################
## Prepare build dir
######################################################################

# for gemeinde in 09/1/61/000 09/1/62/000 09/1/63/000 09/1/71/111 09/1/71/112 09/1/71/113 09/1/71/114 09/1/71/115 09/1/71/116 09/1/71/117 09/1/71/118 09/1/71/119
# for gemeinde in 09/1/62/000 09/1/89/114
for gemeinde in `cut -d / -f 2-5 bayern-ags.csv`
do mkdir -p "${dst}/${gemeinde}" ; done

bash geonames.sh

cat > "${dst}/README.txt" <<FOO
Die Gemeinden in Deutschland nach dem Amtlichen Gemeindeschlüssel
https://de.wikipedia.org/wiki/Amtlicher_Gemeindeschl%C3%BCssel

Derzeit Daten ausschließlich für Bayern.
FOO
echo "Bundesland Bayern" > "${dst}/09/README.txt"

write_about() {
  ags_dir="${1}"
  title="${2}"
  geonames_url="$(cat "${dst}/${ags_dir}/geonames.url")"
  dbpedia_url="http://dbpedia.org/page/${title}"

  cat > "${dst}/${ags_dir}/about.ttl" <<FOO
@prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> .
@prefix dct: <http://purl.org/dc/terms/> .
@prefix owl: <http://www.w3.org/2002/07/owl#> .
@prefix dbp: <http://dbpedia.org/property/> .
@prefix geo: <http://www.w3.org/2003/01/geo/wgs84_pos#> .
@prefix cc: <http://creativecommons.org/ns#> .
@prefix foaf: <http://xmlns.com/foaf/0.1/> .

<${deploy_base_url}${ags_dir}/about.rdf>
  cc:attributionName "🏰"^^<http://www.w3.org/2001/XMLSchema#string> ;
  cc:attributionURL <${deploy_base_url}${ags_dir}/about.rdf> ;
  cc:license <http://creativecommons.org/licenses/by-sa/3.0/> ;
  dct:relation <http://purl.mro.name/denkmaeler> ;
  dct:format <http://purl.org/NET/mediatypes/application/rdf+xml> ;
  dct:source <https://web.archive.org/web/20090411151155/http://www.destatis.de/jetspeed/portal/cms/Sites/destatis/Internet/DE/Content/Statistiken/Regionales/Gemeindeverzeichnis/Administrativ/AdministrativeUebersicht,templateId=renderPrint.psml> ;
  a foaf:Document ;
  foaf:primaryTopic <${deploy_base_url}${ags_dir}/> .

<${deploy_base_url}${ags_dir}/>
  a geo:SpatialThing ;
  rdfs:label """${title}"""@de ;
  dct:identifier "$(echo "${ags_dir}" | tr / ' ')" ;
  dbp:gemeindeschlüssel "$(echo "${ags_dir}" | tr -d /)" ;
  dct:relation <${deploy_base_url}${ags_dir}/denkmal.rdf> ;
FOO
  [ "" != "${geonames_url}" ] && echo "  owl:sameAs <${geonames_url}> ;" >> "${dst}/${ags_dir}/about.ttl"
  [ "" != "${dbpedia_url}" ] && echo "#  owl:sameAs <${dbpedia_url}> ;" >> "${dst}/${ags_dir}/about.ttl"

  echo "." >> "${dst}/${ags_dir}/about.ttl"
  {
    rapper --quiet --input turtle --output rdfxml-abbrev "${dst}/${ags_dir}/about.ttl" > "${dst}/${ags_dir}/about.rdf~"
    diff -q "${dst}/${ags_dir}/about.rdf" "${dst}/${ags_dir}/about.rdf~" 2>/dev/null || cp "${dst}/${ags_dir}/about.rdf~" "${dst}/${ags_dir}/about.rdf"
    rm "${dst}/${ags_dir}/about.rdf~"
  }
}

rsync -aP static/assets "${dst}/"

wait

foot="\n\nSiehe http://purl.mro.name/denkmaeler"
while IFS='/' read ignore bundesland regierungsbezirk landkreis gemeinde name
do
  if [ "" = "${landkreis}" ] ; then
    rme="${dst}/${bundesland}/${regierungsbezirk}/README.txt"
    [ -r "${rme}" ] || echo "${name}${foot}" > "${rme}"
  else
    if [ "" = "${gemeinde}" ] ; then
      rme="${dst}/${bundesland}/${regierungsbezirk}/${landkreis}/README.txt"
      [ -r "${rme}" ] || echo "Landkreis ${name}${foot}" > "${rme}"
    else
      rme="${dst}/${bundesland}/${regierungsbezirk}/${landkreis}/${gemeinde}/README.txt"
      [ -r "${rme}" ] || echo "Gemeinde ${name}${foot}" > "${rme}"
      write_about "${bundesland}/${regierungsbezirk}/${landkreis}/${gemeinde}" "${name}" &
    fi
  fi
done < bayern-ags.csv

wait

######################################################################
## fetch and scrape PDF, turn into ttl and RDF
######################################################################

xml2ttl="denkmaeler-xml2ttl-cmd/denkmaeler-xml2ttl"-*-*-"0.0.1"
[ -x ${xml2ttl} ] || { echo "I need the transformation tool, please run \$ sh build.sh" 1>&2 && exit 1; }

for gemeinde in `ls -d build/??/?/??/??? | cut -d / -f2-`
do
  printf "%s " "${gemeinde}"
  geonames_url="$(cat "${dst}/${gemeinde}/geonames.url")"
  file="${dst}/${gemeinde}/denkmal"
  nummer="$(echo ${gemeinde} | cut -d / -f 2- | tr -d /)"
  bayern_prefix="09"
  url="http://geodaten.bayern.de/denkmal_static_data/externe_denkmalliste/pdf/denkmalliste_merge_${nummer}.pdf"
  deploy_url="${deploy_base_url}${gemeinde}/"
  base_url="http://geodaten.bayern.de/"
  geonames_url="$(cat "${dst}/${gemeinde}/geonames.url")"
  pdf="${file}.pdf"
  http_code="$(curl --user-agent http://github.com/mro/Denkmaeler --silent --write-out "%{http_code}" --location --remote-time --create-dirs --output "${pdf}" --time-cond "${pdf}" --url "${url}")"
  if [ "200" = "${http_code}" ] ; then
    {
      ttl="${file}.ttl"
      cat > "${ttl}" <<FOO
@prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> .
@prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .
@prefix owl: <http://www.w3.org/2002/07/owl#> .
@prefix dct: <http://purl.org/dc/terms/> .
@prefix foaf: <http://xmlns.com/foaf/0.1/> .
@prefix geo: <http://www.w3.org/2003/01/geo/wgs84_pos#> .
@prefix dbp: <http://dbpedia.org/property/> .
@prefix cc: <http://creativecommons.org/ns#> .

<${deploy_url}denkmal.rdf>
  cc:attributionName "🏰"^^<http://www.w3.org/2001/XMLSchema#string> ;
  cc:attributionURL <${deploy_url}denkmal.rdf> ;
  cc:license <http://creativecommons.org/licenses/by-sa/3.0/> ;
  dct:format <http://purl.org/NET/mediatypes/application/rdf+xml> ;
  dct:hasFormat <${url}> ;
  dct:source <${url}> ;
FOO
  [ "" != "${geonames_url}" ] && echo "  dct:spatial <${geonames_url}> ;" >> "${ttl}"

      cat >> "${ttl}" <<FOO
  dct:relation <http://purl.mro.name/denkmaeler> ;
  dct:created "2016-08-31"^^<http://www.w3.org/2001/XMLSchema#date> ;
  a foaf:Document ;
  foaf:primaryTopic <${deploy_url}> .

FOO
      xml="${file}.xml"
      rdf="${file}.rdf"
      pdftohtml -i -stdout -xml "${pdf}" > "${xml}" \
      && ${xml2ttl} < "${xml}" >> "${ttl}" \
      && touch -r "${pdf}" "${ttl}" \
      && rapper --quiet -i turtle -o rdfxml-abbrev "${ttl}" > "${rdf}" \
      && sed -i'~' '1s:<.xml version=.*:<?xml version="1.0" encoding="utf-8"?><?xml-stylesheet type="text/xsl" href="../../../../assets/denkmal2html.xslt"?>:' "${rdf}" \
      && rm "${rdf}~" \
      && touch -r "${pdf}" "${rdf}"
    } &
  else
    echo "ignore http_code ${http_code} ${url}"
  fi
done

wait

sh makeopml.sh

cd "${dst}/.git/.." && {
  git add --all . \
  && git commit -m "update" \
  && git update-server-info
  cd -
}

ls -Altr "${dst}"/??/?/??/*/*.rdf
ls "${dst}"/??/?/??/*/*.rdf 1>&2

######################################################################
## deploy txt, RDF and denkmal.git
######################################################################

rsync -avPz --delete --delete-excluded --exclude .??* --exclude *.pdf --exclude *.url --exclude *.xml --exclude *.ttl "${dst}/" vario:~/"mro.name/linkeddata/open/country/DE/AGS/"
rsync -avPz --delete "${dst}/.git/" vario:~/"mro.name/linkeddata/open/country/DE/AGS/denkmal.git"
