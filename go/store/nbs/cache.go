// Copyright 2016 Attic Labs, Inc. All rights reserved.
// Licensed under the Apache License, version 2.0:
// http://www.apache.org/licenses/LICENSE-2.0

package nbs

import (
	"context"
	"errors"
	"io/ioutil"
	"os"

	"github.com/liquidata-inc/ld/dolt/go/store/types"
	"github.com/liquidata-inc/ld/dolt/go/store/chunks"
	"github.com/liquidata-inc/ld/dolt/go/store/hash"
)

const (
	defaultCacheMemTableSize uint64 = 1 << 27 // 128MiB
)

func NewCache(ctx context.Context) (*NomsBlockCache, error) {
	dir, err := ioutil.TempDir("", "")

	if err != nil {
		return nil, err
	}

	store, err := NewLocalStore(ctx, types.Format_Default.VersionString(), dir, defaultCacheMemTableSize)

	if err != nil {
		return nil, err
	}

	return &NomsBlockCache{store, dir}, nil
}

// NomsBlockCache holds Chunks, allowing them to be retrieved by hash or enumerated in hash order.
type NomsBlockCache struct {
	chunks *NomsBlockStore
	dbDir  string
}

// Insert stores c in the cache.
func (nbc *NomsBlockCache) Insert(ctx context.Context, c chunks.Chunk) error {
	success := nbc.chunks.addChunk(ctx, addr(c.Hash()), c.Data())

	if !success {
		return errors.New("failed to add chunk")
	}

	return nil
}

// Has checks if the chunk referenced by hash is in the cache.
func (nbc *NomsBlockCache) Has(ctx context.Context, hash hash.Hash) (bool, error) {
	return nbc.chunks.Has(ctx, hash)
}

// HasMany returns a set containing the members of hashes present in the
// cache.
func (nbc *NomsBlockCache) HasMany(ctx context.Context, hashes hash.HashSet) (hash.HashSet, error) {
	return nbc.chunks.HasMany(ctx, hashes)
}

// Get retrieves the chunk referenced by hash. If the chunk is not present,
// Get returns the empty Chunk.
func (nbc *NomsBlockCache) Get(ctx context.Context, hash hash.Hash) (chunks.Chunk, error) {
	return nbc.chunks.Get(ctx, hash)
}

// GetMany gets the Chunks with |hashes| from the store. On return,
// |foundChunks| will have been fully sent all chunks which have been
// found. Any non-present chunks will silently be ignored.
func (nbc *NomsBlockCache) GetMany(ctx context.Context, hashes hash.HashSet, foundChunks chan *chunks.Chunk) error {
	return nbc.chunks.GetMany(ctx, hashes, foundChunks)
}

// Count returns the number of items in the cache.
func (nbc *NomsBlockCache) Count() (uint32, error) {
	return nbc.chunks.Count()
}

// Destroy drops the cache and deletes any backing storage.
func (nbc *NomsBlockCache) Destroy() error {
	chunkErr := nbc.chunks.Close()
	remErr := os.RemoveAll(nbc.dbDir)

	if chunkErr != nil {
		return chunkErr
	}

	return remErr
}
