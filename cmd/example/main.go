// What it does:
//
// This example detects motion using a delta threshold from the first frame,
// and then finds contours to determine where the object is located.
//
// Very loosely based on Adrian Rosebrock code located at:
// http://www.pyimagesearch.com/2015/06/01/home-surveillance-and-motion-detection-with-the-raspberry-pi-python-and-opencv/
//
// How to run:
//
// 		go run ./cmd/motion-detect/main.go 0
//

package main

import (
	"flag"
	"io/fs"
	"log"
	"os"
	"path"

	"github.com/andreykyz/gomoveextractor/coder"
)

func main() {

	in := flag.String("in", "", "path to folder")
	out := flag.String("out", "", "path to folder")
	ext := flag.String("ext", ".mp4", "file extension")
	file := flag.String("file", "", "path to single file")
	flag.Parse()
	files := []string{}
	if *out != "" {
		files = listFiles(*in, *ext)
	}
	if *file != "" {
		files = append(files, *file)
	}
	c := coder.NewCoder(coder.CoderArgs{
		InputFiles:       files,
		OutputVideoDir:   "./",
		OutputRectangles: "./",
	})
	if err := c.Generate(); err != nil {
		log.Fatal(err)

	}

}

func listFiles(dir string, ext string) []string {
	root := os.DirFS(dir)

	mdFiles, err := fs.Glob(root, "*"+ext)

	if err != nil {
		log.Fatal(err)
	}

	files := make([]string, 0, len(mdFiles))
	for _, v := range mdFiles {
		files = append(files, path.Join(dir, v))
	}
	return files
}
