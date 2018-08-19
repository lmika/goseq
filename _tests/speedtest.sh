#!/bin/bash
#
#   Speed Test
#

P=500
N=10000

TEST_BIN="./goseq.test.out"

function die()
{
    echo "$@" >&2
    exit 1
}

function genTest()
{
    for n in `seq 1 $N`; do        
        fr="p$((RANDOM % $P))"
        tr="p$((RANDOM % $P))"
        while [ "$fr" = "$tr" ]; do
            fr="p$((RANDOM % $P))"
            tr="p$((RANDOM % $P))"
        done

        echo "$fr->$tr: Line $n"
    done
}

go build -o $TEST_BIN ../. || die "Failed to build goseq"
trap "rm $TEST_BIN" 0

genTest | time $TEST_BIN > /dev/null