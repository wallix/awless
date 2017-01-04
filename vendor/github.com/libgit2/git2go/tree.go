package git

/*
#include <git2.h>

extern int _go_git_treewalk(git_tree *tree, git_treewalk_mode mode, void *ptr);
*/
import "C"

import (
	"runtime"
	"unsafe"
)

type Filemode int

const (
	FilemodeTree           Filemode = C.GIT_FILEMODE_TREE
	FilemodeBlob           Filemode = C.GIT_FILEMODE_BLOB
	FilemodeBlobExecutable Filemode = C.GIT_FILEMODE_BLOB_EXECUTABLE
	FilemodeLink           Filemode = C.GIT_FILEMODE_LINK
	FilemodeCommit         Filemode = C.GIT_FILEMODE_COMMIT
)

type Tree struct {
	Object
	cast_ptr *C.git_tree
}

type TreeEntry struct {
	Name     string
	Id       *Oid
	Type     ObjectType
	Filemode Filemode
}

func newTreeEntry(entry *C.git_tree_entry) *TreeEntry {
	return &TreeEntry{
		C.GoString(C.git_tree_entry_name(entry)),
		newOidFromC(C.git_tree_entry_id(entry)),
		ObjectType(C.git_tree_entry_type(entry)),
		Filemode(C.git_tree_entry_filemode(entry)),
	}
}

func (t Tree) EntryByName(filename string) *TreeEntry {
	cname := C.CString(filename)
	defer C.free(unsafe.Pointer(cname))

	entry := C.git_tree_entry_byname(t.cast_ptr, cname)
	if entry == nil {
		return nil
	}

	return newTreeEntry(entry)
}

// EntryById performs a lookup for a tree entry with the given SHA value.
//
// It returns a *TreeEntry that is owned by the Tree. You don't have to
// free it, but you must not use it after the Tree is freed.
//
// Warning: this must examine every entry in the tree, so it is not fast.
func (t Tree) EntryById(id *Oid) *TreeEntry {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	entry := C.git_tree_entry_byid(t.cast_ptr, id.toC())
	if entry == nil {
		return nil
	}

	return newTreeEntry(entry)
}

// EntryByPath looks up an entry by its full path, recursing into
// deeper trees if necessary (i.e. if there are slashes in the path)
func (t Tree) EntryByPath(path string) (*TreeEntry, error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	var entry *C.git_tree_entry

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_tree_entry_bypath(&entry, t.cast_ptr, cpath)
	if ret < 0 {
		return nil, MakeGitError(ret)
	}

	return newTreeEntry(entry), nil
}

func (t Tree) EntryByIndex(index uint64) *TreeEntry {
	entry := C.git_tree_entry_byindex(t.cast_ptr, C.size_t(index))
	if entry == nil {
		return nil
	}

	return newTreeEntry(entry)
}

func (t Tree) EntryCount() uint64 {
	num := C.git_tree_entrycount(t.cast_ptr)
	return uint64(num)
}

type TreeWalkCallback func(string, *TreeEntry) int

//export CallbackGitTreeWalk
func CallbackGitTreeWalk(_root *C.char, _entry unsafe.Pointer, ptr unsafe.Pointer) C.int {
	root := C.GoString(_root)
	entry := (*C.git_tree_entry)(_entry)

	if callback, ok := pointerHandles.Get(ptr).(TreeWalkCallback); ok {
		return C.int(callback(root, newTreeEntry(entry)))
	} else {
		panic("invalid treewalk callback")
	}
}

func (t Tree) Walk(callback TreeWalkCallback) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ptr := pointerHandles.Track(callback)
	defer pointerHandles.Untrack(ptr)

	err := C._go_git_treewalk(
		t.cast_ptr,
		C.GIT_TREEWALK_PRE,
		ptr,
	)

	if err < 0 {
		return MakeGitError(err)
	}

	return nil
}

type TreeBuilder struct {
	ptr  *C.git_treebuilder
	repo *Repository
}

func (v *TreeBuilder) Free() {
	runtime.SetFinalizer(v, nil)
	C.git_treebuilder_free(v.ptr)
}

func (v *TreeBuilder) Insert(filename string, id *Oid, filemode Filemode) error {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	err := C.git_treebuilder_insert(nil, v.ptr, cfilename, id.toC(), C.git_filemode_t(filemode))
	if err < 0 {
		return MakeGitError(err)
	}

	return nil
}

func (v *TreeBuilder) Remove(filename string) error {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	err := C.git_treebuilder_remove(v.ptr, cfilename)
	if err < 0 {
		return MakeGitError(err)
	}

	return nil
}

func (v *TreeBuilder) Write() (*Oid, error) {
	oid := new(Oid)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	err := C.git_treebuilder_write(oid.toC(), v.ptr)

	if err < 0 {
		return nil, MakeGitError(err)
	}

	return oid, nil
}
