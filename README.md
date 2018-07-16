# upmail

[![Travis CI](https://img.shields.io/travis/genuinetools/upmail.svg?style=for-the-badge)](https://travis-ci.org/genuinetools/upmail)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge)](https://godoc.org/github.com/genuinetools/upmail)
[![Github All Releases](https://img.shields.io/github/downloads/genuinetools/upmail/total.svg?style=for-the-badge)](https://github.com/genuinetools/upmail/releases)

Provides email notifications for [sourcegraph/checkup](https://github.com/sourcegraph/checkup).

 * [Installation](README.md#installation)
      * [Binaries](README.md#binaries)
      * [Via Go](README.md#via-go)
 * [Usage](README.md#usage)

## Installation

#### Binaries

For installation instructions from binaries please visit the [Releases Page](https://github.com/genuinetools/upmail/releases).

#### Via Go

```console
$ go get github.com/genuinetools/upmail
```

## Usage

```console
$ upmail -h
upmail -  Email notification hook for https://github.com/sourcegraph/checkup.

Usage: upmail <command>

Flags:

  --mailgun-domain  Mailgun Domain to use for sending email (optional) (default: <none>)
  --password        SMTP server password (default: <none>)
  --appengine       enable the server for running in Google App Engine (default: false)
  --config          config file location (default: checkup.json)
  -d                enable debug logging (default: false)
  --server          SMTP server for email notifications (default: <none>)
  --username        SMTP server username (default: <none>)
  --interval        check interval (ex. 5ms, 10s, 1m, 3h) (default: 10m0s)
  --mailgun         Mailgun API Key to use for sending email (optional) (default: <none>)
  --recipient       recipient for email notifications (default: <none>)

Commands:

  version  Show the version information.
```
