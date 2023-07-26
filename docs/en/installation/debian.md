---
title: Debian/Ubuntu
subtitle: Installation of SCM-Manager CLI Client for Debian-based linux distributions
displayToc: true
---

## Quickstart

The following code block will configure an apt repository for SCM-Manager CLI Client and install it.

```bash
echo "X-Repolib-Name: SCM-Manager CLI
Enabled: yes
Types: deb
URIs: https://packages.scm-manager.org/repository/apt-v2-releases/
Suites: stable
Components: main
Architectures: $(dpkg --print-architecture)
Signed-By: /etc/apt/keyrings/scmm-archive-keyring.gpg" | sudo tee /etc/apt/sources.list.d/scm-manager.sources
sudo mkdir -p /etc/apt/keyrings
sudo chmod 0755 /etc/apt/keyrings
sudo sh -c 'curl -s https://keys.openpgp.org/vks/v1/by-fingerprint/23D2625B235E25A4719875A2975922F193B07D6E | gpg --dearmor > /etc/apt/keyrings/scmm-archive-keyring.gpg'
sudo chmod +r /etc/apt/keyrings/scmm-archive-keyring.gpg
sudo apt-get update
sudo apt-get install scm-cli
```

## Detailed installation

To install SCM-Manager CLI as a debian package (.deb), we have to configure a [third-party APT repository](https://wiki.debian.org/DebianRepository/UseThirdParty).
Replace `<arch>` with your system architecture. You can find it out by using `dpkg --print-architecture`.
Create a file at `/etc/apt/sources.list.d/scm-manager.sources` with the following content:

```bash
echo "X-Repolib-Name: SCM-Manager CLI
Enabled: yes
Types: deb
URIs: https://packages.scm-manager.org/repository/apt-v2-releases/
Suites: stable
Components: main
Architectures: <arch>
Signed-By: /etc/apt/keyrings/scmm-archive-keyring.gpg" | sudo tee /etc/apt/sources.list.d/scm-manager.sources
```

We are using the deb822 style format supported by apt itself since version 1.1. Previous versions ignore such files with a notice message. In this case, the one-line-style format can be used:

```bash
echo 'deb [arch=<arch> signed-by=/etc/apt/keyrings/scmm-archive-keyring.gpg] https://packages.scm-manager.org/repository/apt-v2-releases/ stable main' | sudo tee /etc/apt/sources.list.d/scm-manager.list
```

This will add the apt repository of the SCM-Manager stable releases to the list of your apt repositories.
To ensure the integrity of the debian packages, we have to import the gpg key for the repository.
In addition, the keyrings directory should be created beforehand, since it does not always exist by default:

```bash
sudo mkdir -p /etc/apt/keyrings
sudo chmod 0755 /etc/apt/keyrings
sudo sh -c 'curl -s https://keys.openpgp.org/vks/v1/by-fingerprint/23D2625B235E25A4719875A2975922F193B07D6E | gpg --dearmor > /etc/apt/keyrings/scmm-archive-keyring.gpg'
sudo chmod +r /etc/apt/keyrings/scmm-archive-keyring.gpg
```

After we have imported the gpg key, we can update the package index and install the SCM-Manager CLI Client:

```bash
sudo apt-get update
sudo apt-get install scm-cli
```
