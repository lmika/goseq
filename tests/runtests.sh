#!/bin/bash
#
#   Run the test producing a webpage that can be viewed by the
#   user.
#


RESULT_SUFFIX="out"
TEST_BIN="./goseq.test.${RESULT_SUFFIX}"

function die()
{
    echo "$@" >&2
    exit 1
}

function runTest()
{
    local inFile="$1"
    local outFile="$2"

    $TEST_BIN < $inFile > $outFile
    echo "$inFile $outFile"
}

function runTests()
{
    for name in *.seq; do
        local outFile="$name.${RESULT_SUFFIX}"

        runTest "$name" "$outFile"
    done
}

function buildResultsPage()
{
    cat << _EOF_
<html>
<head>
  <style>
    table { border: solid thin black; border-collapse: collapse; }
    td { border: solid; }
  </style>
</head>
<body>
_EOF_

    while read inFile outFile; do
        echo "<p>${inFile}</p>"
        echo "<table><tr><td><pre>"
        cat $inFile
        echo "</pre></td><td>"
        cat $outFile | grep -v '^<.xml'
        echo "</td></tr></table>"
    done

    cat << _EOF_
</body>
</html>
_EOF_
}

function cleanUp()
{
    rm *.${RESULT_SUFFIX}
}

go build -o $TEST_BIN ../. || die "Failed to build goseq"

runTests | buildResultsPage > res.html
cleanUp
