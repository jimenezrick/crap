package store

// XXX: create(O_EXCLUSIVE) the file that is going to be renamed, as a lock for the destination

// XXX: Encrypt AES
// XXX: Conn SSL
// XXX: Retrieve(key string) *Blob

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"hash"
	"io/ioutil"
	"os"
	"path"
	"runtime"
)

type Blob struct {
	store  Store
	file   *os.File
	writer *bufio.Writer
	hash   hash.Hash
}

func (s Store) NewBlob() (*Blob, error) {
	file, err := ioutil.TempFile(s.tempPath(), "blob")
	if err != nil {
		return nil, err
	}

	b := Blob{s, file, bufio.NewWriter(file), sha1.New()}
	runtime.SetFinalizer(&b, func(b *Blob) {
		os.Remove(b.file.Name())
	})

	return &b, nil
}

func (b Blob) Write(buf []byte) (int, error) {
	b.hash.Write(buf)
	return b.writer.Write(buf)
}

func (b Blob) Store() (string, error) {
	defer b.file.Close()
	b.writer.Flush()

	if err := b.file.Sync(); err != nil {
		return "", err
	}

	src := b.file.Name()
	dest := b.Path()

	if err := os.MkdirAll(path.Dir(dest), b.store.perm); err != nil {
		return "", err
	}
	if err := os.Rename(src, dest); err != nil {
		return "", err
	}
	if err := syncDir(path.Dir(dest)); err != nil {
		return "", err
	}

	runtime.SetFinalizer(b, nil)
	return b.Key(), nil
}

func (b Blob) Abort() error {
	if err := b.file.Close(); err != nil {
		return err
	}
	if err := os.Remove(b.file.Name()); err != nil {
		return err
	}

	runtime.SetFinalizer(b, nil)
	return nil
}

func (b Blob) Size() (int64, error) {
	info, err := b.file.Stat()
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

func (b Blob) Key() string {
	return fmt.Sprintf("%x", b.hash.Sum(nil))
}

func (b Blob) Path() string {
	hash := b.Key()
	return path.Join(b.store.blobPath(), hash[:2], hash[2:])
}
