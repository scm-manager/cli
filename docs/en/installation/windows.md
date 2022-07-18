---
title: Windows
subtitle: Install SCM-Manager CLI Client on Windows
displayToc: true
---

The following document describes the installation process for SCM-Manager CLI Client on Windows.

## Install SCM-Manager CLI

### Install via Scoop
To install the CLI CLient via Scoop run the following commands:
```
scoop bucket add scm https://github.com/scm-manager/scoop-bucket
scoop install scm-cli
```

### Manual installation

To install SCM-Manager CLI you have to download the latest Windows package from the [download page](/cli/).
After unpacking the archive move the file to a new directory. 
To make it available on your `PATH` you can follow this [instruction](https://stackoverflow.com/questions/1618280/where-can-i-set-path-to-make-exe-on-windows).


## First start
Now we have to open a Terminal (PowerShell, Bash or CMD), in order to run the SCM-Manager CLI Client.
For this to work you must have an SCM-Manager server running and connect your client first. 
You can start with the following command:

```bash
scm.exe login {server_url}
```

### Known problems with mintty (Git Bash, Cygwin)
On some terminals under Windows, problems with stdin or stdout/stderr may occur. For example, when you try to log in,
hitting enter on your username generates a second enter which prevents the possibility to enter your password and the
login fails.

Or, if you have previously logged in e.g. via Powershell, running the CLI won't create any output and the
CLI won't exit.

You can fix these problems by adding the following line to your `.minttyrc` file:

```properties
ConPTY=on
```
