package git

/*
#include <git2.h>
*/
import "C"
import (
	"runtime"
	"unsafe"
)

type BlameOptions struct {
	Flags              BlameOptionsFlag
	MinMatchCharacters uint16
	NewestCommit       *Oid
	OldestCommit       *Oid
	MinLine            uint32
	MaxLine            uint32
}

func DefaultBlameOptions() (BlameOptions, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	opts := C.git_blame_options{}
	ecode := C.git_blame_init_options(&opts, C.GIT_BLAME_OPTIONS_VERSION)
	if ecode < 0 {
		return BlameOptions{}, MakeGitError(ecode)
	}

	return BlameOptions{
		Flags:              BlameOptionsFlag(opts.flags),
		MinMatchCharacters: uint16(opts.min_match_characters),
		NewestCommit:       newOidFromC(&opts.newest_commit),
		OldestCommit:       newOidFromC(&opts.oldest_commit),
		MinLine:            uint32(opts.min_line),
		MaxLine:            uint32(opts.max_line),
	}, nil
}

type BlameOptionsFlag uint32

const (
	BlameNormal                      BlameOptionsFlag = C.GIT_BLAME_NORMAL
	BlameTrackCopiesSameFile         BlameOptionsFlag = C.GIT_BLAME_TRACK_COPIES_SAME_FILE
	BlameTrackCopiesSameCommitMoves  BlameOptionsFlag = C.GIT_BLAME_TRACK_COPIES_SAME_COMMIT_MOVES
	BlameTrackCopiesSameCommitCopies BlameOptionsFlag = C.GIT_BLAME_TRACK_COPIES_SAME_COMMIT_COPIES
	BlameTrackCopiesAnyCommitCopies  BlameOptionsFlag = C.GIT_BLAME_TRACK_COPIES_ANY_COMMIT_COPIES
	BlameFirstParent                 BlameOptionsFlag = C.GIT_BLAME_FIRST_PARENT
)

func (v *Repository) BlameFile(path string, opts *BlameOptions) (*Blame, error) {
	var blamePtr *C.git_blame

	var copts *C.git_blame_options
	if opts != nil {
		copts = &C.git_blame_options{
			version:              C.GIT_BLAME_OPTIONS_VERSION,
			flags:                C.uint32_t(opts.Flags),
			min_match_characters: C.uint16_t(opts.MinMatchCharacters),
			min_line:             C.size_t(opts.MinLine),
			max_line:             C.size_t(opts.MaxLine),
		}
		if opts.NewestCommit != nil {
			copts.newest_commit = *opts.NewestCommit.toC()
		}
		if opts.OldestCommit != nil {
			copts.oldest_commit = *opts.OldestCommit.toC()
		}
	}

	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ecode := C.git_blame_file(&blamePtr, v.ptr, cpath, copts)
	if ecode < 0 {
		return nil, MakeGitError(ecode)
	}

	return newBlameFromC(blamePtr), nil
}

type Blame struct {
	ptr *C.git_blame
}

func (blame *Blame) HunkCount() int {
	return int(C.git_blame_get_hunk_count(blame.ptr))
}

func (blame *Blame) HunkByIndex(index int) (BlameHunk, error) {
	ptr := C.git_blame_get_hunk_byindex(blame.ptr, C.uint32_t(index))
	if ptr == nil {
		return BlameHunk{}, ErrInvalid
	}
	return blameHunkFromC(ptr), nil
}

func (blame *Blame) HunkByLine(lineno int) (BlameHunk, error) {
	ptr := C.git_blame_get_hunk_byline(blame.ptr, C.size_t(lineno))
	if ptr == nil {
		return BlameHunk{}, ErrInvalid
	}
	return blameHunkFromC(ptr), nil
}

func newBlameFromC(ptr *C.git_blame) *Blame {
	if ptr == nil {
		return nil
	}

	blame := &Blame{
		ptr: ptr,
	}

	runtime.SetFinalizer(blame, (*Blame).Free)
	return blame
}

func (blame *Blame) Free() error {
	if blame.ptr == nil {
		return ErrInvalid
	}
	runtime.SetFinalizer(blame, nil)
	C.git_blame_free(blame.ptr)
	blame.ptr = nil
	return nil
}

type BlameHunk struct {
	LinesInHunk          uint16
	FinalCommitId        *Oid
	FinalStartLineNumber uint16
	FinalSignature       *Signature
	OrigCommitId         *Oid
	OrigPath             string
	OrigStartLineNumber  uint16
	OrigSignature        *Signature
	Boundary             bool
}

func blameHunkFromC(hunk *C.git_blame_hunk) BlameHunk {
	return BlameHunk{
		LinesInHunk:          uint16(hunk.lines_in_hunk),
		FinalCommitId:        newOidFromC(&hunk.final_commit_id),
		FinalStartLineNumber: uint16(hunk.final_start_line_number),
		FinalSignature:       newSignatureFromC(hunk.final_signature),
		OrigCommitId:         newOidFromC(&hunk.orig_commit_id),
		OrigPath:             C.GoString(hunk.orig_path),
		OrigStartLineNumber:  uint16(hunk.orig_start_line_number),
		OrigSignature:        newSignatureFromC(hunk.orig_signature),
		Boundary:             hunk.boundary == 1,
	}
}
