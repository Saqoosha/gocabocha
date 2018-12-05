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

// Chunks returns chunks of target string
func Chunks(target string) []string {
	argc := C.int(0)
	argv := make([]*C.char, 0)
	for _, a := range []string{
		os.Args[0],
		// "-b", "/usr/local/etc/mecabrc",
		// "-d", "/usr/local/Cellar/mecab/0.996/lib/mecab/dic/ipadic",
		// "-r", "/usr/local/etc/cabocharc",
		// "-m", "./lib/cabocha/model/dep.ipa.model",
		// "-M", "./lib/cabocha/model/chunk.ipa.model",
		// "-N", "./lib/cabocha/model/ne.ipa.model",
	} {
		argv = append(argv, C.CString(a))
		argc++
	}
	defer func() {
		for _, s := range argv {
			C.free(unsafe.Pointer(s))
		}
	}()

	cabocha := C.cabocha_new(argc, (**C.char)(unsafe.Pointer(&argv[0])))
	if cabocha == nil {
		log.Print("error: cabocha_new failed")
		return nil
	}
	defer C.cabocha_destroy(cabocha)

	ctarget := C.CString(target)
	defer C.free(unsafe.Pointer(ctarget))

	tree := C.cabocha_sparse_totree(cabocha, ctarget)
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
