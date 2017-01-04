package git

/*
#include <git2.h>

extern void _go_git_populate_checkout_cb(git_checkout_options *opts);
*/
import "C"
import (
	"os"
	"runtime"
	"unsafe"
)

type CheckoutNotifyType uint
type CheckoutStrategy uint

const (
	CheckoutNotifyNone      CheckoutNotifyType = C.GIT_CHECKOUT_NOTIFY_NONE
	CheckoutNotifyConflict  CheckoutNotifyType = C.GIT_CHECKOUT_NOTIFY_CONFLICT
	CheckoutNotifyDirty     CheckoutNotifyType = C.GIT_CHECKOUT_NOTIFY_DIRTY
	CheckoutNotifyUpdated   CheckoutNotifyType = C.GIT_CHECKOUT_NOTIFY_UPDATED
	CheckoutNotifyUntracked CheckoutNotifyType = C.GIT_CHECKOUT_NOTIFY_UNTRACKED
	CheckoutNotifyIgnored   CheckoutNotifyType = C.GIT_CHECKOUT_NOTIFY_IGNORED
	CheckoutNotifyAll       CheckoutNotifyType = C.GIT_CHECKOUT_NOTIFY_ALL

	CheckoutNone                      CheckoutStrategy = C.GIT_CHECKOUT_NONE                         // Dry run, no actual updates
	CheckoutSafe                      CheckoutStrategy = C.GIT_CHECKOUT_SAFE                         // Allow safe updates that cannot overwrite uncommitted data
	CheckoutForce                     CheckoutStrategy = C.GIT_CHECKOUT_FORCE                        // Allow all updates to force working directory to look like index
	CheckoutRecreateMissing           CheckoutStrategy = C.GIT_CHECKOUT_RECREATE_MISSING             // Allow checkout to recreate missing files
	CheckoutAllowConflicts            CheckoutStrategy = C.GIT_CHECKOUT_ALLOW_CONFLICTS              // Allow checkout to make safe updates even if conflicts are found
	CheckoutRemoveUntracked           CheckoutStrategy = C.GIT_CHECKOUT_REMOVE_UNTRACKED             // Remove untracked files not in index (that are not ignored)
	CheckoutRemoveIgnored             CheckoutStrategy = C.GIT_CHECKOUT_REMOVE_IGNORED               // Remove ignored files not in index
	CheckoutUpdateOnly                CheckoutStrategy = C.GIT_CHECKOUT_UPDATE_ONLY                  // Only update existing files, don't create new ones
	CheckoutDontUpdateIndex           CheckoutStrategy = C.GIT_CHECKOUT_DONT_UPDATE_INDEX            // Normally checkout updates index entries as it goes; this stops that
	CheckoutNoRefresh                 CheckoutStrategy = C.GIT_CHECKOUT_NO_REFRESH                   // Don't refresh index/config/etc before doing checkout
	CheckoutSkipUnmerged              CheckoutStrategy = C.GIT_CHECKOUT_SKIP_UNMERGED                // Allow checkout to skip unmerged files
	CheckoutUserOurs                  CheckoutStrategy = C.GIT_CHECKOUT_USE_OURS                     // For unmerged files, checkout stage 2 from index
	CheckoutUseTheirs                 CheckoutStrategy = C.GIT_CHECKOUT_USE_THEIRS                   // For unmerged files, checkout stage 3 from index
	CheckoutDisablePathspecMatch      CheckoutStrategy = C.GIT_CHECKOUT_DISABLE_PATHSPEC_MATCH       // Treat pathspec as simple list of exact match file paths
	CheckoutSkipLockedDirectories     CheckoutStrategy = C.GIT_CHECKOUT_SKIP_LOCKED_DIRECTORIES      // Ignore directories in use, they will be left empty
	CheckoutDontOverwriteIgnored      CheckoutStrategy = C.GIT_CHECKOUT_DONT_OVERWRITE_IGNORED       // Don't overwrite ignored files that exist in the checkout target
	CheckoutConflictStyleMerge        CheckoutStrategy = C.GIT_CHECKOUT_CONFLICT_STYLE_MERGE         // Write normal merge files for conflicts
	CheckoutConflictStyleDiff3        CheckoutStrategy = C.GIT_CHECKOUT_CONFLICT_STYLE_DIFF3         // Include common ancestor data in diff3 format files for conflicts
	CheckoutDontRemoveExisting        CheckoutStrategy = C.GIT_CHECKOUT_DONT_REMOVE_EXISTING         // Don't overwrite existing files or folders
	CheckoutDontWriteIndex            CheckoutStrategy = C.GIT_CHECKOUT_DONT_WRITE_INDEX             // Normally checkout writes the index upon completion; this prevents that
	CheckoutUpdateSubmodules          CheckoutStrategy = C.GIT_CHECKOUT_UPDATE_SUBMODULES            // Recursively checkout submodules with same options (NOT IMPLEMENTED)
	CheckoutUpdateSubmodulesIfChanged CheckoutStrategy = C.GIT_CHECKOUT_UPDATE_SUBMODULES_IF_CHANGED // Recursively checkout submodules if HEAD moved in super repo (NOT IMPLEMENTED)
)

type CheckoutNotifyCallback func(why CheckoutNotifyType, path string, baseline, target, workdir DiffFile) ErrorCode
type CheckoutProgressCallback func(path string, completed, total uint) ErrorCode

type CheckoutOpts struct {
	Strategy         CheckoutStrategy   // Default will be a dry run
	DisableFilters   bool               // Don't apply filters like CRLF conversion
	DirMode          os.FileMode        // Default is 0755
	FileMode         os.FileMode        // Default is 0644 or 0755 as dictated by blob
	FileOpenFlags    int                // Default is O_CREAT | O_TRUNC | O_WRONLY
	NotifyFlags      CheckoutNotifyType // Default will be none
	NotifyCallback   CheckoutNotifyCallback
	ProgressCallback CheckoutProgressCallback
	TargetDirectory  string // Alternative checkout path to workdir
	Paths            []string
	Baseline         *Tree
}

func checkoutOptionsFromC(c *C.git_checkout_options) CheckoutOpts {
	opts := CheckoutOpts{}
	opts.Strategy = CheckoutStrategy(c.checkout_strategy)
	opts.DisableFilters = c.disable_filters != 0
	opts.DirMode = os.FileMode(c.dir_mode)
	opts.FileMode = os.FileMode(c.file_mode)
	opts.FileOpenFlags = int(c.file_open_flags)
	opts.NotifyFlags = CheckoutNotifyType(c.notify_flags)
	if c.notify_payload != nil {
		opts.NotifyCallback = pointerHandles.Get(c.notify_payload).(*CheckoutOpts).NotifyCallback
	}
	if c.progress_payload != nil {
		opts.ProgressCallback = pointerHandles.Get(c.progress_payload).(*CheckoutOpts).ProgressCallback
	}
	if c.target_directory != nil {
		opts.TargetDirectory = C.GoString(c.target_directory)
	}
	return opts
}

func (opts *CheckoutOpts) toC() *C.git_checkout_options {
	if opts == nil {
		return nil
	}
	c := C.git_checkout_options{}
	populateCheckoutOpts(&c, opts)
	return &c
}

//export checkoutNotifyCallback
func checkoutNotifyCallback(why C.git_checkout_notify_t, cpath *C.char, cbaseline, ctarget, cworkdir, data unsafe.Pointer) int {
	if data == nil {
		return 0
	}
	path := C.GoString(cpath)
	var baseline, target, workdir DiffFile
	if cbaseline != nil {
		baseline = diffFileFromC((*C.git_diff_file)(cbaseline))
	}
	if ctarget != nil {
		target = diffFileFromC((*C.git_diff_file)(ctarget))
	}
	if cworkdir != nil {
		workdir = diffFileFromC((*C.git_diff_file)(cworkdir))
	}
	opts := pointerHandles.Get(data).(*CheckoutOpts)
	if opts.NotifyCallback == nil {
		return 0
	}
	return int(opts.NotifyCallback(CheckoutNotifyType(why), path, baseline, target, workdir))
}

//export checkoutProgressCallback
func checkoutProgressCallback(path *C.char, completed_steps, total_steps C.size_t, data unsafe.Pointer) int {
	opts := pointerHandles.Get(data).(*CheckoutOpts)
	if opts.ProgressCallback == nil {
		return 0
	}
	return int(opts.ProgressCallback(C.GoString(path), uint(completed_steps), uint(total_steps)))
}

// Convert the CheckoutOpts struct to the corresponding
// C-struct. Returns a pointer to ptr, or nil if opts is nil, in order
// to help with what to pass.
func populateCheckoutOpts(ptr *C.git_checkout_options, opts *CheckoutOpts) *C.git_checkout_options {
	if opts == nil {
		return nil
	}

	C.git_checkout_init_options(ptr, 1)
	ptr.checkout_strategy = C.uint(opts.Strategy)
	ptr.disable_filters = cbool(opts.DisableFilters)
	ptr.dir_mode = C.uint(opts.DirMode.Perm())
	ptr.file_mode = C.uint(opts.FileMode.Perm())
	ptr.notify_flags = C.uint(opts.NotifyFlags)
	if opts.NotifyCallback != nil || opts.ProgressCallback != nil {
		C._go_git_populate_checkout_cb(ptr)
	}
	payload := pointerHandles.Track(opts)
	if opts.NotifyCallback != nil {
		ptr.notify_payload = payload
	}
	if opts.ProgressCallback != nil {
		ptr.progress_payload = payload
	}
	if opts.TargetDirectory != "" {
		ptr.target_directory = C.CString(opts.TargetDirectory)
	}
	if len(opts.Paths) > 0 {
		ptr.paths.strings = makeCStringsFromStrings(opts.Paths)
		ptr.paths.count = C.size_t(len(opts.Paths))
	}

	if opts.Baseline != nil {
		ptr.baseline = opts.Baseline.cast_ptr
	}

	return ptr
}

func freeCheckoutOpts(ptr *C.git_checkout_options) {
	if ptr == nil {
		return
	}
	C.free(unsafe.Pointer(ptr.target_directory))
	if ptr.paths.count > 0 {
		freeStrarray(&ptr.paths)
	}
	if ptr.notify_payload != nil {
		pointerHandles.Untrack(ptr.notify_payload)
	}
}

// Updates files in the index and the working tree to match the content of
// the commit pointed at by HEAD. opts may be nil.
func (v *Repository) CheckoutHead(opts *CheckoutOpts) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	cOpts := opts.toC()
	defer freeCheckoutOpts(cOpts)

	ret := C.git_checkout_head(v.ptr, cOpts)
	if ret < 0 {
		return MakeGitError(ret)
	}

	return nil
}

// Updates files in the working tree to match the content of the given
// index. If index is nil, the repository's index will be used. opts
// may be nil.
func (v *Repository) CheckoutIndex(index *Index, opts *CheckoutOpts) error {
	var iptr *C.git_index = nil
	if index != nil {
		iptr = index.ptr
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	cOpts := opts.toC()
	defer freeCheckoutOpts(cOpts)

	ret := C.git_checkout_index(v.ptr, iptr, cOpts)
	if ret < 0 {
		return MakeGitError(ret)
	}

	return nil
}

func (v *Repository) CheckoutTree(tree *Tree, opts *CheckoutOpts) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	cOpts := opts.toC()
	defer freeCheckoutOpts(cOpts)

	ret := C.git_checkout_tree(v.ptr, tree.ptr, cOpts)
	if ret < 0 {
		return MakeGitError(ret)
	}

	return nil
}
