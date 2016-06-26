#!/bin/bash

# format compress inputzipfilename


tmpdir=`mktemp -d`

for zipfilename in $1/*.zip
do
    echo "Expanding ${zipfilename}"
    unzip -qq -d $tmpdir $zipfilename

    fltfile="${tmpdir}/*.flt"
    xmlfile="${tmpdir}/*meta.xml"

    map_convert ${xmlfile} ${fltfile}

    rm -rf ${tmpdir}/*
done


