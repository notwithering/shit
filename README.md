# shit - `sh`are `it`

shit is a lightweight file server written in go designed for quickly sharing files and directories over http

it supports serving single or multiple files, directories, or a combination of both with features like mime type detection, directory browsing, file uploading, tls, and customizable server settings

```bash
go install github.com/notwithering/shit@latest
```

```bash
yay -S shit-git
```

## example

```cpp
$ shit dir/ a.txt b.txt
$ curl 127.0.0.1:8080
dir/
a.txt
b.txt
$ curl 127.0.0.1:8080/dir
c.txt
```

## licenses

this project uses the following dependencies with the license as noted:

- [github.com/alecthomas/kong](https://github.com/alecthomas/kong) - MIT License

each dependency retains its respective license. for more details refer to their official documentation or source code
