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
                             _ _
 _   _ _ __  _ __ ___   __ _(_) |
| | | | '_ \| '_ ` _ \ / _` | | |
| |_| | |_) | | | | | | (_| | | |
 \__,_| .__/|_| |_| |_|\__,_|_|_|
      |_|

 Email notification hook for https://github.com/sourcegraph/checkup.
 Version: v0.4.2
 Build: 2905d94

  -appengine
        enable the server for running in Google App Engine
  -config string
        config file location (default "checkup.json")
  -d    run in debug mode
  -interval duration
        check interval (ex. 5ms, 10s, 1m, 3h) (default 10m0s)
  -mailgun string
        Mailgun API Key to use for sending email (optional)
  -mailgun-domain string
        Mailgun Domain to use for sending email (optional)
  -password string
        SMTP server password
  -recipient string
        recipient for email notifications
  -sender string
        SMTP default sender email address for email notifications
  -server string
        SMTP server for email notifications
  -username string
        SMTP server username
  -v    print version and exit (shorthand)
  -version
        print version and exit
```
