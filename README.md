# gff
Go file flatting - extract files from sub-directories.

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gonejack/gff)
![Build](https://github.com/gonejack/gff/actions/workflows/go.yml/badge.svg)
[![GitHub license](https://img.shields.io/github/license/gonejack/gff.svg?color=blue)](LICENSE)

## Usage
```
before
.
├── a
│   └── a.flv
├── b
│   └── b.flv
└── c.flv
```

```
> go get github.com/gonejack/gff

> gff -yes '*.flv'
```

```
after
.
├── a
├── a_a.flv
├── b
├── b_b.flv
└── c.flv
```

## Help
```
Usage: gff [<patterns> ...]

Extract files from nested directory https://github.com/gonejack/gff

Arguments:
  [<patterns> ...]

Flags:
  -h, --help             Show context-sensitive help.
  -s, --separator="_"    Filename separator.
  -y, --yes              Make real changes.
```
