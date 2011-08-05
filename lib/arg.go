// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

import (
	"fmt"
	"os"
	"unsafe"
)

func packValues(values []interface{}) ([]unsafe.Pointer, os.Error) {
	cValues := make([]unsafe.Pointer, len(values))
	for i, value := range values {
		switch v := value.(type) {
		case *string:
			cValues[i] = unsafe.Pointer(new(*C.char))
		case *Object:
			cValues[i] = unsafe.Pointer(new(*C.PyObject))
		case *int:
			cValues[i] = unsafe.Pointer(new(C.int))
		case *int8:
			cValues[i] = unsafe.Pointer(new(C.int8_t))
		case *int16:
			cValues[i] = unsafe.Pointer(new(C.int16_t))
		case *int32:
			cValues[i] = unsafe.Pointer(new(C.int32_t))
		case *int64:
			cValues[i] = unsafe.Pointer(new(C.int64_t))
		case *uint:
			cValues[i] = unsafe.Pointer(new(C.uint))
		case *uint8:
			cValues[i] = unsafe.Pointer(new(C.uint8_t))
		case *uint16:
			cValues[i] = unsafe.Pointer(new(C.uint16_t))
		case *uint32:
			cValues[i] = unsafe.Pointer(new(C.uint32_t))
		case *uint64:
			cValues[i] = unsafe.Pointer(new(C.uint64_t))
		case *float32:
			cValues[i] = unsafe.Pointer(new(C.float))
		case *float64:
			cValues[i] = unsafe.Pointer(new(C.double))
		default:
			return nil, fmt.Errorf("Unsupported type: %T", v)
		}
	}
	return cValues, nil
}

func unpackValues(cValues []unsafe.Pointer, values []interface{}) os.Error {
	for i, value := range values {
		switch v := value.(type) {
		case *string:
			*v = C.GoString(*(**C.char)(cValues[i]))
		case *Object:
			*v = newBaseObject(*(**C.PyObject)(cValues[i])).actual()
		case *int:
			*v = int(*(*C.int)(cValues[i]))
		case *int8:
			*v = int8(*(*C.int8_t)(cValues[i]))
		case *int16:
			*v = int16(*(*C.int16_t)(cValues[i]))
		case *int32:
			*v = int32(*(*C.int32_t)(cValues[i]))
		case *int64:
			*v = int64(*(*C.int64_t)(cValues[i]))
		case *uint:
			*v = uint(*(*C.uint)(cValues[i]))
		case *uint8:
			*v = uint8(*(*C.uint8_t)(cValues[i]))
		case *uint16:
			*v = uint16(*(*C.uint16_t)(cValues[i]))
		case *uint32:
			*v = uint32(*(*C.uint32_t)(cValues[i]))
		case *uint64:
			*v = uint64(*(*C.uint64_t)(cValues[i]))
		case *float32:
			*v = float32(*(*C.float)(cValues[i]))
		case *float64:
			*v = float64(*(*C.double)(cValues[i]))
		default:
			return fmt.Errorf("Unsupported type: %T", v)
		}
	}
	return nil
}

func Arg_ParseTuple(args *Tuple, format string, values ...interface{}) os.Error {
	if args == nil {
		return fmt.Errorf("Arg_ParseTuple: args was nil")
	}

	cValues, err := packValues(values)
	if err != nil {
		return nil
	}

	f := C.CString(format)
	defer C.free(unsafe.Pointer(f))

	ret := C.doParseTuple(c(args), f, &cValues[0], C.int(len(cValues)))
	if ret == 0 {
		return exception()
	}

	return unpackValues(cValues, values)
}

func Arg_ParseTupleAndKeywords(args *Tuple, kw *Dict, format string, kwlist []string, values ...interface{}) os.Error {
	if args == nil {
		return fmt.Errorf("Arg_ParseTuple: args was nil")
	}

	cValues, err := packValues(values)
	if err != nil {
		return nil
	}

	f := C.CString(format)
	defer C.free(unsafe.Pointer(f))

	klist := make([]*C.char, len(kwlist)+1)

	for i, k := range kwlist {
		klist[i] = C.CString(k)
		defer C.free(unsafe.Pointer(klist[i]))
	}

	ret := C.doParseTupleKwds(c(args), c(kw), f, &klist[0], &cValues[0], C.int(len(cValues)))
	if ret == 0 {
		return exception()
	}

	return unpackValues(cValues, values)
}

func BuildValue(format string, values ...interface{}) (Object, os.Error) {
	cValues := make([]C.ArgValue, len(values))
	for i, value := range values {
		switch v := value.(type) {
		case string:
			s := C.CString(v)
			defer C.free(unsafe.Pointer(s))
			cValues[i]._type = &C.ffi_type_pointer
			cValues[i].value = unsafe.Pointer(&s)
		case Object:
			b := v.Base()
			cValues[i]._type = &C.ffi_type_pointer
			cValues[i].value = unsafe.Pointer(&b)
		case int:
			iv := C.int(v)
			cValues[i]._type = &C.ffi_type_sint
			cValues[i].value = unsafe.Pointer(&iv)
		case int8:
			iv := C.int8_t(v)
			cValues[i]._type = &C.ffi_type_sint8
			cValues[i].value = unsafe.Pointer(&iv)
		case int16:
			iv := C.int16_t(v)
			cValues[i]._type = &C.ffi_type_sint16
			cValues[i].value = unsafe.Pointer(&iv)
		case int32:
			iv := C.int32_t(v)
			cValues[i]._type = &C.ffi_type_sint32
			cValues[i].value = unsafe.Pointer(&iv)
		case int64:
			iv := C.int64_t(v)
			cValues[i]._type = &C.ffi_type_sint64
			cValues[i].value = unsafe.Pointer(&iv)
		case uint:
			iv := C.uint(v)
			cValues[i]._type = &C.ffi_type_uint
			cValues[i].value = unsafe.Pointer(&iv)
		case uint8:
			iv := C.uint8_t(v)
			cValues[i]._type = &C.ffi_type_uint8
			cValues[i].value = unsafe.Pointer(&iv)
		case uint16:
			iv := C.uint16_t(v)
			cValues[i]._type = &C.ffi_type_uint16
			cValues[i].value = unsafe.Pointer(&iv)
		case uint32:
			iv := C.uint32_t(v)
			cValues[i]._type = &C.ffi_type_uint32
			cValues[i].value = unsafe.Pointer(&iv)
		case uint64:
			iv := C.uint64_t(v)
			cValues[i]._type = &C.ffi_type_uint64
			cValues[i].value = unsafe.Pointer(&iv)
		case float32:
			fv := C.float(v)
			cValues[i]._type = &C.ffi_type_float
			cValues[i].value = unsafe.Pointer(&fv)
		case float64:
			fv := C.double(v)
			cValues[i]._type = &C.ffi_type_double
			cValues[i].value = unsafe.Pointer(&fv)
		default:
			return nil, fmt.Errorf("Unsupported type: %T", v)
		}
	}
	f := C.CString(format)
	defer C.free(unsafe.Pointer(f))
	ret := C.doBuildValue(f, &cValues[0], C.int(len(cValues)))
	if ret == nil {
		return nil, exception()
	}
	return newBaseObject(ret).actual(), nil
}
