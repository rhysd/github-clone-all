#! /bin/bash

set -e

rm -rf release
gox -verbose
mkdir -p release
mv github-clone-all_* release/
cd release
for bin in *; do
    if [[ "$bin" == *windows* ]]; then
        command="github-clone-all.exe"
    else
        command="github-clone-all"
    fi
    mv "$bin" "$command"
    zip "${bin}.zip" "$command"
    rm "$command"
done
