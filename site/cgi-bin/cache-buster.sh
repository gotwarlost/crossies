#!/bin/bash

if [[ -z "${GATEWAY_INTERFACE}" ]]
then
  true
else
  exit 1
fi

d=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
for file in $(find . -type f \( -name "*.html" -o -name "*.js" \))
do
  sed -i "s/?bust=1/?bust=$d/g" "${file}"
done
