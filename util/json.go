package util

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func WriteJSONToFile(path string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if dir != "." && dir != ".." {
		err = os.Mkdir(dir, 0o777)
		if err != nil && !os.IsExist(err) {
			return fmt.Errorf("failed to create log directory: %v", err)
		}
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.WithError(err).Errorf("unable to open file %q", path)
		return err
	}

	defer file.Close()

	_, err = file.Write(b)
	if err != nil {
		return err
	}

	return nil
}
