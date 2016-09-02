#!/bin/sh
cd "$(dirname "$0")"

pdftohtml 2>/dev/null
[ 1 -eq $? ] || { echo "Please install pdftohtml." 1>&2 && exit 1 ; }
rapper --version >/dev/null || { echo "Please install raptor." 1>&2 && exit 1 ; }

dst="build"
deploy_base_url="http://linkeddata.mro.name/open/country/DE/AGS/"

######################################################################
## Prepare build dir
######################################################################

# for gemeinde in 09/4/77/148 09/1/62/000 09/1/89/159 ; do mkdir -p "${dst}/${gemeinde}" ; done
for gemeinde in `cut -d / -f 2-5 bayern-ags.csv` ; do mkdir -p "${dst}/${gemeinde}" ; done
bash geonames.sh

cat > "${dst}/README.txt" <<FOO
Die Gemeinden in Deutschland nach dem Amtlichen Gemeindeschlüssel
https://de.wikipedia.org/wiki/Amtlicher_Gemeindeschl%C3%BCssel

Derzeit Daten ausschließlich für Bayern.
FOO
echo "Bundesland Bayern" > "${dst}/09/README.txt"
echo "Regierungsbezirk Oberbayern" > "${dst}/09/1/README.txt"

write_about() {
  ags_dir="${1}"
  title="${2}"
  geonames_url="$(cat "${dst}/${ags_dir}/geonames.url")"

  cat > "${dst}/${ags_dir}/about.ttl" <<FOO
@prefix dct: <http://purl.org/dc/terms/> .
@prefix gn: <http://www.geonames.org/ontology#> .
<${deploy_base_url}${ags_dir}/>
  dct:title """${title}""" ;
  gn:admin4Code "$(echo "${ags_dir}" | tr / ' ')" ;
FOO
  [ "" != "${geonames_url}" ] && echo "  dct:spatial <${geonames_url}> ;" >> "${dst}/${ags_dir}/about.ttl"
  echo "." >> "${dst}/${ags_dir}/about.ttl"

  rapper --quiet --input turtle --output rdfxml-abbrev "${dst}/${ags_dir}/about.ttl" > "${dst}/${ags_dir}/about.rdf"
}

while IFS='/' read ignore bundesland regierungsbezirk landkreis gemeinde name
do
  if [ "" = "${landkreis}" ] ; then
    echo "${name}" > "${dst}/${bundesland}/${regierungsbezirk}/README.txt"
  else
    if [ "" = "${gemeinde}" ] ; then
      echo "Landkreis ${name}" > "${dst}/${bundesland}/${regierungsbezirk}/${landkreis}/README.txt"
      # write_about "${bundesland}/${regierungsbezirk}/${landkreis}/${gemeinde}" "${name}" &
    else
      echo "Gemeinde ${name}" > "${dst}/${bundesland}/${regierungsbezirk}/${landkreis}/${gemeinde}/README.txt"
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
  file="${dst}/${gemeinde}/denkmal-liste"
  nummer="$(echo ${gemeinde} | cut -d / -f 2- | tr -d /)"
  bayern_prefix="09"
  url="http://geodaten.bayern.de/denkmal_static_data/externe_denkmalliste/pdf/denkmalliste_merge_${nummer}.pdf"
  deploy_url="${deploy_base_url}${gemeinde}/"
  base_url="http://geodaten.bayern.de/"
  pdf="${file}.pdf"
  http_code="$(curl --user-agent http://github.com/mro/Denkmaeler --silent --write-out "%{http_code}" --location --remote-time --create-dirs --output "${pdf}" --time-cond "${pdf}" --url "${url}")"
  if [ "200" = "${http_code}" ] ; then
    {
      ttl="${file}.ttl"
      cat > "${ttl}" <<FOO
@prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .
@prefix cc: <http://creativecommons.org/ns#> .
@prefix foaf: <http://xmlns.com/foaf/0.1/> .
@prefix dct: <http://purl.org/dc/terms/> .
@prefix gn: <http://www.geonames.org/ontology#> .

# <${deploy_url}denkmal-liste.rdf>
#    dct:modified "$(date +%F)"^^<http://www.w3.org/2001/XMLSchema#date> .

<${deploy_url}about.rdf>
  cc:attributionName "🏰"^^<http://www.w3.org/2001/XMLSchema#string> ;
  cc:attributionURL <${deploy_url}about.rdf> ;
  cc:license <http://creativecommons.org/licenses/by/3.0/> ;
  dct:spatial <${geonames_url}> ;
  gn:admin4Code "${bayern_prefix}${nummer}" ;
  dct:created "2016-08-31"^^<http://www.w3.org/2001/XMLSchema#date> ;
  dct:source <https://web.archive.org/web/20090411151155/http://www.destatis.de/jetspeed/portal/cms/Sites/destatis/Internet/DE/Content/Statistiken/Regionales/Gemeindeverzeichnis/Administrativ/AdministrativeUebersicht,templateId=renderPrint.psml> ;
  a foaf:Document ;
  foaf:primaryTopic <${deploy_url}> .

<${deploy_url}denkmal-liste.rdf>
  cc:attributionName "🏰"^^<http://www.w3.org/2001/XMLSchema#string> ;
  cc:attributionURL <${deploy_url}denkmal-liste.rdf> ;
  cc:license <http://creativecommons.org/licenses/by/3.0/> ;
  dct:source <${url}> ;
  dct:relation <http://github.com/mro/Denkmaeler> ;
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
      && touch -r "${pdf}" "${rdf}"
    } &
  else
    echo "ignore http_code ${http_code} ${url}"
  fi
done

wait

ls -Altr "${dst}"/??/?/??/*/*.rdf
ls "${dst}"/??/?/??/*/*.rdf 1>&2

######################################################################
## deploy txt and RDF.
######################################################################

rsync -avPz --delete --delete-excluded --exclude .??* --exclude *.pdf --exclude *.url --exclude *.xml --exclude *.ttl "build/" vario:~/"mro.name/linkeddata/open/country/DE/AGS/"
