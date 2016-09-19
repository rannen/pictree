package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
)

func getExifDate(file string) (tm time.Time, err error) {
	exif.RegisterParsers(mknote.All...)

	fexif, err := os.Open(file)
	defer fexif.Close()
	if err != nil {
		return tm, err
	}
	x, err := exif.Decode(fexif)
	if err != nil {
		return tm, err
	}

	tm, err = x.DateTime()
	if err != nil {
		return tm, fmt.Errorf("NO_TIME_TAKEN")
	}
	return tm, nil
}
