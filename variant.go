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

	Flags = Int + shift // uint32+

	shift VariantType = 1 << 7

	Undef VariantType = Bool + shift
)

const (
	VTrue  = uintptr(1 << 6)
	VFalse = uintptr(^(1 << 6) & 0xFFFFFF)
	VDunno = 1
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
	return VariantType(v.typInfo & typeFieldMask)
}

func NewInt(x int64) Variant {
	return Variant{uintptr(x), 1, uint32(Int)}
}
func (v Variant) Int() int64 {
	if uint32(Int) == v.typInfo { // fastpath
		return int64(v.value)
	}
	return 0
}

func NewFloat(x float64) Variant {
	return Variant{uintptr(x), 1, uint32(Float)}
}
func (v Variant) Float() float64 {
	if uint32(Float) == v.typInfo { // fastpath
		return float64(v.value)
	}
	return 0
}

func NewBool(x bool) Variant {
	var val uintptr
	switch x {
	case true:
		val = VTrue
	case false:
		val = VFalse
	}
	return Variant{val, 1, uint32(Bool)}
}
func NewTristate() Variant {
	return Variant{0, 1, uint32(Bool)}
}

func (v Variant) Bool(defval bool) bool {
	if uint32(Bool) == v.typInfo { // fastpath
		switch v.value {
		case VTrue:
			return true
		case VFalse:
			return false
		default:
			return defval
		}
	}
	return defval // upon VDunno
}
