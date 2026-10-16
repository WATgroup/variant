// SPDX-FileCopyrightText: © 2026 W-A-T EU Operations Oü
// SPDX-License-Identifier: MPL-2.0
// SPDX-FileContributor: Created by Jose Luis Tallon <jltallon@w-a-t.group>

package variant

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
	return Variant{0, 2, uint32(Bool)} // default "dunno"
}

func (v Variant) IsDunno() bool {
	// assume we can rely on shortcircuit evaluation
	return (Bool == VariantType(v.typInfo) && 2 == v.length) //(0 == v.value)) || (0 != (v.typInfo & uint32(shift)))
}

func (v Variant) Bool() bool {
	switch v.value {
	case VTrue:
		return true
	case VFalse:
		return false
	}
	// XXX: TODO:  panic upon dunno ??
	return false
}

func (v Variant) Bool3(defval bool) bool {
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
	// XXX: slightly "surprising" behaviour...
	return defval
}

func (v Variant) LogicAnd(input bool) bool {
	if (uint32(Bool) == v.typInfo) && (2 == v.length || VDunno == v.value) {
		// Dunno!
		return input
	}
	switch v.value {
	case VTrue:
		return input
	case VFalse:
		return false
	default:
		break
	}
	return input && (0 != v.value) // shouldn't happen...
}

func (v Variant) LogicOr(input bool) bool {
	if (uint32(Bool) == v.typInfo) && (2 == v.length || VDunno == v.value) {
		// Dunno!
		return input
	}
	switch v.value {
	case VTrue:
		return true
	case VFalse:
		return input
	default:
		break
	}
	return input || (0 != v.value) // shouldn't happen...
}

///////////////////////////////////////////////////////////////////////////////

func NewFlags(val uint32) Variant {
	return Variant{uintptr(val), 1, uint32(Flags)}
}
