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
	// Special handling for install command (no "run" prefix needed)
	isInstall := script == "install"

	if fileExists(filepath.Join(dir, "bun.lockb")) || fileExists(filepath.Join(dir, "bun.lock")) {
		if isInstall {
			return exec.Command("bun", "install")
		}
		return exec.Command("bun", "run", script)
	}
	if fileExists(filepath.Join(dir, "pnpm-lock.yaml")) {
		if isInstall {
			return exec.Command("pnpm", "install")
		}
		return exec.Command("pnpm", script)
	}
	if fileExists(filepath.Join(dir, "yarn.lock")) {
		if isInstall {
			return exec.Command("yarn", "install")
		}
		return exec.Command("yarn", script)
	}
	// Default to npm
	if isInstall {
		return exec.Command("npm", "install")
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

// FilteredWriter filters output based on keywords
type FilteredWriter struct {
	prefix   string
	writer   *os.File
	buffer   []byte
	keywords []string
}

// NewFilteredWriter creates a new FilteredWriter that only shows lines containing keywords
func NewFilteredWriter(writer *os.File, prefix string, keywords []string) *FilteredWriter {
	return &FilteredWriter{
		prefix:   prefix,
		writer:   writer,
		buffer:   make([]byte, 0),
		keywords: keywords,
	}
}

// Write implements io.Writer interface with filtering
func (fw *FilteredWriter) Write(p []byte) (n int, err error) {
	fw.buffer = append(fw.buffer, p...)

	// Process complete lines
	for {
		idx := strings.IndexByte(string(fw.buffer), '\n')
		if idx == -1 {
			break
		}

		line := string(fw.buffer[:idx])
		fw.buffer = fw.buffer[idx+1:]

		// Check if line contains any keyword
		shouldShow := false
		for _, keyword := range fw.keywords {
			if strings.Contains(line, keyword) {
				shouldShow = true
				break
			}
		}

		if shouldShow {
			fw.writer.WriteString(fw.prefix)
			fw.writer.WriteString(line)
			fw.writer.WriteString("\n")
		}
	}

	return len(p), nil
}

// Unzip extracts a zip archive to a destination directory
func Unzip(src, dest string) error {
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

// pluralize converts a word to plural form (simple implementation)
func pluralize(word string) string {
	if strings.HasSuffix(word, "s") || strings.HasSuffix(word, "x") ||
		strings.HasSuffix(word, "ch") || strings.HasSuffix(word, "sh") {
		return word + "es"
	}
	if strings.HasSuffix(word, "y") {
		return word[:len(word)-1] + "ies"
	}
	return word + "s"
}

// titleCase converts first letter to uppercase
func titleCase(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}