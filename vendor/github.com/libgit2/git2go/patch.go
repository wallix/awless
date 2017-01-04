package git

/*
#include <git2.h>
*/
import "C"
import (
	"runtime"
	"unsafe"
)

type Patch struct {
	ptr *C.git_patch
}

func newPatchFromC(ptr *C.git_patch) *Patch {
	if ptr == nil {
		return nil
	}

	patch := &Patch{
		ptr: ptr,
	}

	runtime.SetFinalizer(patch, (*Patch).Free)
	return patch
}

func (patch *Patch) Free() error {
	if patch.ptr == nil {
		return ErrInvalid
	}
	runtime.SetFinalizer(patch, nil)
	C.git_patch_free(patch.ptr)
	patch.ptr = nil
	return nil
}

func (patch *Patch) String() (string, error) {
	if patch.ptr == nil {
		return "", ErrInvalid
	}
	var buf C.git_buf

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ecode := C.git_patch_to_buf(&buf, patch.ptr)
	if ecode < 0 {
		return "", MakeGitError(ecode)
	}
	return C.GoString(buf.ptr), nil
}

func toPointer(data []byte) (ptr unsafe.Pointer) {
	if len(data) > 0 {
		ptr = unsafe.Pointer(&data[0])
	} else {
		ptr = unsafe.Pointer(nil)
	}
	return
}

func (v *Repository) PatchFromBuffers(oldPath, newPath string, oldBuf, newBuf []byte, opts *DiffOptions) (*Patch, error) {
	var patchPtr *C.git_patch

	oldPtr := toPointer(oldBuf)
	newPtr := (*C.char)(toPointer(newBuf))

	cOldPath := C.CString(oldPath)
	defer C.free(unsafe.Pointer(cOldPath))

	cNewPath := C.CString(newPath)
	defer C.free(unsafe.Pointer(cNewPath))

	copts, _ := diffOptionsToC(opts)
	defer freeDiffOptions(copts)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ecode := C.git_patch_from_buffers(&patchPtr, oldPtr, C.size_t(len(oldBuf)), cOldPath, newPtr, C.size_t(len(newBuf)), cNewPath, copts)
	if ecode < 0 {
		return nil, MakeGitError(ecode)
	}
	return newPatchFromC(patchPtr), nil
}
