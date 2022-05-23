package capture

// #include <stdlib.h>
import "C"
import (
	"unsafe"
)

//export ReturnFrame
func ReturnFrame(d *C.uint8_t, w, h C.int) {
	ww := int(w)
	hh := int(h)
	pix := (*[1 << 30]uint8)(unsafe.Pointer(d))[: ww*hh*4 : ww*hh*4]
	callbackFrameWriter(pix)
}
