package main

import "bazil.org/fuse/fs"

type FS struct {
	root *Dir
}

type Node struct {
	inode uint64
	name  string
	path string
}

func (f *FS) Root() (fs.Node, error) {
	return f.root, nil
}
