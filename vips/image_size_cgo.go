package vips

/*
#cgo pkg-config: vips
#include <vips/vips.h>
#include <stdlib.h>

// wrapper: return VipsImage* (NULL on error). autorotate_flag 0 or 1.
static VipsImage* vips_image_new_from_file_autorotate(const char* name, int autorotate_flag) {
    if (autorotate_flag) {
        return vips_image_new_from_file(name, "autorotate", TRUE, NULL);
    } else {
        return vips_image_new_from_file(name, NULL);
    }
}
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// GetSizeParams holds optional parameters for getting image size.
type GetSizeParams struct {
	Autorotate bool
}

// Convenience
func GetImageSizeFromFile(file string) (width, height int, err error) {
	return GetImageSizeFromFileWithParams(file, nil)
}

// Main implementation
func GetImageSizeFromFileWithParams(file string, params *GetSizeParams) (width, height int, err error) {
	// ensure libvips initialized
	startupIfNeeded()

	cFilePath := C.CString(file)
	defer C.free(unsafe.Pointer(cFilePath))

	autorotateFlag := 0
	if params != nil && params.Autorotate {
		autorotateFlag = 1
	}

	// call C wrapper: returns VipsImage* or NULL on error
	img := C.vips_image_new_from_file_autorotate(cFilePath, C.int(autorotateFlag))
	if img == nil {
		// read libvips error buffer
		errStr := C.GoString(C.vips_error_buffer())
		C.vips_error_clear()
		return 0, 0, fmt.Errorf("vips: %s", errStr)
	}
	// ensure we free the VipsImage when done
	defer C.g_object_unref(C.gpointer(img))

	w := int(C.vips_image_get_width(img))
	h := int(C.vips_image_get_height(img))

	return w, h, nil
}
