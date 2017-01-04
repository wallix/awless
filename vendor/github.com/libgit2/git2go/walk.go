package git

/*
#include <git2.h>
*/
import "C"

import (
	"runtime"
	"unsafe"
)

// RevWalk

type SortType uint

const (
	SortNone        SortType = C.GIT_SORT_NONE
	SortTopological SortType = C.GIT_SORT_TOPOLOGICAL
	SortTime        SortType = C.GIT_SORT_TIME
	SortReverse     SortType = C.GIT_SORT_REVERSE
)

type RevWalk struct {
	ptr  *C.git_revwalk
	repo *Repository
}

func revWalkFromC(repo *Repository, c *C.git_revwalk) *RevWalk {
	v := &RevWalk{ptr: c, repo: repo}
	runtime.SetFinalizer(v, (*RevWalk).Free)
	return v
}

func (v *RevWalk) Reset() {
	C.git_revwalk_reset(v.ptr)
}

func (v *RevWalk) Push(id *Oid) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ecode := C.git_revwalk_push(v.ptr, id.toC())
	if ecode < 0 {
		return MakeGitError(ecode)
	}
	return nil
}

func (v *RevWalk) PushGlob(glob string) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	cstr := C.CString(glob)
	defer C.free(unsafe.Pointer(cstr))

	ecode := C.git_revwalk_push_glob(v.ptr, cstr)
	if ecode < 0 {
		return MakeGitError(ecode)
	}
	return nil
}

func (v *RevWalk) PushRange(r string) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	cstr := C.CString(r)
	defer C.free(unsafe.Pointer(cstr))

	ecode := C.git_revwalk_push_range(v.ptr, cstr)
	if ecode < 0 {
		return MakeGitError(ecode)
	}
	return nil
}

func (v *RevWalk) PushRef(r string) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	cstr := C.CString(r)
	defer C.free(unsafe.Pointer(cstr))

	ecode := C.git_revwalk_push_ref(v.ptr, cstr)
	if ecode < 0 {
		return MakeGitError(ecode)
	}
	return nil
}

func (v *RevWalk) PushHead() (err error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ecode := C.git_revwalk_push_head(v.ptr)
	if ecode < 0 {
		err = MakeGitError(ecode)
	}
	return nil
}

func (v *RevWalk) Hide(id *Oid) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ecode := C.git_revwalk_hide(v.ptr, id.toC())
	if ecode < 0 {
		return MakeGitError(ecode)
	}
	return nil
}

func (v *RevWalk) HideGlob(glob string) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	cstr := C.CString(glob)
	defer C.free(unsafe.Pointer(cstr))

	ecode := C.git_revwalk_hide_glob(v.ptr, cstr)
	if ecode < 0 {
		return MakeGitError(ecode)
	}
	return nil
}

func (v *RevWalk) HideRef(r string) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	cstr := C.CString(r)
	defer C.free(unsafe.Pointer(cstr))

	ecode := C.git_revwalk_hide_ref(v.ptr, cstr)
	if ecode < 0 {
		return MakeGitError(ecode)
	}
	return nil
}

func (v *RevWalk) HideHead() (err error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ecode := C.git_revwalk_hide_head(v.ptr)
	if ecode < 0 {
		err = MakeGitError(ecode)
	}
	return nil
}

func (v *RevWalk) Next(id *Oid) (err error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_revwalk_next(id.toC(), v.ptr)
	switch {
	case ret < 0:
		err = MakeGitError(ret)
	}

	return
}

type RevWalkIterator func(commit *Commit) bool

func (v *RevWalk) Iterate(fun RevWalkIterator) (err error) {
	oid := new(Oid)
	for {
		err = v.Next(oid)
		if IsErrorCode(err, ErrIterOver) {
			return nil
		}
		if err != nil {
			return err
		}

		commit, err := v.repo.LookupCommit(oid)
		if err != nil {
			return err
		}

		cont := fun(commit)
		if !cont {
			break
		}
	}

	return nil
}

func (v *RevWalk) SimplifyFirstParent() {
	C.git_revwalk_simplify_first_parent(v.ptr)
}

func (v *RevWalk) Sorting(sm SortType) {
	C.git_revwalk_sorting(v.ptr, C.uint(sm))
}

func (v *RevWalk) Free() {

	runtime.SetFinalizer(v, nil)
	C.git_revwalk_free(v.ptr)
}
