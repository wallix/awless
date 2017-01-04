package git

/*
#include <git2.h>
*/
import "C"
import "unsafe"

type CredType uint

const (
	CredTypeUserpassPlaintext CredType = C.GIT_CREDTYPE_USERPASS_PLAINTEXT
	CredTypeSshKey            CredType = C.GIT_CREDTYPE_SSH_KEY
	CredTypeSshCustom         CredType = C.GIT_CREDTYPE_SSH_CUSTOM
	CredTypeDefault           CredType = C.GIT_CREDTYPE_DEFAULT
)

type Cred struct {
	ptr *C.git_cred
}

func (o *Cred) HasUsername() bool {
	if C.git_cred_has_username(o.ptr) == 1 {
		return true
	}
	return false
}

func (o *Cred) Type() CredType {
	return (CredType)(o.ptr.credtype)
}

func credFromC(ptr *C.git_cred) *Cred {
	return &Cred{ptr}
}

func NewCredUserpassPlaintext(username string, password string) (int, Cred) {
	cred := Cred{}
	cusername := C.CString(username)
	defer C.free(unsafe.Pointer(cusername))
	cpassword := C.CString(password)
	defer C.free(unsafe.Pointer(cpassword))
	ret := C.git_cred_userpass_plaintext_new(&cred.ptr, cusername, cpassword)
	return int(ret), cred
}

// NewCredSshKey creates new ssh credentials reading the public and private keys
// from the file system.
func NewCredSshKey(username string, publicKeyPath string, privateKeyPath string, passphrase string) (int, Cred) {
	cred := Cred{}
	cusername := C.CString(username)
	defer C.free(unsafe.Pointer(cusername))
	cpublickey := C.CString(publicKeyPath)
	defer C.free(unsafe.Pointer(cpublickey))
	cprivatekey := C.CString(privateKeyPath)
	defer C.free(unsafe.Pointer(cprivatekey))
	cpassphrase := C.CString(passphrase)
	defer C.free(unsafe.Pointer(cpassphrase))
	ret := C.git_cred_ssh_key_new(&cred.ptr, cusername, cpublickey, cprivatekey, cpassphrase)
	return int(ret), cred
}

// NewCredSshKeyFromMemory creates new ssh credentials using the publicKey and privateKey
// arguments as the values for the public and private keys.
func NewCredSshKeyFromMemory(username string, publicKey string, privateKey string, passphrase string) (int, Cred) {
	cred := Cred{}
	cusername := C.CString(username)
	defer C.free(unsafe.Pointer(cusername))
	cpublickey := C.CString(publicKey)
	defer C.free(unsafe.Pointer(cpublickey))
	cprivatekey := C.CString(privateKey)
	defer C.free(unsafe.Pointer(cprivatekey))
	cpassphrase := C.CString(passphrase)
	defer C.free(unsafe.Pointer(cpassphrase))
	ret := C.git_cred_ssh_key_memory_new(&cred.ptr, cusername, cpublickey, cprivatekey, cpassphrase)
	return int(ret), cred
}

func NewCredSshKeyFromAgent(username string) (int, Cred) {
	cred := Cred{}
	cusername := C.CString(username)
	defer C.free(unsafe.Pointer(cusername))
	ret := C.git_cred_ssh_key_from_agent(&cred.ptr, cusername)
	return int(ret), cred
}

func NewCredDefault() (int, Cred) {
	cred := Cred{}
	ret := C.git_cred_default_new(&cred.ptr)
	return int(ret), cred
}
