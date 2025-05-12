# shit - **sh**are **it**

shit is a lightweight file server written in go designed for quickly sharing files and directories over http

it supports serving single or multiple files, directories, or a combination of both with  features like mime type detection, directory browsing, file uploading, tls, and customizable server settings

shit is a successor to my previous project [vmshare](https://github.com/notwithering/vmshare)

```bash
go install github.com/notwithering/shit@latest
```

## features

- smart
    - share single files, multiple files, and/or directories
    - auto-redirects for single items
    - built-in directory browsing with HTML interface
    - proper MIME type detection and handling
    - detects if using cURL and gives raw text output
- many options
    - custom host/port (`--host`, `--port`)
    - TLS support (`--tls`, `--cert`, `--key`)
    - file upload support (`--upload`)
    - index file serving (`--index`)
    - go's `http.FileServer` mode (`--go`)
    - fine tuned control (`--max-upload-memory`, `--read-timeout`, `--write-timeout`, `--upload-timeout`, `--permanent-redirect`)
- security & reliability
    - path validation & traversal protection
    - error logging with timestamps
    - environment variable support

## examples

<!-- code examples labeled as C++ for better bash formatting -->

```cpp
$ echo "hello" > a.txt
$ echo "hi" > b.txt
$ mkdir dir
$ echo "hey" > dir/c.txt
```

host current directory

```cpp
$ shit
$ curl 127.0.0.1:8080
a.txt
b.txt
dir/
$ curl 127.0.0.1:8080/a.txt
hello
```

host a single file

```cpp
$ shit b.txt
$ curl -L 127.0.0.1:8080
hi
```

host a few files

```cpp
$ shit a.txt b.txt
$ curl 127.0.0.1:8080
a.txt
b.txt
```

host a file and a directory

```cpp
$ shit a.txt dir/
$ curl 127.0.0.1:8080
a.txt
dir/
$ curl 127.0.0.1:8080/dir
c.txt
$ curl 127.0.0.1:8080/dir/c.txt
hey
```

etc.

## licenses

this project uses the following dependencies with the license as noted:

- [github.com/alecthomas/kingpin](https://github.com/alecthomas/kingpin) - MIT License

each dependency retains its respective license. for more details refer to their official documentation or source code
