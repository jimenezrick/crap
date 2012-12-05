package store

//
// XXX: Copiar metadatos en JSON en blobs/xx/xxxxxxxxxxxxxxxxx.meta
//

// XXX: Encrypt AES
// XXX: Conn SSL
// XXX: Retrieve(key string) *Blob

import (
	"io/ioutil"
	"os"
	"path"
	"runtime"
)

import (
	"crap/hashed"
	"crap/util"
)

type Blob struct {
	store Store
	file  *os.File
	hashed.HashedWriter
}

func (s Store) NewBlobSize(size uint64) (*Blob, error) {
	b, err := s.NewBlob()
	if err != nil {
		return nil, err
	}

	if err := util.Fallocate(b.file, size); err != nil {
		return nil, err
	}
	return b, nil
}

func (s Store) NewBlob() (*Blob, error) {
	file, err := ioutil.TempFile(s.tempPath(), "blob")
	if err != nil {
		return nil, err
	}

	b := Blob{s, file, hashed.NewSha1FileWriter(file)}
	runtime.SetFinalizer(&b, func(b *Blob) {
		os.Remove(b.file.Name())
	})

	return &b, nil
}

func (b *Blob) Store(sync bool) ([]byte, error) {
	defer b.file.Close()
	b.HashedWriter.Flush()

	if err := b.file.Chmod(b.store.filePerm); err != nil {
		return nil, err
	}

	if sync {
		if err := util.Fdatasync(b.file); err != nil {
			return nil, err
		}
	}

	src := b.file.Name()
	dest := b.path()

	if err := os.MkdirAll(path.Dir(dest), b.store.dirPerm); err != nil {
		return nil, err
	}
	if err := b.lock(); err != nil {
		return b.Sum(nil), err
	}
	if err := os.Rename(src, dest); err != nil {
		return nil, err
	}

	if sync {
		if err := util.Datasync(path.Dir(dest)); err != nil {
			return nil, err
		}
	}

	runtime.SetFinalizer(b, nil)
	return b.Sum(nil), nil
}

func (b *Blob) Abort() error {
	err1 := b.file.Close()
	err2 := os.Remove(b.file.Name())
	runtime.SetFinalizer(b, nil)

	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}

	return nil
}

func (b *Blob) lock() error {
	return util.CreateLockFile(b.path(), b.store.filePerm)
}

func (b *Blob) path() string {
	hash := util.HexHash(b)
	return path.Join(b.store.blobPath(), hash[:2], hash[2:])
}
