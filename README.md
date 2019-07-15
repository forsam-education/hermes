# Hermes

[![CircleCI](https://circleci.com/gh/forsam-education/hermes/tree/master.svg?style=svg)](https://circleci.com/gh/forsam-education/hermes/tree/master)
[![GoDoc](https://godoc.org/github.com/forsam-education/hermes?status.svg)](https://godoc.org/github.com/forsam-education/hermes)
[![Go Report Card](https://goreportcard.com/badge/github.com/forsam-education/hermes)](https://goreportcard.com/report/github.com/forsam-education/hermes)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A simple project that handles sending e-mails through SMTP from a template storage system, using SQS messages, meant to be run on an AWS Lambda.

## Quality assurance

To fix the basics of code format, you can run `go fmt`.

For a bit more advanced code style checks, you can run `golint $(go list ./... | grep -v /vendor/)`. You'll have to run `go get -u golang.org/x/lint/golint` before.

## Storage Connectors

Storage connectors just have to implement the StorageConnector interface. 

At the moment, only the S3 Bucket connector is available but feel free to implement any other storage connector and make a pull request.

## Templates naming

You need to have both HTML and plain text versions of a template, and store them using `templatename.html.template` and `templatename.txt.template` naming system.

You then only have to pass the template name in the SQS message, and it will get both versions.

## Templates format

The templates are in the basic [Go HTML Template](https://golang.org/pkg/html/template/) and [Go TEXT Template](https://golang.org/pkg/text/template/) formats, and therefor you must use the `{{.myVar}}` notation, the var_name being the key of your data in the `template_context` json object.

## Environment Variables

You have to configure the SMTP server connection details and the S3 template bucket using environment variables.

You can customise their names in the `config` structure, in `main.go`, specifically if you implement a new storage connector.

## Call process

When deployed, this lambda has to subscribe to an SQS queue that will transport the messages containing the informations about the mails to send.

Here is an example of message body to send:

```json
{
  "from_address": "test-lambda@forsam.education",
  "from_name": "Example Name",
  "reply_to": "reply-to@forsam.education",
  "to_address": "cto@forsam.education",
  "subject": "This is my subject",
  "template_name": "template-example",
  "template_context": {
    "myVar": "value"
  },
  "bcc": ["sneaky@yourmanager.com"],
  "cc": ["not-so-sneaky@example.com"]
}
```
