package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	appVersion = "0.1"
)

var (
	version = flag.Bool("version", false, "Print the version number.")
	v       = flag.Bool("verbose", false, "Print the detailed log.")
	r       = flag.Bool("r", false, "Rename the file with the extracted date <Year>-<Month>-<Day>_<Originale name>.jpg")
	src     = flag.String("src", "", "Source folder that contains the files to process.")
	dst     = flag.String("dst", "", "Destination folder.")
)

func main() {
	flag.Parse() // Scan the arguments list

	if *version {
		fmt.Println("Version:", appVersion)
		return
	}

	if *src == "" {
		log.Println("[ERROR] Source folder is missing")
		return

	}
	if *dst == "" {
		log.Println("[ERROR] Destination folder is missing")
		return

	}
	numScanned := 0

	var scan = func(path string, f os.FileInfo, err error) error {
		if err != nil {
			log.Println("[ERROR] ", path, " does not exist or is not accessible on the filesystem:", err.Error())
			return err
		}
		if !f.IsDir() {
			numScanned++
			if *v {
				log.Println("[INFO] Processing file:", path)
			}

			mimeType, err := GetMIME(path)
			if err != nil {
				log.Println("[ERROR] Can't read MIME type:", err.Error())
				return err
			}
			if *v {
				log.Println("[INFO] MIME type of the file:", mimeType)
			}
			var timeTaken time.Time
			switch {
			case strings.Contains(mimeType, "image/"):
				timeTaken, err = getExifDate(path)
				if err != nil {
					log.Println("[ERROR] Can't read date time:", err.Error())
					// In case of NO_TIME_TAKEN
					timeTaken = f.ModTime()
				}

			case strings.Contains(mimeType, "video/"):
				// Waiting for a better solution...
				timeTaken = f.ModTime()

			default:
				timeTaken = f.ModTime()
			}
			if *v {
				log.Println("[INFO] Detected Date Time:", timeTaken.Format(time.RFC3339))
			}
			//Destination folder
			fLvl1 := strconv.Itoa(timeTaken.Year())
			fLvl2 := timeTaken.Format("2006-01")
			fLvl3 := timeTaken.Format("2006-01-02")

			dstFolder := filepath.Join(*dst, fLvl1, fLvl2, fLvl3)
			dstFile := f.Name()
			if *r {
				dstFile = fLvl3 + "_" + f.Name()
			}
			if *v {
				log.Println("[INFO] Move the file to:", dstFolder, dstFile)
			}

			// Is the destination file exists?
			if _, err := os.Stat(filepath.Join(dstFolder, dstFile)); err == nil {
				if *v {
					log.Println("[INFO] The file already exists in destination folder:", dstFolder, dstFile)
				}
			} else {
				if err := MoveFile(path, dstFolder, dstFile); err != nil {
					log.Println("[ERROR] Can't move the file", path, "to:", dstFolder)
				}
			}
		}
		return nil
	}
	filepath.Walk(fmt.Sprint(*src), scan)
	fmt.Println("[INFO] Total scanned:", numScanned)
}
