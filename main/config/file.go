package config

import (
	"fmt"
	"os"
)

// FromFile loads the configuration from the specified file.
func FromFile(filename string) (err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return
	}
	if info.Size() == 0 {
		err = fmt.Errorf("empty %s", filename)
		return
	}

	Cfg, err = parseConfig(file)
	if err != nil {
		return
	}

	err = validate(Cfg)

	return
}
