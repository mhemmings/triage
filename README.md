# Triage

Provide a list of repositories, and get back all the issues that haven't been triaged. Currently only supports Github, but easily extensible to add providers.

## Installation

```console
$ go get github.com/mhemmings/triage
```

## Usage

```console
$ triage repolist
```

Where repolist is a file containing a list of repositories, for example:

```
owner/reponame
github.com/owner/reponame
https://github.com/owner/reponame
```

Once all the issues are collected, a simple webpage will be served displaying them all.

An individual repository can also be passed without the need to use a separate repo list file:

```console
$ triage -r "owner/reponame"
```

By default, only unlabeled issues are shown. A comma-separated list of label names can be provided with `-l/--labels` to override this and match only those label names.


```console
$ triage --labels bug,todo,critical repolist
```

Or all issues can be gathered, regardless of label status using `--all`.

To filter by time since issue was created, use `-s`. For example:
```
$ triage repolist -s 6h  # Last 6 hours
$ triage repolist -s 2d  # Last 2 days
$ triage repolist -s 1w  # Last week
```

## Github Auth

The Github API has strict rate limits for unauthenticated requests. It is recommended you [generate a token](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/) and set this to the `TRIAGE_GITHUB_TOKEN` environment variable to authenticate your requests. This will also give you access to private repositories.

## TODO

- Better tests
- More outputs (static, JSON etc)
- More providers (Bitbucket etc)
- Ability to triage without visiting Github
