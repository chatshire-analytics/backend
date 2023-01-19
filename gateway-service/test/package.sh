#!/usr/local/bin/bash

rm -rf package/*

for d in */; do
  base=$(basename "$d")

  if [ $base == "package" ] || [ $base == "compressed" ]; then
    continue
  fi

  gotouch package --source "./$base" --target ./package/
done
