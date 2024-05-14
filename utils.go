package main

import (
	"log"
	"os"
)

func LoadFile(path string) ([]byte, error) {
	file, err := os.ReadFile("static/" + path)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return file, nil
}

func ReadDir(path string) ([]os.DirEntry, error) {
	files, err := os.ReadDir("static/" + path)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return files, nil
}
