# Pinbak 🐧

Pinbak is a small backup manager. It will save and restore file or directory with git with and works with any git 
repository manager.

## Getting started

Default utilisation
```
$ pinbak init
$ pinbak add repo configuration git@github.com:Pinbak/configuration-backup.git
$ pinbak add repo ssh git@github.com:Pinbak/ssh-backup.git
$ pinbak add configuration ~/.bashrc ~/.gitconfig ~/.config/Code/User/settings.json
$ pinbak add ssh ~/.ssh/id_rsa.pub
$ pinbak list
```

To restore all the items from a fresh install
```
$ pinbak init
$ pinbak add repository configuration git@github.com:Pinbak/configuration-backup.git
$ pinbak restore all
```

To update all the backed items
```
$ pinbak update
```

To remove an item
```
$ pinbak list
configuration : git@github.com:Pinbak/configuration.git
    -  bm1nrku1nn09qvbipro0  :  {HOME}/.gitconfig
$ pinbak remove bm1nrku1nn09qvbipro0
Done.
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
  remove      Remove a file from the backup.
  restore     Restore all file in repository.
  update      Update all files in all repositories.

Flags:
  -h, --help   help for pinbak

Use "pinbak [command] --help" for more information about a command.
```

## Disclamer

* Pinbak will write all the directory with the UNIX decimal write `drwxr-xr-x` and file with `-rw-r--r--`.
* All directories and files will be writed with as the current user.
* When Pinbak detect an home path it will be restore in the home directory of the current user.
