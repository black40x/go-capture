package capture

/*
#cgo LDFLAGS: -framework CoreGraphics -framework CoreFoundation -framework IOSurface

#include <CoreGraphics/CGDisplayStream.h>
#include <CoreGraphics/CoreGraphics.h>
#include <IOSurface/IOSurface.h>
#include "./capture_darwin.h"

CGDisplayStreamRef CaptureInit(int width, int height) {
	size_t output_width = width;
	size_t output_height = height;
	uint32_t pixel_format = '720f';
	pixel_format = 'BGRA';
	dispatch_queue_t dq = dispatch_queue_create("com.black40x.screencapture", DISPATCH_QUEUE_SERIAL);
	CGDirectDisplayID display_ids[5];
	uint32_t found_displays = 0;
	CGError err = CGGetActiveDisplayList(5, display_ids, &found_displays);
	CGDisplayStreamRef sref;

	sref = CGDisplayStreamCreateWithDispatchQueue(
		display_ids[0], output_width, output_height, pixel_format, NULL, dq,
		^(CGDisplayStreamFrameStatus status, uint64_t time, IOSurfaceRef frame, CGDisplayStreamUpdateRef ref) {
			if (kCGDisplayStreamFrameStatusFrameComplete == status && NULL != frame) {
				IOSurfaceLock(frame, 0x00000001, NULL);
				uint8_t* pix = (uint8_t*)IOSurfaceGetBaseAddress(frame);
				if (NULL != pix) {
					ReturnFrame(pix, width, height);
				}
				IOSurfaceUnlock(frame, 0x00000001, NULL);
			}
		});

	return sref;
}

void CaptureStart(CGDisplayStreamRef sref) {
	CGError err;
	err = CGDisplayStreamStart(sref);

	if (kCGErrorSuccess != err) {
		exit(EXIT_FAILURE);
  	}
}

void CaptureStop(CGDisplayStreamRef sref) {
	CGDisplayStreamStop(sref);
}
*/
import "C"

var callbackFrameWriter func([]uint8)

func GetDisplayRect() *DisplayRect {
	dr := C.CGDisplayBounds(C.CGMainDisplayID())
	return &DisplayRect{int(dr.size.width), int(dr.size.height)}
}
