package git

/*
#include <git2.h>
#include <git2/pack.h>
#include <stdlib.h>

extern int _go_git_packbuilder_foreach(git_packbuilder *pb, void *payload);
*/
import "C"
import (
	"io"
	"os"
	"runtime"
	"unsafe"
)

type Packbuilder struct {
	ptr *C.git_packbuilder
}

func (repo *Repository) NewPackbuilder() (*Packbuilder, error) {
	builder := &Packbuilder{}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_packbuilder_new(&builder.ptr, repo.ptr)
	if ret != 0 {
		return nil, MakeGitError(ret)
	}
	runtime.SetFinalizer(builder, (*Packbuilder).Free)
	return builder, nil
}

func (pb *Packbuilder) Free() {
	runtime.SetFinalizer(pb, nil)
	C.git_packbuilder_free(pb.ptr)
}

func (pb *Packbuilder) Insert(id *Oid, name string) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_packbuilder_insert(pb.ptr, id.toC(), cname)
	if ret != 0 {
		return MakeGitError(ret)
	}
	return nil
}

func (pb *Packbuilder) InsertCommit(id *Oid) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_packbuilder_insert_commit(pb.ptr, id.toC())
	if ret != 0 {
		return MakeGitError(ret)
	}
	return nil
}

func (pb *Packbuilder) InsertTree(id *Oid) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_packbuilder_insert_tree(pb.ptr, id.toC())
	if ret != 0 {
		return MakeGitError(ret)
	}
	return nil
}

func (pb *Packbuilder) ObjectCount() uint32 {
	return uint32(C.git_packbuilder_object_count(pb.ptr))
}

func (pb *Packbuilder) WriteToFile(name string, mode os.FileMode) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_packbuilder_write(pb.ptr, cname, C.uint(mode.Perm()), nil, nil)
	if ret != 0 {
		return MakeGitError(ret)
	}
	return nil
}

func (pb *Packbuilder) Write(w io.Writer) error {
	return pb.ForEach(func(slice []byte) error {
		_, err := w.Write(slice)
		return err
	})
}

func (pb *Packbuilder) Written() uint32 {
	return uint32(C.git_packbuilder_written(pb.ptr))
}

type PackbuilderForeachCallback func([]byte) error
type packbuilderCbData struct {
	callback PackbuilderForeachCallback
	err      error
}

//export packbuilderForEachCb
func packbuilderForEachCb(buf unsafe.Pointer, size C.size_t, handle unsafe.Pointer) int {
	payload := pointerHandles.Get(handle)
	data, ok := payload.(*packbuilderCbData)
	if !ok {
		panic("could not get packbuilder CB data")
	}

	slice := C.GoBytes(buf, C.int(size))

	err := data.callback(slice)
	if err != nil {
		data.err = err
		return C.GIT_EUSER
	}

	return 0
}

// ForEach repeatedly calls the callback with new packfile data until
// there is no more data or the callback returns an error
func (pb *Packbuilder) ForEach(callback PackbuilderForeachCallback) error {
	data := packbuilderCbData{
		callback: callback,
		err:      nil,
	}
	handle := pointerHandles.Track(&data)
	defer pointerHandles.Untrack(handle)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	err := C._go_git_packbuilder_foreach(pb.ptr, handle)
	if err == C.GIT_EUSER {
		return data.err
	}
	if err < 0 {
		return MakeGitError(err)
	}

	return nil
}
