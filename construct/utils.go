package construct

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// findProjectRoot looks for main.go to determine project root
func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if fileExists(filepath.Join(dir, "main.go")) {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("project root not found (no main.go)")
		}
		dir = parent
	}
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// detectPackageManager detects which package manager to use
func detectPackageManager(dir, script string) *exec.Cmd {
	if fileExists(filepath.Join(dir, "bun.lockb")) {
		return exec.Command("bun", "run", script)
	}
	if fileExists(filepath.Join(dir, "pnpm-lock.yaml")) {
		return exec.Command("pnpm", script)
	}
	if fileExists(filepath.Join(dir, "yarn.lock")) {
		return exec.Command("yarn", script)
	}
	return exec.Command("npm", "run", script)
}

// PrefixWriter adds a prefix to each line of output
type PrefixWriter struct {
	prefix string
	writer *os.File
	buffer []byte
}

// NewPrefixWriter creates a new PrefixWriter
func NewPrefixWriter(writer *os.File, prefix string) *PrefixWriter {
	return &PrefixWriter{
		prefix: prefix,
		writer: writer,
		buffer: make([]byte, 0),
	}
}

// Write implements io.Writer interface
func (pw *PrefixWriter) Write(p []byte) (n int, err error) {
	pw.buffer = append(pw.buffer, p...)

	// Process complete lines
	for {
		idx := strings.IndexByte(string(pw.buffer), '\n')
		if idx == -1 {
			break
		}

		line := pw.buffer[:idx+1]
		pw.writer.WriteString(pw.prefix)
		pw.writer.Write(line)
		pw.buffer = pw.buffer[idx+1:]
	}

	return len(p), nil
}

// unzip extracts a zip archive to a destination directory
func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		// Construct the file path
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip vulnerability
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			// Create directory
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Create parent directories if needed
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		// Create the file
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		// Extract the file
		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}