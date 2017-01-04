git2go
======
[![GoDoc](https://godoc.org/github.com/libgit2/git2go?status.svg)](http://godoc.org/github.com/libgit2/git2go) [![Build Status](https://travis-ci.org/libgit2/git2go.svg?branch=master)](https://travis-ci.org/libgit2/git2go)


Go bindings for [libgit2](http://libgit2.github.com/). The `master` branch follows the latest libgit2 release. The versioned branches indicate which libgit2 version they work against.

Installing
----------

This project wraps the functionality provided by libgit2. If you're using a stable version, install it to your system via your system's package manager and then install git2go as usual.

Otherwise (`next` which tracks an unstable version), we need to build libgit2 as well. In order to build it, you need `cmake`, `pkg-config` and a C compiler. You will also need the development packages for OpenSSL and LibSSH2 installed if you want libgit2 to support HTTPS and SSH respectively.

### Stable version

git2go has `master` which tracks the latest release of libgit2, and versioned branches which indicate which version of libgit2 they work against. Install the development package on your system via your favourite package manager or from source and you can use a service like gopkg.in to use the appropriate version. For the libgit2 v0.22 case, you can use

    import "gopkg.in/libgit2/git2go.v22"

to use a version of git2go which will work against libgit2 v0.22 and dynamically link to the library. You can use

    import "github.com/libgit2/git2go"

to use the version which works against the latest release.

### From `next`

The `next` branch follows libgit2's master branch, which means there is no stable API or ABI to link against. git2go can statically link against a vendored version of libgit2.

Run `go get -d github.com/libgit2/git2go` to download the code and go to your `$GOPATH/src/github.com/libgit2/git2go` directory. From there, we need to build the C code and put it into the resulting go binary.

    git checkout next
    git submodule update --init # get libgit2
    make install

will compile libgit2. Run `go install` so that it's statically linked to the git2go package.

Parallelism and network operations
----------------------------------

libgit2 uses OpenSSL and LibSSH2 for performing encrypted network connections. For now, git2go asks libgit2 to set locking for OpenSSL. This makes HTTPS connections thread-safe, but it is fragile and will likely stop doing it soon. This may also make SSH connections thread-safe if your copy of libssh2 is linked against OpenSSL. Check libgit2's `THREADSAFE.md` for more information.

Running the tests
-----------------

For the stable version, `go test` will work as usual. For the `next` branch, similarly to installing, running the tests requires linking against the local libgit2 library, so the Makefile provides a wrapper

    make test

Alternatively, if you want to pass arguments to `go test`, you can use the script that sets it all up

    ./script/with-static.sh go test -v

which will run the specified arguments with the correct environment variables.

License
-------

M to the I to the T. See the LICENSE file if you've never seen a MIT license before.

Authors
-------

- Carlos Martín (@carlosmn)
- Vicent Martí (@vmg)

