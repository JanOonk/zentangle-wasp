// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package wasmclient

import (
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/dict"
)

// The Results struct is used to gather all arguments for a smart
// contract function call and encode it into a deterministic byte array
type Results struct {
	res dict.Dict
}

func (r Results) Exists(key string) bool {
	_, ok := r.res[kv.Key(key)]
	return ok
}

func (r Results) get(key string, typeID int32) []byte {
	size := TypeSizes[typeID]
	bytes, ok := r.res[kv.Key(key)]
	if ok {
		if size != 0 && len(bytes) != int(size) {
			panic("invalid type size")
		}
		return bytes
	}
	// return default all-zero bytes value
	return make([]byte, size)
}

func (r Results) getBase58(key string, typeID int32) string {
	return Base58Encode(r.get(key, typeID))
}

func (r Results) GetAddress(key string) Address {
	return Address(r.getBase58(key, TYPE_ADDRESS))
}

func (r Results) GetAgentID(key string) AgentID {
	return AgentID(r.getBase58(key, TYPE_AGENT_ID))
}

func (r Results) GetBytes(key string) []byte {
	return r.get(key, TYPE_BYTES)
}

func (r Results) GetBool(key string) bool {
	return r.get(key, TYPE_BOOL)[0] != 0
}

func (r Results) GetChainID(key string) ChainID {
	return ChainID(r.getBase58(key, TYPE_CHAIN_ID))
}

func (r Results) GetColor(key string) Color {
	return Color(r.getBase58(key, TYPE_COLOR))
}

func (r Results) GetHash(key string) Hash {
	return Hash(r.getBase58(key, TYPE_HASH))
}

func (r Results) GetHname(key string) Hname {
	return Hname(r.getUint64(key, TYPE_HNAME))
}

func (r Results) GetInt8(key string) int8 {
	return int8(r.get(key, TYPE_INT8)[0])
}

func (r Results) GetInt16(key string) int16 {
	return int16(r.getUint64(key, TYPE_INT16))
}

func (r Results) GetInt32(key string) int32 {
	return int32(r.getUint64(key, TYPE_INT32))
}

func (r Results) GetInt64(key string) int64 {
	return int64(r.getUint64(key, TYPE_INT64))
}

func (r Results) GetRequestID(key string) RequestID {
	return RequestID(r.getBase58(key, TYPE_REQUEST_ID))
}

func (r Results) GetString(key string) string {
	return string(r.get(key, TYPE_STRING))
}

func (r Results) GetUint8(key string) uint8 {
	return r.get(key, TYPE_INT8)[0]
}

func (r Results) GetUint16(key string) uint16 {
	return uint16(r.getUint64(key, TYPE_INT16))
}

func (r Results) GetUint32(key string) uint32 {
	return uint32(r.getUint64(key, TYPE_INT32))
}

func (r Results) GetUint64(key string) uint64 {
	return r.getUint64(key, TYPE_INT64)
}

func (r Results) getUint64(key string, typeID int32) uint64 {
	b := r.get(key, typeID)
	v := uint64(0)
	for i := len(b) - 1; i >= 0; i-- {
		v = (v << 8) | uint64(b[i])
	}
	return v
}

// TODO Decode() from view call response into map
