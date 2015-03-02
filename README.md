# goseq - text based sequence diagrams

A small command line utility used to generate UML sequence diagrams from a text-base definition file.

Inspired by [js-sequence-diagram](http://bramp.github.io/js-sequence-diagrams/) and
[websequencediagram](http://www.websequencediagrams.com/).

## Install

To install it:

    go get bitbucket.org/lmika/goseq

To allow automatic generation of PNG files, ImageMagick is required.  To disable this feature, add the `noim` tag:

    go get -tags noim bitbucket.org/lmika/goseq

## Usage

    goseq [FLAGS] FILES ...

Supported flags:

* `-o filename`: Specify output filename (either .svg or, if supported, .png)