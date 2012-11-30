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

func (sw *SHA1FileWriter) Write(b []byte) (int, error) {
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

func (sr *SHA1FileReader) Read(b []byte) (int, error) {
	n, err := sr.Reader.Read(b)
	sr.Hash.Write(b[:n])
	return n, err
}
