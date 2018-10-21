# INI

The ini parser and operate


## TODO

- [X] read ini config from `reader` , return ini sturct
- [X] check config format
- [X] default section operate: set, get(default session sould not able remove)
- [X] section operate: set, get, del, update
    - [X] get section key's value
    - [X] set section key's value
    - [X] update section using map
- [X] default section key operate: set, get, del
    - should be same as section operate, except default section could not delete

## Install

```
go install github.com/kmollee/ini/....
```


## Usage

```
cfg -f {filename}  get {section_name} {key}
cfg -f {filename}  set {section name} {key} {value}
```

