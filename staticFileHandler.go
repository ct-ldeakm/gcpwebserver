package gcpwebserv

import (
	"fmt"
	"io/fs"
	"net/http"
	"strings"
)

// Setup Static
func SetupStaticFileHandler(httpRoute string, fileSystemPath string) {

	fsys := dotFileHidingFileSystem{http.Dir("/static/")}
	//http.FileServer(fsys)
	// handler := func(w http.ResponseWriter, r *http.Request) {
	// 	log.Printf("In static Handler:%s %v", path, fsys)
	// 	http.FileServer(fsys)
	// }
	//RouteHandler(httpRoute, http.FileServer(fsys))
	//mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(fsys)))
	mux.Handle("/static/", http.FileServer(fsys))
	//return handler
}

// These Functions are copied directly from the net/http package

// containsDotFile reports whether name contains a path element starting with a period.
// The name is assumed to be a delimited by forward slashes, as guaranteed
// by the http.FileSystem interface.
func containsDotFile(name string) bool {
	parts := strings.Split(name, "/")
	for _, part := range parts {
		if strings.HasPrefix(part, ".") {
			return true
		}
	}
	return false
}

// dotFileHidingFile is the http.File use in dotFileHidingFileSystem.
// It is used to wrap the Readdir method of http.File so that we can
// remove files and directories that start with a period from its output.
type dotFileHidingFile struct {
	http.File
}

// Readdir is a wrapper around the Readdir method of the embedded File
// that filters out all files that start with a period in their name.
func (f dotFileHidingFile) Readdir(n int) (fis []fs.FileInfo, err error) {
	fmt.Printf("%s", "In Func")
	files, err := f.File.Readdir(n)
	for _, file := range files { // Filters out the dot files
		fmt.Printf("%s", file.Name())
		if !strings.HasPrefix(file.Name(), ".") {
			fis = append(fis, file)
		}

	}
	return
}

// dotFileHidingFileSystem is an http.FileSystem that hides
// hidden "dot files" from being served.
type dotFileHidingFileSystem struct {
	http.FileSystem
}

// Open is a wrapper around the Open method of the embedded FileSystem
// that serves a 403 permission error when name has a file or directory
// with whose name starts with a period in its path.
func (fsys dotFileHidingFileSystem) Open(name string) (http.File, error) {

	//fmt.Printf("%s\n", name)
	if containsDotFile(name) { // If dot file, return 403 response
		return nil, fs.ErrPermission
	}

	file, err := fsys.FileSystem.Open(name)
	if err != nil {
		fmt.Printf("%v", name)
		return nil, err
	}

	// Prevents directories from being listed
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		return nil, fs.ErrPermission
	}

	return dotFileHidingFile{file}, err
}
