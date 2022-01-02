package fs

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Local implements local file storage
type Local struct {
	rootDir string
}

func NewLocal(rootDir string) (*Local, error) {
	return &Local{rootDir: rootDir}, nil
}

func (l *Local) Open(name string) (http.File, error) {
	path := filepath.Join(l.rootDir, name)
	return os.Open(path)
}

func (l *Local) Delete(_ctx context.Context, name string) error {
	path := filepath.Join(l.rootDir, name)
	return os.Remove(path)
}

func (l *Local) Create(_ctx context.Context, name string, reader io.Reader) (int64, error) {
	var (
		logger = log.WithField("name", name)
		path   = filepath.Join(l.rootDir, name)
	)

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return 0, errors.Wrapf(err, "failed to mkdir: %s", path)
	}

	logger.Infof("creating file: %s", path)
	written, err := l.copyFile(reader, path)
	if err != nil {
		return 0, errors.Wrap(err, "failed to copy file")
	}

	logger.Debugf("written %d bytes", written)
	return written, nil
}

func (l *Local) copyFile(source io.Reader, destinationPath string) (int64, error) {
	dest, err := os.Create(destinationPath)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create destination file")
	}

	defer dest.Close()

	written, err := io.Copy(dest, source)
	if err != nil {
		return 0, errors.Wrap(err, "failed to copy data")
	}

	return written, nil
}
