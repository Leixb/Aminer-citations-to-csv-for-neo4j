#!/usr/bin/env bash

cd ~/Documents/json_stream_parser/

neo4j-admin import \
    --database=test.db \
    --delimiter ";" \
    --array-delimiter "|" \
    --id-type INTEGER \
    --nodes=article="articles_header.csv,articles.csv" \
    --nodes=author="authors_header.csv,authors.csv" \
    --nodes=journal="journal_header.csv,journal.csv" \
    --nodes=workshop="journal_header.csv,workshop.csv" \
    --nodes=conference="journal_header.csv,conference.csv" \
    --nodes=keyword="keywords_header.csv,keywords.csv" \
    --nodes=edition="edition_header.csv,edition.csv" \
    --nodes=volume="volume_header.csv,volume.csv" \
    --relationships=published_in="rel_published_header.csv,rel_published.csv" \
    --relationships=authored_by="rel_authored_header.csv,rel_authored.csv" \
    --relationships=has_topic="rel_keywords_header.csv,rel_keywords.csv" \
    --relationships=cites="rel_cites_header.csv,rel_cites.csv" \
    --relationships=reviews="rel_reviews_header.csv,rel_reviews.csv" \
    --relationships=belongs="rel_belongs_header.csv,rel_belongs.csv" \
    "$@"
