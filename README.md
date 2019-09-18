# Pinbak 🐧

Pinbak is a small backup manager. It will save and restore file with git with any git repository manager.

## Getting started


```
pinbak init
pinbak add repository configuration git@github.com:Pinbak/configuration-backup.git
pinbak add repository ssh git@github.com:Pinbak/ssh-backup.git
pinbak add configuration ~/.bashrc ~/.gitconfig ~/.config/Code/User/settings.json
pinbak add ssh ~/.ssh/id_rsa.pub
pinbak list
```

To restore the file from a fresh install
```
pinbak init
pinbak add repository configuration git@github.com:Pinbak/configuration-backup.git
pinbak restore all
```


## Usage

```raw
Pinbak is a simple backup manager that store files in a git repository.

Usage:
  pinbak [flags]
  pinbak [command]

Available Commands:
  add         Add a file to backup.
  help        Help about any command
  init        Init pinbak.
  list        List all file in repository.
  restore     Restore all file in repository.

Flags:
  -h, --help   help for pinbak

Use "pinbak [command] --help" for more information about a command.
```


