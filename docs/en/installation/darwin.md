---
title: Mac OS X
subtitle: SCM-Manager CLI Client installation on OS X using homebrew
displayToc: true
---

# Homebrew
To install the SCM-Manager CLI Client on OS X we offer a [Homebrew](https://brew.sh/) tap.
This CLI Client does not work standalone but need a running SCM-Manager server.

## Quickstart

```bash
brew install scm-manager/tap/scm-cli
```

## Detailed installation

To install SCM-Manager with homebrew we had to add the SCM-Manager tap:
```bash
brew tap scm-manager/tap
```
After the tap was added, we can install the SCM-Manager CLI Client:
```bash
brew install scm-cli
```

Now the CLI Client can be started:

```bash
scm login {server_url}
```

