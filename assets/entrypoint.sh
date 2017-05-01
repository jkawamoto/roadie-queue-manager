#!/bin/bash
#
# entrypoint.sh
#
# Copyright (c) 2017 Junpei Kawamoto
#
# This file is part of Roadie queue manager.
#
# Roadie queue manager is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Roadie queue manager is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with Roadie queue manager. If not, see <http://www.gnu.org/licenses/>.
#

# This template is an entrypoint of a docker container to execute run steps.
#
if [[ $# != 0 ]]; then
  exec $@
fi

extract_zip(){
  echo "Unzipping $1"
  unzip -o -d $(dirname $1) $1
  rm $1
}

unpack_targz(){
  echo "Unpacking $1"
  $(cd $(dirname $1) && tar -zxvf $1)
  rm $1
}

unpack_tar(){
  echo "Unpacking $1"
  $(cd $(dirname $1) && tar -xvf $1)
  rm $1
}

{{with .Git}}
  echo "Cloning git repository {{.}}"
  git clone {{.}} .
{{end}}

{{range .Downloads}}
  echo "Downloading {{.Src}}"
  if [[$(curl -I -H 'Accept-Encoding: gzip,deflate' {{.Src}} 2>/dev/null | grep "Content-Encoding" | grep "gzip" | wc -l) == 1]]; then
    curl -L -o /tmp/gzippedfile {{.Src}}
    gzip -dc /tmp/gzipppedfile > {{.Dest}}
  else
    curl -L -o {{.Dest}} {{.Src}}
  fi
  {{if .Zip}}
    extract_zip {{.Dest}}
  {{else if .TarGz}}
    unpack_targz {{.Dest}}
  {{else if .Tar}}
    unpack_tar {{.Dest}}
  {{end}}
{{end}}

{{range .GSFiles}}
  echo "Downloading {{.Src}}"
  gsutil cp {{.Src}} {{.Dest}}
  {{if .Zip}}
    extract_zip {{.Dest}}
  {{else if .TarGz}}
    unpack_targz {{.Dest}}
  {{else if .Tar}}
    unpack_tar {{.Dest}}
  {{end}}
{{end}}

if [[ -e requirements.txt ]]; then
  echo "Installing required python packages defined in requirements.txt"
  pip install --exists-action i -r requirements.txt
fi

export LC_ALL=C
{{range $index, $elements := .Run}}
echo "{{.}}"
sh -c "{{.}}" > /tmp/stdout{{$index}}.txt
{{end}}

echo "Uploading stdouts"
{{$result := .Result}}
gsutil -m cp "/tmp/stdout*.txt" {{$result}}
{{range .Uploads}}
  echo "Uploading {{.}}"
  gsutil -m cp "{{.}}" {{$result}}
{{end}}
