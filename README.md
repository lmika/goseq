# goseq - text based sequence diagrams

A small command line utility used to generate UML sequence diagrams from a text-base definition file.

Inspired by [js-sequence-diagram](http://bramp.github.io/js-sequence-diagrams/) and
[websequencediagram](http://www.websequencediagrams.com/).

## Install

To install it:

    go get github.org/lmika/goseq

To allow automatic generation of PNG files, build with the `im` tag (ImageMagick is required):

    go get -tags im bitbucket.org/lmika/goseq

## Usage

    goseq [FLAGS] FILES ...

Supported flags:

* `-o filename`: Specify output filename (either .svg or, if supported, .png)

## Licence

Released under the MIT Licence.