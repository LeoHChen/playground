#!/bin/bash

FILE=$1

if [ -z $FILE ]; then
   echo "$0 consensus.s?.log"
   exit
fi

OUT=${FILE##consensus.}

jq -c '{time:.time, block:.blockNum, msg:.message}' $FILE > $OUT
