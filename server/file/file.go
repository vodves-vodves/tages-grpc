package file

import (
	"bytes"
	"io/ioutil"
)

type File struct {
	name   string
	buffer *bytes.Buffer
}

func (f *File) GetName() string {
	return f.name
}

func (f *File) EditName(name string) {
	f.name = name
}

func NewFile() *File {
	return &File{
		name:   "",
		buffer: &bytes.Buffer{},
	}
}

func (f *File) Write(chunk []byte) error {
	_, err := f.buffer.Write(chunk)

	return err
}

func Store(file *File) error {
	if err := ioutil.WriteFile("server/files/"+file.name, file.buffer.Bytes(), 0644); err != nil {
		return err
	}

	return nil
}
