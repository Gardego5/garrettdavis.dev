package middleware

import (
	"context"
	"io/fs"
	"net/http"
	"path/filepath"
)

type fileserver struct{ fs http.FileSystem }
type filesystem interface {
	FileSystem(fs http.FileSystem) http.FileSystem
}

func (fs fileserver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s := r.Context().Value(fileserverKey); s != nil {
		if s, ok := s.(filesystem); ok {
			http.FileServer(s.FileSystem(fs.fs)).ServeHTTP(w, r)
			return
		}
	}

	http.NotFound(w, r)
}

func FileServer(fs http.FileSystem) http.Handler { return fileserver{fs} }
func FileServerFS(fs fs.FS) http.Handler         { return fileserver{http.FS(fs)} }

type standardFileSystem struct{}

func FileSystem(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), fileserverKey, standardFileSystem{})))
	})
}
func (sfs standardFileSystem) FileSystem(fs http.FileSystem) http.FileSystem { return fs }

type neuteredFileSystem struct{ fs http.FileSystem }

func NeuteredFileSystem(next http.Handler) http.Handler {
	fs := neuteredFileSystem{}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), fileserverKey, fs)))
	})
}
func (nfs neuteredFileSystem) FileSystem(fs http.FileSystem) http.FileSystem {
	return neuteredFileSystem{fs}
}
func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}
