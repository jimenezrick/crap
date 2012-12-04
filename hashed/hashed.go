package hashed

import (
	"bufio"
	"crypto/sha1"
	"hash"
	"io"
	"os"
)

type HashedReader interface {
	io.Reader
	hash.Hash
}

type HashedWriter interface {
	hash.Hash
	Flush() error
}

type sha1FileWriter struct {
	*bufio.Writer
	hash.Hash
}

func NewSha1FileWriter(file *os.File) HashedWriter {
	return &sha1FileWriter{bufio.NewWriter(file), sha1.New()}
}

func (sw *sha1FileWriter) Write(b []byte) (int, error) {
	sw.Hash.Write(b)
	return sw.Writer.Write(b)
}

type sha1FileReader struct {
	*bufio.Reader
	hash.Hash
}

func NewSha1FileReader(file *os.File) HashedReader {
	return &sha1FileReader{bufio.NewReader(file), sha1.New()}
}

func (sr *sha1FileReader) Read(b []byte) (int, error) {
	n, err := sr.Reader.Read(b)
	sr.Hash.Write(b[:n])
	return n, err
}
