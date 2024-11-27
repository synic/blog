#!/bin/bash

# This script will use git commit times to determine if markdown articles
# are newer than the json output. This is because git doens't actually restore
# the actual modified time, so if an article is edited directly on github, the
# build system won't know that it is newer than the output json, and it won't
# be converted. Here we use the commit time to check both files, and if the
# markdown file is newer, delete the json file, which will cause the build
# system to re-create it.

for mdfn in ./articles/*.md; do
  mdct=$(git log -1 --format="%ct" -- "${mdfn}")
  jsonfn=$(echo "${mdfn}" | sed 's/\.md$/\.json/' | sed 's#^\./articles#\./assets/articles#')
  jsonct=$(git log -1 --format="%ct" -- "${jsonfn}")

  if [[ ! -f "$jsonfn" ]]; then
    continue
  fi

  if [[ $mdct -gt $jsonct ]]; then
    echo "removing ${jsonfn} because it is older..."
    rm "$jsonfn"
  fi
done
