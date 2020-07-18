#!/bin/bash -eu

target_directory=_examples
for d in $target_directory/*/ ; do
    (
        cd $d
        go build -o main
    )
done
