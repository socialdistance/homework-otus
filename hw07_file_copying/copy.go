package main

import (
	"errors"
	"io"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrInvalidOffset         = errors.New("invalid offset")
	ErrInvalidLimit          = errors.New("invalid limit")
	ErrOpenFile              = errors.New("failed open file")
	ErrCreateFile            = errors.New("failed create file")
	ErrSeek                  = errors.New("failed seek")
	ErrCopyFile              = errors.New("failed copy file")
	ErrInvalidFile           = errors.New("error invalid file")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if offset < 0 {
		return ErrInvalidOffset
	}

	if limit < 0 {
		return ErrInvalidLimit
	}

	bytes, err := os.Stat(fromPath)
	if err != nil {
		return ErrInvalidFile
	}

	if bytes.IsDir() {
		return ErrInvalidFile
	}
	if bytes.Size() < 1 {
		return ErrUnsupportedFile
	}

	if offset > bytes.Size() {
		return ErrOffsetExceedsFileSize
	}
	if limit == 0 {
		limit = bytes.Size()
	}

	file, err := os.OpenFile(fromPath, os.O_RDONLY, 0o777)
	if err != nil {
		return ErrOpenFile
	}

	defer file.Close()

	fileWrite, err := os.Create(toPath)
	if err != nil {
		return ErrCreateFile
	}
	defer fileWrite.Close()

	_, err = file.Seek(offset, 0)
	if err != nil {
		return ErrSeek
	}

	bar := pb.New(int(limit))
	bar.SetRefreshRate(time.Millisecond * 100)
	bar.Start()
	defer bar.Finish()
	proxy := bar.NewProxyReader(file)

	_, err = io.CopyN(fileWrite, proxy, limit)
	if errors.Is(err, io.EOF) {
		return nil
	}
	if errors.Is(err, nil) {
		return nil
	}

	return ErrCopyFile
}
