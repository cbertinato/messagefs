package main

import (
	"flag"
	"log"
	"os"
	"time"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

var filesys *FS
var inode uint64

var Usage = func() {
	log.Printf("Usage of %s:\n", os.Args[0])
	log.Printf("  %s MOUNTPOINT\n", os.Args[0])
	flag.PrintDefaults()
}

func run(mountpoint string, done chan bool) error {
	// TO DO: Strip off trailing / from mountpoint
	c, err := fuse.Mount(
		mountpoint,
		fuse.FSName("messagefs"),
		fuse.Subtype("messagefs"),
		fuse.LocalVolume(),
		fuse.VolumeName("Message filesystem"),
	)

	if err != nil {
		return err
	}

	defer c.Close()

	if p := c.Protocol(); !p.HasInvalidate() {
		log.Panicln("kernel FUSE support is too old to have invalidations: version %v", p)
	}
	srv := fs.New(c, nil)

	n := Node{
		name: "head",
		inode: NewInode(),
		path: mountpoint + "/",
	}

	r := &Dir{
		Node: n,
		files: &[]*File{},
		directories: &[]*Dir{},
	}

	filesys = &FS{r}

	log.Println("About to serve fs")
	if err := srv.Serve(filesys); err != nil {
		return err
	}

	done<-true

	// Check if the mount process has an error to report.
	<-c.Ready
	if err := c.MountError; err != nil {
		return err
	}

	return nil
} // func run

func main() {
	flag.Usage = Usage
	flag.Parse()

	if flag.NArg() != 1 {
		Usage()
		os.Exit(2)
	}

	mountpoint := flag.Arg(0)

	done := make(chan bool)

	go run(mountpoint, done)

	time.Sleep(1*time.Second)

	if err := os.Mkdir(mountpoint + "/all/", 0777); err != nil {
    log.Panicln(err)
  }

	go Chat(mountpoint)

	<-done

} // func main
