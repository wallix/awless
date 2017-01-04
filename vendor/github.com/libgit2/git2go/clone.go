package git

/*
#include <git2.h>

extern void _go_git_populate_remote_cb(git_clone_options *opts);
*/
import "C"
import (
	"runtime"
	"unsafe"
)

type RemoteCreateCallback func(repo *Repository, name, url string) (*Remote, ErrorCode)

type CloneOptions struct {
	*CheckoutOpts
	*FetchOptions
	Bare                 bool
	CheckoutBranch       string
	RemoteCreateCallback RemoteCreateCallback
}

func Clone(url string, path string, options *CloneOptions) (*Repository, error) {
	curl := C.CString(url)
	defer C.free(unsafe.Pointer(curl))

	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	copts := (*C.git_clone_options)(C.calloc(1, C.size_t(unsafe.Sizeof(C.git_clone_options{}))))
	populateCloneOptions(copts, options)
	defer freeCloneOptions(copts)

	if len(options.CheckoutBranch) != 0 {
		copts.checkout_branch = C.CString(options.CheckoutBranch)
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var ptr *C.git_repository
	ret := C.git_clone(&ptr, curl, cpath, copts)

	if ret < 0 {
		return nil, MakeGitError(ret)
	}

	return newRepositoryFromC(ptr), nil
}

//export remoteCreateCallback
func remoteCreateCallback(cremote unsafe.Pointer, crepo unsafe.Pointer, cname, curl *C.char, payload unsafe.Pointer) C.int {
	name := C.GoString(cname)
	url := C.GoString(curl)
	repo := newRepositoryFromC((*C.git_repository)(crepo))
	// We don't own this repository, so make sure we don't try to free it
	runtime.SetFinalizer(repo, nil)

	if opts, ok := pointerHandles.Get(payload).(CloneOptions); ok {
		remote, err := opts.RemoteCreateCallback(repo, name, url)
		// clear finalizer as the calling C function will
		// free the remote itself
		runtime.SetFinalizer(remote, nil)

		if err == ErrOk && remote != nil {
			cptr := (**C.git_remote)(cremote)
			*cptr = remote.ptr
		} else if err == ErrOk && remote == nil {
			panic("no remote created by callback")
		}

		return C.int(err)
	} else {
		panic("invalid remote create callback")
	}
}

func populateCloneOptions(ptr *C.git_clone_options, opts *CloneOptions) {
	C.git_clone_init_options(ptr, C.GIT_CLONE_OPTIONS_VERSION)

	if opts == nil {
		return
	}
	populateCheckoutOpts(&ptr.checkout_opts, opts.CheckoutOpts)
	populateFetchOptions(&ptr.fetch_opts, opts.FetchOptions)
	ptr.bare = cbool(opts.Bare)

	if opts.RemoteCreateCallback != nil {
		// Go v1.1 does not allow to assign a C function pointer
		C._go_git_populate_remote_cb(ptr)
		ptr.remote_cb_payload = pointerHandles.Track(*opts)
	}
}

func freeCloneOptions(ptr *C.git_clone_options) {
	if ptr == nil {
		return
	}

	freeCheckoutOpts(&ptr.checkout_opts)

	if ptr.remote_cb_payload != nil {
		pointerHandles.Untrack(ptr.remote_cb_payload)
	}

	C.free(unsafe.Pointer(ptr.checkout_branch))
	C.free(unsafe.Pointer(ptr))
}
