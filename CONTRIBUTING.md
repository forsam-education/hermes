# Contribution Guidelines

## Pull request

Use the pull request template available when you start your PR, or in [PULL_REQUEST_TEMPLATE.md](PULL_REQUEST_TEMPLATE.md).

## Commit Convention

Commit message summaries must follow this basic format:

```none
Tag: Message
```

`Tag` should not be confused with git tag.
`Message` should not be confused with git commit message.

The `Tag` is one of the following:

- `Fix` - for a bug fix.
- `Update` - for a backwards-compatible enhancement of an existing feature.
- `Breaking` - for a backwards-incompatible enhancement.
- `Docs` - changes to documentation only.
- `Build` - changes to build / CI / CD process only.
- `New` - implemented a new feature.
- `Upgrade` - for a dependency upgrade.

The message summary should be a one-sentence description of the change.

Here are some good commit message summary examples:

```none
Build: Update Travis to only test Node 0.10
Fix: Semi rule incorrectly flagging extra semicolon
Upgrade: Esprima to 1.2, switch to using Esprima comment attachment
```
