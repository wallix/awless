package git

/*
#include <git2.h>

extern void _go_git_revspec_free(git_revspec *revspec);
*/
import "C"
import (
	"runtime"
	"unsafe"
)

type RevparseFlag int

const (
	RevparseSingle    RevparseFlag = C.GIT_REVPARSE_SINGLE
	RevparseRange     RevparseFlag = C.GIT_REVPARSE_RANGE
	RevparseMergeBase RevparseFlag = C.GIT_REVPARSE_MERGE_BASE
)

type Revspec struct {
	to    *Object
	from  *Object
	flags RevparseFlag
}

func (rs *Revspec) To() *Object {
	return rs.to
}

func (rs *Revspec) From() *Object {
	return rs.from
}

func (rs *Revspec) Flags() RevparseFlag {
	return rs.flags
}

func newRevspecFromC(ptr *C.git_revspec, repo *Repository) *Revspec {
	var to *Object
	var from *Object

	if ptr.to != nil {
		to = allocObject(ptr.to, repo)
	}

	if ptr.from != nil {
		from = allocObject(ptr.from, repo)
	}

	return &Revspec{
		to:    to,
		from:  from,
		flags: RevparseFlag(ptr.flags),
	}
}

func (r *Repository) Revparse(spec string) (*Revspec, error) {
	cspec := C.CString(spec)
	defer C.free(unsafe.Pointer(cspec))

	var crevspec C.git_revspec

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ecode := C.git_revparse(&crevspec, r.ptr, cspec)
	if ecode != 0 {
		return nil, MakeGitError(ecode)
	}

	return newRevspecFromC(&crevspec, r), nil
}

func (v *Repository) RevparseSingle(spec string) (*Object, error) {
	cspec := C.CString(spec)
	defer C.free(unsafe.Pointer(cspec))

	var ptr *C.git_object

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ecode := C.git_revparse_single(&ptr, v.ptr, cspec)
	if ecode < 0 {
		return nil, MakeGitError(ecode)
	}

	return allocObject(ptr, v), nil
}

func (r *Repository) RevparseExt(spec string) (*Object, *Reference, error) {
	cspec := C.CString(spec)
	defer C.free(unsafe.Pointer(cspec))

	var obj *C.git_object
	var ref *C.git_reference

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ecode := C.git_revparse_ext(&obj, &ref, r.ptr, cspec)
	if ecode != 0 {
		return nil, nil, MakeGitError(ecode)
	}

	if ref == nil {
		return allocObject(obj, r), nil, nil
	}

	return allocObject(obj, r), newReferenceFromC(ref, r), nil
}
