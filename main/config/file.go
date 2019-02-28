package config

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const maxFileSize = 1 << 20

// Load loads the configuration from the specified file.
func Load(filename string) {
	setFilename(filename)

	info, err := os.Stat(filename)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Println(err)
			return
		}
		if _, err = os.Create(filename); err != nil {
			log.Println(err)
			return
		}
		log.Printf("INFO: '%s' empty configuration file created. No config loaded\n", filename)
		return
	}
	size := info.Size()
	if size == 0 {
		log.Printf("WARNING! '%s' configuration file is empty. No config loaded\n", filename)
		return
	}
	if size > maxFileSize {
		log.Printf("WARNING! '%s' size %v exceeded the limit of 1 Mb. No config loaded\n", filename, size)
		return
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	c, err := parseConfig(file)
	if err != nil {
		log.Println(err)
		return
	}

	_, errs := Update(c)
	if errs != nil {
		for _, e := range errs {
			log.Println(e)
		}
	}

	return
}

// Save saves the config to the file.
func Save() error {
	cfg.Lock()
	defer cfg.Unlock()

	if len(cfg.filename) == 0 {
		err := errors.New("Error saving the config to the file: missing filename")
		log.Println(err)
		return err
	}

	content := bytes.Buffer{}
	content.WriteByte(comment)
	content.WriteByte(' ')
	content.WriteString("Updated at ")
	content.WriteString(time.Now().Format(time.RFC850))
	content.WriteString("\n")

	for sect, opts := range cfg.config {
		content.WriteString("\n")
		content.WriteByte(namestart)
		content.WriteString(sect)
		content.WriteByte(namestop)
		content.WriteString("\n")
		for opt, val := range opts {
			content.WriteString(opt)
			content.WriteByte(sep)
			content.WriteString(val)
			content.WriteString("\n")
		}
	}

	err := ioutil.WriteFile(cfg.filename, content.Bytes(), 0644)
	if err != nil {
		log.Println("Error saving the config to the file:", err)
	}

	return err
}

func setFilename(filename string) {
	cfg.Lock()
	cfg.filename = filename
	cfg.Unlock()
}
