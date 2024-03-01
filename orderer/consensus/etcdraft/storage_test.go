/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package etcdraft

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.etcd.io/etcd/pkg/fileutil"
	"go.etcd.io/etcd/raft"
	"go.etcd.io/etcd/raft/raftpb"
	"go.etcd.io/etcd/wal"
	"go.uber.org/zap"
)

var (
	err    error
	logger *flogging.FabricLogger

	dataDir, walDir, snapDir string

	ram   *raft.MemoryStorage
	store *RaftStorage
)

func setup(t *testing.T) {
	logger = flogging.NewFabricLogger(zap.NewExample())
	ram = raft.NewMemoryStorage()
	dataDir, err = ioutil.TempDir("", "etcdraft-")
	assert.NoError(t, err)
	walDir, snapDir = path.Join(dataDir, "wal"), path.Join(dataDir, "snapshot")
	store, err = CreateStorage(logger, walDir, snapDir, ram)
	assert.NoError(t, err)
}

func clean(t *testing.T) {
	err = store.Close()
	assert.NoError(t, err)
	err = os.RemoveAll(dataDir)
	assert.NoError(t, err)
}

func fileCount(files []string, suffix string) (c int) {
	for _, f := range files {
		if strings.HasSuffix(f, suffix) {
			c++
		}
	}
	return
}

func assertFileCount(t *testing.T, wal, snap int) {
	files, err := fileutil.ReadDir(walDir)
	assert.NoError(t, err)
	assert.Equal(t, wal, fileCount(files, ".wal"), "WAL file count mismatch")

	files, err = fileutil.ReadDir(snapDir)
	assert.NoError(t, err)
	assert.Equal(t, snap, fileCount(files, ".snap"), "Snap file count mismatch")
}

func TestOpenWAL(t *testing.T) {
	t.Run("Last WAL file is broken", func(t *testing.T) {
		setup(t)
		defer clean(t)

		// create 10 new wal files
		for i := 0; i < 10; i++ {
			store.Store(
				[]raftpb.Entry{{Index: uint64(i), Data: make([]byte, 10)}},
				raftpb.HardState{Commit: uint64(i)},
				raftpb.Snapshot{},
			)
		}
		assertFileCount(t, 1, 0)
		lasti, _ := store.ram.LastIndex() // it never returns err

		// close current storage
		err = store.Close()
		assert.NoError(t, err)

		// truncate wal file
		w := func() string {
			files, err := fileutil.ReadDir(walDir)
			assert.NoError(t, err)
			for _, f := range files {
				if strings.HasSuffix(f, ".wal") {
					return path.Join(walDir, f)
				}
			}
			t.FailNow()
			return ""
		}()
		err = os.Truncate(w, 200)
		assert.NoError(t, err)

		// create new storage
		ram = raft.NewMemoryStorage()
		store, err = CreateStorage(logger, walDir, snapDir, ram)
		require.NoError(t, err)
		lastI, _ := store.ram.LastIndex()
		assert.True(t, lastI > 0)     // we are still able to read some entries
		assert.True(t, lasti > lastI) // but less than before because some are broken
	})
}

func TestTakeSnapshot(t *testing.T) {
	// To make this test more understandable, here's a list
	// of expected wal files:
	// (wal file name format: seq-index.wal, index is the first index in this file)
	//
	// 0000000000000000-0000000000000000.wal (this is created initially by etcd/wal)
	// 0000000000000001-0000000000000001.wal
	// 0000000000000002-0000000000000002.wal
	// 0000000000000003-0000000000000003.wal
	// 0000000000000004-0000000000000004.wal
	// 0000000000000005-0000000000000005.wal
	// 0000000000000006-0000000000000006.wal
	// 0000000000000007-0000000000000007.wal
	// 0000000000000008-0000000000000008.wal
	// 0000000000000009-0000000000000009.wal
	// 000000000000000a-000000000000000a.wal

	t.Run("Good", func(t *testing.T) {
		t.Run("MaxSnapshotFiles==1", func(t *testing.T) {
			backup := MaxSnapshotFiles
			MaxSnapshotFiles = 1
			defer func() { MaxSnapshotFiles = backup }()

			setup(t)
			defer clean(t)

			// set SegmentSizeBytes to a small value so that
			// every entry persisted to wal would result in
			// a new wal being created.
			oldSegmentSizeBytes := wal.SegmentSizeBytes
			wal.SegmentSizeBytes = 10
			defer func() {
				wal.SegmentSizeBytes = oldSegmentSizeBytes
			}()

			// create 10 new wal files
			for i := 0; i < 10; i++ {
				store.Store(
					[]raftpb.Entry{{Index: uint64(i), Data: make([]byte, 100)}},
					raftpb.HardState{Commit: uint64(i)},
					raftpb.Snapshot{},
				)
			}

			assertFileCount(t, 11, 0)

			err = store.TakeSnapshot(uint64(3), raftpb.ConfState{Nodes: []uint64{1}}, make([]byte, 10))
			assert.NoError(t, err)
			// Snapshot is taken at index 3, which releases lock up to 2 (excl.).
			// This results in wal files with index [0, 1] being purged (2 files)
			assertFileCount(t, 9, 1)

			err = store.TakeSnapshot(uint64(5), raftpb.ConfState{Nodes: []uint64{1}}, make([]byte, 10))
			assert.NoError(t, err)
			// Snapshot is taken at index 5, which releases lock up to 4 (excl.).
			// This results in wal files with index [2, 3] being purged (2 files)
			assertFileCount(t, 7, 1)

			t.Logf("Close the storage and create a new one based on existing files")

			err = store.Close()
			assert.NoError(t, err)
			ram := raft.NewMemoryStorage()
			store, err = CreateStorage(logger, walDir, snapDir, ram)
			assert.NoError(t, err)

			err = store.TakeSnapshot(uint64(7), raftpb.ConfState{Nodes: []uint64{1}}, make([]byte, 10))
			assert.NoError(t, err)
			// Snapshot is taken at index 7, which releases lock up to 6 (excl.).
			// This results in wal files with index [4, 5] being purged (2 file)
			assertFileCount(t, 5, 1)

			err = store.TakeSnapshot(uint64(9), raftpb.ConfState{Nodes: []uint64{1}}, make([]byte, 10))
			assert.NoError(t, err)
			// Snapshot is taken at index 9, which releases lock up to 8 (excl.).
			// This results in wal files with index [6, 7] being purged (2 file)
			assertFileCount(t, 3, 1)
		})

		t.Run("MaxSnapshotFiles==2", func(t *testing.T) {
			backup := MaxSnapshotFiles
			MaxSnapshotFiles = 2
			defer func() { MaxSnapshotFiles = backup }()

			setup(t)
			defer clean(t)

			// set SegmentSizeBytes to a small value so that
			// every entry persisted to wal would result in
			// a new wal being created.
			oldSegmentSizeBytes := wal.SegmentSizeBytes
			wal.SegmentSizeBytes = 10
			defer func() {
				wal.SegmentSizeBytes = oldSegmentSizeBytes
			}()

			// create 10 new wal files
			for i := 0; i < 10; i++ {
				store.Store(
					[]raftpb.Entry{{Index: uint64(i), Data: make([]byte, 100)}},
					raftpb.HardState{Commit: uint64(i)},
					raftpb.Snapshot{},
				)
			}

			assertFileCount(t, 11, 0)

			// Only one snapshot is taken, no wal pruning happened
			err = store.TakeSnapshot(uint64(3), raftpb.ConfState{Nodes: []uint64{1}}, make([]byte, 10))
			assert.NoError(t, err)
			assertFileCount(t, 11, 1)

			// Two snapshots at index 3, 5. And we keep one extra wal file prior to oldest snapshot.
			// So we should have pruned wal file with index [0, 1]
			err = store.TakeSnapshot(uint64(5), raftpb.ConfState{Nodes: []uint64{1}}, make([]byte, 10))
			assert.NoError(t, err)
			assertFileCount(t, 9, 2)

			t.Logf("Close the storage and create a new one based on existing files")

			err = store.Close()
			assert.NoError(t, err)
			ram := raft.NewMemoryStorage()
			store, err = CreateStorage(logger, walDir, snapDir, ram)
			assert.NoError(t, err)

			// Two snapshots at index 5, 7. And we keep one extra wal file prior to oldest snapshot.
			// So we should have pruned wal file with index [2, 3]
			err = store.TakeSnapshot(uint64(7), raftpb.ConfState{Nodes: []uint64{1}}, make([]byte, 10))
			assert.NoError(t, err)
			assertFileCount(t, 7, 2)

			// Two snapshots at index 7, 9. And we keep one extra wal file prior to oldest snapshot.
			// So we should have pruned wal file with index [4, 5]
			err = store.TakeSnapshot(uint64(9), raftpb.ConfState{Nodes: []uint64{1}}, make([]byte, 10))
			assert.NoError(t, err)
			assertFileCount(t, 5, 2)
		})
	})

	t.Run("Bad", func(t *testing.T) {

		t.Run("MaxSnapshotFiles==2", func(t *testing.T) {
			// If latest snapshot file is corrupted, storage should be able
			// to recover from an older one.

			backup := MaxSnapshotFiles
			MaxSnapshotFiles = 2
			defer func() { MaxSnapshotFiles = backup }()

			setup(t)
			defer clean(t)

			// set SegmentSizeBytes to a small value so that
			// every entry persisted to wal would result in
			// a new wal being created.
			oldSegmentSizeBytes := wal.SegmentSizeBytes
			wal.SegmentSizeBytes = 10
			defer func() {
				wal.SegmentSizeBytes = oldSegmentSizeBytes
			}()

			// create 10 new wal files
			for i := 0; i < 10; i++ {
				store.Store(
					[]raftpb.Entry{{Index: uint64(i), Data: make([]byte, 100)}},
					raftpb.HardState{Commit: uint64(i)},
					raftpb.Snapshot{},
				)
			}

			assertFileCount(t, 11, 0)

			// Only one snapshot is taken, no wal pruning happened
			err = store.TakeSnapshot(uint64(3), raftpb.ConfState{Nodes: []uint64{1}}, make([]byte, 10))
			assert.NoError(t, err)
			assertFileCount(t, 11, 1)

			// Two snapshots at index 3, 5. And we keep one extra wal file prior to oldest snapshot.
			// So we should have pruned wal file with index [0, 1]
			err = store.TakeSnapshot(uint64(5), raftpb.ConfState{Nodes: []uint64{1}}, make([]byte, 10))
			assert.NoError(t, err)
			assertFileCount(t, 9, 2)

			d, err := os.Open(snapDir)
			assert.NoError(t, err)
			defer d.Close()
			names, err := d.Readdirnames(-1)
			assert.NoError(t, err)
			sort.Sort(sort.Reverse(sort.StringSlice(names)))

			corrupted := filepath.Join(snapDir, names[0])
			t.Logf("Corrupt latest snapshot file: %s", corrupted)
			f, err := os.OpenFile(corrupted, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
			assert.NoError(t, err)
			_, err = f.WriteString("Corrupted Snapshot")
			assert.NoError(t, err)
			f.Close()

			// Corrupted snapshot file should've been renamed by ListSnapshots
			_ = ListSnapshots(logger, snapDir)
			assertFileCount(t, 9, 1)

			// Rollback the rename
			broken := corrupted + ".broken"
			err = os.Rename(broken, corrupted)
			assert.NoError(t, err)
			assertFileCount(t, 9, 2)

			t.Logf("Close the storage and create a new one based on existing files")

			err = store.Close()
			assert.NoError(t, err)
			ram := raft.NewMemoryStorage()
			store, err = CreateStorage(logger, walDir, snapDir, ram)
			assert.NoError(t, err)

			// Corrupted snapshot file should've been renamed by CreateStorage
			assertFileCount(t, 9, 1)

			files, err := fileutil.ReadDir(snapDir)
			assert.NoError(t, err)
			assert.Equal(t, 1, fileCount(files, ".broken"))

			err = store.TakeSnapshot(uint64(7), raftpb.ConfState{Nodes: []uint64{1}}, make([]byte, 10))
			assert.NoError(t, err)
			assertFileCount(t, 9, 2)

			err = store.TakeSnapshot(uint64(9), raftpb.ConfState{Nodes: []uint64{1}}, make([]byte, 10))
			assert.NoError(t, err)
			assertFileCount(t, 5, 2)
		})
	})
}

func TestApplyOutOfDateSnapshot(t *testing.T) {
	t.Run("Apply out of date snapshot", func(t *testing.T) {
		setup(t)
		defer clean(t)

		// set SegmentSizeBytes to a small value so that
		// every entry persisted to wal would result in
		// a new wal being created.
		oldSegmentSizeBytes := wal.SegmentSizeBytes
		wal.SegmentSizeBytes = 10
		defer func() {
			wal.SegmentSizeBytes = oldSegmentSizeBytes
		}()

		// create 10 new wal files
		for i := 0; i < 10; i++ {
			store.Store(
				[]raftpb.Entry{{Index: uint64(i), Data: make([]byte, 100)}},
				raftpb.HardState{Commit: uint64(i)},
				raftpb.Snapshot{},
			)
		}
		assertFileCount(t, 11, 0)

		err = store.TakeSnapshot(uint64(3), raftpb.ConfState{Nodes: []uint64{1}}, make([]byte, 10))
		assert.NoError(t, err)
		assertFileCount(t, 11, 1)

		snapshot := store.Snapshot()
		assert.NotNil(t, snapshot)

		// Applying old snapshot should have no effect
		store.ApplySnapshot(snapshot)

		// Storing old snapshot gets no error
		err := store.Store(
			[]raftpb.Entry{{Index: uint64(10), Data: make([]byte, 100)}},
			raftpb.HardState{},
			snapshot,
		)
		assert.NoError(t, err)
		assertFileCount(t, 12, 1)
	})
}

func TestAbortWhenWritingSnapshot(t *testing.T) {
	t.Run("Abort when writing snapshot", func(t *testing.T) {
		setup(t)
		defer clean(t)

		// set SegmentSizeBytes to a small value so that
		// every entry persisted to wal would result in
		// a new wal being created.
		oldSegmentSizeBytes := wal.SegmentSizeBytes
		wal.SegmentSizeBytes = 10
		defer func() {
			wal.SegmentSizeBytes = oldSegmentSizeBytes
		}()

		// create 5 new entry
		for i := 0; i < 5; i++ {
			store.Store(
				[]raftpb.Entry{{Index: uint64(i), Data: make([]byte, 100)}},
				raftpb.HardState{Commit: uint64(i)},
				raftpb.Snapshot{},
			)
		}
		assertFileCount(t, 6, 0)

		// Assume an orderer missed some records due to exceptions and receives a new snapshot from other orderers.
		commit := 10
		store.Store(
			[]raftpb.Entry{},
			raftpb.HardState{Commit: uint64(commit)},
			raftpb.Snapshot{
				Metadata: raftpb.SnapshotMetadata{
					Index: uint64(commit),
				},
				Data: make([]byte, 100),
			},
		)
		err = store.Close()
		assert.NoError(t, err)

		// In old logic, it will use rs.wal.Save(hardstate, entries) to save the state firstly, so we remove the snapshot files.
		// sd, err := os.Open(snapDir)
		// assert.NoError(t, err)
		// defer sd.Close()
		// names, err := sd.Readdirnames(-1)
		// assert.NoError(t, err)
		// sort.Sort(sort.Reverse(sort.StringSlice(names)))
		// os.Remove(filepath.Join(snapDir, names[0]))
		// wd, err := os.Open(walDir)
		// assert.NoError(t, err)
		// defer wd.Close()
		// names, err = wd.Readdirnames(-1)
		// assert.NoError(t, err)
		// sort.Sort(sort.Reverse(sort.StringSlice(names)))
		// os.Remove(filepath.Join(walDir, names[0]))

		// But in the new logic, it will use rs.saveSnap(snapshot) to save the snapshot firstly, so we remove the WAL files.
		wd, err := os.Open(walDir)
		assert.NoError(t, err)
		defer wd.Close()
		names, err := wd.Readdirnames(-1)
		assert.NoError(t, err)
		sort.Sort(sort.Reverse(sort.StringSlice(names)))
		os.Remove(filepath.Join(walDir, names[0]))

		// Then restart the orderer.
		ram := raft.NewMemoryStorage()
		store, err = CreateStorage(logger, walDir, snapDir, ram)
		assert.NoError(t, err)

		// Check the state from go.etcd.io/etcd/raft/raft.go
		// func (r *raft) loadState(state pb.HardState)
		hd, _, err := store.ram.InitialState()
		assert.NoError(t, err)
		lastIndex, err := store.ram.LastIndex()
		assert.NoError(t, err)
		assert.False(t, hd.Commit > lastIndex)
	})
}
