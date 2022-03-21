# DBLP-Citation-network V13 converter for neo4j

Small idiomatic go program to convert the data from [aminer][1] to csv
files which can be imported with neo4j.

This is tailored for v11 and v13, but it can be adapted easily to other
versions. It can process the whole 17Gb file using minimal memory (it only processes
an article at a time) and at ~12.000 entries per second.

Currently it only takes the relevant fields which I needed for a graph algorithm
analysis (and it adds some fake data such as venue location) but this can be
disabled easily in the code.

Note: the original JSON file from aminer contains mongodb JSON v2 type annotations which
should be dropped. You can use: `sed -i 's/NumberInt(\([0-9]*\))/\1/' input.json` or
any other method.

[1]: https://www.aminer.org/citation
