package git

/*
#include <git2.h>

extern int _go_git_visit_submodule(git_repository *repo, void *fct);
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// SubmoduleUpdateOptions
type SubmoduleUpdateOptions struct {
	*CheckoutOpts
	*FetchOptions
	CloneCheckoutStrategy CheckoutStrategy
}

// Submodule
type Submodule struct {
	ptr *C.git_submodule
}

type SubmoduleUpdate int

const (
	SubmoduleUpdateCheckout SubmoduleUpdate = C.GIT_SUBMODULE_UPDATE_CHECKOUT
	SubmoduleUpdateRebase   SubmoduleUpdate = C.GIT_SUBMODULE_UPDATE_REBASE
	SubmoduleUpdateMerge    SubmoduleUpdate = C.GIT_SUBMODULE_UPDATE_MERGE
	SubmoduleUpdateNone     SubmoduleUpdate = C.GIT_SUBMODULE_UPDATE_NONE
)

type SubmoduleIgnore int

const (
	SubmoduleIgnoreNone      SubmoduleIgnore = C.GIT_SUBMODULE_IGNORE_NONE
	SubmoduleIgnoreUntracked SubmoduleIgnore = C.GIT_SUBMODULE_IGNORE_UNTRACKED
	SubmoduleIgnoreDirty     SubmoduleIgnore = C.GIT_SUBMODULE_IGNORE_DIRTY
	SubmoduleIgnoreAll       SubmoduleIgnore = C.GIT_SUBMODULE_IGNORE_ALL
)

type SubmoduleStatus int

const (
	SubmoduleStatusInHead          SubmoduleStatus = C.GIT_SUBMODULE_STATUS_IN_HEAD
	SubmoduleStatusInIndex         SubmoduleStatus = C.GIT_SUBMODULE_STATUS_IN_INDEX
	SubmoduleStatusInConfig        SubmoduleStatus = C.GIT_SUBMODULE_STATUS_IN_CONFIG
	SubmoduleStatusInWd            SubmoduleStatus = C.GIT_SUBMODULE_STATUS_IN_WD
	SubmoduleStatusIndexAdded      SubmoduleStatus = C.GIT_SUBMODULE_STATUS_INDEX_ADDED
	SubmoduleStatusIndexDeleted    SubmoduleStatus = C.GIT_SUBMODULE_STATUS_INDEX_DELETED
	SubmoduleStatusIndexModified   SubmoduleStatus = C.GIT_SUBMODULE_STATUS_INDEX_MODIFIED
	SubmoduleStatusWdUninitialized SubmoduleStatus = C.GIT_SUBMODULE_STATUS_WD_UNINITIALIZED
	SubmoduleStatusWdAdded         SubmoduleStatus = C.GIT_SUBMODULE_STATUS_WD_ADDED
	SubmoduleStatusWdDeleted       SubmoduleStatus = C.GIT_SUBMODULE_STATUS_WD_DELETED
	SubmoduleStatusWdModified      SubmoduleStatus = C.GIT_SUBMODULE_STATUS_WD_MODIFIED
	SubmoduleStatusWdIndexModified SubmoduleStatus = C.GIT_SUBMODULE_STATUS_WD_INDEX_MODIFIED
	SubmoduleStatusWdWdModified    SubmoduleStatus = C.GIT_SUBMODULE_STATUS_WD_WD_MODIFIED
	SubmoduleStatusWdUntracked     SubmoduleStatus = C.GIT_SUBMODULE_STATUS_WD_UNTRACKED
)

type SubmoduleRecurse int

const (
	SubmoduleRecurseNo       SubmoduleRecurse = C.GIT_SUBMODULE_RECURSE_NO
	SubmoduleRecurseYes      SubmoduleRecurse = C.GIT_SUBMODULE_RECURSE_YES
	SubmoduleRecurseOndemand SubmoduleRecurse = C.GIT_SUBMODULE_RECURSE_ONDEMAND
)

type SubmoduleCollection struct {
	repo *Repository
}

func SubmoduleStatusIsUnmodified(status int) bool {
	o := SubmoduleStatus(status) & ^(SubmoduleStatusInHead | SubmoduleStatusInIndex |
		SubmoduleStatusInConfig | SubmoduleStatusInWd)
	return o == 0
}

func (c *SubmoduleCollection) Lookup(name string) (*Submodule, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	sub := new(Submodule)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_submodule_lookup(&sub.ptr, c.repo.ptr, cname)
	if ret < 0 {
		return nil, MakeGitError(ret)
	}

	return sub, nil
}

type SubmoduleCbk func(sub *Submodule, name string) int

//export SubmoduleVisitor
func SubmoduleVisitor(csub unsafe.Pointer, name *C.char, handle unsafe.Pointer) C.int {
	sub := &Submodule{(*C.git_submodule)(csub)}

	if callback, ok := pointerHandles.Get(handle).(SubmoduleCbk); ok {
		return (C.int)(callback(sub, C.GoString(name)))
	} else {
		panic("invalid submodule visitor callback")
	}
}

func (c *SubmoduleCollection) Foreach(cbk SubmoduleCbk) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	handle := pointerHandles.Track(cbk)
	defer pointerHandles.Untrack(handle)

	ret := C._go_git_visit_submodule(c.repo.ptr, handle)
	if ret < 0 {
		return MakeGitError(ret)
	}
	return nil
}

func (c *SubmoduleCollection) Add(url, path string, use_git_link bool) (*Submodule, error) {
	curl := C.CString(url)
	defer C.free(unsafe.Pointer(curl))
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	sub := new(Submodule)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_submodule_add_setup(&sub.ptr, c.repo.ptr, curl, cpath, cbool(use_git_link))
	if ret < 0 {
		return nil, MakeGitError(ret)
	}
	return sub, nil
}

func (sub *Submodule) FinalizeAdd() error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_submodule_add_finalize(sub.ptr)
	if ret < 0 {
		return MakeGitError(ret)
	}
	return nil
}

func (sub *Submodule) AddToIndex(write_index bool) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_submodule_add_to_index(sub.ptr, cbool(write_index))
	if ret < 0 {
		return MakeGitError(ret)
	}
	return nil
}

func (sub *Submodule) Name() string {
	n := C.git_submodule_name(sub.ptr)
	return C.GoString(n)
}

func (sub *Submodule) Path() string {
	n := C.git_submodule_path(sub.ptr)
	return C.GoString(n)
}

func (sub *Submodule) Url() string {
	n := C.git_submodule_url(sub.ptr)
	return C.GoString(n)
}

func (c *SubmoduleCollection) SetUrl(submodule, url string) error {
	csubmodule := C.CString(submodule)
	defer C.free(unsafe.Pointer(csubmodule))
	curl := C.CString(url)
	defer C.free(unsafe.Pointer(curl))

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_submodule_set_url(c.repo.ptr, csubmodule, curl)
	if ret < 0 {
		return MakeGitError(ret)
	}
	return nil
}

func (sub *Submodule) IndexId() *Oid {
	idx := C.git_submodule_index_id(sub.ptr)
	if idx == nil {
		return nil
	}
	return newOidFromC(idx)
}

func (sub *Submodule) HeadId() *Oid {
	idx := C.git_submodule_head_id(sub.ptr)
	if idx == nil {
		return nil
	}
	return newOidFromC(idx)
}

func (sub *Submodule) WdId() *Oid {
	idx := C.git_submodule_wd_id(sub.ptr)
	if idx == nil {
		return nil
	}
	return newOidFromC(idx)
}

func (sub *Submodule) Ignore() SubmoduleIgnore {
	o := C.git_submodule_ignore(sub.ptr)
	return SubmoduleIgnore(o)
}

func (c *SubmoduleCollection) SetIgnore(submodule string, ignore SubmoduleIgnore) error {
	csubmodule := C.CString(submodule)
	defer C.free(unsafe.Pointer(csubmodule))

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_submodule_set_ignore(c.repo.ptr, csubmodule, C.git_submodule_ignore_t(ignore))
	if ret < 0 {
		return MakeGitError(ret)
	}

	return nil
}

func (sub *Submodule) UpdateStrategy() SubmoduleUpdate {
	o := C.git_submodule_update_strategy(sub.ptr)
	return SubmoduleUpdate(o)
}

func (c *SubmoduleCollection) SetUpdate(submodule string, update SubmoduleUpdate) error {
	csubmodule := C.CString(submodule)
	defer C.free(unsafe.Pointer(csubmodule))

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_submodule_set_update(c.repo.ptr, csubmodule, C.git_submodule_update_t(update))
	if ret < 0 {
		return MakeGitError(ret)
	}

	return nil
}

func (sub *Submodule) FetchRecurseSubmodules() SubmoduleRecurse {
	return SubmoduleRecurse(C.git_submodule_fetch_recurse_submodules(sub.ptr))
}

func (c *SubmoduleCollection) SetFetchRecurseSubmodules(submodule string, recurse SubmoduleRecurse) error {
	csubmodule := C.CString(submodule)
	defer C.free(unsafe.Pointer(csubmodule))

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_submodule_set_fetch_recurse_submodules(c.repo.ptr, csubmodule, C.git_submodule_recurse_t(recurse))
	if ret < 0 {
		return MakeGitError(C.int(ret))
	}
	return nil
}

func (sub *Submodule) Init(overwrite bool) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_submodule_init(sub.ptr, cbool(overwrite))
	if ret < 0 {
		return MakeGitError(ret)
	}
	return nil
}

func (sub *Submodule) Sync() error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_submodule_sync(sub.ptr)
	if ret < 0 {
		return MakeGitError(ret)
	}
	return nil
}

func (sub *Submodule) Open() (*Repository, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var ptr *C.git_repository
	ret := C.git_submodule_open(&ptr, sub.ptr)
	if ret < 0 {
		return nil, MakeGitError(ret)
	}
	return newRepositoryFromC(ptr), nil
}

func (sub *Submodule) Update(init bool, opts *SubmoduleUpdateOptions) error {
	var copts C.git_submodule_update_options
	err := populateSubmoduleUpdateOptions(&copts, opts)
	if err != nil {
		return err
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_submodule_update(sub.ptr, cbool(init), &copts)
	if ret < 0 {
		return MakeGitError(ret)
	}

	return nil
}

func populateSubmoduleUpdateOptions(ptr *C.git_submodule_update_options, opts *SubmoduleUpdateOptions) error {
	C.git_submodule_update_init_options(ptr, C.GIT_SUBMODULE_UPDATE_OPTIONS_VERSION)

	if opts == nil {
		return nil
	}

	populateCheckoutOpts(&ptr.checkout_opts, opts.CheckoutOpts)
	populateFetchOptions(&ptr.fetch_opts, opts.FetchOptions)
	ptr.clone_checkout_strategy = C.uint(opts.CloneCheckoutStrategy)

	return nil
}
