#! /bin/bash

set -e

rm -rf release
gox -verbose
mkdir -p release
mv github-clone-all_* release/
cd release
for bin in *; do
    mv "$bin" github-clone-all
    zip "${bin}.zip" github-clone-all
    rm github-clone-all
done
