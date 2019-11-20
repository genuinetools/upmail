# upmail

![make-all](https://github.com/genuinetools/upmail/workflows/make%20all/badge.svg)
![make-image](https://github.com/genuinetools/upmail/workflows/make%20image/badge.svg)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge)](https://godoc.org/github.com/genuinetools/upmail)
[![Github All Releases](https://img.shields.io/github/downloads/genuinetools/upmail/total.svg?style=for-the-badge)](https://github.com/genuinetools/upmail/releases)

Provides email notifications for [sourcegraph/checkup](https://github.com/sourcegraph/checkup).

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [Installation](#installation)
    - [Binaries](#binaries)
    - [Via Go](#via-go)
- [Usage](#usage)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

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

  --appengine       enable the server for running in Google App Engine (default: false)
  --password        SMTP server password (default: <none>)
  --recipient       recipient for email notifications (default: <none>)
  --mailgun-domain  Mailgun Domain to use for sending email (optional) (default: <none>)
  --sender          SMTP default sender email address for email notifications (default: <none>)
  --server          SMTP server for email notifications (default: <none>)
  --username        SMTP server username (default: <none>)
  --config          config file location (default: checkup.json)
  -d                enable debug logging (default: false)
  --interval        check interval (ex. 5ms, 10s, 1m, 3h) (default: 10m0s)
  --mailgun         Mailgun API Key to use for sending email (optional) (default: <none>)

Commands:

  version  Show the version information.
```
