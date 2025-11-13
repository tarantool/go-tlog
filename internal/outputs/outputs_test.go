package outputs_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tarantool/go-tlog/internal/outputs"
)

// Assert *outputs.Outputs implements io.WriteCloser.
var _ io.WriteCloser = (*outputs.Outputs)(nil)

func Test_New_BadPath(t *testing.T) {
	require := require.New(t)

	_, err := outputs.New("/not/exist")
	require.ErrorContains(err, "open /not/exist: no such file or directory")
}

func Test_New_MultipleWithBadPath(t *testing.T) {
	require := require.New(t)

	_, err := outputs.New("/dev/null,/not/exist")
	require.ErrorContains(err, "open /not/exist: no such file or directory")
}

func Test_New_EmptyPaths(t *testing.T) {
	require := require.New(t)

	_, err := outputs.New("")
	require.ErrorContains(err, "empty paths")
}

func Test_New_MultipleWithEmptyPath(t *testing.T) {
	require := require.New(t)

	_, err := outputs.New("/dev/null,")
	require.ErrorContains(err, "empty path")
}

func Test_Outputs_Std(t *testing.T) {
	for _, tc := range []string{"stdout", "stderr"} {
		t.Run(tc, func(t *testing.T) {
			require := require.New(t)

			r, w, _ := os.Pipe()

			switch tc {
			case "stdout":
				orig := os.Stdout
				os.Stdout = w

				defer func() {
					os.Stdout = orig
				}()
			case "stderr":
				orig := os.Stderr
				os.Stderr = w

				defer func() {
					os.Stderr = orig
				}()
			}

			outputs, err := outputs.New(tc)
			require.NoError(err)

			_, err = outputs.Write([]byte("log_message"))
			require.NoError(err)

			_ = w.Close()

			out, err := io.ReadAll(r)
			require.NoError(err)
			require.Contains(string(out), "log_message")
		})
	}
}

func Test_Outputs_File(t *testing.T) {
	require := require.New(t)

	filename := filepath.Join(os.TempDir(), "Test_Outputs_File.log")
	defer func(name string) {
		_ = os.Remove(name)
	}(filename)

	outputs, err := outputs.New(filename)
	require.NoError(err)

	_, err = outputs.Write([]byte("log_message"))
	require.NoError(err)

	out, err := os.ReadFile(filename)
	require.NoError(err)
	require.Contains(string(out), "log_message")
}

func Test_Outputs_Multiple(t *testing.T) {
	require := require.New(t)

	// Prepare file 1.
	filename1 := filepath.Join(os.TempDir(), "Test_Outputs_Multiple1.log")
	defer func(name string) {
		_ = os.Remove(name)
	}(filename1)

	// Prepare stdout.
	r, w, _ := os.Pipe()

	origStdout := os.Stdout
	os.Stdout = w

	defer func() {
		os.Stdout = origStdout
	}()

	// Prepare file 2.
	filename2 := filepath.Join(os.TempDir(), "Test_Outputs_Multiple2.log")
	defer func(name string) {
		_ = os.Remove(name)
	}(filename2)

	outputs, err := outputs.New(filename1 + ",stdout," + filename2)
	require.NoError(err)

	_, err = outputs.Write([]byte("log_message"))
	require.NoError(err)

	// Assert file 1 contents.
	file1Out, err := os.ReadFile(filename1)
	require.NoError(err)
	require.Contains(string(file1Out), "log_message")

	// Assert stdout contents.
	_ = w.Close()

	stdOut, err := io.ReadAll(r)
	require.NoError(err)
	require.Contains(string(stdOut), "log_message")

	// Assert file 2 contents.
	file2Out, err := os.ReadFile(filename2)
	require.NoError(err)
	require.Contains(string(file2Out), "log_message")
}
