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

### AMD64 architecture
```bash
wget https://packages.scm-manager.org/repository/scm-cli-releases/<version>/scm-cli_<version>_Linux_x86_64.tar.gz
echo "<checksum> *scm-cli_<version>_Linux_x86_64.tar.gz" | sha1sum -c -
```
Extract the archive:
```bash
tar xvfz scm-cli_<version>_Linux_x86_64.tar.gz -C /opt
```

### ARM64 architecture
```bash
wget https://packages.scm-manager.org/repository/scm-cli-releases/<version>/scm-cli_<version>_Linux_arm64.tar.gz
echo "<checksum> *scm-cli_<version>_Linux_arm64.tar.gz" | sha1sum -c -
```
Extract the archive:
```bash
tar xvfz scm-cli_<version>_Linux_arm64.tar.gz -C /opt
```

## First start
The cli client can be started by using `scm` in your terminal.
```bash
scm
```
