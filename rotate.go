package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/abibby/uptime/config"
	"github.com/davecgh/go-spew/spew"
)

func rotate(t time.Time) error {
	files, err := os.ReadDir(config.LogPath)
	if err != nil {
		return err
	}

	for _, f := range files {
		name := f.Name()
		date, err := time.Parse("2006-01-02.csv", name)
		if err != nil {
			continue
		}
		startOfDay := t.Truncate(24 * time.Hour)
		if !date.Equal(startOfDay) {
			err := rotateFile(path.Join(config.LogPath, name), startOfDay)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func rotateFile(p string, t time.Time) error {
	f, err := os.Open(p)
	if err != nil {
		return err
	}
	defer f.Close()

	gzF, err := OpenFile(p + ".gz")
	if err != nil {
		return err
	}
	defer gzF.Close()

	agF, err := OpenFile(config.AggregatePath, fmt.Sprintf("%s.csv", t.Format("2006-01")))
	if err != nil {
		return err
	}
	defer gzF.Close()

	w := gzip.NewWriter(gzF)
	defer w.Close()

	s := bufio.NewScanner(f)

	good := 0
	bad := 0

	for s.Scan() {
		line := s.Text()

		w.Write([]byte(line))

		parts := strings.SplitN(line, "\t", 2)
		if len(parts) < 2 {
			continue
		}
		succeeded, err := strconv.Atoi(parts[1])
		if err != nil {
			spew.Dump(parts)
			return err
		}
		if succeeded > 0 {
			good++
		} else {
			bad++
		}
	}

	err = s.Err()
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(agF, "%s\t%0.2f%%\t%d\t%d\n", t.Format("2006-01-02"), float32(bad)/float32(good+bad)*100, good, bad)
	if err != nil {
		return err
	}

	err = os.Remove(p)
	if err != nil {
		return err
	}

	return nil
}
