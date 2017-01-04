package git

/*
#include <git2.h>
*/
import "C"
import "runtime"

type ResetType int

const (
	ResetSoft  ResetType = C.GIT_RESET_SOFT
	ResetMixed ResetType = C.GIT_RESET_MIXED
	ResetHard  ResetType = C.GIT_RESET_HARD
)

func (r *Repository) ResetToCommit(commit *Commit, resetType ResetType, opts *CheckoutOpts) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	ret := C.git_reset(r.ptr, commit.ptr, C.git_reset_t(resetType), opts.toC())

	if ret < 0 {
		return MakeGitError(ret)
	}
	return nil
}

func (r *Repository) ResetDefaultToCommit(commit *Commit, pathspecs []string) error {
	cpathspecs := C.git_strarray{}
	cpathspecs.count = C.size_t(len(pathspecs))
	cpathspecs.strings = makeCStringsFromStrings(pathspecs)
	defer freeStrarray(&cpathspecs)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	ret := C.git_reset_default(r.ptr, commit.ptr, &cpathspecs)

	if ret < 0 {
		return MakeGitError(ret)
	}
	return nil
}
