package cabocha

/*
#cgo LDFLAGS: -lcabocha
#include <stdlib.h>
#include <cabocha.h>
*/
import "C"
import (
	"log"
	"os"
	"unsafe"
)

// Cabocha struct
type Cabocha struct {
	h *C.cabocha_t
}

// NewCabocha creates new cabocha instance
func NewCabocha(options []string) *Cabocha {
	argc := C.int(1)
	argv := []*C.char{C.CString(os.Args[0])}
	if options != nil {
		for _, o := range options {
			argv = append(argv, C.CString(o))
			argc++
		}
	}
	defer func() {
		for _, s := range argv {
			C.free(unsafe.Pointer(s))
		}
	}()

	h := C.cabocha_new(argc, (**C.char)(unsafe.Pointer(&argv[0])))
	if h == nil {
		log.Print("error: cabocha_new failed")
		return nil
	}
	return &Cabocha{h}
}

// Destroy destroy and cleanup cabocha instance
func (c *Cabocha) Destroy() {
	C.cabocha_destroy(c.h)
}

// Chunks returns chunks of target string
func (c *Cabocha) Chunks(target string) []string {
	ctarget := C.CString(target)
	defer C.free(unsafe.Pointer(ctarget))

	tree := C.cabocha_sparse_totree(c.h, ctarget)
	size := C.cabocha_tree_token_size(tree)

	runes := []rune(target)
	chunks := make([]string, 0)

	var chunk string
	r := 0
	for i := 0; i < int(size); i++ {
		token := C.cabocha_tree_token(tree, C.ulong(i))
		if token.chunk != nil {
			if chunk != "" {
				chunks = append(chunks, chunk)
			}
			chunk = ""
		}
		surface := C.GoString(token.surface)
		for _, s := range []rune(surface) {
			for s != runes[r] {
				if chunk == "" {
					chunks[len(chunks)-1] += string(runes[r])
				} else {
					chunk += string(runes[r])
				}
				r++
			}
			chunk += string(runes[r])
			r++
		}
	}
	if chunk != "" {
		chunks = append(chunks, chunk)
	}

	return chunks
}

// Chunks returns chunks of target string
func Chunks(target string) []string {
	cabocha := NewCabocha(nil)
	defer cabocha.Destroy()
	return cabocha.Chunks(target)
}
