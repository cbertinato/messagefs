# messagefs
A messaging file system written in Go.

*messagefs* is a userspace file system that enables communications between nodes
with file operations. In this file system, files represent nodes (or users) identified by the last 6 characters of a randomly assigned UUID,
and directories represent groups.

Messages can be sent to a single node by writing to
the file represent that nodes.

Messages can be sent to a group by writing to the
hidden *.all* file that exists in each directory.

This is currently a proof-of-concept, but is being actively developed to implement
all functionality of a typical file system as well as that expected of
a secure communication protocol.

The implementation relies on two packages:
* bazil.org/fuse: A native Go implementation of the kernel-userspace
communication protocol. It does not depend upon the FUSE C library.
* zeromq/gyre: A Golang port of Zyre 2.0, which is an open-source framework for
proximity-based peer-to-peer applications.

Here's how to get going:
* Install the dependencies:
```shell
go get github.com/bazil.org/fuse
go get github.com/zeromq/gyre
```
* Open a terminal and create an empty directory for a mountpoint.
* Start the filesystem:
```shell
go run comms.go dir.go file.go fs.go main.go <mountpoint>
```
* Open another terminal and run the test chat:
```shell
go run chat.go
```

## License

MIT License (MIT)

Copyright (c) 2017 Chris Bertinato

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
