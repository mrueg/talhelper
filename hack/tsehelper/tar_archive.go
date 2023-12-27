// tar_archive.go
package main

import (
	"archive/tar"
	"bytes"
	"io"
	"strings"
)

// processTarArchive processes a tar archive and returns a slice of strings or an error.
func processTarArchive(tarData []byte) ([]string, error) {
	// Create a new reader from the tar data
	tarReader := tar.NewReader(bytes.NewReader(tarData))
	// Iterate through the files in the tar archive
	var exts []string
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			// End of tar archive
			break
		}
		if err != nil {
			return nil, err
		}

		// Check if the current entry matches the target filename
		targetFilename := "image-digests"
		if header.Name == targetFilename {
			// Read the content of the file
			content := make([]byte, header.Size)
			_, err := tarReader.Read(content)
			if err != io.EOF {
				return nil, err
			}
			// Process the file content, split by new line
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				// Skip empty lines
				if line == "" {
					continue
				}
				exts = append(exts, line)
			}
		}
	}

	return exts, nil
}
