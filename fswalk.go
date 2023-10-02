package tepidreload

import (
	"errors"
	"io/fs"
	"log"
	"strings"
	"time"
)

var (
	WalkErr = errors.New("failed dir walk")
)

func CheckFileMods(config Config) (bool, error) {
	isChanged := false

	timeAgo := time.Now().UTC().Add(time.Millisecond * time.Duration(config.TickIntervalMS*-1))

	err := fs.WalkDir(config.WatchPath, ".", func(path string, dirInfo fs.DirEntry, err error) error {

		if err != nil {
			log.Printf("error accessing file for reload monitoring: %v", dirInfo.Name())
			return err
		}

		if dirInfo.IsDir() {
			for _, exclusionDir := range config.ExcludeDirs {
				if dirInfo.Name() == exclusionDir {
					return fs.SkipDir
				}
			}
			return nil
		}

		for _, exclusionExt := range config.ExcludeExtensions {
			if strings.HasSuffix(dirInfo.Name(), exclusionExt) {
				return nil
			}
		}

		for _, exclusionFilename := range config.ExcludeFiles {
			if exclusionFilename == dirInfo.Name() {
				return nil
			}
		}

		fileInfo, err := dirInfo.Info()
		if err != nil {
			log.Printf("error accessing file for reload monitoring: %v", dirInfo.Name())
			return err
		}

		if fileInfo.ModTime().UTC().After(timeAgo) {
			isChanged = true
			return fs.SkipAll
		}

		return nil
	}) // end fs.WalkDirFunc

	if err != nil {
		return false, WalkErr
	}

	return isChanged, nil
}
