package git

/*
#include <git2.h>
*/
import "C"
import (
	"runtime"
	"time"
	"unsafe"
)

type Signature struct {
	Name  string
	Email string
	When  time.Time
}

func newSignatureFromC(sig *C.git_signature) *Signature {
	// git stores minutes, go wants seconds
	loc := time.FixedZone("", int(sig.when.offset)*60)
	return &Signature{
		C.GoString(sig.name),
		C.GoString(sig.email),
		time.Unix(int64(sig.when.time), 0).In(loc),
	}
}

// the offset in mintes, which is what git wants
func (v *Signature) Offset() int {
	_, offset := v.When.Zone()
	return offset / 60
}

func (sig *Signature) toC() (*C.git_signature, error) {
	if sig == nil {
		return nil, nil
	}

	var out *C.git_signature

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	name := C.CString(sig.Name)
	defer C.free(unsafe.Pointer(name))

	email := C.CString(sig.Email)
	defer C.free(unsafe.Pointer(email))

	ret := C.git_signature_new(&out, name, email, C.git_time_t(sig.When.Unix()), C.int(sig.Offset()))
	if ret < 0 {
		return nil, MakeGitError(ret)
	}

	return out, nil
}

func (repo *Repository) DefaultSignature() (*Signature, error) {
	var out *C.git_signature

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	cErr := C.git_signature_default(&out, repo.ptr)
	if cErr < 0 {
		return nil, MakeGitError(cErr)
	}

	defer C.git_signature_free(out)

	return newSignatureFromC(out), nil
}
