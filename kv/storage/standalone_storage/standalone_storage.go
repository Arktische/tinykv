package standalone_storage

import (
	"github.com/Connor1996/badger"
	"github.com/pingcap-incubator/tinykv/kv/config"
	"github.com/pingcap-incubator/tinykv/kv/storage"
	"github.com/pingcap-incubator/tinykv/kv/util/engine_util"
	"github.com/pingcap-incubator/tinykv/proto/pkg/kvrpcpb"
)

// StandAloneStorage is an implementation of `Storage` for a single-node TinyKV instance. It does not
// communicate with other nodes and all data is stored locally.
type StandAloneStorage struct {
	// Your Data Here (1).
	// storage engine
	engine *engine_util.Engines
}

// NewStandAloneStorage returns a standalone storage instance
func NewStandAloneStorage(conf *config.Config) *StandAloneStorage {
	// Your Code Here (1).
	return &StandAloneStorage{
		engine: engine_util.NewEngines(engine_util.CreateDB(conf.DBPath, false),
			nil, conf.DBPath, ""),
	}
}

// Start starts the standalone storage
func (s *StandAloneStorage) Start() error {
	// Your Code Here (1).

	return nil
}

// Stop stops the standalone storage
func (s *StandAloneStorage) Stop() error {
	// Your Code Here (1).
	return s.engine.Kv.Close()
}

// Reader returns reader for standalone storage
func (s *StandAloneStorage) Reader(ctx *kvrpcpb.Context) (storage.StorageReader, error) {
	// Your Code Here (1).
	txn := s.engine.Kv.NewTransaction(false)
	reader := NewStandAloneStorageReader(txn)
	return reader, nil
}

// Write writes batch to standalone storage with transaction
func (s *StandAloneStorage) Write(ctx *kvrpcpb.Context, batch []storage.Modify) error {
	// Your Code Here (1).
	kvWB := new(engine_util.WriteBatch)
	for _, m := range batch {
		switch m.Data.(type) {
		case storage.Put:
			put := m.Data.(storage.Put)
			kvWB.SetCF(put.Cf, put.Key, put.Value)
			break
		case storage.Delete:
			del := m.Data.(storage.Delete)
			kvWB.DeleteCF(del.Cf, del.Key)
			break
		}
	}
	err := kvWB.WriteToDB(s.engine.Kv)
	if err != nil {
		return err
	}
	return nil
}

// StandAloneStorageReader reader
type StandAloneStorageReader struct {
	txn *badger.Txn
}

// NewStandAloneStorageReader returns a new reader
func NewStandAloneStorageReader(txn *badger.Txn) *StandAloneStorageReader {
	return &StandAloneStorageReader{
		txn: txn,
	}
}

// GetCF returns column family
func (s *StandAloneStorageReader) GetCF(cf string, key []byte) ([]byte, error) {
	val, err := engine_util.GetCFFromTxn(s.txn, cf, key)
	if err == badger.ErrKeyNotFound {
		return nil, nil
	}
	return val, err
}

// IterCF returns DBIterator for given column family
func (s *StandAloneStorageReader) IterCF(cf string) engine_util.DBIterator {
	return engine_util.NewCFIterator(cf, s.txn)
}

// Close closes the reader and discards the internal transaction
func (s *StandAloneStorageReader) Close() {
	s.txn.Discard()
}
