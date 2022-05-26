#!/usr/bin/env bash

split articles.csv  papers_ -n l/4 -a 1 --additional-suffix .csv
split conference.csv  conferences_ -n l/3 -a 1 --additional-suffix .csv

for f in papers_*; do
    # if header not set, set it
    if [ -z "$HEADER" ]; then
        HEADER="$(head -n 1 $f)"
    else
        # add header to start of file
        sed -i "1s/^/$HEADER\n/" $f
    fi
done

# unset variable
unset HEADER

for f in conferences_*; do
    # if header not set, set it
    if [ -z "$HEADER" ]; then
        HEADER="$(head -n 1 $f)"
    else
        # add header to start of file
        sed -i "1s/^/$HEADER\n/" $f
    fi
done
