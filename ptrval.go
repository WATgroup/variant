// SPDX-FileCopyrightText: © 2025,2026 W-A-T EU Operations Oü
// SPDX-License-Identifier: MPL-2.0
// SPDX-FileContributor: Created by Jose Luis Tallon <jltallon@w-a-t.group>

package variant

import (
	"runtime"
	"sync"
	"sync/atomic"
)

var (
	map_strings stringMap
	map_bytevec bytevecMap

	//map_boolvec
	//map_intvec
	//map_floatvec
)

var handleCnt atomic.Uintptr
var vecHeuristic uint32 = 13

func init() {
	// initialize the "pointer base" with a large and "distinct" value
	// ..to help with diagnostics and make it fail clearly if misused ;)
	handleCnt.Store(ptrAddxMark)
}

/** IMPLEMENTATION NOTES:
 * This is similar to runtime/cgo's approach... only specialized
 * We only support a few specific types... hence pre-specialize it
 *
 * Note well: this is *the only correct way* to do it in Go
 * Messing with pointers "à la C union" is guaranteed to panic once Go's GC
 * does its thing -- specially as soon as it starts moving mem blocks around(!)
 *
 * Type-specific maps are demand-
 *
 * Until go's generics implementation matures / supports partial specialization,
 * this will be faster and more efficient
 */

type stringMap struct {
	mx     sync.RWMutex
	data   map[uintptr]string
	inited uint32
	_      noCopy
}

// The relative overhead of a slice (2 64bit words) is negligible when compared
type bytevec = []byte

type bytevecMap struct {
	mx     sync.RWMutex
	data   map[uintptr]bytevec
	inited uint32
	_      noCopy
}

///////////////////////////////////////////////////////////////////////////////

func NewString(str string) (ret Variant) {
	h, l := regString(str)
	ret = Variant{h, l, uint32(String & typeFieldMask)}
	if !DISABLE_FREEMEM {
		runtime.SetFinalizer(&ret, delString)
	}
	return
}

func NewStringEx(str []byte) (ret Variant) {
	h, l := regString(string(str))
	ret = Variant{h, l, uint32(String & typeFieldMask)}
	if !DISABLE_FREEMEM {
		runtime.SetFinalizer(&ret, delString)
	}
	return
}

func (v Variant) String() string {
	if 0 == v.value {
		return ""
	}
	return stringdata(v.value)
}

func stringdata(v uintptr) string {
	map_strings.mx.RLock()
	defer map_strings.mx.RUnlock()

	return map_strings.data[v]
}

//////////////////////////////////////////////////////////////////////

func NewBlob(x []byte) (ret Variant) {
	h, l := regBlob(x)
	ret = Variant{h, l, uint32(Blob & typeFieldMask)}
	runtime.SetFinalizer(&ret, delBlob) // XXX: TODO: check validity of supposition
	return
}

// /////////////////////////////////////////////////////////////////////////////
func initStrings() {
	if atomic.CompareAndSwapUint32(&map_strings.inited, 0, 1) {
		// we are the ones initializing: do it
		map_strings.data = make(map[uintptr]string, vecHeuristic)
	}
}

func regString(x string) (uintptr, uint32) {
	l := uint32(len(x) & maxVectorLen)

	initStrings()
	map_strings.mx.Lock()
	defer map_strings.mx.Unlock()

	c := handleCnt.Add(1)
	map_strings.data[c] = x

	return c, l
}

func delString(o *Variant) {
	delete(map_strings.data, o.value)
}

///////////////////////////////////////////////////////////

func initBlobs() {
	if atomic.CompareAndSwapUint32(&map_bytevec.inited, 0, 1) {
		// we are the ones initializing: do it
		map_bytevec.data = make(map[uintptr]bytevec, vecHeuristic)
	}
}

func regBlob(x []byte) (uintptr, uint32) {

	l := uint32(len(x) & maxVectorLen)

	initBlobs()
	map_bytevec.mx.Lock()
	defer map_bytevec.mx.Unlock()

	c := handleCnt.Add(1)
	map_bytevec.data[c] = x

	return c, l
}

func delBlob(o *Variant) {
	delete(map_bytevec.data, o.value)
}

// /////////////////////////////////////////////////////////////////////////////
type noCopy struct{}

func (_ *noCopy) Lock()   {}
func (_ *noCopy) Unlock() {}
