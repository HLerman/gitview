# gitview

Gitview is a tool which scan directories to find all of your git repositories and check if they are up-to-date.

## Installation

Use the [release](https://github.com/HLerman/gitview/releases) page to download the lastest version

Or you can get and build the tool yourself

```bash
git clone https://github.com/HLerman/gitview.git
cd gitview
go build
```

## Usage

Display all of your git repositories and check if it's up to date.
```bash
$ ./gitview
/home/user/go/src/github.com/karrick/godirwalk/      GIT[master] outdated
/home/user/go/src/github.com/sqs/goreturns/          GIT[master] up-to-date
/home/user/go/src/github.com/uudashr/gopkgs/         GIT[master] up-to-date
/home/user/go/src/github.com/alexflint/go-arg/       GIT[master] up-to-date
/home/user/go/src/github.com/fatih/color/            GIT[master] up-to-date
/home/user/go/src/github.com/go-delve/delve/         GIT[master] outdated
/home/user/go/src/github.com/stamblerre/gocode/      GIT[master] up-to-date
/home/user/go/src/golang.org/x/sys/                  GIT[master] outdated
/home/user/go/src/github.com/acroca/go-symbols/      GIT[master] up-to-date
/home/user/go/src/github.com/google/uuid/            GIT[master] up-to-date
/home/user/go/src/github.com/mdempsky/gocode/        GIT[master] up-to-date
/home/user/go/src/github.com/pkg/errors/             GIT[master] up-to-date
/home/user/go/src/github.com/ramya-rao-a/go-outline/ GIT[master] up-to-date
/home/user/go/src/github.com/sirupsen/logrus/        GIT[master] outdated
/home/user/go/src/github.com/rogpeppe/godef/         GIT[master] up-to-date
/home/user/go/src/gitview/                           GIT[master] up-to-date
/home/user/go/src/golang.org/x/lint/                 GIT[master] up-to-date
/home/user/go/src/golang.org/x/tools/                GIT[master] outdated
```

Create a json file (gitview.json in the same directory than gitview binary) to sav repositories paths. Improve next standard execution speed (path are already known).
```bash
$ ./gitview --refresh
```

Equivalent to a git pull in all of your git repositories which use the master branch.
```bash
$ ./gitview --pull
/home/user/go/src/github.com/go-delve/delve/         GIT[master] up-to-date
/home/user/go/src/github.com/alexflint/go-arg/       GIT[master] up-to-date
/home/user/go/src/github.com/pkg/errors/             GIT[master] up-to-date
/home/user/go/src/golang.org/x/tools/                GIT[master] up-to-date
/home/user/go/src/github.com/fatih/color/            GIT[master] up-to-date
/home/user/go/src/github.com/karrick/godirwalk/      GIT[master] up-to-date
/home/user/go/src/github.com/mdempsky/gocode/        GIT[master] up-to-date
/home/user/go/src/github.com/uudashr/gopkgs/         GIT[master] up-to-date
/home/user/go/src/golang.org/x/sys/                  GIT[master] up-to-date
/home/user/go/src/gitview/                           GIT[master] up-to-date
/home/user/go/src/github.com/rogpeppe/godef/         GIT[master] up-to-date
/home/user/go/src/github.com/google/uuid/            GIT[master] up-to-date
/home/user/go/src/github.com/ramya-rao-a/go-outline/ GIT[master] up-to-date
/home/user/go/src/golang.org/x/lint/                 GIT[master] up-to-date
/home/user/go/src/github.com/acroca/go-symbols/      GIT[master] up-to-date
/home/user/go/src/github.com/stamblerre/gocode/      GIT[master] up-to-date
/home/user/go/src/github.com/sqs/goreturns/          GIT[master] up-to-date
/home/user/go/src/github.com/sirupsen/logrus/        GIT[master] up-to-date
```