package git

/*
#include <git2.h>
*/
import "C"
import (
	"runtime"
	"unsafe"
)

func (v *Repository) AddIgnoreRule(rules string) error {
	crules := C.CString(rules)
	defer C.free(unsafe.Pointer(crules))

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_ignore_add_rule(v.ptr, crules)
	if ret < 0 {
		return MakeGitError(ret)
	}
	return nil
}

func (v *Repository) ClearInternalIgnoreRules() error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_ignore_clear_internal_rules(v.ptr)
	if ret < 0 {
		return MakeGitError(ret)
	}
	return nil
}

func (v *Repository) IsPathIgnored(path string) (bool, error) {
	var ignored C.int

	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_ignore_path_is_ignored(&ignored, v.ptr, cpath)
	if ret < 0 {
		return false, MakeGitError(ret)
	}
	return ignored == 1, nil
}
