package gcpwebserv

import (
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

// SetupStaticFileHandler sets up a static file server based on the
// route and path provided. The diectory must exist to succeed. Using
// the default route of / will not work currently
func SetupStaticFileHandler(httpRoute string, fileSystemPath string) error {

	// Check to see if the directory exists. These issues can be hard
	// to troubleshoot so prechecking is in order.
	d, err := os.Stat(fileSystemPath)
	if err != nil {
		slog.Error("Error with static file server", "error", err)
		return err
	}
	if !d.IsDir() {
		slog.Error("Static Path is not a directory", "path", fileSystemPath)
		return fmt.Errorf("The path %s in not a directory", fileSystemPath)
	}

	fsys := dotFileHidingFileSystem{http.Dir(fileSystemPath)}

	mux.Handle(httpRoute, http.StripPrefix("/"+fileSystemPath, http.FileServer(fsys)))
	return nil
}

// These Functions are copied directly from the net/http package
//
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

	fmt.Printf("n:%v\n", n)
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

	fmt.Printf("requested %s\n", name)
	if containsDotFile(name) { // If dot file, return 403 response
		return nil, fs.ErrPermission
	}

	file, err := fsys.FileSystem.Open(name)
	if err != nil {
		fmt.Printf("%v", name)
		return nil, err
	}

	// Updated to prevents directories from being listed
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		return nil, fs.ErrPermission
	}

	return dotFileHidingFile{file}, err
}
