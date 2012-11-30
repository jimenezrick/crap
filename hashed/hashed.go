package hashed

import (
	"os"
	"io"
	"bufio"
	"hash"
	"crypto/sha1"
)

type HashedReader interface {
	io.Reader
	hash.Hash
}

type HashedWriter interface {
	hash.Hash
}

type SHA1FileWriter struct {
	*bufio.Writer
	hash.Hash
}

func NewSHA1FileWriter(file *os.File) *SHA1FileWriter {
	return &SHA1FileWriter{bufio.NewWriter(file), sha1.New()}
}

func (sw *SHA1FileWriter) Write(b []byte) (n int, err error) {
	sw.Hash.Write(b)
	return sw.Writer.Write(b)
}

type SHA1FileReader struct {
	*bufio.Reader
	hash.Hash
}

func NewSHA1FileReader(file *os.File) *SHA1FileReader {
	return &SHA1FileReader{bufio.NewReader(file), sha1.New()}
}

func (sr *SHA1FileReader) Read(b []byte) (n int, err error) {
	sr.Hash.Write(b)
	return sr.Reader.Read(b)
}
