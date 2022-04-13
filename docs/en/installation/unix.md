---
title: Unix
subtitle: General unix installation
displayToc: true
---

## Requirements
This tool does not work standalone but need a running SCM-Manager server.

## Installation
Grab the latest version and checksum from [download page](/cli) and replace `<version>` and `<checksum>` in the code blocks below.
Download and verify the checksum.

```bash
wget https://packages.scm-manager.org/repository/scm-cli-releases/<version>/scm-cli_<version>_<os>_<arch>.tar.gz
echo "<checksum> *scm-cli_<version>_<os>_<arch>.tar.gz" | sha256sum -c -
```
Extract the archive:
```bash
sudo tar xvfz scm-cli_<version>_<os>_<arch>.tar.gz -C /usr/local/bin
```

## First start
The cli client can be started by using `scm` in your terminal.
```bash
scm
```
