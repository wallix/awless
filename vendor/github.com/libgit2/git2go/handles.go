package git

/*
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"sync"
	"unsafe"
)

type HandleList struct {
	sync.RWMutex
	// stores the Go pointers
	handles map[unsafe.Pointer]interface{}
}

func NewHandleList() *HandleList {
	return &HandleList{
		handles: make(map[unsafe.Pointer]interface{}),
	}
}

// Track adds the given pointer to the list of pointers to track and
// returns a pointer value which can be passed to C as an opaque
// pointer.
func (v *HandleList) Track(pointer interface{}) unsafe.Pointer {
	handle := C.malloc(1)

	v.Lock()
	v.handles[handle] = pointer
	v.Unlock()

	return handle
}

// Untrack stops tracking the pointer given by the handle
func (v *HandleList) Untrack(handle unsafe.Pointer) {
	v.Lock()
	delete(v.handles, handle)
	C.free(handle)
	v.Unlock()
}

// Get retrieves the pointer from the given handle
func (v *HandleList) Get(handle unsafe.Pointer) interface{} {
	v.RLock()
	defer v.RUnlock()

	ptr, ok := v.handles[handle]
	if !ok {
		panic(fmt.Sprintf("invalid pointer handle: %p", handle))
	}

	return ptr
}
