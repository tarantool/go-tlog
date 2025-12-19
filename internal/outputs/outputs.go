package outputs

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
)

// Outputs is io.WriteCloser for multiple output paths.
type Outputs struct {
	files []*os.File
	w     io.Writer
}

// New creates Outputs from comma-separated string of paths.
// Use "stdout" and "stderr" for os streams and file paths for files.
func New(paths string) (*Outputs, error) {
	if paths == "" {
		return nil, errors.New("empty paths")
	}

	slice := splitPaths(paths)

	files := make([]*os.File, 0, len(slice))
	writers := make([]io.Writer, 0, len(slice))

	for _, path := range slice {
		file, err := openFile(path)
		if err != nil {
			_ = multiClose(files)

			return nil, fmt.Errorf("failed to open path %q: %w", path, err)
		}

		files = append(files, file)
		writers = append(writers, file)
	}

	return &Outputs{
		files: files,
		w:     io.MultiWriter(writers...),
	}, nil
}

func splitPaths(paths string) []string {
	if paths == "" {
		return []string{}
	}

	split := strings.Split(paths, ",")

	for i, path := range split {
		split[i] = strings.TrimSpace(path)
	}

	return split
}

// https://github.com/uber-go/zap/blob/6d482535bdd97f4d97b2f9573ac308f1cf9b574e/sink.go#L158
var defaultFilePerms uint32 = 0o666

func openFile(path string) (*os.File, error) {
	switch path {
	case "stdout":
		// https://github.com/uber-go/zap/blob/6d482535bdd97f4d97b2f9573ac308f1cf9b574e/sink.go#L153-L154
		return os.Stdout, nil
		// https://github.com/uber-go/zap/blob/6d482535bdd97f4d97b2f9573ac308f1cf9b574e/sink.go#L155-L156
	case "stderr":
		return os.Stderr, nil
	case "":
		return nil, errors.New("empty path")
	default:
		return os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, fs.FileMode(defaultFilePerms))
	}
}

func multiClose(files []*os.File) error {
	errs := make([]error, 0, len(files))

	for _, file := range files {
		switch file {
		case os.Stdout, os.Stderr, nil:
			continue
		default:
			errs = append(errs, file.Close())
		}
	}

	return errors.Join(errs...)
}

// Write writes p to all configured output destinations.
// It implements io.Writer and is used by slog handlers.
func (o *Outputs) Write(p []byte) (int, error) {
	return o.w.Write(p)
}

// Close closes all file outputs except stdout and stderr.
func (o *Outputs) Close() error {
	return multiClose(o.files)
}
