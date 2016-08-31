#!/bin/sh
cd "$(dirname "${0}")"

curl --silent --location --remote-time --time-cond DE.zip --output DE.zip --url "http://download.geonames.org/export/dump/DE.zip"
unzip -ud DE DE.zip

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
  echo "http://sws.geonames.org/${geoname}/" > "${dir}/geonames.url"
done

# ls -Al build/??/?/??/???/geonames.url
