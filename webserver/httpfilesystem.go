// How to Use //go:embed
// https://blog.carlmjohnson.net/post/2021/how-to-use-go-embed/
package webserver

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
)

func GetHttpFileSystem(embededFiles embed.FS, rootDir string, useEmbedded bool) http.FileSystem {
	if useEmbedded {
		log.Printf("using embedded files for %s", rootDir)
		fsys, err := fs.Sub(embededFiles, rootDir)
		if err != nil {
			panic(err)
		}
		return http.FS(fsys)
	} else {
		log.Printf("using files on disk for %s", rootDir)
		return http.FS(os.DirFS(rootDir))
	}
}

func GetFileSystem(embededFiles embed.FS, rootDir string, useEmbedded bool) fs.FS {
	if useEmbedded {
		log.Printf("using embedded files for %s", rootDir)
		fsys, err := fs.Sub(embededFiles, rootDir)
		if err != nil {
			panic(err)
		}
		return fsys
	} else {
		log.Printf("using files on disk for %s", rootDir)
		return os.DirFS(rootDir)
	}
}
