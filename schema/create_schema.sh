#!/bin/sh

for sqlfile in `ls *.sql`
do
    psql postgres < $sqlfile
done
