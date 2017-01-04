package git

/*
#include <git2.h>
*/
import "C"

import (
	"runtime"
	"unsafe"
)

type Status int

const (
	StatusCurrent         Status = C.GIT_STATUS_CURRENT
	StatusIndexNew        Status = C.GIT_STATUS_INDEX_NEW
	StatusIndexModified   Status = C.GIT_STATUS_INDEX_MODIFIED
	StatusIndexDeleted    Status = C.GIT_STATUS_INDEX_DELETED
	StatusIndexRenamed    Status = C.GIT_STATUS_INDEX_RENAMED
	StatusIndexTypeChange Status = C.GIT_STATUS_INDEX_TYPECHANGE
	StatusWtNew           Status = C.GIT_STATUS_WT_NEW
	StatusWtModified      Status = C.GIT_STATUS_WT_MODIFIED
	StatusWtDeleted       Status = C.GIT_STATUS_WT_DELETED
	StatusWtTypeChange    Status = C.GIT_STATUS_WT_TYPECHANGE
	StatusWtRenamed       Status = C.GIT_STATUS_WT_RENAMED
	StatusIgnored         Status = C.GIT_STATUS_IGNORED
	StatusConflicted      Status = C.GIT_STATUS_CONFLICTED
)

type StatusEntry struct {
	Status         Status
	HeadToIndex    DiffDelta
	IndexToWorkdir DiffDelta
}

func statusEntryFromC(statusEntry *C.git_status_entry) StatusEntry {
	var headToIndex DiffDelta = DiffDelta{}
	var indexToWorkdir DiffDelta = DiffDelta{}

	// Based on the libgit2 status example, head_to_index can be null in some cases
	if statusEntry.head_to_index != nil {
		headToIndex = diffDeltaFromC(statusEntry.head_to_index)
	}
	if statusEntry.index_to_workdir != nil {
		indexToWorkdir = diffDeltaFromC(statusEntry.index_to_workdir)
	}

	return StatusEntry{
		Status:         Status(statusEntry.status),
		HeadToIndex:    headToIndex,
		IndexToWorkdir: indexToWorkdir,
	}
}

type StatusList struct {
	ptr *C.git_status_list
}

func newStatusListFromC(ptr *C.git_status_list) *StatusList {
	if ptr == nil {
		return nil
	}

	statusList := &StatusList{
		ptr: ptr,
	}

	runtime.SetFinalizer(statusList, (*StatusList).Free)
	return statusList
}

func (statusList *StatusList) Free() {
	if statusList.ptr == nil {
		return
	}
	runtime.SetFinalizer(statusList, nil)
	C.git_status_list_free(statusList.ptr)
	statusList.ptr = nil
}

func (statusList *StatusList) ByIndex(index int) (StatusEntry, error) {
	if statusList.ptr == nil {
		return StatusEntry{}, ErrInvalid
	}
	ptr := C.git_status_byindex(statusList.ptr, C.size_t(index))
	return statusEntryFromC(ptr), nil
}

func (statusList *StatusList) EntryCount() (int, error) {
	if statusList.ptr == nil {
		return -1, ErrInvalid
	}
	return int(C.git_status_list_entrycount(statusList.ptr)), nil
}

type StatusOpt int

const (
	StatusOptIncludeUntracked      StatusOpt = C.GIT_STATUS_OPT_INCLUDE_UNTRACKED
	StatusOptIncludeIgnored        StatusOpt = C.GIT_STATUS_OPT_INCLUDE_IGNORED
	StatusOptIncludeUnmodified     StatusOpt = C.GIT_STATUS_OPT_INCLUDE_UNMODIFIED
	StatusOptExcludeSubmodules     StatusOpt = C.GIT_STATUS_OPT_EXCLUDE_SUBMODULES
	StatusOptRecurseUntrackedDirs  StatusOpt = C.GIT_STATUS_OPT_RECURSE_UNTRACKED_DIRS
	StatusOptDisablePathspecMatch  StatusOpt = C.GIT_STATUS_OPT_DISABLE_PATHSPEC_MATCH
	StatusOptRecurseIgnoredDirs    StatusOpt = C.GIT_STATUS_OPT_RECURSE_IGNORED_DIRS
	StatusOptRenamesHeadToIndex    StatusOpt = C.GIT_STATUS_OPT_RENAMES_HEAD_TO_INDEX
	StatusOptRenamesIndexToWorkdir StatusOpt = C.GIT_STATUS_OPT_RENAMES_INDEX_TO_WORKDIR
	StatusOptSortCaseSensitively   StatusOpt = C.GIT_STATUS_OPT_SORT_CASE_SENSITIVELY
	StatusOptSortCaseInsensitively StatusOpt = C.GIT_STATUS_OPT_SORT_CASE_INSENSITIVELY
	StatusOptRenamesFromRewrites   StatusOpt = C.GIT_STATUS_OPT_RENAMES_FROM_REWRITES
	StatusOptNoRefresh             StatusOpt = C.GIT_STATUS_OPT_NO_REFRESH
	StatusOptUpdateIndex           StatusOpt = C.GIT_STATUS_OPT_UPDATE_INDEX
)

type StatusShow int

const (
	StatusShowIndexAndWorkdir StatusShow = C.GIT_STATUS_SHOW_INDEX_AND_WORKDIR
	StatusShowIndexOnly       StatusShow = C.GIT_STATUS_SHOW_INDEX_ONLY
	StatusShowWorkdirOnly     StatusShow = C.GIT_STATUS_SHOW_WORKDIR_ONLY
)

type StatusOptions struct {
	Show     StatusShow
	Flags    StatusOpt
	Pathspec []string
}

func (v *Repository) StatusList(opts *StatusOptions) (*StatusList, error) {
	var ptr *C.git_status_list
	var copts *C.git_status_options

	if opts != nil {
		cpathspec := C.git_strarray{}
		if opts.Pathspec != nil {
			cpathspec.count = C.size_t(len(opts.Pathspec))
			cpathspec.strings = makeCStringsFromStrings(opts.Pathspec)
			defer freeStrarray(&cpathspec)
		}

		copts = &C.git_status_options{
			version:  C.GIT_STATUS_OPTIONS_VERSION,
			show:     C.git_status_show_t(opts.Show),
			flags:    C.uint(opts.Flags),
			pathspec: cpathspec,
		}
	} else {
		copts = &C.git_status_options{}
		ret := C.git_status_init_options(copts, C.GIT_STATUS_OPTIONS_VERSION)
		if ret < 0 {
			return nil, MakeGitError(ret)
		}
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_status_list_new(&ptr, v.ptr, copts)
	if ret < 0 {
		return nil, MakeGitError(ret)
	}

	return newStatusListFromC(ptr), nil
}

func (v *Repository) StatusFile(path string) (Status, error) {
	var statusFlags C.uint
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_status_file(&statusFlags, v.ptr, cPath)
	if ret < 0 {
		return 0, MakeGitError(ret)
	}
	return Status(statusFlags), nil
}
