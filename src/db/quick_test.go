// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package db

import (
	"flag"
	"math/rand"
	"reflect"
	"testing/quick"
	"time"
)

// testing/quick defaults to 5 iterations and a random seed.
// You can override these settings from the command line:
//
//   -quick.count     The number of iterations to perform.
//   -quick.seed      The seed to use for randomizing.
//   -quick.maxitems  The maximum number of items to insert into a DB.
//   -quick.maxksize  The maximum size of a key.
//   -quick.maxvsize  The maximum size of a value.
//

var qcount, qseed, qmaxitems, qmaxksize, qmaxvsize int

func init() {
	flag.IntVar(&qcount, "quick.count", 5, "")
	flag.IntVar(&qseed, "quick.seed", int(time.Now().UnixNano())%100000, "")
	flag.IntVar(&qmaxitems, "quick.maxitems", 1000, "")
	flag.IntVar(&qmaxksize, "quick.maxksize", 1024, "")
	flag.IntVar(&qmaxvsize, "quick.maxvsize", 1024, "")
	flag.Parse()
	warn("seed:", qseed)
}

func qconfig() *quick.Config {
	return &quick.Config{
		MaxCount: qcount,
		Rand:     rand.New(rand.NewSource(int64(qseed))),
	}
}

type testdata []testdataitem

func (t testdata) Generate(rand *rand.Rand, size int) reflect.Value {
	n := rand.Intn(qmaxitems-1) + 1
	items := make(testdata, n)
	for i := 0; i < n; i++ {
		item := &items[i]
		item.Key = randByteSlice(rand, 1, qmaxksize)
		item.Value = randByteSlice(rand, 0, qmaxvsize)
	}
	return reflect.ValueOf(items)
}

type testdataitem struct {
	Key   []byte
	Value []byte
}

func randByteSlice(rand *rand.Rand, minSize, maxSize int) []byte {
	n := rand.Intn(maxSize-minSize) + minSize
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte(rand.Intn(255))
	}
	return b
}
