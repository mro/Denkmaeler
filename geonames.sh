#!/bin/bash
cd "$(dirname "${0}")"
#  Copyright (c) 2016-2016 Marcus Rohrmoser, http://mro.name/~me

# use geonames data for correlation rather than dbpedia which is less accurate.
#
# See e.g.
# - wrong gemeindeschlüssel of http://dbpedia.org/page/Reit_im_Winkl
# - geoname of http://dbpedia.org/page/Wildpoldsried does not point to ADM4
#
# # SPARQL endpoint http://dbpedia.org/sparql
# PREFIX dbp:<http://dbpedia.org/property/>
# PREFIX owl:<http://www.w3.org/2002/07/owl#>
# select distinct ?url,?ags,?geoname where {
# ?url
#   dbp:state "Bayern"@en ;
#   dbp:gemeindeschlüssel ?ags ;
#   owl:sameAs ?geoname .
#  FILTER regex(str(?geoname), "^http://sws\\.geonames\\.org/", 'i')
# } ORDER BY ?ags
#

curl --silent --location --remote-time --time-cond DE.zip --output DE.zip --url "http://download.geonames.org/export/dump/DE.zip"
unzip -u -d DE DE.zip

fgrep ADM4 DE/DE.txt | cut -f 1,14,8 | while read geoname level gemeinde name
do
  [ "" = "${gemeinde}" ] && continue
  [ "ADM4" = "${level}" ] || continue
  land="${gemeinde:0:2}"
  reg="${gemeinde:2:1}"
  kreis="${gemeinde:3:2}"
  gem="${gemeinde:5:3}"
  dir="build/${land}/${reg}/${kreis}/${gem}"

  [ -d "${dir}" ] || continue
  # echo "${dir} http://sws.geonames.org/${geoname}/ ${name} ${level}"
  printf "."

  [ -r "${dir}/geonames.url" ] || echo "http://sws.geonames.org/${geoname}/" > "${dir}/geonames.url"
done

# ls -Al build/??/?/??/???/geonames.url
