package boltbrowser

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/boltdb/bolt"
)

/*
boltDB A Database, basically a collection of buckets
*/
type boltDB struct {
	db       *bolt.DB
	readOnly bool
	buckets  []boltBucket
}

/*
boltBucket is just a struct representation of a Bucket in the Bolt DB
*/
type boltBucket struct {
	name      string
	pairs     []boltPair
	buckets   []boltBucket
	parent    *boltBucket
	expanded  bool
	errorFlag bool
}

/*
boltPair is just a struct representation of a Pair in the Bolt DB
*/
type boltPair struct {
	parent *boltBucket
	key    string
	val    string
}

func (bd *boltDB) getGenericFromPath(path []string) (*boltBucket, *boltPair, error) {
	// Check if 'path' leads to a pair
	p, err := bd.getPairFromPath(path)
	if err == nil {
		return nil, p, nil
	}
	// Nope, check if it leads to a bucket
	b, err := bd.getBucketFromPath(path)
	if err == nil {
		return b, nil, nil
	}
	// Nope, error
	return nil, nil, errors.New("Invalid Path")
}

func (bd *boltDB) getBucketFromPath(path []string) (*boltBucket, error) {
	if len(path) > 0 {
		// Find the BoltBucket with a path == path
		var b *boltBucket
		var err error
		// Find the root bucket
		b, err = bd.getBucket(path[0])
		if err != nil {
			return nil, err
		}
		if len(path) > 1 {
			for p := 1; p < len(path); p++ {
				b, err = b.getBucket(path[p])
				if err != nil {
					return nil, err
				}
			}
		}
		return b, nil
	}
	return nil, errors.New("Invalid Path")
}

func (bd *boltDB) getPairFromPath(path []string) (*boltPair, error) {
	if len(path) <= 0 {
		return nil, errors.New("No Path")
	}
	b, err := bd.getBucketFromPath(path[:len(path)-1])
	if err != nil {
		return nil, err
	}
	// Found the bucket, pull out the pair
	p, err := b.getPair(path[len(path)-1])
	return p, err
}

func (bd *boltDB) getVisibleItemCount(path []string) (int, error) {
	vis := 0
	var retErr error
	if len(path) == 0 {
		for i := range bd.buckets {
			n, err := bd.getVisibleItemCount(bd.buckets[i].GetPath())
			if err != nil {
				return 0, err
			}
			vis += n
		}
	} else {
		b, err := bd.getBucketFromPath(path)
		if err != nil {
			return 0, err
		}
		// 1 for the bucket
		vis++
		if b.expanded {
			// This bucket is expanded, add up it's children
			// * 1 for each pair
			vis += len(b.pairs)
			// * recurse for buckets
			for i := range b.buckets {
				n, err := bd.getVisibleItemCount(b.buckets[i].GetPath())
				if err != nil {
					return 0, err
				}
				vis += n
			}
		}
	}
	return vis, retErr
}

func (bd *boltDB) buildVisiblePathSlice() ([][]string, error) {
	var retSlice [][]string
	var retErr error
	// The root path, recurse for root buckets
	for i := range bd.buckets {
		bktS, bktErr := bd.buckets[i].buildVisiblePathSlice([]string{})
		if bktErr == nil {
			retSlice = append(retSlice, bktS...)
		} else {
			// Something went wrong, set the error flag
			bd.buckets[i].errorFlag = true
		}
	}
	return retSlice, retErr
}

func (bd *boltDB) getPrevVisiblePath(path []string) []string {
	visPaths, err := bd.buildVisiblePathSlice()
	if path == nil {
		if len(visPaths) > 0 {
			return visPaths[len(visPaths)-1]
		}
		return nil
	}
	if err == nil {
		for idx, pth := range visPaths {
			isCurPath := true
			for i := range path {
				if len(pth) <= i || path[i] != pth[i] {
					isCurPath = false
					break
				}
			}
			if isCurPath && idx > 0 {
				return visPaths[idx-1]
			}
		}
	}
	return nil
}
func (bd *boltDB) getNextVisiblePath(path []string) []string {
	visPaths, err := bd.buildVisiblePathSlice()
	if path == nil {
		if len(visPaths) > 0 {
			return visPaths[0]
		}
		return nil
	}
	if err == nil {
		for idx, pth := range visPaths {
			isCurPath := true
			for i := range path {
				if len(pth) <= i || path[i] != pth[i] {
					isCurPath = false
					break
				}
			}
			if isCurPath && len(visPaths) > idx+1 {
				return visPaths[idx+1]
			}
		}
	}
	return nil
}

func (bd *boltDB) toggleOpenBucket(path []string) error {
	// Find the BoltBucket with a path == path
	b, err := bd.getBucketFromPath(path)
	if err == nil {
		b.expanded = !b.expanded
	}
	return err
}

func (bd *boltDB) closeBucket(path []string) error {
	// Find the BoltBucket with a path == path
	b, err := bd.getBucketFromPath(path)
	if err == nil {
		b.expanded = false
	}
	return err
}

func (bd *boltDB) openBucket(path []string) error {
	// Find the BoltBucket with a path == path
	b, err := bd.getBucketFromPath(path)
	if err == nil {
		b.expanded = true
	}
	return err
}

func (bd *boltDB) getBucket(k string) (*boltBucket, error) {
	for i := range bd.buckets {
		if bd.buckets[i].name == k {
			return &bd.buckets[i], nil
		}
	}
	return nil, errors.New("Bucket Not Found")
}

func (bd *boltDB) openAllBuckets() {
	for i := range bd.buckets {
		bd.buckets[i].openAllBuckets()
		bd.buckets[i].expanded = true
	}
}

func (bd *boltDB) syncOpenBuckets(shadow *boltDB) {
	// First test this bucket
	for i := range bd.buckets {
		for j := range shadow.buckets {
			if bd.buckets[i].name == shadow.buckets[j].name {
				bd.buckets[i].syncOpenBuckets(&shadow.buckets[j])
			}
		}
	}
}

func (bd *boltDB) refreshDatabase() {
	newDB := newModel(bd.db, bd.readOnly)
	*bd = *newDB
}

func newModel(db *bolt.DB, readOnly bool) *boltDB {
	memBolt := &boltDB{
		db:       db,
		readOnly: readOnly,
	}
	db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(nm []byte, b *bolt.Bucket) error {
			bb, err := readBucket(b)
			if err == nil {
				bb.name = string(nm)
				bb.expanded = false
				memBolt.buckets = append(memBolt.buckets, *bb)
				return nil
			}
			return err
		})
	})
	return memBolt
}

/*
GetPath returns the database path leading to this BoltBucket
*/
func (b *boltBucket) GetPath() []string {
	if b.parent != nil {
		return append(b.parent.GetPath(), b.name)
	}
	return []string{b.name}
}

/*
buildVisiblePathSlice builds a slice of string slices containing all visible paths in this bucket
The passed prefix is the path leading to the current bucket
*/
func (b *boltBucket) buildVisiblePathSlice(prefix []string) ([][]string, error) {
	var retSlice [][]string
	var retErr error
	retSlice = append(retSlice, append(prefix, b.name))
	if b.expanded {
		// Add subbuckets
		for i := range b.buckets {
			bktS, bktErr := b.buckets[i].buildVisiblePathSlice(append(prefix, b.name))
			if bktErr != nil {
				return retSlice, bktErr
			}
			retSlice = append(retSlice, bktS...)
		}
		// Add pairs
		for i := range b.pairs {
			retSlice = append(retSlice, append(append(prefix, b.name), b.pairs[i].key))
		}
	}
	return retSlice, retErr
}

func (b *boltBucket) syncOpenBuckets(shadow *boltBucket) {
	// First test this bucket
	b.expanded = shadow.expanded
	for i := range b.buckets {
		for j := range shadow.buckets {
			if b.buckets[i].name == shadow.buckets[j].name {
				b.buckets[i].syncOpenBuckets(&shadow.buckets[j])
			}
		}
	}
}

func (b *boltBucket) openAllBuckets() {
	for i := range b.buckets {
		b.buckets[i].openAllBuckets()
		b.buckets[i].expanded = true
	}
}

func (b *boltBucket) getBucket(k string) (*boltBucket, error) {
	for i := range b.buckets {
		if b.buckets[i].name == k {
			return &b.buckets[i], nil
		}
	}
	return nil, errors.New("Bucket Not Found")
}

func (b *boltBucket) getPair(k string) (*boltPair, error) {
	for i := range b.pairs {
		if b.pairs[i].key == k {
			return &b.pairs[i], nil
		}
	}
	return nil, errors.New("Pair Not Found")
}

/*
GetPath Returns the path of the BoltPair
*/
func (p *boltPair) GetPath() []string {
	return append(p.parent.GetPath(), p.key)
}

/* This is a go-between function (between the boltbrowser structs
 * above, and the bolt convenience functions below)
 * for taking a boltbrowser bucket and recursively adding it
 * and all of it's children into the database.
 * Mainly used for moving a bucket from one path to another
 * as in the 'renameBucket' function below.
 */
func (bd *boltDB) addBucketFromBoltBucket(path []string, bb *boltBucket) error {
	if err := bd.insertBucket(path, bb.name); err == nil {
		bucketPath := append(path, bb.name)
		for i := range bb.pairs {
			if err = bd.insertPair(bucketPath, bb.pairs[i].key, bb.pairs[i].val); err != nil {
				return err
			}
		}
		for i := range bb.buckets {
			if err = bd.addBucketFromBoltBucket(bucketPath, &bb.buckets[i]); err != nil {
				return err
			}
		}
	}
	return nil
}

func (bd *boltDB) deleteKey(path []string) error {
	if bd.readOnly {
		return errors.New("DB is in Read-Only Mode")
	}
	err := bd.db.Update(func(tx *bolt.Tx) error {
		// len(b.path)-1 is the key we need to delete,
		// the rest are buckets leading to that key
		if len(path) == 1 {
			// Deleting a root bucket
			return tx.DeleteBucket([]byte(path[0]))
		}
		b := tx.Bucket([]byte(path[0]))
		if b != nil {
			if len(path) > 1 {
				for i := range path[1 : len(path)-1] {
					b = b.Bucket([]byte(path[i+1]))
					if b == nil {
						return errors.New("deleteKey: Invalid Path")
					}
				}
			}
			// Now delete the last key in the path
			var err error
			if deleteBkt := b.Bucket([]byte(path[len(path)-1])); deleteBkt == nil {
				// Must be a pair
				err = b.Delete([]byte(path[len(path)-1]))
			} else {
				err = b.DeleteBucket([]byte(path[len(path)-1]))
			}
			return err
		}
		return errors.New("deleteKey: Invalid Path")
	})
	return err
}

func readBucket(b *bolt.Bucket) (*boltBucket, error) {
	bb := new(boltBucket)
	b.ForEach(func(k, v []byte) error {
		if v == nil {
			tb, err := readBucket(b.Bucket(k))
			tb.parent = bb
			if err == nil {
				tb.name = string(k)
				bb.buckets = append(bb.buckets, *tb)
			}
		} else {
			tp := boltPair{key: string(k), val: string(v)}
			tp.parent = bb
			bb.pairs = append(bb.pairs, tp)
		}
		return nil
	})
	return bb, nil
}

func (bd *boltDB) renameBucket(path []string, name string) error {
	if name == path[len(path)-1] {
		// No change requested
		return nil
	}
	var bb *boltBucket // For caching the current bucket
	err := bd.db.View(func(tx *bolt.Tx) error {
		// len(b.path)-1 is the key we need to delete,
		// the rest are buckets leading to that key
		b := tx.Bucket([]byte(path[0]))
		if b != nil {
			if len(path) > 1 {
				for i := range path[1:] {
					b = b.Bucket([]byte(path[i+1]))
					if b == nil {
						return errors.New("renameBucket: Invalid Path")
					}
				}
			}
			var err error
			// Ok, cache b
			bb, err = readBucket(b)
			if err != nil {
				return err
			}
		} else {
			return errors.New("renameBucket: Invalid Bucket")
		}
		return nil
	})
	if err != nil {
		return err
	}
	if bb == nil {
		return errors.New("renameBucket: Couldn't find Bucket")
	}

	// Ok, we have the bucket cached, now delete the current instance
	if err = bd.deleteKey(path); err != nil {
		return err
	}
	// Rechristen our cached bucket
	bb.name = name
	// And re-add it

	parentPath := path[:len(path)-1]
	if err = bd.addBucketFromBoltBucket(parentPath, bb); err != nil {
		return err
	}
	return nil
}

func (bd *boltDB) updatePairKey(path []string, k string) error {
	if bd.readOnly {
		return errors.New("DB is in Read-Only Mode")
	}
	err := bd.db.Update(func(tx *bolt.Tx) error {
		// len(b.path)-1 is the key for the pair we're updating,
		// the rest are buckets leading to that key
		b := tx.Bucket([]byte(path[0]))
		if b != nil {
			if len(path) > 0 {
				for i := range path[1 : len(path)-1] {
					b = b.Bucket([]byte(path[i+1]))
					if b == nil {
						return errors.New("updatePairValue: Invalid Path")
					}
				}
			}
			bk := []byte(path[len(path)-1])
			v := b.Get(bk)
			err := b.Delete(bk)
			if err == nil {
				// Old pair has been deleted, now add the new one
				err = b.Put([]byte(k), v)
			}
			// Now update the last key in the path
			return err
		}
		return errors.New("updatePairValue: Invalid Path")
	})
	return err
}

func (bd *boltDB) updatePairValue(path []string, v string) error {
	if bd.readOnly {
		return errors.New("DB is in Read-Only Mode")
	}
	err := bd.db.Update(func(tx *bolt.Tx) error {
		// len(b.GetPath())-1 is the key for the pair we're updating,
		// the rest are buckets leading to that key
		b := tx.Bucket([]byte(path[0]))
		if b != nil {
			if len(path) > 0 {
				for i := range path[1 : len(path)-1] {
					b = b.Bucket([]byte(path[i+1]))
					if b == nil {
						return errors.New("updatePairValue: Invalid Path")
					}
				}
			}
			// Now update the last key in the path
			err := b.Put([]byte(path[len(path)-1]), []byte(v))
			return err
		}
		return errors.New("updatePairValue: Invalid Path")
	})
	return err
}

func (bd *boltDB) insertBucket(path []string, n string) error {
	if bd.readOnly {
		return errors.New("DB is in Read-Only Mode")
	}
	// Inserts a new bucket named 'n' at 'path'
	err := bd.db.Update(func(tx *bolt.Tx) error {
		if len(path) == 0 {
			// insert at root
			_, err := tx.CreateBucket([]byte(n))
			if err != nil {
				return fmt.Errorf("insertBucket: %s", err)
			}
		} else {
			rootBucket, path := path[0], path[1:]
			b := tx.Bucket([]byte(rootBucket))
			if b != nil {
				for len(path) > 0 {
					tstBucket := ""
					tstBucket, path = path[0], path[1:]
					nB := b.Bucket([]byte(tstBucket))
					if nB == nil {
						// Not a bucket, if we're out of path, just move on
						if len(path) != 0 {
							// Out of path, error
							return errors.New("insertBucket: Invalid Path 1")
						}
					} else {
						b = nB
					}
				}
				_, err := b.CreateBucket([]byte(n))
				return err
			}
			return fmt.Errorf("insertBucket: Invalid Path %s", rootBucket)
		}
		return nil
	})
	return err
}

func (bd *boltDB) insertPair(path []string, k string, v string) error {
	if bd.readOnly {
		return errors.New("DB is in Read-Only Mode")
	}
	// Insert a new pair k => v at path
	err := bd.db.Update(func(tx *bolt.Tx) error {
		if len(path) == 0 {
			// We cannot insert a pair at root
			return errors.New("insertPair: Cannot insert pair at root")
		}
		var err error
		b := tx.Bucket([]byte(path[0]))
		if b != nil {
			if len(path) > 0 {
				for i := 1; i < len(path); i++ {
					b = b.Bucket([]byte(path[i]))
					if b == nil {
						return fmt.Errorf("insertPair: %s", err)
					}
				}
			}
			err := b.Put([]byte(k), []byte(v))
			if err != nil {
				return fmt.Errorf("insertPair: %s", err)
			}
		}
		return nil
	})
	return err
}

func (bd *boltDB) exportValue(path []string, fName string) error {
	return bd.db.View(func(tx *bolt.Tx) error {
		// len(b.path)-1 is the key whose value we want to export
		// the rest are buckets leading to that key
		b := tx.Bucket([]byte(path[0]))
		if b != nil {
			if len(path) > 1 {
				for i := range path[1 : len(path)-1] {
					b = b.Bucket([]byte(path[i+1]))
					if b == nil {
						return errors.New("exportValue: Invalid Path: " + strings.Join(path, "/"))
					}
				}
			}
			bk := []byte(path[len(path)-1])
			v := b.Get(bk)
			return writeToFile(fName, string(v)+"\n", os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
		} else {
			return errors.New("exportValue: Invalid Bucket")
		}
		return nil
	})
}

func (bd *boltDB) exportJSON(path []string, fName string) error {
	return bd.db.View(func(tx *bolt.Tx) error {
		// len(b.path)-1 is the key whose value we want to export
		// the rest are buckets leading to that key
		b := tx.Bucket([]byte(path[0]))
		if b != nil {
			if len(path) > 1 {
				for i := range path[1 : len(path)-1] {
					b = b.Bucket([]byte(path[i+1]))
					if b == nil {
						return errors.New("exportValue: Invalid Path: " + strings.Join(path, "/"))
					}
				}
			}
			bk := []byte(path[len(path)-1])
			if v := b.Get(bk); v != nil {
				return writeToFile(fName, "{\""+string(bk)+"\":\""+string(v)+"\"}", os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
			}
			if b.Bucket(bk) != nil {
				return writeToFile(fName, genJSONString(b.Bucket(bk)), os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
			}
			return writeToFile(fName, genJSONString(b), os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
		} else {
			return errors.New("exportValue: Invalid Bucket")
		}
		return nil
	})
}

func genJSONString(b *bolt.Bucket) string {
	ret := "{"
	b.ForEach(func(k, v []byte) error {
		ret = fmt.Sprintf("%s\"%s\":", ret, string(k))
		if v == nil {
			ret = fmt.Sprintf("%s%s,", ret, genJSONString(b.Bucket(k)))
		} else {
			ret = fmt.Sprintf("%s\"%s\",", ret, string(v))
		}
		return nil
	})
	ret = fmt.Sprintf("%s}", ret[:len(ret)-1])
	return ret
}

var f *os.File

func logToFile(s string) error {
	return writeToFile("bolt-log", s+"\n", os.O_RDWR|os.O_APPEND)
}

func writeToFile(fn, s string, mode int) error {
	var err error
	if f == nil {
		f, err = os.OpenFile(fn, mode, 0660)
	}
	if err != nil {
		return err
	}
	if _, err = f.WriteString(s); err != nil {
		return err
	}
	if err = f.Sync(); err != nil {
		return err
	}
	return nil
}
