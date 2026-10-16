// SPDX-FileCopyrightText: © 2026 W-A-T EU Operations Oü
// SPDX-License-Identifier: MPL-2.0
// SPDX-FileContributor: Created by Jose Luis Tallon <jltallon@w-a-t.group>

package variant_test

import (
	vv "github.com/WATgroup/variant"

	"math"
	"testing"
)

func TestNewEmptyIsEmpty(t *testing.T) {
	if !vv.NewEmpty().IsEmpty() {
		t.Error("NewEmpty() should be empty")
	}
	if vv.NewEmpty().Type() != vv.Empty {
		t.Errorf("NewEmpty().Type() = %v, want Empty", vv.NewEmpty().Type())
	}
}

func TestZeroValueIsEmpty(t *testing.T) {
	var v vv.Variant
	if !v.IsEmpty() {
		t.Error("zero-value Variant should be empty")
	}
}

func TestIntRoundTrip(t *testing.T) {
	cases := []int64{0, 1, -1, 42, -42, math.MaxInt32, math.MinInt32}
	for _, c := range cases {
		v := vv.NewInt(c)
		if got := v.Int(); got != c {
			t.Errorf("NewInt(%d).Int() = %d, want %d", c, got, c)
		}
		if v.Type() != vv.Int {
			t.Errorf("NewInt(%d).Type() = %v, want Int", c, v.Type())
		}
		if v.IsEmpty() {
			t.Errorf("NewInt(%d) should not be Empty", c)
		}
	}
}

// NewInt stores its value via uintptr(x). On 32-bit plataforms uintptr is only 32 bits wide, so values outside that range would silently truncate.
// This will only fail on a 32-bit build
func TestIntFull64BitRange(t *testing.T) {
	cases := []int64{math.MaxInt64, math.MinInt64}
	for _, c := range cases {
		v := vv.NewInt(c)
		if got := v.Int(); got != c {
			t.Errorf("NewInt(%d).Int() = %d, want %d (expected failure on 32-bit plataforms only)", c, got, c)
		}
	}
}

// NewFloat/Float() are expected to round-trip a float64. As implemented, NewFloat converts via uintptr(x),
// a numeric conversion that truncates towards 0 rather than preserving the value. This test documents
// the intended contract. if this is right, it will fail for every non-integer case.
func TestFloatRoundTrip(t *testing.T) {
	cases := []float64{0, 1, -1, math.Phi, -math.E, math.Pi, -3.14, 0.5}
	for _, c := range cases {
		v := vv.NewFloat(c)
		if got := v.Float(); got != c {
			t.Errorf("NewFloat(%v).Float() = %v, want %v", c, got, c)
		}
	}
}

// --- Bool / Tristate
func TestBoolRoundTrip(t *testing.T) {
	if got := vv.NewBool(true).Bool(); got != true {
		t.Errorf("NewBool(true).Bool(false) = %v, want true", got)
	}
	if got := vv.NewBool(false).Bool(); got != false {
		t.Errorf("NewBool(false).Bool(true) = %v, want false", got)
	}
}

func TestTristateReturnsDefault(t *testing.T) {
	tri := vv.NewTristate()
	if got := tri.Bool3(true); got != true {
		t.Errorf("NewTristate().Bool(true) = %v, want true (default)", got)
	}
	if got := tri.Bool3(false); got != false {
		t.Errorf("NewTristate().Bool(false) = %v, want false (default)", got)
	}
	if tri.IsEmpty() {
		t.Errorf("NewTristate() should not be Empty (length = 1 distinguishes it from the zero value)")
	}
}

func TestTristateLogic(t *testing.T) {
	tri := vv.NewTristate() // dunno
	if got := tri.LogicAnd(true); got != true {
		t.Errorf("NewTristate().LogicAnd(true) = %v, want true -- triLogic", got)
	}
	if got := tri.LogicOr(false); got != false {
		t.Errorf("NewTristate().LogicOr(false) = %v, want false -- triLogic", got)
	}

	val := vv.NewInt(0)
	if got := val.LogicOr(false); got != false {
		t.Errorf("NewInt().LogicOr(false) = %v, want false (default)", got)
	}
}

////////////////////////////////////////////////////////////////////////////////

func TestTypeAndValues(t *testing.T) {
	cases := []struct {
		name   string
		v      vv.Variant
		want   vv.VariantType
		strrep string
	}{
		{"Empty", vv.NewEmpty(), vv.Empty, "<empty>"},
		{"Int", vv.NewInt(1), vv.Int, "1"},
		{"Float", vv.NewFloat(4), vv.Float, "4"}, // value 4 chosen deliberately ---> an integer float, so it should pass even if TestFloatRoundTrip fails. Interesting bit-pattern
		{"Bool", vv.NewBool(true), vv.Bool, "true"},
		{"Flags", vv.NewFlags(0x55AA), vv.Flags, "000055AA"},
		{"String", vv.NewString("x"), vv.String, "x"},
		{"Blob", vv.NewBlob([]byte("xXx")), vv.Blob, "eFh4"},
	}
	for _, c := range cases {
		if got, val := c.v.Type(), c.v.String(); got != c.want || val != c.strrep {
			t.Errorf(`%s: Type() = %v, want %v; String()->"%s" (expected "%s")`, c.name, got, c.want, val, c.strrep)
		}
	}
}
