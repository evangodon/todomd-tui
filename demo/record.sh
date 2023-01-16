#!/usr/bin/env bash

set -o errexit 
set -o nounset

rm todo.md
cp template.md todo.md

TAPE_FILE=cassette.tape

echo "Set Width 1920" > $TAPE_FILE 

vhs record >> $TAPE_FILE

rm todo.md
cp template.md todo.md
