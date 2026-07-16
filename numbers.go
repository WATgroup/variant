// SPDX-FileCopyrightText: © 2025,2026 W-A-T EU Operations Oü
// SPDX-License-Identifier: MPL-2.0
// SPDX-FileContributor: Created by Jose Luis Tallon <jltallon@w-a-t.group>

package variant

import "math"

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
	return Variant{uintptr(math.Float64bits(x)), 1, uint32(Float)}
}

func (v Variant) Float() float64 {
	if uint32(Float) == v.typInfo { // fastpath
		return math.Float64frombits(uint64(v.value))
	}
	if uint32(Float|shift) == v.typInfo {
		return math.NaN()
	}
	return 0
}

func opaque2float(v uintptr) float64 {
	return math.Float64frombits(uint64(v))
}
