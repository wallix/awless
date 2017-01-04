package git

/*
#include <git2.h>
*/
import "C"
import (
	"runtime"
)

type CherrypickOptions struct {
	Version      uint
	Mainline     uint
	MergeOpts    MergeOptions
	CheckoutOpts CheckoutOpts
}

func cherrypickOptionsFromC(c *C.git_cherrypick_options) CherrypickOptions {
	opts := CherrypickOptions{
		Version:      uint(c.version),
		Mainline:     uint(c.mainline),
		MergeOpts:    mergeOptionsFromC(&c.merge_opts),
		CheckoutOpts: checkoutOptionsFromC(&c.checkout_opts),
	}
	return opts
}

func (opts *CherrypickOptions) toC() *C.git_cherrypick_options {
	if opts == nil {
		return nil
	}
	c := C.git_cherrypick_options{}
	c.version = C.uint(opts.Version)
	c.mainline = C.uint(opts.Mainline)
	c.merge_opts = *opts.MergeOpts.toC()
	c.checkout_opts = *opts.CheckoutOpts.toC()
	return &c
}

func freeCherrypickOpts(ptr *C.git_cherrypick_options) {
	if ptr == nil {
		return
	}
	freeCheckoutOpts(&ptr.checkout_opts)
}

func DefaultCherrypickOptions() (CherrypickOptions, error) {
	c := C.git_cherrypick_options{}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ecode := C.git_cherrypick_init_options(&c, C.GIT_CHERRYPICK_OPTIONS_VERSION)
	if ecode < 0 {
		return CherrypickOptions{}, MakeGitError(ecode)
	}
	defer freeCherrypickOpts(&c)
	return cherrypickOptionsFromC(&c), nil
}

func (v *Repository) Cherrypick(commit *Commit, opts CherrypickOptions) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	cOpts := opts.toC()
	defer freeCherrypickOpts(cOpts)

	ecode := C.git_cherrypick(v.ptr, commit.cast_ptr, cOpts)
	if ecode < 0 {
		return MakeGitError(ecode)
	}
	return nil
}
