#!/usr/bin/env bash

# Name of neo4j database to use
DATABASE_NAME="${DATABASE_NAME:-dblp.db}"

# Output direcotry containing all the csv files
DATA_DIR="${DATA_DIR:-"$(dirname "$0")/data"}"
cd "${DATA_DIR}"

neo4j-admin import \
    --database="$DATABASE_NAME" \
    --delimiter ";" \
    --array-delimiter "|" \
    --id-type INTEGER \
    --multiline-fields=true \
    --nodes=paper="articles_header.csv,articles.csv" \
    --nodes=author="authors_header.csv,authors.csv" \
    --nodes=journal:venue="journal_header.csv,journal.csv" \
    --nodes=workshop:venue="journal_header.csv,workshop.csv" \
    --nodes=conference:venue="journal_header.csv,conference.csv" \
    --nodes=keyword="keywords_header.csv,keywords.csv" \
    --nodes=edition:publication="edition_header.csv,edition.csv" \
    --nodes=volume:publication="volume_header.csv,volume.csv" \
    --relationships=published_in="rel_published_header.csv,rel_published.csv" \
    --relationships=authored_by="rel_authored_header.csv,rel_authored.csv" \
    --relationships=has_topic="rel_keywords_header.csv,rel_keywords.csv" \
    --relationships=cites="rel_cites_header.csv,rel_cites.csv" \
    --relationships=from="rel_belongs_header.csv,rel_belongs.csv" \
    --relationships=reviews="rel_reviews_header.csv,rel_reviews.csv" \
    "$@"
