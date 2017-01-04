package git

/*
#include <git2.h>
*/
import "C"

import (
	"runtime"
	"unsafe"
)

type BranchType uint

const (
	BranchAll    BranchType = C.GIT_BRANCH_ALL
	BranchLocal  BranchType = C.GIT_BRANCH_LOCAL
	BranchRemote BranchType = C.GIT_BRANCH_REMOTE
)

type Branch struct {
	*Reference
}

func (r *Reference) Branch() *Branch {
	return &Branch{Reference: r}
}

type BranchIterator struct {
	ptr  *C.git_branch_iterator
	repo *Repository
}

type BranchIteratorFunc func(*Branch, BranchType) error

func newBranchIteratorFromC(repo *Repository, ptr *C.git_branch_iterator) *BranchIterator {
	i := &BranchIterator{repo: repo, ptr: ptr}
	runtime.SetFinalizer(i, (*BranchIterator).Free)
	return i
}

func (i *BranchIterator) Next() (*Branch, BranchType, error) {

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var refPtr *C.git_reference
	var refType C.git_branch_t

	ecode := C.git_branch_next(&refPtr, &refType, i.ptr)

	if ecode < 0 {
		return nil, BranchLocal, MakeGitError(ecode)
	}

	branch := newReferenceFromC(refPtr, i.repo).Branch()

	return branch, BranchType(refType), nil
}

func (i *BranchIterator) Free() {
	runtime.SetFinalizer(i, nil)
	C.git_branch_iterator_free(i.ptr)
}

func (i *BranchIterator) ForEach(f BranchIteratorFunc) error {
	b, t, err := i.Next()

	for err == nil {
		err = f(b, t)
		if err == nil {
			b, t, err = i.Next()
		}
	}

	if err != nil && IsErrorCode(err, ErrIterOver) {
		return nil
	}

	return err
}

func (repo *Repository) NewBranchIterator(flags BranchType) (*BranchIterator, error) {
	refType := C.git_branch_t(flags)
	var ptr *C.git_branch_iterator

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ecode := C.git_branch_iterator_new(&ptr, repo.ptr, refType)
	if ecode < 0 {
		return nil, MakeGitError(ecode)
	}

	return newBranchIteratorFromC(repo, ptr), nil
}

func (repo *Repository) CreateBranch(branchName string, target *Commit, force bool) (*Branch, error) {

	var ptr *C.git_reference
	cBranchName := C.CString(branchName)
	defer C.free(unsafe.Pointer(cBranchName))
	cForce := cbool(force)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_branch_create(&ptr, repo.ptr, cBranchName, target.cast_ptr, cForce)
	if ret < 0 {
		return nil, MakeGitError(ret)
	}
	return newReferenceFromC(ptr, repo).Branch(), nil
}

func (b *Branch) Delete() error {

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	ret := C.git_branch_delete(b.Reference.ptr)
	if ret < 0 {
		return MakeGitError(ret)
	}
	return nil
}

func (b *Branch) Move(newBranchName string, force bool) (*Branch, error) {
	var ptr *C.git_reference
	cNewBranchName := C.CString(newBranchName)
	defer C.free(unsafe.Pointer(cNewBranchName))
	cForce := cbool(force)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_branch_move(&ptr, b.Reference.ptr, cNewBranchName, cForce)
	if ret < 0 {
		return nil, MakeGitError(ret)
	}
	return newReferenceFromC(ptr, b.repo).Branch(), nil
}

func (b *Branch) IsHead() (bool, error) {

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_branch_is_head(b.Reference.ptr)
	switch ret {
	case 1:
		return true, nil
	case 0:
		return false, nil
	}
	return false, MakeGitError(ret)

}

func (repo *Repository) LookupBranch(branchName string, bt BranchType) (*Branch, error) {
	var ptr *C.git_reference

	cName := C.CString(branchName)
	defer C.free(unsafe.Pointer(cName))

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_branch_lookup(&ptr, repo.ptr, cName, C.git_branch_t(bt))
	if ret < 0 {
		return nil, MakeGitError(ret)
	}
	return newReferenceFromC(ptr, repo).Branch(), nil
}

func (b *Branch) Name() (string, error) {
	var cName *C.char
	defer C.free(unsafe.Pointer(cName))

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_branch_name(&cName, b.Reference.ptr)
	if ret < 0 {
		return "", MakeGitError(ret)
	}

	return C.GoString(cName), nil
}

func (repo *Repository) RemoteName(canonicalBranchName string) (string, error) {
	cName := C.CString(canonicalBranchName)
	defer C.free(unsafe.Pointer(cName))

	nameBuf := C.git_buf{}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_branch_remote_name(&nameBuf, repo.ptr, cName)
	if ret < 0 {
		return "", MakeGitError(ret)
	}
	defer C.git_buf_free(&nameBuf)

	return C.GoString(nameBuf.ptr), nil
}

func (b *Branch) SetUpstream(upstreamName string) error {
	cName := C.CString(upstreamName)
	defer C.free(unsafe.Pointer(cName))

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_branch_set_upstream(b.Reference.ptr, cName)
	if ret < 0 {
		return MakeGitError(ret)
	}
	return nil
}

func (b *Branch) Upstream() (*Reference, error) {

	var ptr *C.git_reference
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_branch_upstream(&ptr, b.Reference.ptr)
	if ret < 0 {
		return nil, MakeGitError(ret)
	}
	return newReferenceFromC(ptr, b.repo), nil
}

func (repo *Repository) UpstreamName(canonicalBranchName string) (string, error) {
	cName := C.CString(canonicalBranchName)
	defer C.free(unsafe.Pointer(cName))

	nameBuf := C.git_buf{}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_branch_upstream_name(&nameBuf, repo.ptr, cName)
	if ret < 0 {
		return "", MakeGitError(ret)
	}
	defer C.git_buf_free(&nameBuf)

	return C.GoString(nameBuf.ptr), nil
}
