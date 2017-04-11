package main

import (
	"log"
	"os"
	"strings"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"bazil.org/fuse/fuseutil"
	"golang.org/x/net/context"
)

type File struct {
	Node
	data []byte
	channel chan string
}

var newData string

// TO DO: Recycle inode numbers.
func NewInode() uint64 {
	inode += 1
	return inode
}

func (f *File) Attr(ctx context.Context, a *fuse.Attr) error {
	log.Println("Requested Attr for File", f.name, "has data size", len(f.data))
	a.Inode = f.inode
	a.Mode = 0777
	a.Size = uint64(len(f.data))
	return nil
}

func (f *File) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	log.Println("Requested Read on File", f.name)
	fuseutil.HandleRead(req, resp, f.data)
	return nil
}

func (f *File) ReadAll(ctx context.Context) ([]byte, error) {
	log.Println("Reading all of file", f.name)
	return []byte(f.data), nil
}

func (f *File) Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error {
	log.Println("Trying to write to ", f.name, "offset", req.Offset, "dataSize:", len(req.Data), "data: ", string(req.Data))

	if f.name == ".all" {
		newData = string(req.Data)
		// get current directory
		pwd := (*f).Node.path

		// get list of files and subdirectories
		p, err := os.Open(pwd)
		if err != nil {
			return err
		}
		defer p.Close()

		files, err := p.Readdirnames(-1)
		if err != nil {
			return err
		}

		// write to all files and directories in the directory
		for _, fname := range files {
			if fname == ".all" {
				continue
			} else {
				fp, _ := os.OpenFile(pwd + fname, os.O_APPEND|os.O_RDWR, 0777)
				finfo, _ := fp.Stat()
				fn := finfo.Name()

				// if directory, write to it's .all file
				if finfo.IsDir() {
					fp.Close()
					fp, _ = os.OpenFile(pwd + fn + "/.all", os.O_RDWR, 0777)
				}

				_, err = fp.Write(req.Data)
				if err != nil {
					log.Panicln(err)
					return err
				}
				fp.Close()
			} // if
		} // for
	} else {
		if f.name != Myname {
			msg := Msg{
				Type: "DIRECT_MSG",
				// Payload: string(req.Data),
				Payload: newData,
				Target: f.name,
			}
			input <- msg
		}

		f.data = req.Data
		log.Println("Wrote to file", f.name)
	} // if

	resp.Size = len(req.Data)
	return nil
}

func (f *File) Flush(ctx context.Context, req *fuse.FlushRequest) error {
	log.Println("Flushing file", f.name)
	return nil
}

func (f *File) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {
	log.Println("Open call on file", f.name)
	return f, nil
}

func (f *File) Release(ctx context.Context, req *fuse.ReleaseRequest) error {
	log.Println("Release requested on file", f.name)
	return nil
}

func (f *File) Fsync(ctx context.Context, req *fuse.FsyncRequest) error {
	log.Println("Fsync call on file", f.name)
	return nil
}

func (f *File) Rename(ctx context.Context, req *fuse.RenameRequest, newDir Node) error {
	// not yet implemented
}
