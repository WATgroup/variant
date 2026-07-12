// SPDX-FileCopyrightText: © 2025,2026 W-A-T EU Operations Oü
// SPDX-License-Identifier: MPL-2.0
// SPDX-FileContributor: Created by Jose Luis Tallon <jltallon@w-a-t.group>

package variant

/**
 * A Variant is implemented as a struct containing the following information:
 * `value`, `length`, `typInfo`
 * ...hence, it is a BY-VALUE TYPE, which can be copied and sent around :)
 * .
 * What gets stored in the struct fields depends on the type of the Variant.
 *
 * This is version TWO of the implementation: the first one was based on an
 * incorrect [C-like] memory model, which would have panicked under even a
 * slight memory pressure --- with a relocating GC, which is the case with
 * GreenTea (go 1.25+), it's more probable than ever
 */
