# nhexport

[![Build Status](https://api.travis-ci.org/mrd0ll4r/nhexport.svg?branch=master)](https://travis-ci.org/mrd0ll4r/nhexport)
[![Go Report Card](https://goreportcard.com/badge/github.com/mrd0ll4r/nhexport)](https://goreportcard.com/report/github.com/mrd0ll4r/nhexport)
[![GoDoc](https://godoc.org/github.com/mrd0ll4r/nhexport?status.svg)](https://godoc.org/github.com/mrd0ll4r/nhexport)
![Lines of Code](https://tokei.rs/b1/github/mrd0ll4r/nhexport)

A quick-and-dirty nicehash-to-CSV exporter.

It exports either the payment history (set `-payment`), or the history of hashrates and unpaid balances per algorithm (default).

This will print out CSV (including headers) to a file named `FROM-TO-ADDR-MODE.csv`.

## Installation

Make sure you have a working Go installation.
Then just

```bash
go get -u github.com/mrd0ll4r/nhexport/cmd/nhexport
```

Will fetch, compile, and install `nhexport` and all its dependencies.
If you added `$GOPATH/bin` to your `$PATH`, the `nhexport` binary should be on your path.

## Usage

```
  -addr string
        your bitcoin address (mandatory)
  -from string
        begin date (inclusive) (default YESTERDAY)
  -payments
        whether to export hashrates+history (default) or payments
  -to string
        end date (exclusive) (default TODAY)
```

If ran without setting `from` and `to`, it will export the stats of yesterday.

Note that there is an API call rate limit of ~one request/30 seconds.

## Exit status

Exit status will be zero in case of success, one in case of runtime errors and two in case of malformed arguments.
