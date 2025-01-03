# shit - **sh**are **it**

simple and effective fileserver written in go to share files quickly

shit when i first created it was just a go script that ran an http.FileServer on the current directory but i kept adding more features as i needed them over time (hosting a single file, hosting multiple files, hosting both files and directories, etc.) but now i think its ready for other people to use

shit is a successor to my previous project [vmshare](https://github.com/notwithering/vmshare)

```
go install github.com/notwithering/shit@latest
```

## features

- file and directory sharing
	+ allows serving multiple files or directories specified via command-line arguments
	+ defaults to serving the current working directory if no arguments are provided
- port specification
	+ customizable server port using `--port` or `-p` flag
	+ defaults to port 8080 if not specified
- dynamic file and directory browsing
	+ automatically lists directory contents as clickable links in an html interface
	+ redirects to a file or directory if the root directory contains a single export target
- file handling
	+ serves files with their correct MIME types
	+ reads and streams file content directly 
- error handling
	+ provides error logging in the console including timestamps for debugging
	+ returns appropriate http error response to the client for issues like missing files or directories
- path validation and security
	+ ensured safe file and directory resolution to prevent directory traversal attacks
	+ optional flags for enabling tls
- ease of use
	+ supports simple command line usage with the [`github.com/alecthomas/kingpin`](https://github.com/alecthomas/kingpin) library for argument parsing
	+ automatically displays the local server address (`http://127.0.0.1:<port>`)

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