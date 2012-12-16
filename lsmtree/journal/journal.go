package journal

//
// XXX: Usar el termino journal en vez de log
//

import (
	"os"
	"hash/fnv"
)

// XXX: Transacciones, insertar en varios indices atomicamente (un log unico)
//      Usar un unico log con un identificador de indice en cada record
//      En una insercion multi-indice, crear un solo record con todas las inserciones

/* XXX
   Insert tombstone (delete) in lsmtree
   Copy blob to blobs/
   Insert blob key in lsmtree

   When a log replay happens, look for with the Iterator()
   tombstone without an insert, delete those blobs from blobs/.
*/

// Opciones para recibir por red:
//
// XXX: {embed-index: true} -> blob dentro del index
//      {async: true} -> no fsync ni en blob ni en index

//
// XXX: Sync() el log cada 5 segundo desde la gorutina?
//      Las escrituras en el log no hacen Sync() directamente
//      O usar O_DIRECT sin fsync ni nada?
//

// XXX XXX XXX
// -- index/
//     |-- foo.log (current log)                         (in memtable)
//     |-- foo.log-1 ... foo.log-N (to be mered with L0) (in memtable)
//     |-- foo.sst-0_1 ... foo.sst-0_4 2MB (young level) (in memtable)
//     |-- foo.sst-1 8MB
//     +-- foo.sst-2 64MB (growth rate: 8^i)
// XXX XXX XXX

const (
	magicNum = 0x7d3afedc4cb9752d // XXX: No se pone cuando un record pisa el final de un segmento
	blockSize = 32 * 1024
	maxBlockTailSize = 256
)

const logSize = 4 * 1024 * 1024 // XXX: Config in MB

type journal struct {
	*os.File
	offset int64
}

func Open(config config.Config) (*journal, error) {
	// TODO: Sync dir
	// TODO: fallocate, fadvise

	journal = new(journal)

	file, err := os.Open("/tmp/journal")
	if err != nil {
		return nil, err
	}
	journal.File = file

	offset, err := l.Seek(0, os.SEEK_END)
	if err != nil {
		return nil, err
	}
	journal.offset = offset

	return journal, nil
}

func (l *journal) Append(r record) error {
	bytes := r.Encode()

	if offset % blockSize == 0 {
		l.AppendMagicNum()
	} else if l.BlockFree() <= maxBlockTailSize && len(bytes) > l.BlockFree() {
		if err := l.PadBlock(); err != nil {
			return err
		}
		return l.Append()
	}

	n, err := l.Write(bytes)
	if err != nil {
		return err
	}

	l.offset += n
	return nil
}

func (l *journal) BlockFree() int64 {
	return (l.offset / blockSize + 1) * blockSize - offset
}

func (l *journal) PadBlock() error {
	offset, err := l.Seek(l.BlockFree(), os.SEEK_CUR)
	if err != nil {
		return err
	}
	l.offset = offset
}

func (l *journal) Sync() error {
	return util.Fdatasync(l.File)
}

//func Iterator() {
//}

func (l *journal) Close() error {
	defer l.Close()
	return util.Fdatasync(l.File)
}




func (l *journal) AppendMagicNum() error {
	l.offset += 8
	return binary.Write(l, binary.BigEndian, magicNum)
}
















type record struct {
	tombstone, compress bool
	key []byte
	value []byte
}






func newTombstone() *record {
	return &record{tombstone: true}
}


func newRecord(key, value []byte) *record {
	r := &record{false, false, info, key, value}
}

func newCompressedRecord(key, value []byte) *record {
	r := &record{false, true, info, key, value}
}


func (r *record) Encode() []byte {
	bytesCrc := make([]byte, 4 + binary.MaxVarintLen64 * 3 + 1 + len(key) + len(value))
	bytes := bytesCrc[4 + binary.MaxVarintLen64:]

	n := binary.PutUvarint(bytes[1:] // XXX

	keyLen := make([]byte, binary.MaxVarintLen64)
	n = binary.PutUvarint(keyLen, len(key))
	keyLen = keyLen[:n]

	valueLen := make([]byte, binary.MaxVarintLen64)
	n = binary.PutUvarint(valueLen, len(value))
	valueLen = valueLen[:n]

	bytes = append(bytes, info)
	bytes = append(bytes, keyLen)
	bytes = append(bytes, key)
	bytes = append(bytes, valueLen)
	bytes = append(bytes, value)

	recordLen := make([]byte, binary.MaxVarintLen64)
	n = binary.PutUvarint(recordLen, len(bytes))
	recordLen = recordLen[:n]

	crc := crcBytes(bytes)

	offset := (4 + binary.MaxVarintLen64) - (len(crc) + len(recordLen))
	bytesCrc = bytesCrc[offset:]

	copy(bytesCrc, crc)
	copy(bytesCrc[4:], recordLen)
}

func crcBytes(bytes []byte) []byte {
	crc := fnv.New32a()
	crc.Write(bytes)

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, crc.Sum32()); err != nil {
		panic(err)
	}
	return buf.Bytes()
}
