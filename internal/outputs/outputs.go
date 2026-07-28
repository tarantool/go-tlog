package outputs

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
	"sync"
)

// Errors reported when a path list cannot be turned into outputs.
var (
	// ErrEmptyPaths is returned by New when the path list is empty.
	ErrEmptyPaths = errors.New("empty paths")
	// ErrEmptyPath is returned when the path list contains an empty element,
	// as in "stdout,,stderr".
	ErrEmptyPath = errors.New("empty path")
)

// Outputs is io.WriteCloser for multiple output paths.
type Outputs struct {
	mu    sync.RWMutex
	paths []string
	files []*os.File
	w     io.Writer
}

// New creates Outputs from comma-separated string of paths.
// Use "stdout" and "stderr" for os streams and file paths for files.
func New(paths string) (*Outputs, error) {
	if paths == "" {
		return nil, ErrEmptyPaths
	}

	pathsSlice := splitPaths(paths)

	files, w, err := openFiles(pathsSlice)
	if err != nil {
		return nil, err
	}

	return &Outputs{
		paths: pathsSlice,
		files: files,
		w:     w,
	}, nil
}

func openFiles(paths []string) ([]*os.File, io.Writer, error) {
	files := make([]*os.File, 0, len(paths))
	writers := make([]io.Writer, 0, len(paths))

	for _, path := range paths {
		file, err := openFile(path)
		if err != nil {
			_ = multiClose(files)

			return nil, nil, fmt.Errorf("failed to open path %q: %w", path, err)
		}

		files = append(files, file)
		writers = append(writers, file)
	}

	return files, io.MultiWriter(writers...), nil
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
		return nil, ErrEmptyPath
	default:
		// #nosec G304 -- the log destination is chosen by the application that
		// configures the logger, so a variable path is the intended behavior.
		//nolint:wrapcheck // New wraps this error with the offending path.
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
	o.mu.RLock()
	defer o.mu.RUnlock()

	//nolint:wrapcheck // io.Writer callers match on the underlying error
	// (os.ErrClosed, syscall.EPIPE); wrapping it would break that.
	return o.w.Write(p)
}

// Reopen closes current file outputs and opens them once again.
// Can be used for logrotate.
func (o *Outputs) Reopen() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	files, w, err := openFiles(o.paths)
	if err != nil {
		return fmt.Errorf("failed to reopen outputs: %w", err)
	}

	oldFiles := o.files
	o.files = files
	o.w = w

	if err := multiClose(oldFiles); err != nil {
		return fmt.Errorf("failed to close old outputs (reopen succeeded): %w", err)
	}

	return nil
}

// Close closes all file outputs except stdout and stderr.
func (o *Outputs) Close() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	return multiClose(o.files)
}
