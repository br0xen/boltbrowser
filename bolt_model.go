package main

import (
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"os"
	"strings"
)

// A Database, basically a collection of buckets
type BoltDB struct {
	buckets []BoltBucket
}

type BoltBucket struct {
	name       string
	pairs      []BoltPair
	buckets    []BoltBucket
	parent     *BoltBucket
	expanded   bool
	error_flag bool
}

type BoltPair struct {
	parent *BoltBucket
	key    string
	val    string
}

func (bd *BoltDB) getGenericFromPath(path []string) (*BoltBucket, *BoltPair, error) {
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

func (bd *BoltDB) getBucketFromPath(path []string) (*BoltBucket, error) {
	if len(path) > 0 {
		// Find the BoltBucket with a path == path
		var b *BoltBucket
		var err error
		// Find the root bucket
		b, err = memBolt.getBucket(path[0])
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

func (bd *BoltDB) getPairFromPath(path []string) (*BoltPair, error) {
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

func (bd *BoltDB) getVisibleItemCount(path []string) (int, error) {
	vis := 0
	var ret_err error
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
		vis += 1
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
	return vis, ret_err
}

func (bd *BoltDB) buildVisiblePathSlice() ([]string, error) {
	var ret_slice []string
	var ret_err error
	// The root path, recurse for root buckets
	for i := range bd.buckets {
		bkt_s, bkt_err := bd.buckets[i].buildVisiblePathSlice([]string{})
		if bkt_err == nil {
			ret_slice = append(ret_slice, bkt_s...)
		} else {
			// Something went wrong, set the error flag
			bd.buckets[i].error_flag = true
		}
	}
	return ret_slice, ret_err
}

func (bd *BoltDB) getPrevVisiblePath(path []string) []string {
	vis_paths, err := bd.buildVisiblePathSlice()
	if path == nil {
		if len(vis_paths) > 0 {
			return strings.Split(vis_paths[len(vis_paths)-1], "/")
		} else {
			return nil
		}
	}
	if err == nil {
		find_path := strings.Join(path, "/")
		for i := range vis_paths {
			if vis_paths[i] == find_path && i > 0 {
				return strings.Split(vis_paths[i-1], "/")
			}
		}
	}
	return nil
}
func (bd *BoltDB) getNextVisiblePath(path []string) []string {
	vis_paths, err := bd.buildVisiblePathSlice()
	if path == nil {
		if len(vis_paths) > 0 {
			return strings.Split(vis_paths[0], "/")
		} else {
			return nil
		}
	}
	if err == nil {
		find_path := strings.Join(path, "/")
		for i := range vis_paths {
			if vis_paths[i] == find_path && i < len(vis_paths)-1 {
				return strings.Split(vis_paths[i+1], "/")
			}
		}
	}
	return nil
}

func (bd *BoltDB) toggleOpenBucket(path []string) error {
	// Find the BoltBucket with a path == path
	b, err := bd.getBucketFromPath(path)
	if err == nil {
		b.expanded = !b.expanded
	}
	return err
}

func (bd *BoltDB) closeBucket(path []string) error {
	// Find the BoltBucket with a path == path
	b, err := bd.getBucketFromPath(path)
	if err == nil {
		b.expanded = false
	}
	return err
}

func (bd *BoltDB) openBucket(path []string) error {
	// Find the BoltBucket with a path == path
	b, err := bd.getBucketFromPath(path)
	if err == nil {
		b.expanded = true
	}
	return err
}

func (bd *BoltDB) getBucket(k string) (*BoltBucket, error) {
	for i := range bd.buckets {
		if bd.buckets[i].name == k {
			return &bd.buckets[i], nil
		}
	}
	return nil, errors.New("Bucket Not Found")
}

func (bd *BoltDB) syncOpenBuckets(shadow *BoltDB) {
	// First test this bucket
	for i := range bd.buckets {
		for j := range shadow.buckets {
			if bd.buckets[i].name == shadow.buckets[j].name {
				bd.buckets[i].syncOpenBuckets(&shadow.buckets[j])
			}
		}
	}
}

func (bd *BoltDB) refreshDatabase() *BoltDB {
	// Reload the database into memBolt
	memBolt = new(BoltDB)
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

func (b *BoltBucket) GetPath() []string {
	if b.parent != nil {
		return append(b.parent.GetPath(), b.name)
	} else {
		return []string{b.name}
	}
}

func (b *BoltBucket) buildVisiblePathSlice(prefix []string) ([]string, error) {
	var ret_slice []string
	var ret_err error
	// Add this bucket to the slice
	prefix = append(prefix, b.name)
	ret_slice = append(ret_slice, strings.Join(prefix, "/"))
	if b.expanded {
		// Add subbuckets
		for i := range b.buckets {
			bkt_s, bkt_err := b.buckets[i].buildVisiblePathSlice(prefix)
			if bkt_err == nil {
				ret_slice = append(ret_slice, bkt_s...)
			} else {
				// Something went wrong, set the error flag
				b.buckets[i].error_flag = true
			}
		}
		// Add Pairs
		for i := range b.pairs {
			ret_slice = append(ret_slice, strings.Join(prefix, "/")+"/"+b.pairs[i].key)
		}
	}
	return ret_slice, ret_err
}

func (b *BoltBucket) syncOpenBuckets(shadow *BoltBucket) {
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

func (b *BoltBucket) getBucket(k string) (*BoltBucket, error) {
	for i := range b.buckets {
		if b.buckets[i].name == k {
			return &b.buckets[i], nil
		}
	}
	return nil, errors.New("Bucket Not Found")
}

func (b *BoltBucket) getPair(k string) (*BoltPair, error) {
	for i := range b.pairs {
		if b.pairs[i].key == k {
			return &b.pairs[i], nil
		}
	}
	return nil, errors.New("Pair Not Found")
}

func (p *BoltPair) GetPath() []string {
	return append(p.parent.GetPath(), p.key)
}

/* This is a go-between function (between the boltbrowser structs
 * above, and the bolt convenience functions below)
 * for taking a boltbrowser bucket and recursively adding it
 * and all of it's children into the database.
 * Mainly used for moving a bucket from one path to another
 * as in the 'renameBucket' function below.
 */
func addBucketFromBoltBucket(path []string, bb *BoltBucket) error {
	if err := insertBucket(path, bb.name); err == nil {
		bucket_path := append(path, bb.name)
		for i := range bb.pairs {
			if err = insertPair(bucket_path, bb.pairs[i].key, bb.pairs[i].val); err != nil {
				return err
			}
		}
		for i := range bb.buckets {
			if err = addBucketFromBoltBucket(bucket_path, &bb.buckets[i]); err != nil {
				return err
			}
		}
	}
	return nil
}

func deleteKey(path []string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		// len(b.path)-1 is the key we need to delete,
		// the rest are buckets leading to that key
		if len(path) == 1 {
			// Deleting a root bucket
			return tx.DeleteBucket([]byte(path[0]))
		} else {
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
				if delete_bkt := b.Bucket([]byte(path[len(path)-1])); delete_bkt == nil {
					// Must be a pair
					err = b.Delete([]byte(path[len(path)-1]))
				} else {
					err = b.DeleteBucket([]byte(path[len(path)-1]))
				}
				return err
			} else {
				return errors.New("deleteKey: Invalid Path")
			}
		}
	})
	return err
}

func readBucket(b *bolt.Bucket) (*BoltBucket, error) {
	bb := new(BoltBucket)
	b.ForEach(func(k, v []byte) error {
		if v == nil {
			tb, err := readBucket(b.Bucket(k))
			tb.parent = bb
			if err == nil {
				tb.name = string(k)
				bb.buckets = append(bb.buckets, *tb)
			}
		} else {
			tp := BoltPair{key: string(k), val: string(v)}
			tp.parent = bb
			bb.pairs = append(bb.pairs, tp)
		}
		return nil
	})
	return bb, nil
}

func renameBucket(path []string, name string) error {
	if name == path[len(path)-1] {
		// No change requested
		return nil
	}
	var bb *BoltBucket // For caching the current bucket
	err := db.View(func(tx *bolt.Tx) error {
		// len(b.path)-1 is the key we need to delete,
		// the rest are buckets leading to that key
		b := tx.Bucket([]byte(path[0]))
		if b != nil {
			if len(path) > 1 {
				for i := range path[1:len(path)] {
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
	if err = deleteKey(path); err != nil {
		return err
	}
	// Rechristen our cached bucket
	bb.name = name
	// And re-add it

	parent_path := path[:len(path)-1]
	if err = addBucketFromBoltBucket(parent_path, bb); err != nil {
		return err
	}
	return nil
}

func updatePairKey(path []string, k string) error {
	err := db.Update(func(tx *bolt.Tx) error {
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
		} else {
			return errors.New("updatePairValue: Invalid Path")
		}
	})
	return err
}

func updatePairValue(path []string, v string) error {
	err := db.Update(func(tx *bolt.Tx) error {
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
		} else {
			return errors.New("updatePairValue: Invalid Path")
		}
	})
	return err
}

func insertBucket(path []string, n string) error {
	// Inserts a new bucket named 'n' at 'path'
	err := db.Update(func(tx *bolt.Tx) error {
		if len(path) == 0 {
			// insert at root
			_, err := tx.CreateBucket([]byte(n))
			if err != nil {
				return fmt.Errorf("insertBucket: %s", err)
			}
		} else {
			root_bucket, path := path[0], path[1:]
			b := tx.Bucket([]byte(root_bucket))
			if b != nil {
				for len(path) > 0 {
					tst_bucket := ""
					tst_bucket, path = path[0], path[1:]
					n_b := b.Bucket([]byte(tst_bucket))
					if n_b == nil {
						// Not a bucket, if we're out of path, just move on
						if len(path) != 0 {
							// Out of path, error
							return errors.New("insertBucket: Invalid Path 1")
						}
					} else {
						b = n_b
					}
				}
				_, err := b.CreateBucket([]byte(n))
				return err
			}
			return fmt.Errorf("insertBucket: Invalid Path %s", root_bucket)
		}
		return nil
	})
	return err
}

func insertPair(path []string, k string, v string) error {
	// Insert a new pair k => v at path
	err := db.Update(func(tx *bolt.Tx) error {
		if len(path) == 0 {
			// We cannot insert a pair at root
			return errors.New("insertPair: Cannot insert pair at root.")
		} else {
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
		}
		return nil
	})
	return err
}

var f *os.File

func logToFile(s string) error {
	var err error
	if f == nil {
		f, err = os.OpenFile("bolt-log", os.O_RDWR|os.O_APPEND, 0660)
	}
	if err != nil {
		return err
	}
	if _, err = f.WriteString(s + "\n"); err != nil {
		return err
	}
	if err = f.Sync(); err != nil {
		return err
	}
	return nil
}
