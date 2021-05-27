# changelogger
A changelog file generator heavily inspired on [towncrier](https://github.com/twisted/towncrier).

## Why?
I'm looking for an alternative to towncrier that doesn't require a Python environment,
so I decided to build one with only basic features.

## How it works?
```
$ changelogger -h
Usage of changelogger:
  -dir string
        Directory of changelog entries
  -path string
        Path of the changelog file
  -version string
        The release version

$ changelogger -dir ./.changelog -path ./CHANGELOG.md -version v1.1.2
```
changelogger will search for files on the directory of changelog entries (flag `-dir`) that contains the follow extensions:

* `.added`
* `.changed`
* `.removed`
* `.fixed`

Each file represents an entry on changelog and the extension defines which session the entry will be added to.

After collecting all entries, changelogger remove the entry file and edit the changelog file (flag `-path`) replacing the placeholder
`[changelogger-notes]::` for the content of the entries with a version header (flag `-version`). F.ex:

```diff
# Changelog

- [changelogger-notes]::
+ [changelogger-notes]::
+
+ ## v0.0.1 (2021-05-22)
+ ### Added
+ * A new cool feature
+
+ ## Removed
+ * The deprecated feature
```
