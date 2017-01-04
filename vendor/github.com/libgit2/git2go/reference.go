package git

/*
#include <git2.h>
*/
import "C"
import (
	"runtime"
	"unsafe"
)

type ReferenceType int

const (
	ReferenceSymbolic ReferenceType = C.GIT_REF_SYMBOLIC
	ReferenceOid      ReferenceType = C.GIT_REF_OID
)

type Reference struct {
	ptr  *C.git_reference
	repo *Repository
}

type ReferenceCollection struct {
	repo *Repository
}

func (c *ReferenceCollection) Lookup(name string) (*Reference, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	var ptr *C.git_reference

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ecode := C.git_reference_lookup(&ptr, c.repo.ptr, cname)
	if ecode < 0 {
		return nil, MakeGitError(ecode)
	}

	return newReferenceFromC(ptr, c.repo), nil
}

func (c *ReferenceCollection) Create(name string, id *Oid, force bool, msg string) (*Reference, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var cmsg *C.char
	if msg == "" {
		cmsg = nil
	} else {
		cmsg = C.CString(msg)
		defer C.free(unsafe.Pointer(cmsg))
	}

	var ptr *C.git_reference

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ecode := C.git_reference_create(&ptr, c.repo.ptr, cname, id.toC(), cbool(force), cmsg)
	if ecode < 0 {
		return nil, MakeGitError(ecode)
	}

	return newReferenceFromC(ptr, c.repo), nil
}

func (c *ReferenceCollection) CreateSymbolic(name, target string, force bool, msg string) (*Reference, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	ctarget := C.CString(target)
	defer C.free(unsafe.Pointer(ctarget))

	var cmsg *C.char
	if msg == "" {
		cmsg = nil
	} else {
		cmsg = C.CString(msg)
		defer C.free(unsafe.Pointer(cmsg))
	}

	var ptr *C.git_reference

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ecode := C.git_reference_symbolic_create(&ptr, c.repo.ptr, cname, ctarget, cbool(force), cmsg)
	if ecode < 0 {
		return nil, MakeGitError(ecode)
	}

	return newReferenceFromC(ptr, c.repo), nil
}

// EnsureLog ensures that there is a reflog for the given reference
// name and creates an empty one if necessary.
func (c *ReferenceCollection) EnsureLog(name string) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_reference_ensure_log(c.repo.ptr, cname)
	if ret < 0 {
		return MakeGitError(ret)
	}

	return nil
}

// HasLog returns whether there is a reflog for the given reference
// name
func (c *ReferenceCollection) HasLog(name string) (bool, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_reference_has_log(c.repo.ptr, cname)
	if ret < 0 {
		return false, MakeGitError(ret)
	}

	return ret == 1, nil
}

// Dwim looks up a reference by DWIMing its short name
func (c *ReferenceCollection) Dwim(name string) (*Reference, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var ptr *C.git_reference
	ret := C.git_reference_dwim(&ptr, c.repo.ptr, cname)
	if ret < 0 {
		return nil, MakeGitError(ret)
	}

	return newReferenceFromC(ptr, c.repo), nil
}

func newReferenceFromC(ptr *C.git_reference, repo *Repository) *Reference {
	ref := &Reference{ptr: ptr, repo: repo}
	runtime.SetFinalizer(ref, (*Reference).Free)
	return ref
}

func (v *Reference) SetSymbolicTarget(target string, msg string) (*Reference, error) {
	var ptr *C.git_reference

	ctarget := C.CString(target)
	defer C.free(unsafe.Pointer(ctarget))

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var cmsg *C.char
	if msg == "" {
		cmsg = nil
	} else {
		cmsg = C.CString(msg)
		defer C.free(unsafe.Pointer(cmsg))
	}

	ret := C.git_reference_symbolic_set_target(&ptr, v.ptr, ctarget, cmsg)
	if ret < 0 {
		return nil, MakeGitError(ret)
	}

	return newReferenceFromC(ptr, v.repo), nil
}

func (v *Reference) SetTarget(target *Oid, msg string) (*Reference, error) {
	var ptr *C.git_reference

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var cmsg *C.char
	if msg == "" {
		cmsg = nil
	} else {
		cmsg = C.CString(msg)
		defer C.free(unsafe.Pointer(cmsg))
	}

	ret := C.git_reference_set_target(&ptr, v.ptr, target.toC(), cmsg)
	if ret < 0 {
		return nil, MakeGitError(ret)
	}

	return newReferenceFromC(ptr, v.repo), nil
}

func (v *Reference) Resolve() (*Reference, error) {
	var ptr *C.git_reference

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_reference_resolve(&ptr, v.ptr)
	if ret < 0 {
		return nil, MakeGitError(ret)
	}

	return newReferenceFromC(ptr, v.repo), nil
}

func (v *Reference) Rename(name string, force bool, msg string) (*Reference, error) {
	var ptr *C.git_reference
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var cmsg *C.char
	if msg == "" {
		cmsg = nil
	} else {
		cmsg = C.CString(msg)
		defer C.free(unsafe.Pointer(cmsg))
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_reference_rename(&ptr, v.ptr, cname, cbool(force), cmsg)

	if ret < 0 {
		return nil, MakeGitError(ret)
	}

	return newReferenceFromC(ptr, v.repo), nil
}

func (v *Reference) Target() *Oid {
	return newOidFromC(C.git_reference_target(v.ptr))
}

func (v *Reference) SymbolicTarget() string {
	cstr := C.git_reference_symbolic_target(v.ptr)
	if cstr == nil {
		return ""
	}

	return C.GoString(cstr)
}

func (v *Reference) Delete() error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_reference_delete(v.ptr)

	if ret < 0 {
		return MakeGitError(ret)
	}

	return nil
}

func (v *Reference) Peel(t ObjectType) (*Object, error) {
	var cobj *C.git_object

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if err := C.git_reference_peel(&cobj, v.ptr, C.git_otype(t)); err < 0 {
		return nil, MakeGitError(err)
	}

	return allocObject(cobj, v.repo), nil
}

// Owner returns a weak reference to the repository which owns this
// reference.
func (v *Reference) Owner() *Repository {
	return &Repository{
		ptr: C.git_reference_owner(v.ptr),
	}
}

// Cmp compares both references, retursn 0 on equality, otherwise a
// stable sorting.
func (v *Reference) Cmp(ref2 *Reference) int {
	return int(C.git_reference_cmp(v.ptr, ref2.ptr))
}

// Shorthand returns a "human-readable" short reference name
func (v *Reference) Shorthand() string {
	return C.GoString(C.git_reference_shorthand(v.ptr))
}

func (v *Reference) Name() string {
	return C.GoString(C.git_reference_name(v.ptr))
}

func (v *Reference) Type() ReferenceType {
	return ReferenceType(C.git_reference_type(v.ptr))
}

func (v *Reference) IsBranch() bool {
	return C.git_reference_is_branch(v.ptr) == 1
}

func (v *Reference) IsRemote() bool {
	return C.git_reference_is_remote(v.ptr) == 1
}

func (v *Reference) IsTag() bool {
	return C.git_reference_is_tag(v.ptr) == 1
}

// IsNote checks if the reference is a note.
func (v *Reference) IsNote() bool {
	return C.git_reference_is_note(v.ptr) == 1
}

func (v *Reference) Free() {
	runtime.SetFinalizer(v, nil)
	C.git_reference_free(v.ptr)
}

type ReferenceIterator struct {
	ptr  *C.git_reference_iterator
	repo *Repository
}

type ReferenceNameIterator struct {
	*ReferenceIterator
}

// NewReferenceIterator creates a new iterator over reference names
func (repo *Repository) NewReferenceIterator() (*ReferenceIterator, error) {
	var ptr *C.git_reference_iterator

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_reference_iterator_new(&ptr, repo.ptr)
	if ret < 0 {
		return nil, MakeGitError(ret)
	}

	iter := &ReferenceIterator{ptr: ptr, repo: repo}
	runtime.SetFinalizer(iter, (*ReferenceIterator).Free)
	return iter, nil
}

// NewReferenceIterator creates a new branch iterator over reference names
func (repo *Repository) NewReferenceNameIterator() (*ReferenceNameIterator, error) {
	var ptr *C.git_reference_iterator

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_reference_iterator_new(&ptr, repo.ptr)
	if ret < 0 {
		return nil, MakeGitError(ret)
	}

	iter := &ReferenceIterator{ptr: ptr, repo: repo}
	runtime.SetFinalizer(iter, (*ReferenceIterator).Free)
	return iter.Names(), nil
}

// NewReferenceIteratorGlob creates an iterator over reference names
// that match the speicified glob. The glob is of the usual fnmatch
// type.
func (repo *Repository) NewReferenceIteratorGlob(glob string) (*ReferenceIterator, error) {
	cstr := C.CString(glob)
	defer C.free(unsafe.Pointer(cstr))
	var ptr *C.git_reference_iterator

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_reference_iterator_glob_new(&ptr, repo.ptr, cstr)
	if ret < 0 {
		return nil, MakeGitError(ret)
	}

	iter := &ReferenceIterator{ptr: ptr}
	runtime.SetFinalizer(iter, (*ReferenceIterator).Free)
	return iter, nil
}

func (i *ReferenceIterator) Names() *ReferenceNameIterator {
	return &ReferenceNameIterator{i}
}

// NextName retrieves the next reference name. If the iteration is over,
// the returned error is git.ErrIterOver
func (v *ReferenceNameIterator) Next() (string, error) {
	var ptr *C.char

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_reference_next_name(&ptr, v.ptr)
	if ret < 0 {
		return "", MakeGitError(ret)
	}

	return C.GoString(ptr), nil
}

// Next retrieves the next reference. If the iterationis over, the
// returned error is git.ErrIterOver
func (v *ReferenceIterator) Next() (*Reference, error) {
	var ptr *C.git_reference

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_reference_next(&ptr, v.ptr)
	if ret < 0 {
		return nil, MakeGitError(ret)
	}

	return newReferenceFromC(ptr, v.repo), nil
}

// Free the reference iterator
func (v *ReferenceIterator) Free() {
	runtime.SetFinalizer(v, nil)
	C.git_reference_iterator_free(v.ptr)
}

// ReferenceIsValidName ensures the reference name is well-formed.
//
// Valid reference names must follow one of two patterns:
//
// 1. Top-level names must contain only capital letters and underscores,
// and must begin and end with a letter. (e.g. "HEAD", "ORIG_HEAD").
//
// 2. Names prefixed with "refs/" can be almost anything. You must avoid
// the characters '~', '^', ':', ' \ ', '?', '[', and '*', and the sequences
// ".." and " @ {" which have special meaning to revparse.
func ReferenceIsValidName(name string) bool {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	if C.git_reference_is_valid_name(cname) == 1 {
		return true
	}
	return false
}
