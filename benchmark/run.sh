#!/bin/bash

module load golang/1.19

cd "$(dirname "$0")"

for size in xsmall small medium large xlarge
do
    mkdir -p ./$size
    touch ./$size/sequential.txt
    for i in 2 4 6 8 12
    do
        touch ./$size/$i-thread.txt
    done
    
    echo "Running $size test"
    for i in {1..5}
    do
        go run ../benchmark.go s $size >> ./$size/sequential.txt
    done

    for i in 2 4 6 8 12
    do
        for j in {1..5}
        do
            go run ../benchmark.go p $size $i >> ./$size/$i-thread.txt
        done
    done
    echo ""
done

python3 ./analysis.py
