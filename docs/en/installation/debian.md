---
title: Debian/Ubuntu
subtitle: Installation of SCM-Manager CLI Client for Debian-based linux distributions
displayToc: true
---

## Quickstart

The following code block will configure an apt repository for SCM-Manager CLI Client and install it.

```bash
echo 'deb [arch=all] https://packages.scm-manager.org/repository/apt-v2-releases/ stable main' | sudo tee /etc/apt/sources.list.d/scm-manager.list
sudo apt-key adv --recv-keys --keyserver hkps://keys.openpgp.org 0x975922F193B07D6E
sudo apt-get update
sudo apt-get install scm-cli
```

## Detailed installation

To install SCM-Manager as a debian package (.deb), we have to configure an apt repository.
Create a file at `/etc/apt/sources.list.d/scm-manager.list` with the following content:

```text
deb [arch=all] https://packages.scm-manager.org/repository/apt-v2-releases/ stable main
```

This will add the apt repository of the SCM-Manager stable releases to the list of your apt repositories.
To ensure the integrity of the debian packages we have to import the gpg key for the repository.

```bash
sudo apt-key adv --recv-keys --keyserver hkps://keys.openpgp.org 0x975922F193B07D6E
```

After we have imported the gpg key, we can update the package index and install the SCM-Manager CLI Client:

```bash
sudo apt-get update
sudo apt-get install scm-cli
```
