# DBLP-Citation-network V13 converter for neo4j

This is a small go program to convert the data from [aminer][1] to CSV
so that they can be imported with `neo4j-admin import`. This was used for
a course exercise, so some assumptions and dummy data is added in some cases,
the `.sh` scripts provided can easily be modified to only import the data that
is relevant for your analysis.

This is tailored for v11 and v13, but it can be adapted easily to other
versions. It can process the whole 17Gb file using minimal memory (it only processes
an article at a time) and at ~12.000 entries per second.

## Usage

Build the program with `go build`, run the program `./parser`, you can
optionally pass arguments to configure the input and output folder, (see
`./parser -help` for details)

Once the parser finishes (it should take around 7 minutes, depending on you
computer specs), all the CSV data will be in the `data` folder. You can then run
`./neo4j_import_A2.sh` (with `--force` if necessary) to load the data into
neo4j. This should take around 4 minutes, depending on your Neo4J
configuration.

There are 2 import scripts, `./neo4j_import_A2.sh` and `./neo4j_import_A3.sh`
corresponding to the solutions of the tasks proposed in the course.

However, the parser only needs to run once, since creating all the CSV files at
once was more efficient than parsing the file twice and generating slightly
different outputs.

### Note

If you use the original JSON file from aminer contains mongodb extended JSON v2
annotations which should be dropped. You can use:
`sed -i 's/NumberInt(\([0-9]*\))/\1/' input.json` or any tool to clean it.
Additionally, the last entry is corrupted and contains decimal points, so you
have to remove them or discard it manually.

Since this may not be a trivial task on such a big file, we provide a compressed
zip file with the fixed JSON file and another with the computed CSV files
[here (only upc.edu accounts)][2]

[1]: https://www.aminer.org/citation
[2]: https://drive.google.com/drive/folders/1Pz00DOnqoGlOUfqALr5R7EP9-p2EdROm
