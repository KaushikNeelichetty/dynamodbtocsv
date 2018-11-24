# dynamodbtocsv

[![Build Status](https://travis-ci.org/KaushikNeelichetty/dynamodbtocsv.svg?branch=master)](https://travis-ci.org/KaushikNeelichetty/dynamodbtocsv)
[![standard-readme compliant](https://img.shields.io/badge/standard--readme-OK-green.svg?style=flat-square)](https://github.com/RichardLitt/standard-readme)
[![Go Report Card](https://goreportcard.com/badge/github.com/KaushikNeelichetty/dynamodbtocsv)](https://goreportcard.com/report/github.com/KaushikNeelichetty/dynamodbtocsv)
[![Downloads](https://img.shields.io/github/downloads/KaushikNeelichetty/dynamodbtocsv/latest/total.svg)](https://github.com/KaushikNeelichetty/dynamodbtocsv/releases)
[![Latest release](https://img.shields.io/github/release/KaushikNeelichetty/dynamodbtocsv.svg)](https://github.com/KaushikNeelichetty/dynamodbtocsv/releases)

Dump your DynamoDB table as CSV

## Table of Contents

- [Install](#install)
- [Example](#example)
- [Usage](#usage)
- [Maintainers](#maintainers)
- [Contribute](#contribute)
- [License](#license)

## Install

#### Homebrew

```
brew tap KaushikNeelichetty/homebrew-tap
brew install dynamodbtocsv
```

#### Direct Binary Download

```
# binary will be $GOPATH/bin/dynamodbtocsv
curl -sfL https://raw.githubusercontent.com/KaushikNeelichetty/dynamodbtocsv/master/install.sh | sh -s -- -b $GOPATH/bin

# or install it into ./bin/
curl -sfL https://raw.githubusercontent.com/KaushikNeelichetty/dynamodbtocsv/master/install.sh | sh -s

# In alpine linux (as it does not come with curl by default)
wget -O - -q https://raw.githubusercontent.com/KaushikNeelichetty/dynamodbtocsv/master/install.sh | sh -s
```

## Example 

```bash
$ AWS_REGION=ap-southeast-1 dynamodbtocsv -t timesheet-cop -f timestamp,count > dump.csv
$ head dump.csv
timestamp,count
2018-09-05 09:41:20,3
2018-10-05 06:41:20,1
2018-09-10 00:41:24,95
2018-07-31 05:41:21,5
2018-07-28 14:41:14
2018-08-01 12:41:20,1
2018-09-19 22:11:21,9
2018-10-12 04:11:21,4
2018-09-18 20:41:21,11
```
To load the shared config, `export AWS_SDK_LOAD_CONFIG="true"` before running the above command. This is useful to load the shared credentials file from `~/.aws/credentials`.

## Usage

```bash
dynamodbtocsv -h
Usage:
  dynamodbtocsv [OPTIONS]

Application Options:
  -t, --table-name=       Name of the dynamo db table
  -f, --fields-name=      List of comma separated fieldName to be output to CSV [Example "timestamp,count"] [nested structures will not be flattened]

Help Options:
  -h, --help              Show this help message
  
```

## Maintainers

- [@KaushikNeelichetty](https://github.com/KaushikNeelichetty)
- [@kishaningithub](https://github.com/kishaningithub)

## Contribute

PRs accepted.

Small note: If editing the README, please conform to the [standard-readme](https://github.com/RichardLitt/standard-readme) specification.

## License

MIT © 2018 Kaushik Neelichetty
