package vips

/*
#cgo pkg-config: vips
#include <vips/vips.h>
#include <stdlib.h>

// wrapper: return VipsImage* (NULL on error).
static VipsImage* vips_image_new_from_file_simple(const char* name) {
    return vips_image_new_from_file(name, NULL);
}

// wrapper: return VipsImage* (NULL on error) from buffer.
static VipsImage* vips_image_new_from_buffer_simple(const void* buf, size_t len) {
    // Option string NULL, varargs terminated with NULL
    return vips_image_new_from_buffer(buf, len, NULL, NULL);
}

// get orientation (thin wrapper not strictly necessary but keeps C declarations together)
static int vips_image_get_orientation_wrapper(VipsImage* img) {
    return vips_image_get_orientation(img);
}
*/
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

// GetSizeParams holds optional parameters for getting image size.
type GetSizeParams struct {
	// Autorotate: if true, return width/height as if image were autorotated to upright.
	Autorotate bool
}

// GetImageSizeFromFile is a convenience wrapper without params.
func GetImageSizeFromFile(file string) (width, height int, err error) {
	return GetImageSizeFromFileWithParams(file, nil)
}

// GetImageSizeFromFileWithParams loads image from file and returns width/height.
func GetImageSizeFromFileWithParams(file string, params *GetSizeParams) (width, height int, err error) {
	startupIfNeeded() // 确保 libvips 初始化（你包里已有实现）

	if file == "" {
		return 0, 0, errors.New("govips: filename cannot be empty")
	}

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

// GetImageSizeFromBuffer is a convenience wrapper without params.
func GetImageSizeFromBuffer(buf []byte) (width, height int, err error) {
	return GetImageSizeFromBufferWithParams(buf, nil)
}

// GetImageSizeFromBufferWithParams loads image from memory and returns width/height.
// Respects Autorotate in GetSizeParams.
func GetImageSizeFromBufferWithParams(buf []byte, params *GetSizeParams) (width, height int, err error) {
	startupIfNeeded()

	if len(buf) == 0 {
		return 0, 0, errors.New("govips: buffer cannot be empty")
	}

	p := unsafe.Pointer(&buf[0])
	length := C.size_t(len(buf))

	img := C.vips_image_new_from_buffer_simple(p, length)
	if img == nil {
		errStr := C.GoString(C.vips_error_buffer())
		C.vips_error_clear()
		return 0, 0, fmt.Errorf("vips: %s", errStr)
	}
	defer C.g_object_unref(C.gpointer(img))

	w := int(C.vips_image_get_width(img))
	h := int(C.vips_image_get_height(img))

	if params != nil && params.Autorotate {
		orient := int(C.vips_image_get_orientation_wrapper(img))
		if orient >= 5 && orient <= 8 {
			w, h = h, w
		}
	}

	return w, h, nil
}
