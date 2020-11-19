package utils

import (
	"io"
	"os"
	"path"
)

// StaticPrefix is prefix for path to files
const StaticPrefix = "/static/"

// StaticDir is folder of static files
func StaticDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return path.Join(cwd, "web", StaticPrefix), nil
}

func ServeStaticFile(filePath string, w io.Writer) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	buf := make([]byte, 4*1024) // 4Kb
	if _, err = io.CopyBuffer(w, f, buf); err != nil {
		return err
	}

	return nil
}
