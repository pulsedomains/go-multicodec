// Copyright © 2019 Weald Technology Trading
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package multicodec

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestID(t *testing.T) {
	tests := []struct {
		name string
		id   uint64
		err  error
	}{
		{
			name: "not present",
			err:  errors.New("unknown name not present"),
		},
		{
			name: "ipfs-ns",
			id:   0xe3,
		},
		{
			name: "ipfs",
			id:   0x1a5,
		},
		{
			name: "identity",
			id:   0x0,
		},
	}

	for i, test := range tests {
		id, err := ID(test.name)
		if test.err != nil {
			assert.Equal(t, test.err.Error(), err.Error(), fmt.Sprintf("incorrect error at test %d", i))
		} else {
			assert.Nil(t, err, fmt.Sprintf("unexpected error at test %d", i))
			assert.Equal(t, test.id, id, fmt.Sprintf("unexpected result at test %d", i))
		}
	}
}

func TestMustID(t *testing.T) {
	tests := []struct {
		name   string
		id     uint64
		panics bool
	}{
		{
			name:   "not present",
			panics: true,
		},
		{
			name: "ipfs-ns",
			id:   0xe3,
		},
		{
			name: "ipfs",
			id:   0x1a5,
		},
		{
			name: "identity",
			id:   0x0,
		},
	}

	for _, test := range tests {
		if test.panics {
			require.Panics(t, func() { MustID(test.name) })
		} else {
			require.NotPanics(t, func() { MustID(test.name) })
		}
	}
}

func TestName(t *testing.T) {
	tests := []struct {
		id   uint64
		name string
		err  error
	}{
		{
			id:  0x123fe78bc9a,
			err: errors.New("unknown ID 0x123fe78bc9a"),
		},
		{
			id:   0xe3,
			name: "ipfs-ns",
		},
		{
			id:   0x1a5,
			name: "p2p",
		},
		{
			id:   0x0,
			name: "identity",
		},
	}

	for i, test := range tests {
		name, err := Name(test.id)
		if test.err != nil {
			assert.Equal(t, test.err.Error(), err.Error(), fmt.Sprintf("incorrect error at test %d", i))
		} else {
			assert.Nil(t, err, fmt.Sprintf("unexpected error at test %d", i))
			assert.Equal(t, test.name, name, fmt.Sprintf("unexpected result at test %d", i))
		}
	}
}

func TestMustName(t *testing.T) {
	tests := []struct {
		id     uint64
		name   string
		panics bool
	}{
		{
			id:     0x123fe78bc9a,
			panics: true,
		},
		{
			id:   0xe3,
			name: "ipfs-ns",
		},
		{
			id:   0x1a5,
			name: "p2p",
		},
		{
			id:   0x0,
			name: "identity",
		},
	}

	for _, test := range tests {
		if test.panics {
			require.Panics(t, func() { MustName(test.id) })
		} else {
			require.NotPanics(t, func() { MustName(test.id) })
		}
	}
}

func TestAddCodec(t *testing.T) {
	tests := []struct {
		data          []byte
		name          string
		dataWithCodec []byte
		err           error
	}{
		{
			data: []byte{0x74, 0x65, 0x73, 0x74},
			name: "not present",
			err:  errors.New("unknown name not present"),
		},
		{
			data:          []byte{0x74, 0x65, 0x73, 0x74},
			name:          "ipfs-ns",
			dataWithCodec: []byte{0xe3, 0x01, 0x74, 0x65, 0x73, 0x74},
		},
	}

	for i, test := range tests {
		dataWithCodec, err := AddCodec(test.name, test.data)
		if test.err != nil {
			assert.Equal(t, test.err.Error(), err.Error(), fmt.Sprintf("incorrect error at test %d", i))
		} else {
			assert.Nil(t, err, fmt.Sprintf("unexpected error at test %d", i))
			assert.Equal(t, test.dataWithCodec, dataWithCodec, fmt.Sprintf("unexpected result at test %d", i))
		}
	}
}

func TestRemoveCodec(t *testing.T) {
	tests := []struct {
		data             []byte
		id               uint64
		dataWithoutCodec []byte
		err              error
	}{
		{
			dataWithoutCodec: []byte{},
			err:              errors.New("failed to find codec prefix to remove"),
		},
		{
			data:             []byte{0xe3, 0x01, 0x74, 0x65, 0x73, 0x74},
			dataWithoutCodec: []byte{0x74, 0x65, 0x73, 0x74},
			id:               0xe3,
		},
	}

	for i, test := range tests {
		dataWithoutCodec, id, err := RemoveCodec(test.data)
		if test.err != nil {
			assert.Equal(t, test.err.Error(), err.Error(), fmt.Sprintf("incorrect error at test %d", i))
		} else {
			assert.Nil(t, err, fmt.Sprintf("unexpected error at test %d", i))
			assert.Equal(t, test.dataWithoutCodec, dataWithoutCodec, fmt.Sprintf("unexpected result at test %d", i))
			assert.Equal(t, test.id, id, fmt.Sprintf("unexpected ID at test %d", i))
		}
	}
}

func TestIsCodec(t *testing.T) {
	tests := []struct {
		data   []byte
		codec  string
		result bool
	}{
		{ // 0
			data:   []byte{},
			codec:  "ipfs-ns",
			result: false,
		},
		{ // 1
			data:   []byte{0xe3, 0x01, 0x74, 0x65, 0x73, 0x74},
			codec:  "does not match",
			result: false,
		},
		{ // 2
			data:   []byte{0xef, 0x01, 0x74, 0x65, 0x73, 0x74},
			codec:  "ipfs-ns",
			result: false,
		},
		{ // 3
			data:   []byte{0xe3, 0x01, 0x74, 0x65, 0x73, 0x74},
			codec:  "ipfs-ns",
			result: true,
		},
	}

	for i, test := range tests {
		result := IsCodec(test.codec, test.data)
		assert.Equal(t, test.result, result, fmt.Sprintf("unexpected result at test %d", i))
	}
}

func TestGetCodec(t *testing.T) {
	tests := []struct {
		data  []byte
		codec string
		err   error
	}{
		{ // 0
			data: []byte{},
			err:  errors.New("failed to find codec prefix to remove"),
		},
		{ // 1
			data:  []byte{0xe3, 0x01, 0x74, 0x65, 0x73, 0x74},
			codec: "ipfs-ns",
		},
	}

	for i, test := range tests {
		codecID, err := GetCodec(test.data)
		if test.err != nil {
			assert.Equal(t, test.err.Error(), err.Error(), fmt.Sprintf("incorrect error at test %d", i))
		} else {
			assert.Nil(t, err, fmt.Sprintf("unexpected error at test %d", i))
			codec, err := Name(codecID)
			require.Nil(t, err, fmt.Sprintf("failed with error at test %d", i))
			assert.Equal(t, test.codec, codec, fmt.Sprintf("unexpected result at test %d", i))
		}
	}
}
