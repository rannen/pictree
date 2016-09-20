package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

//GetMIME return the MIME Type of the media file
func GetMIME(file string) (string, error) {
	buffMIME, err := os.Open(file)
	if err != nil {
		return "", err
	}

	buff := make([]byte, 512)
	// why 512 bytes ? see http://golang.org/pkg/net/http/#DetectContentType
	_, err = buffMIME.Read(buff)
	if err != nil {
		return "", err
	}
	fileType := http.DetectContentType(buff)
	buffMIME.Close()
	return fileType, nil
}

//MoveFile moves file to a destination folder
func MoveFile(src, dstFolder, dstFile string) error {

	if err := CopyFile(src, dstFolder, dstFile); err != nil {
		return fmt.Errorf("[ERROR] Failed to create blob in blobstore: %s", err.Error())
	}
	if err := os.Remove(src); err != nil {
		return fmt.Errorf("[ERROR] Failed to remove temp file: %s %s", src, err.Error())
	}
	return nil
}

//CopyFile copies file to a destination folder
func CopyFile(src, dstFolder, dstFile string) error {
	sfi, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}

	// Test if the file already exists in destination
	if dfi, err := os.Stat(filepath.Join(dstFolder, dstFile)); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("The file already exists in destination: %s %d", dfi.Name(), dfi.Size())
		}
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	// Create folder if it not exists
	if err := os.MkdirAll(dstFolder, os.ModePerm); err != nil {
		return err
	}

	out, err := os.Create(filepath.Join(dstFolder, dstFile))
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	err = out.Close()
	if err != nil {
		return err
	}

	// Check if the copy is ok
	if dfi, err := os.Stat(filepath.Join(dstFolder, dstFile)); err != nil {
		if !os.SameFile(sfi, dfi) {
			return fmt.Errorf("Error during the copy of the file: %s", src)
		}
	}
	return nil
}

//RemoveContents avoid the deletion of the folder in order to avoid problem with permission access denied error on Windows
func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
