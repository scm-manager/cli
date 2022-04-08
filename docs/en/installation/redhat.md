---
title: Redhat/CentOS/Fedora
subtitle: Installation of SCM-Manager CLI Client for RedHat-based linux distributions
displayToc: true
---

## Quickstart

The following code block will configure a yum repository for scm-manager and install it.

```bash
cat << EOF | sudo tee /etc/yum.repos.d/SCM-Manager.repo
[scm-manager]
name=SCM-Manager Repository
baseurl=https://packages.scm-manager.org/repository/yum-v2-releases/
enabled=1
gpgcheck=1
priority=1
gpgkey=file:///etc/pki/rpm-gpg/SCM-Manager
EOF
sudo curl -o /etc/pki/rpm-gpg/SCM-Manager https://packages.scm-manager.org/repository/keys/gpg/oss-cloudogu-com.pub
sudo yum install scm-cli
```

## Detailed installation

To install SCM-Manager as a redhat package (.rpm), we have to configure a yum repository.
Create a file at `/etc/yum.repos.d/SCM-Manager.repo` with the following content:

```ini
[scm-manager]
name=SCM-Manager Repository
baseurl=https://packages.scm-manager.org/repository/yum-v2-releases/
enabled=1
gpgcheck=1
priority=1
gpgkey=file:///etc/pki/rpm-gpg/SCM-Manager
```

This will add the yum repository of the scm-manager stable releases to the list of your yum repositories.
To ensure the integrity of the rpm packages we have to import the gpg key for the repository.

```bash
sudo curl -o /etc/pki/rpm-gpg/SCM-Manager https://packages.scm-manager.org/repository/keys/gpg/oss-cloudogu-com.pub
```

After we have imported the gpg key, we can install scm-manager:

```bash
sudo yum install scm-cli
```
