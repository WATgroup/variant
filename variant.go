// SPDX-FileCopyrightText: © 2025,2026 W-A-T EU Operations Oü
// SPDX-License-Identifier: MPL-2.0
// SPDX-FileContributor: Created by Jose Luis Tallon <jltallon@w-a-t.group>

package variant

// Variant represents a *by-value* variant [scalar] variant
// * value contains the relevant data
// * length contains the value's length OR flags
// * typInfo contains the type information
// Alternative implementations would use a "lenInfo" field
// which contains the length and type information ...and fits in 32bits.
// But in 64bit platforms, due to alignment constraints, it'll be 64bits anyway
// so we simplify the implementation
type Variant struct {
	value   uintptr
	length  uint32
	typInfo uint32
}

type VariantType uint16

const (
	Empty  VariantType = iota
	Bool               // yes-or-no-or-maybe
	Int                // int64
	Float              // float64
	vector             // ptr-like storage
	String             // slice-of-runes ;)
	Blob               // slice-of-bytes
	Reserved

	shift VariantType = 1 << (typeFieldBitCount + 2)

	Flags             = Int + shift // uint32+
	Undef VariantType = shift

	mutable VariantType = 0x8000
)

const (
	VTrue  = uintptr(3 << 8)
	VFalse = uintptr(2 << 8) // & 0xFFFFFF)
	VDunno = 0
)

func (v Variant) IsEmpty() bool {
	return 0 == (uint64(v.value) | uint64(v.length) | uint64(v.typInfo))
}

// NewEmpty creates a Variant of type Empty. Equivalent to Variant{}.
func NewEmpty() Variant {
	return Variant{}
}

func (v Variant) Type() VariantType {
	if v.IsEmpty() {
		return Empty
	}
	return VariantType(v.typInfo)
}

func (v Variant) Len() int {
	if Bool == VariantType(v.typInfo&typeFieldMask) {
		return 1 // "fix" trick with encoding
	}
	return int(v.length)
}
