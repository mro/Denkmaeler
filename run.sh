#!/bin/sh
cd "$(dirname "$0")"

pdftohtml 2>/dev/null
[ 1 -eq $? ] || { echo "Please install pdftohtml." 1>&2 && exit 1 ; }

dst="build"

xml2ttl="denkmaeler-xml2ttl-cmd/denkmaeler-xml2ttl"-*-*-"0.0.1"
[ -x ${xml2ttl} ] || { echo "I need the transformation tool, please run \$ sh build.sh" 1>&2 && exit 1; }

for gemeinde in 162000 189159 ; do mkdir -p build/${gemeinde} ; done

for gemeinde in `ls build`
do
  pdf="${dst}/${gemeinde}/liste.pdf"
  xml="${dst}/${gemeinde}/liste.xml"
  ttl="${dst}/${gemeinde}/liste.ttl"
  url="http://geodaten.bayern.de/denkmal_static_data/externe_denkmalliste/pdf/denkmalliste_merge_${gemeinde}.pdf"
  http_code="$(curl --silent --write-out "%{http_code}" --location --remote-time --create-dirs --output "${pdf}" --time-cond "${pdf}" --url "${url}")"
  if [ "200" = "${http_code}" ] ; then
    pdftohtml -i -stdout -xml "${pdf}" > "${xml}" \
    && ${xml2ttl} < "${xml}" > "${ttl}" \
    && touch -r "${pdf}" "${ttl}"
  else
    echo "ignore http_code ${http_code} ${url}"
  fi
done

ls -Altr "${dst}"/*/*.ttl
ls "${dst}"/*/*.ttl 1>&2

rsync -avPz --delete --delete-excluded --exclude *.pdf --exclude *.xml build/ vario:~/mro.name/linkeddata/open/denkmal/DE/
