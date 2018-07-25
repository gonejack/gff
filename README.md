# gff
Go file flatting - extract files from directories.

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

