package git

/*
#include <git2.h>
*/
import "C"
import (
	"runtime"
	"unsafe"
)

type ConfigLevel int

const (
	// System-wide on Windows, for compatibility with portable git
	ConfigLevelProgramdata ConfigLevel = C.GIT_CONFIG_LEVEL_PROGRAMDATA

	// System-wide configuration file; /etc/gitconfig on Linux systems
	ConfigLevelSystem ConfigLevel = C.GIT_CONFIG_LEVEL_SYSTEM

	// XDG compatible configuration file; typically ~/.config/git/config
	ConfigLevelXDG ConfigLevel = C.GIT_CONFIG_LEVEL_XDG

	// User-specific configuration file (also called Global configuration
	// file); typically ~/.gitconfig
	ConfigLevelGlobal ConfigLevel = C.GIT_CONFIG_LEVEL_GLOBAL

	// Repository specific configuration file; $WORK_DIR/.git/config on
	// non-bare repos
	ConfigLevelLocal ConfigLevel = C.GIT_CONFIG_LEVEL_LOCAL

	// Application specific configuration file; freely defined by applications
	ConfigLevelApp ConfigLevel = C.GIT_CONFIG_LEVEL_APP

	// Represents the highest level available config file (i.e. the most
	// specific config file available that actually is loaded)
	ConfigLevelHighest ConfigLevel = C.GIT_CONFIG_HIGHEST_LEVEL
)

type ConfigEntry struct {
	Name  string
	Value string
	Level ConfigLevel
}

func newConfigEntryFromC(centry *C.git_config_entry) *ConfigEntry {
	return &ConfigEntry{
		Name:  C.GoString(centry.name),
		Value: C.GoString(centry.value),
		Level: ConfigLevel(centry.level),
	}
}

type Config struct {
	ptr *C.git_config
}

// NewConfig creates a new empty configuration object
func NewConfig() (*Config, error) {
	config := new(Config)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if ret := C.git_config_new(&config.ptr); ret < 0 {
		return nil, MakeGitError(ret)
	}

	return config, nil
}

// AddFile adds a file-backed backend to the config object at the specified level.
func (c *Config) AddFile(path string, level ConfigLevel, force bool) error {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_config_add_file_ondisk(c.ptr, cpath, C.git_config_level_t(level), cbool(force))
	if ret < 0 {
		return MakeGitError(ret)
	}

	return nil
}

func (c *Config) LookupInt32(name string) (int32, error) {
	var out C.int32_t
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_config_get_int32(&out, c.ptr, cname)
	if ret < 0 {
		return 0, MakeGitError(ret)
	}

	return int32(out), nil
}

func (c *Config) LookupInt64(name string) (int64, error) {
	var out C.int64_t
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_config_get_int64(&out, c.ptr, cname)
	if ret < 0 {
		return 0, MakeGitError(ret)
	}

	return int64(out), nil
}

func (c *Config) LookupString(name string) (string, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	valBuf := C.git_buf{}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if ret := C.git_config_get_string_buf(&valBuf, c.ptr, cname); ret < 0 {
		return "", MakeGitError(ret)
	}
	defer C.git_buf_free(&valBuf)

	return C.GoString(valBuf.ptr), nil
}

func (c *Config) LookupBool(name string) (bool, error) {
	var out C.int
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_config_get_bool(&out, c.ptr, cname)
	if ret < 0 {
		return false, MakeGitError(ret)
	}

	return out != 0, nil
}

func (c *Config) NewMultivarIterator(name, regexp string) (*ConfigIterator, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var cregexp *C.char
	if regexp == "" {
		cregexp = nil
	} else {
		cregexp = C.CString(regexp)
		defer C.free(unsafe.Pointer(cregexp))
	}

	iter := new(ConfigIterator)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_config_multivar_iterator_new(&iter.ptr, c.ptr, cname, cregexp)
	if ret < 0 {
		return nil, MakeGitError(ret)
	}

	runtime.SetFinalizer(iter, (*ConfigIterator).Free)
	return iter, nil
}

// NewIterator creates an iterator over each entry in the
// configuration
func (c *Config) NewIterator() (*ConfigIterator, error) {
	iter := new(ConfigIterator)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_config_iterator_new(&iter.ptr, c.ptr)
	if ret < 0 {
		return nil, MakeGitError(ret)
	}

	return iter, nil
}

// NewIteratorGlob creates an iterator over each entry in the
// configuration whose name matches the given regular expression
func (c *Config) NewIteratorGlob(regexp string) (*ConfigIterator, error) {
	iter := new(ConfigIterator)
	cregexp := C.CString(regexp)
	defer C.free(unsafe.Pointer(cregexp))

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_config_iterator_glob_new(&iter.ptr, c.ptr, cregexp)
	if ret < 0 {
		return nil, MakeGitError(ret)
	}

	return iter, nil
}

func (c *Config) SetString(name, value string) (err error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_config_set_string(c.ptr, cname, cvalue)
	if ret < 0 {
		return MakeGitError(ret)
	}

	return nil
}

func (c *Config) Free() {
	runtime.SetFinalizer(c, nil)
	C.git_config_free(c.ptr)
}

func (c *Config) SetInt32(name string, value int32) (err error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_config_set_int32(c.ptr, cname, C.int32_t(value))
	if ret < 0 {
		return MakeGitError(ret)
	}

	return nil
}

func (c *Config) SetInt64(name string, value int64) (err error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_config_set_int64(c.ptr, cname, C.int64_t(value))
	if ret < 0 {
		return MakeGitError(ret)
	}

	return nil
}

func (c *Config) SetBool(name string, value bool) (err error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_config_set_bool(c.ptr, cname, cbool(value))
	if ret < 0 {
		return MakeGitError(ret)
	}

	return nil
}

func (c *Config) SetMultivar(name, regexp, value string) (err error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	cregexp := C.CString(regexp)
	defer C.free(unsafe.Pointer(cregexp))

	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_config_set_multivar(c.ptr, cname, cregexp, cvalue)
	if ret < 0 {
		return MakeGitError(ret)
	}

	return nil
}

func (c *Config) Delete(name string) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_config_delete_entry(c.ptr, cname)

	if ret < 0 {
		return MakeGitError(ret)
	}

	return nil
}

// OpenLevel creates a single-level focused config object from a multi-level one
func (c *Config) OpenLevel(parent *Config, level ConfigLevel) (*Config, error) {
	config := new(Config)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_config_open_level(&config.ptr, parent.ptr, C.git_config_level_t(level))
	if ret < 0 {
		return nil, MakeGitError(ret)
	}

	return config, nil
}

// OpenOndisk creates a new config instance containing a single on-disk file
func OpenOndisk(parent *Config, path string) (*Config, error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	config := new(Config)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if ret := C.git_config_open_ondisk(&config.ptr, cpath); ret < 0 {
		return nil, MakeGitError(ret)
	}

	return config, nil
}

type ConfigIterator struct {
	ptr *C.git_config_iterator
}

// Next returns the next entry for this iterator
func (iter *ConfigIterator) Next() (*ConfigEntry, error) {
	var centry *C.git_config_entry

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_config_next(&centry, iter.ptr)
	if ret < 0 {
		return nil, MakeGitError(ret)
	}

	return newConfigEntryFromC(centry), nil
}

func (iter *ConfigIterator) Free() {
	runtime.SetFinalizer(iter, nil)
	C.free(unsafe.Pointer(iter.ptr))
}

func ConfigFindGlobal() (string, error) {
	var buf C.git_buf
	defer C.git_buf_free(&buf)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_config_find_global(&buf)
	if ret < 0 {
		return "", MakeGitError(ret)
	}

	return C.GoString(buf.ptr), nil
}

func ConfigFindSystem() (string, error) {
	var buf C.git_buf
	defer C.git_buf_free(&buf)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_config_find_system(&buf)
	if ret < 0 {
		return "", MakeGitError(ret)
	}

	return C.GoString(buf.ptr), nil
}

func ConfigFindXDG() (string, error) {
	var buf C.git_buf
	defer C.git_buf_free(&buf)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_config_find_xdg(&buf)
	if ret < 0 {
		return "", MakeGitError(ret)
	}

	return C.GoString(buf.ptr), nil
}

// ConfigFindProgramdata locate the path to the configuration file in ProgramData.
//
// Look for the file in %PROGRAMDATA%\Git\config used by portable git.
func ConfigFindProgramdata() (string, error) {
	var buf C.git_buf
	defer C.git_buf_free(&buf)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_config_find_programdata(&buf)
	if ret < 0 {
		return "", MakeGitError(ret)
	}

	return C.GoString(buf.ptr), nil
}
