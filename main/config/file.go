package config

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// Load loads the configuration from the specified file.
func Load(filename string) {
	cfg.Lock()
	defer cfg.Unlock()

	cfg.filename = filename

	info, err := os.Stat(filename)
	if err != nil {
		return
	}
	if info.Size() == 0 {
		return
	}

	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	privcfg, err := parseConfig(file)
	if err != nil {
		return
	}

	cfg.config = privcfg

	return
}

// Save saves the config to the file.
func Save() {
	cfg.Lock()
	defer cfg.Unlock()

	if len(cfg.filename) == 0 {
		log.Println("Error saving the config to the file: missing filename")
		return
	}

	content := bytes.Buffer{}
	content.WriteByte(comment)
	content.WriteByte(' ')
	content.WriteString("Updated at ")
	content.WriteString(time.Now().Format(time.RFC850))

	for sect, opts := range cfg.config {
		content.WriteString("\n")
		content.WriteByte(namestart)
		content.WriteString(sect)
		content.WriteByte(namestop)
		for opt, val := range opts {
			content.WriteString(opt)
			content.WriteByte(sep)
			content.WriteString(val)
		}
	}

	err := ioutil.WriteFile(cfg.filename, content.Bytes(), 0644)
	if err != nil {
		log.Println("Error saving the config to the file:", err)
	}
}
