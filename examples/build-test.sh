#!/bin/bash -eu

target_directory=examples
for d in $target_directory/*/ ; do
    (
        cd $d
        go build -o main
    )
done
