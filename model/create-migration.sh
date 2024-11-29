#!/usr/bin/env bash

cd $(dirname $0)

timestamp=$(date +%s)

usage="Usage: $0 <migration-name>"
if [ $# -ne 1 ]; then
	echo "$usage"
	exit 1
fi

touch "./migrations/${timestamp}_${1}.up.sql"
touch "./migrations/${timestamp}_${1}.down.sql"
