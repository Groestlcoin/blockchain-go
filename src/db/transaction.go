// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package db

// Transaction represents a read-only transaction on the database.
// It can be used for retrieving values for keys as well as creating cursors for
// iterating over the data.
//
// IMPORTANT: You must close transactions when you are done with them. Pages
// can not be reclaimed by the writer until no more transactions are using them.
// A long running read transaction can cause the database to quickly grow.
type Transaction struct {
	db            *DB
	rwtransaction *RWTransaction
	meta          *meta
	buckets       *buckets
	pages         map[pgid]*page
}

// txnid represents the internal transaction identifier.
type txnid uint64

// init initializes the transaction and associates it with a database.
func (t *Transaction) init(db *DB) {
	t.db = db
	t.pages = nil

	// Copy the meta page since it can be changed by the writer.
	t.meta = &meta{}
	db.meta().copy(t.meta)

	// Read in the buckets page.
	t.buckets = &buckets{}
	t.buckets.read(t.page(t.meta.buckets))
}

// id returns the transaction id.
func (t *Transaction) id() txnid {
	return t.meta.txnid
}

// Close closes the transaction and releases any pages it is using.
func (t *Transaction) Close() {
	if t.rwtransaction != nil {
		t.rwtransaction.Rollback()
	}
	t.db.removeTransaction(t)
}

// DB returns a reference to the database that created the transaction.
func (t *Transaction) DB() *DB {
	return t.db
}

// Bucket retrieves a bucket by name.
// Returns nil if the bucket does not exist.
func (t *Transaction) Bucket(name string) *Bucket {
	b := t.buckets.get(name)
	if b == nil {
		return nil
	}

	return &Bucket{
		bucket:        b,
		name:          name,
		transaction:   t,
		rwtransaction: t.rwtransaction,
	}
}

// Buckets retrieves a list of all buckets.
func (t *Transaction) Buckets() []*Bucket {
	buckets := make([]*Bucket, 0, len(t.buckets.items))
	for name, b := range t.buckets.items {
		bucket := &Bucket{
			bucket:        b,
			name:          name,
			transaction:   t,
			rwtransaction: t.rwtransaction,
		}
		buckets = append(buckets, bucket)
	}
	return buckets
}

// page returns a reference to the page with a given id.
// If page has been written to then a temporary bufferred page is returned.
func (t *Transaction) page(id pgid) *page {
	// Check the dirty pages first.
	if t.pages != nil {
		if p, ok := t.pages[id]; ok {
			return p
		}
	}

	// Otherwise return directly from the mmap.
	return t.db.page(id)
}

// forEachPage iterates over every page within a given page and executes a function.
func (t *Transaction) forEachPage(pgid pgid, depth int, fn func(*page, int)) {
	p := t.page(pgid)

	// Execute function.
	fn(p, depth)

	// Recursively loop over children.
	if (p.flags & branchPageFlag) != 0 {
		for i := 0; i < int(p.count); i++ {
			elem := p.branchPageElement(uint16(i))
			t.forEachPage(elem.pgid, depth+1, fn)
		}
	}
}
