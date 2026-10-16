// SPDX-FileCopyrightText: © 2026 W-A-T EU Operations Oü
// SPDX-License-Identifier: MPL-2.0
// SPDX-FileContributor: Created by Jose Luis Tallon <jltallon@w-a-t.group>

package variant

import (
	"encoding/base64"
	"fmt"
	"strconv"
)

func (v Variant) String() string {
	return v.toString()
}

func (v Variant) toString() string {

	switch VariantType(v.typInfo) {
	case Empty:
		return `<empty>`

	case Bool:
		switch v.value {
		case VTrue:
			return `true`
		case VFalse:
			return `false`
		case VDunno:
			return `neutral`
		}
	case Int:
		return strconv.FormatInt(int64(v.value), 10)
	case Float:
		return strconv.FormatFloat(opaque2float(v.value), 'g', 3, 64)
	case String:
		return stringdata(v.value)
	case Blob:
		return base64.StdEncoding.EncodeToString(blobdata(v.value))
	case Flags:
		return fmt.Sprintf("%08X", v.value)
	case Undef:
		// if 0 != (Undef & VariantType(v.typInfo)) {
		// 	return "undef"
		// }
		return `<undef>`
	default:
		break // XXX: ...
	}
	panic(sCORRUPTED)
	return "<??>"
}

func (v Variant) Blob() []byte {
	if Blob != VariantType(v.typInfo&typeFieldMask) {
		return nil
	}
	return blobdata(v.value)
}
