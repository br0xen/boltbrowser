package main

import (
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"strings"
)

// A Database, basically a collection of buckets
type BoltDB struct {
	buckets []BoltBucket
}

func (bd *BoltDB) getDepthFromPath(path []string) int {
	depth := 0
	b, p, err := bd.getGenericFromPath(path)
	if err != nil {
		return -1
	}
	if p != nil {
		b = p.parent
		depth += 1
	}
	for b != nil {
		b = b.parent
		depth += 1
	}
	return depth
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
				return b.getBucket(path[p])
			}
		}
		return b, nil
	}
	return nil, errors.New("Invalid Path")
}

func (bd *BoltDB) getPairFromPath(path []string) (*BoltPair, error) {
	b, err := memBolt.getBucketFromPath(path[:len(path)-1])
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
			n, err := bd.getVisibleItemCount(bd.buckets[i].path)
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
				n, err := bd.getVisibleItemCount(b.buckets[i].path)
				if err != nil {
					return 0, err
				}
				vis += n
			}
		}
	}
	return vis, ret_err
}

func (bd *BoltDB) buildVisiblePathSlice(path []string) ([]string, error) {
	var ret_slice []string
	var ret_err error
	if len(path) == 0 {
		for i := range bd.buckets {
			n, err := bd.buildVisiblePathSlice(bd.buckets[i].path)
			if err != nil {
				return nil, err
			}
			ret_slice = append(ret_slice, n...)
		}
	} else {
		b, err := bd.getBucketFromPath(path)
		if err != nil {
			return nil, err
		}
		// Add the bucket's path
		ret_slice = append(ret_slice, strings.Join(b.path, "/"))
		if b.expanded {
			// This bucket is expanded, add up it's children
			// * recurse for buckets
			for i := range b.buckets {
				n, err := bd.buildVisiblePathSlice(b.buckets[i].path)
				if err != nil {
					return nil, err
				}
				ret_slice = append(ret_slice, n...)
			}
			// * one path for each pair
			for i := range b.pairs {
				ret_slice = append(ret_slice, strings.Join(b.pairs[i].path, "/"))
			}
		}
	}
	return ret_slice, ret_err
}

func (bd *BoltDB) getPrevVisiblePath(path []string) []string {
	vis_paths, err := bd.buildVisiblePathSlice(nil)
	if path == nil {
		return strings.Split(vis_paths[len(vis_paths)-1], "/")
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
	vis_paths, err := bd.buildVisiblePathSlice(nil)
	if path == nil {
		return strings.Split(vis_paths[0], "/")
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

func (bd *BoltDB) getBucket(k string) (*BoltBucket, error) {
	for i := range bd.buckets {
		if bd.buckets[i].name == k {
			return &bd.buckets[i], nil
		}
	}
	return nil, errors.New("Bucket Not Found")
}

type BoltBucket struct {
	name     string
	path     []string
	pairs    []BoltPair
	buckets  []BoltBucket
	parent   *BoltBucket
	expanded bool
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

type BoltPair struct {
	path   []string
	parent *BoltBucket
	key    string
	val    string
}

func toggleOpenBucket(path []string) error {
	// Find the BoltBucket with a path == path
	b, err := memBolt.getBucketFromPath(path)
	if err == nil {
		b.expanded = !b.expanded
	}
	return err
}

func closeBucket(path []string) error {
	// Find the BoltBucket with a path == path
	b, err := memBolt.getBucketFromPath(path)
	if err == nil {
		b.expanded = false
	}
	return err
}

func openBucket(path []string) error {
	// Find the BoltBucket with a path == path
	b, err := memBolt.getBucketFromPath(path)
	if err == nil {
		b.expanded = true
	}
	return err
}

func deleteKey(path []string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		// len(b.path)-1 is the key we need to delete, the rest are buckets leading to that key
		b := tx.Bucket([]byte(path[0]))
		if b != nil {
			if len(path) > 1 {
				for i := range path[2 : len(path)-1] {
					b = b.Bucket([]byte(path[i+1]))
					if b == nil {
						return errors.New("deleteKey: Invalid Path")
					}
				}
			}
			// Now delete the last key in the path
			err := b.Delete([]byte(path[len(path)-1]))
			return err
		} else {
			return errors.New("deleteKey: Invalid Path")
		}
	})
	return err
}

func refreshDatabase() *BoltDB {
	// Reload the database into memBolt
	memBolt = new(BoltDB)
	db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(nm []byte, b *bolt.Bucket) error {
			bb, err := readBucket(b)
			if err == nil {
				bb.name = string(nm)
				bb.path = []string{bb.name}
				bb.expanded = false
				memBolt.buckets = append(memBolt.buckets, *bb)
				updatePaths(bb)
				return nil
			}
			return err
		})
	})
	return memBolt
}

func readBucket(b *bolt.Bucket) (*BoltBucket, error) {
	bb := new(BoltBucket)
	b.ForEach(func(k, v []byte) error {
		if v == nil {
			tb, err := readBucket(b.Bucket(k))
			tb.parent = bb
			if err == nil {
				tb.name = string(k)
				tb.path = append(bb.path, tb.name)
				bb.buckets = append(bb.buckets, *tb)
			}
		} else {
			tp := BoltPair{key: string(k), val: string(v)}
			tp.parent = bb
			tp.path = append(bb.path, tp.key)
			bb.pairs = append(bb.pairs, tp)
		}
		return nil
	})
	return bb, nil
}

func updatePaths(b *BoltBucket) {
	for i := range b.buckets {
		b.buckets[i].path = append(b.path, b.buckets[i].name)
		updatePaths(&b.buckets[i])
	}
	for i := range b.pairs {
		b.pairs[i].path = append(b.path, b.pairs[i].key)
	}
}

/*
func renameBucket(path []string, name string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		// len(b.path)-1 is the key we need to delete, the rest are buckets leading to that key
		b := tx.Bucket([]byte(path[0]))
		if b != nil {
			if len(path) > 1 {
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
			return errors.New("renameBucket: Invalid Path")
		}
	})
	refreshDatabase()
	return err
}
*/
func updatePairValue(path []string, v string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		// len(b.path)-1 is the key we need to delete, the rest are buckets leading to that key
		b := tx.Bucket([]byte(path[0]))
		if b != nil {
			if len(path) > 1 {
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
	// Inserts a new bucket named 'n' at 'path[len(path)-2]
	err := db.Update(func(tx *bolt.Tx) error {
		if len(path) == 1 {
			// insert at root
			_, err := tx.CreateBucket([]byte(n))
			if err != nil {
				return fmt.Errorf("insertBucket: %s", err)
			}
		} else if len(path) > 1 {
			var err error
			b := tx.Bucket([]byte(path[0]))
			if b != nil {
				if len(path) > 2 {
					for i := range path[1 : len(path)-2] {
						b = b.Bucket([]byte(path[i+1]))
						if b == nil {
							return fmt.Errorf("insertBucket: %s", err)
						}
					}
				}
				_, err := b.CreateBucket([]byte(n))
				if err != nil {
					return fmt.Errorf("insertBucket: %s", err)
				}
			}
		}
		return nil
	})
	return err
}
