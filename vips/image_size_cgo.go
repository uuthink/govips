package vips

/*
#cgo pkg-config: vips
#include <vips/vips.h>
#include <stdlib.h>

// wrapper: return VipsImage* (NULL on error).
static VipsImage* vips_image_new_from_file_simple(const char* name) {
    return vips_image_new_from_file(name, NULL);
}

// get orientation (thin wrapper not strictly necessary but keeps C declarations together)
static int vips_image_get_orientation_wrapper(VipsImage* img) {
    return vips_image_get_orientation(img);
}
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// GetSizeParams holds optional parameters for getting image size.
type GetSizeParams struct {
	// Autorotate: if true, return width/height as if image were autorotated to upright.
	Autorotate bool
}

// Convenience
func GetImageSizeFromFile(file string) (width, height int, err error) {
	return GetImageSizeFromFileWithParams(file, nil)
}

// Main implementation
func GetImageSizeFromFileWithParams(file string, params *GetSizeParams) (width, height int, err error) {
	startupIfNeeded() // 确保 libvips 初始化（你包里已有实现）

	cFilePath := C.CString(file)
	defer C.free(unsafe.Pointer(cFilePath))

	// always load without loader autorotate to avoid loader-specific option errors
	img := C.vips_image_new_from_file_simple(cFilePath)
	if img == nil {
		errStr := C.GoString(C.vips_error_buffer())
		C.vips_error_clear()
		return 0, 0, fmt.Errorf("vips: %s", errStr)
	}
	defer C.g_object_unref(C.gpointer(img))

	w := int(C.vips_image_get_width(img))
	h := int(C.vips_image_get_height(img))

	// if caller wants autorotate semantics, inspect orientation metadata
	if params != nil && params.Autorotate {
		orient := int(C.vips_image_get_orientation_wrapper(img))
		// orientation 5,6,7,8 imply a 90/270-degree rotation (swap dims)
		if orient >= 5 && orient <= 8 {
			w, h = h, w
		}
	}

	return w, h, nil
}
