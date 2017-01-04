package git

/*
#include <git2.h>
*/
import "C"
import (
	"runtime"
)

func (repo *Repository) DescendantOf(commit, ancestor *Oid) (bool, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_graph_descendant_of(repo.ptr, commit.toC(), ancestor.toC())
	if ret < 0 {
		return false, MakeGitError(ret)
	}

	return (ret > 0), nil
}

func (repo *Repository) AheadBehind(local, upstream *Oid) (ahead, behind int, err error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var aheadT C.size_t
	var behindT C.size_t

	ret := C.git_graph_ahead_behind(&aheadT, &behindT, repo.ptr, local.toC(), upstream.toC())
	if ret < 0 {
		return 0, 0, MakeGitError(ret)
	}

	return int(aheadT), int(behindT), nil
}
