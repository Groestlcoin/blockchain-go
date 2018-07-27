// Copyright (c) 2013 Ben Johnson
// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Ensure that the page type can be returned in human readable format.
func TestPageTyp(t *testing.T) {
	assert.Equal(t, (&page{flags: branchPageFlag}).typ(), "branch")
	assert.Equal(t, (&page{flags: leafPageFlag}).typ(), "leaf")
	assert.Equal(t, (&page{flags: metaPageFlag}).typ(), "meta")
	assert.Equal(t, (&page{flags: bucketsPageFlag}).typ(), "buckets")
	assert.Equal(t, (&page{flags: freelistPageFlag}).typ(), "freelist")
	assert.Equal(t, (&page{flags: 20000}).typ(), "unknown<4e20>")
}

// Ensure that the hexdump debugging function doesn't blow up.
func TestPageDump(t *testing.T) {
	(&page{id: 256}).hexdump(16)
}
