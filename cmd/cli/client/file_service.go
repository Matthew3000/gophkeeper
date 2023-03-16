package client

import (
	"os"
	"sync"
)

type SafeFile struct {
	file  *os.File
	mutex *sync.Mutex
}

func NewSafeFile(filename string) (*SafeFile, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return &SafeFile{
		file:  file,
		mutex: &sync.Mutex{},
	}, nil
}

func (sf *SafeFile) Read(p []byte) (int, error) {
	sf.mutex.Lock()
	defer sf.mutex.Unlock()
	return sf.file.Read(p)
}

func (sf *SafeFile) Write(p []byte) (int, error) {
	sf.mutex.Lock()
	defer sf.mutex.Unlock()
	return sf.file.Write(p)
}

func (sf *SafeFile) Close() error {
	sf.mutex.Lock()
	defer sf.mutex.Unlock()
	return sf.file.Close()
}
