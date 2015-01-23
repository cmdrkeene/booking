# Booking

![Build Status](https://travis-ci.org/cmdrkeene/booking.svg)

An apartment booking service.

The goal is to explore new techniques for building rich applications in
[Golang](http://golang.org).

## Dependencies

These should be kept to an absolute minimum and the standard library should be
used as much as possible.

	glog
	inject

## Setup

    $ go version
    go version go1.4.1 darwin/amd64
    $ go get ./...

## Testing

Use the standard `go test ./...` tool, or even better:

    $ ./autotest.sh

## TODO

* Finish one-step booking form
* Web interface to cancel
* Send email on book and cancel to admin and guest

## License

THIS IS NOT OPEN SOURCE YET!

Copyright 2015 Brandon Keene - All Rights Reserved
