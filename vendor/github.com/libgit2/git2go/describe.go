package git

/*
#include <git2.h>
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// DescribeOptions represents the describe operation configuration.
//
// You can use DefaultDescribeOptions() to get default options.
type DescribeOptions struct {
	// How many tags as candidates to consider to describe the input commit-ish.
	// Increasing it above 10 will take slightly longer but may produce a more
	// accurate result. 0 will cause only exact matches to be output.
	MaxCandidatesTags uint // default: 10

	// By default describe only shows annotated tags. Change this in order
	// to show all refs from refs/tags or refs/.
	Strategy DescribeOptionsStrategy // default: DescribeDefault

	// Only consider tags matching the given glob(7) pattern, excluding
	// the "refs/tags/" prefix. Can be used to avoid leaking private
	// tags from the repo.
	Pattern string

	// When calculating the distance from the matching tag or
	// reference, only walk down the first-parent ancestry.
	OnlyFollowFirstParent bool

	// If no matching tag or reference is found, the describe
	// operation would normally fail. If this option is set, it
	// will instead fall back to showing the full id of the commit.
	ShowCommitOidAsFallback bool
}

// DefaultDescribeOptions returns default options for the describe operation.
func DefaultDescribeOptions() (DescribeOptions, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	opts := C.git_describe_options{}
	ecode := C.git_describe_init_options(&opts, C.GIT_DESCRIBE_OPTIONS_VERSION)
	if ecode < 0 {
		return DescribeOptions{}, MakeGitError(ecode)
	}

	return DescribeOptions{
		MaxCandidatesTags: uint(opts.max_candidates_tags),
		Strategy:          DescribeOptionsStrategy(opts.describe_strategy),
	}, nil
}

// DescribeFormatOptions can be used for formatting the describe string.
//
// You can use DefaultDescribeFormatOptions() to get default options.
type DescribeFormatOptions struct {
	// Size of the abbreviated commit id to use. This value is the
	// lower bound for the length of the abbreviated string.
	AbbreviatedSize uint // default: 7

	// Set to use the long format even when a shorter name could be used.
	AlwaysUseLongFormat bool

	// If the workdir is dirty and this is set, this string will be
	// appended to the description string.
	DirtySuffix string
}

// DefaultDescribeFormatOptions returns default options for formatting
// the output.
func DefaultDescribeFormatOptions() (DescribeFormatOptions, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	opts := C.git_describe_format_options{}
	ecode := C.git_describe_init_format_options(&opts, C.GIT_DESCRIBE_FORMAT_OPTIONS_VERSION)
	if ecode < 0 {
		return DescribeFormatOptions{}, MakeGitError(ecode)
	}

	return DescribeFormatOptions{
		AbbreviatedSize:     uint(opts.abbreviated_size),
		AlwaysUseLongFormat: opts.always_use_long_format == 1,
	}, nil
}

// DescribeOptionsStrategy behaves like the --tags and --all options
// to git-describe, namely they say to look for any reference in
// either refs/tags/ or refs/ respectively.
//
// By default it only shows annotated tags.
type DescribeOptionsStrategy uint

// Describe strategy options.
const (
	DescribeDefault DescribeOptionsStrategy = C.GIT_DESCRIBE_DEFAULT
	DescribeTags    DescribeOptionsStrategy = C.GIT_DESCRIBE_TAGS
	DescribeAll     DescribeOptionsStrategy = C.GIT_DESCRIBE_ALL
)

// Describe performs the describe operation on the commit.
func (c *Commit) Describe(opts *DescribeOptions) (*DescribeResult, error) {
	var resultPtr *C.git_describe_result

	var cDescribeOpts *C.git_describe_options
	if opts != nil {
		var cpattern *C.char
		if len(opts.Pattern) > 0 {
			cpattern = C.CString(opts.Pattern)
			defer C.free(unsafe.Pointer(cpattern))
		}

		cDescribeOpts = &C.git_describe_options{
			version:                     C.GIT_DESCRIBE_OPTIONS_VERSION,
			max_candidates_tags:         C.uint(opts.MaxCandidatesTags),
			describe_strategy:           C.uint(opts.Strategy),
			pattern:                     cpattern,
			only_follow_first_parent:    cbool(opts.OnlyFollowFirstParent),
			show_commit_oid_as_fallback: cbool(opts.ShowCommitOidAsFallback),
		}
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ecode := C.git_describe_commit(&resultPtr, c.ptr, cDescribeOpts)
	if ecode < 0 {
		return nil, MakeGitError(ecode)
	}

	return newDescribeResultFromC(resultPtr), nil
}

// DescribeWorkdir describes the working tree. It means describe HEAD
// and appends <mark> (-dirty by default) if the working tree is dirty.
func (repo *Repository) DescribeWorkdir(opts *DescribeOptions) (*DescribeResult, error) {
	var resultPtr *C.git_describe_result

	var cDescribeOpts *C.git_describe_options
	if opts != nil {
		var cpattern *C.char
		if len(opts.Pattern) > 0 {
			cpattern = C.CString(opts.Pattern)
			defer C.free(unsafe.Pointer(cpattern))
		}

		cDescribeOpts = &C.git_describe_options{
			version:                     C.GIT_DESCRIBE_OPTIONS_VERSION,
			max_candidates_tags:         C.uint(opts.MaxCandidatesTags),
			describe_strategy:           C.uint(opts.Strategy),
			pattern:                     cpattern,
			only_follow_first_parent:    cbool(opts.OnlyFollowFirstParent),
			show_commit_oid_as_fallback: cbool(opts.ShowCommitOidAsFallback),
		}
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ecode := C.git_describe_workdir(&resultPtr, repo.ptr, cDescribeOpts)
	if ecode < 0 {
		return nil, MakeGitError(ecode)
	}

	return newDescribeResultFromC(resultPtr), nil
}

// DescribeResult represents the output from the 'git_describe_commit'
// and 'git_describe_workdir' functions in libgit2.
//
// Use Format() to get a string out of it.
type DescribeResult struct {
	ptr *C.git_describe_result
}

func newDescribeResultFromC(ptr *C.git_describe_result) *DescribeResult {
	result := &DescribeResult{
		ptr: ptr,
	}
	runtime.SetFinalizer(result, (*DescribeResult).Free)
	return result
}

// Format prints the DescribeResult as a string.
func (result *DescribeResult) Format(opts *DescribeFormatOptions) (string, error) {
	resultBuf := C.git_buf{}

	var cFormatOpts *C.git_describe_format_options
	if opts != nil {
		cDirtySuffix := C.CString(opts.DirtySuffix)
		defer C.free(unsafe.Pointer(cDirtySuffix))

		cFormatOpts = &C.git_describe_format_options{
			version:                C.GIT_DESCRIBE_FORMAT_OPTIONS_VERSION,
			abbreviated_size:       C.uint(opts.AbbreviatedSize),
			always_use_long_format: cbool(opts.AlwaysUseLongFormat),
			dirty_suffix:           cDirtySuffix,
		}
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ecode := C.git_describe_format(&resultBuf, result.ptr, cFormatOpts)
	if ecode < 0 {
		return "", MakeGitError(ecode)
	}
	defer C.git_buf_free(&resultBuf)

	return C.GoString(resultBuf.ptr), nil
}

// Free cleans up the C reference.
func (result *DescribeResult) Free() {
	runtime.SetFinalizer(result, nil)
	C.git_describe_result_free(result.ptr)
	result.ptr = nil
}
