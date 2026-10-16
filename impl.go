// SPDX-FileCopyrightText: © 2025,2026 W-A-T EU Operations Oü
// SPDX-License-Identifier: MPL-2.0
// SPDX-FileContributor: Created by Jose Luis Tallon <jltallon@w-a-t.group>

package variant

import (
// "fmt"
)

// ** for "merged" len/typeInfo fields -- superseeded by v2 (THIS) implementation
const (
	// Number of bits to use for Type field. n >= log2(max(VariantType))
	typeFieldBitCount = 4

	// Bit mask for Type part of lenInfo field.
	// ...or up to 2^27 == 128MB size
	typeFieldMask = (1 << typeFieldBitCount) - 1

	// Maximum length of a slice-type that can be stored in Variant.
	maxVectorLen = int((^uint(0))>>1) >> typeFieldBitCount

	ptrAddxLen  = 48
	ptrAddxMask = (1 << ptrAddxLen) - 1
	ptrAddxMark = (0xFC80 << ptrAddxLen)
)

// const (
// 	DISABLE_FREEMEM = false
// )

const (
	sCORRUPTED = "CORRUPT Variant!"
)

// func (v *Variant) Coerce(t VariantType) {
// 	switch t {
// 	case Empty: case Bool:
// 	case Int: case Float:
// 	case String: case Blob:
// 	case Undef:
// 		v.typInfo = uint32(t)
// 	default:
// 		return // 0xFFFF	// invalid
// 	}
// }

// implement GoStringer
func (v Variant) GoString() string {
	if v.IsEmpty() {
		return "empty"
	}
	if 0 != (Undef & VariantType(v.typInfo)) {
		return "undef"
	}
	return v.toString()
}
