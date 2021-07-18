// Code generated by "stringer -type=OS -linecomment"; DO NOT EDIT.

package appcast

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Unknown-0]
	_ = x[MacOS-1]
	_ = x[Windows64-2]
	_ = x[Windows32-3]
}

const _OS_name = "macoswindows-x64windows-x86"

var _OS_index = [...]uint8{0, 0, 5, 16, 27}

func (i OS) String() string {
	if i >= OS(len(_OS_index)-1) {
		return "OS(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _OS_name[_OS_index[i]:_OS_index[i+1]]
}
