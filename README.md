# Hermes

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

The templates are in the basic (Go HTML Template)[https://golang.org/pkg/html/template/] and (Go TEXT Template)[https://golang.org/pkg/text/template/] formats, and therefor you must use the `{{.myVar}}` notation, the var_name being the key of your data in the `template_context` json object.

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
