// Shamelessly copied (and slightly altered) from
// https://github.com/client9/misspell/blob/9ce5d97/mime.go.
package ded

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// The number of possible binary formats is very large items that might be
// checked into a repo or be an artifact of a build. Additions welcome.
var binary = map[string]bool{
	".a":     true,
	".bin":   true,
	".bz2":   true,
	".class": true,
	".dll":   true,
	".exe":   true,
	".gif":   true,
	".gpg":   true,
	".gz":    true,
	".ico":   true,
	".jar":   true,
	".jpeg":  true,
	".jpg":   true,
	".mp3":   true,
	".mp4":   true,
	".mpeg":  true,
	".o":     true,
	".pdf":   true,
	".png":   true,
	".pyc":   true,
	".pyo":   true,
	".so":    true,
	".swp":   true,
	".tar":   true,
	".tiff":  true,
	".woff":  true,
	".woff2": true,
	".xz":    true,
	".z":     true,
	".zip":   true,
}

// isBinaryFilename returns true if the file is likely to be binary.
func isBinaryFilename(s string) bool {
	return binary[strings.ToLower(filepath.Ext(s))]
}

var scm = map[string]bool{
	".bzr": true,
	".git": true,
	".hg":  true,
	".svn": true,
	"CVS":  true,
}

// isSCMPath returns true if the path is likely part of a (private) SCM
// directory.
func isSCMPath(s string) bool {
	if strings.Contains(filepath.Base(s), "EDITMSG") {
		return false
	}
	parts := strings.Split(filepath.Clean(s), string(filepath.Separator))
	for _, dir := range parts {
		if scm[dir] {
			return true
		}
	}
	return false
}

var magicHeaders = [][]byte{
	// PGP messages and signatures are "text" but really just
	// blobs of base64-text and should not be checked.
	[]byte("-----BEGIN PGP MESSAGE-----"),
	[]byte("-----BEGIN PGP SIGNATURE-----"),

	// ELF
	{0x7f, 0x45, 0x4c, 0x46},

	// Postscript
	{0x25, 0x21, 0x50, 0x53},

	// PDF
	{0x25, 0x50, 0x44, 0x46},

	// Java class file
	// https://en.wikipedia.org/wiki/Java_class_file
	{0xCA, 0xFE, 0xBA, 0xBE},

	// PNG
	// https://en.wikipedia.org/wiki/Portable_Network_Graphics
	{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a},

	// ZIP, JAR, ODF, OOXML
	{0x50, 0x4B, 0x03, 0x04},
	{0x50, 0x4B, 0x05, 0x06},
	{0x50, 0x4B, 0x07, 0x08},
}

func isTextFile(raw []byte) bool {
	for _, magic := range magicHeaders {
		if bytes.HasPrefix(raw, magic) {
			return false
		}
	}

	// Allow any text/type with utf-8 encoding.
	// DetectContentType sometimes returns charset=utf-16 for XML stuff, in which
	// case we should ignore.
	mime := http.DetectContentType(raw)
	return strings.HasPrefix(mime, "text/") && strings.HasSuffix(mime, "charset=utf-8")
}

// ReadTextFile returns the contents of a file, first testing if it is a text file
//  returns ("", nil) if not a text file
//  returns ("", error) if error
//  returns (string, nil) if text
//
// unfortunately, in worse case, this does
//   1 stat
//   1 open,read,close of 512 bytes
//   1 more stat,open, read everything, close (via ioutil.ReadAll)
//  This could be kinder to the filesystem.
//
// This uses some heuristics of the file's extension (e.g. .zip, .txt) and uses
// a sniffer to determine if the file is text or not. Using file extensions
// isn't great, but probably good enough for real-world use.
//
// Go's built-in sniffer is problematic for different reasons. It's optimized
// for HTML, and is very limited in detection.  It would be good to explicitly
// add some tests for ELF/DWARF formats to make sure we never corrupt binary
// files.
func ReadTextFile(filename string) (string, error) {
	if isBinaryFilename(filename) {
		return "", nil
	}

	if isSCMPath(filename) {
		return "", nil
	}

	fstat, err := os.Stat(filename)
	if err != nil {
		return "", fmt.Errorf("unable to stat %q: %s", filename, err)
	}

	if fstat.IsDir() {
		return "", nil
	}

	// Avoid reading _large_ files.
	// If input is large, read the first 512 bytes to sniff type. Return iff
	// non-text.
	isText := false
	if fstat.Size() > 50000 {
		fin, err := os.Open(filename)
		if err != nil {
			return "", fmt.Errorf("unable to open large file %q: %s", filename, err)
		}
		defer fin.Close()
		buf := make([]byte, 512)
		if _, err = io.ReadFull(fin, buf); err != nil {
			return "", fmt.Errorf("unable to read 512 bytes from %q: %s", filename, err)
		}
		if !isTextFile(buf) {
			return "", nil
		}

		// Avoid double-checking this file.
		isText = true
	}

	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("unable to read all %q: %s", filename, err)
	}
	if !isText && !isTextFile(raw) {
		return "", nil
	}
	return string(raw), nil
}
