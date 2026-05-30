package main

import (
	"embed"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path"
	"strings"
)

//go:embed all:frontend_dist
var embeddedFrontend embed.FS

func embeddedFrontendHandler() http.Handler {
	sub, err := fs.Sub(embeddedFrontend, "frontend_dist")
	if err != nil {
		return nil
	}
	if _, err := fs.Stat(sub, "index.html"); err != nil {
		return nil
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clean := strings.TrimPrefix(r.URL.Path, "/")
		if clean == "" {
			if serveEmbeddedFile(w, r, sub, "index.html") {
				return
			}
			http.NotFound(w, r)
			return
		}
		if strings.HasSuffix(clean, "/") {
			if serveEmbeddedFile(w, r, sub, clean+"index.html") {
				return
			}
		}
		if info, err := fs.Stat(sub, clean); err == nil {
			if info.IsDir() {
				if serveEmbeddedFile(w, r, sub, clean+"/index.html") {
					return
				}
			} else {
				if serveEmbeddedFile(w, r, sub, clean) {
					return
				}
			}
		}
		trimmed := strings.TrimSuffix(clean, "/")
		if serveEmbeddedFile(w, r, sub, trimmed+".html") {
			return
		}
		if serveEmbeddedFile(w, r, sub, trimmed+"/index.html") {
			return
		}
		serveEmbeddedFile(w, r, sub, "index.html")
	})
}

func serveEmbeddedFile(w http.ResponseWriter, r *http.Request, sub fs.FS, name string) bool {
	f, err := sub.Open(name)
	if err != nil {
		return false
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil || info.IsDir() {
		return false
	}
	contentType := mime.TypeByExtension(path.Ext(name))
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "no-cache")
	if r.Method == http.MethodHead {
		w.WriteHeader(http.StatusOK)
		return true
	}
	_, _ = io.Copy(w, f)
	return true
}
