package archive_extractor

import (
	"errors"
	"fmt"
	archive "github.com/gen2brain/go-unarr"
	"github.com/go-archive-extractor/utils"
	"io"
)

type SevenZipArchvier struct {
}

func (za SevenZipArchvier) ExtractArchive(path string, processingFunc func(header *ArchiveHeader, params map[string]interface{}) error,
	params map[string]interface{}) error {
	r, err := archive.NewArchive(path)
	if err != nil {
		if err.Error() == "unarr: No valid RAR, ZIP, 7Z or TAR archive" {
			return nil
		}
		return err
	}
	allFiles, err := r.List()
	if err != nil {
		if err.Error() == "unarr: No valid RAR, ZIP, 7Z or TAR archive" {
			return nil
		}
		return err
	}
	defer r.Close()
	if err == nil {
		for _, archiveEntry := range allFiles {
			err := r.EntryFor(archiveEntry)
			if err != nil {
				return err
			}
			if !utils.IsFolder(archiveEntry) {
				rc := &SevenZipReader{Archive: r, Size: r.Size()}
				archiveHeader := NewArchiveHeader(rc, r.Name(), r.ModTime().Unix(), int64(r.Size()))
				err = processingFunc(archiveHeader, params)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

type SevenZipReader struct {
	Archive *archive.Archive
	Size    int
}

func (a *SevenZipReader) Read(p []byte) (n int, err error) {
	if a.Size <= 0 {
		return 0, io.EOF
	}
	size := len(p)
	if len(p) > a.Size {
		size = a.Size
	}
	b := make([]byte, size)
	n, err = a.Archive.Read(b)
	if err != nil && err != io.EOF {
		return 0, err
	}
	copied := copy(p, b)
	if copied != n {
		return 0, errors.New(fmt.Sprintf("copy arrays failed, copied only %v from %v bytes", copied, n))
	}
	a.Size -= n
	return n, nil
}
