package main

import (
	"errors"
	"flag"
	"log"
)

/*
#cgo pkg-config: python-3.8
#cgo LDFLAGS: -L/usr/lib/python3.8/config-3.8-x86_64-linux-gnu -L/usr/lib -lpython3.8 -lcrypt -lpthread -ldl  -lutil -lm -lm
#define PY_SSIZE_T_CLEAN
#define SIZEOF_WCHAR_T 4
#include <Python.h>

static int Run(char *name){
PyObject *obj = Py_BuildValue("s", name);
FILE *f = _Py_fopen_obj(obj, "r+");
	int res = PyRun_SimpleFileEx(f,name,1);
	return res;
}
*/
import "C"

var (
	file *string
	str  *string
)

func argparse() error {
	file = flag.String("f", "", "File to use")
	str = flag.String("c", "", "String to run")
	flag.Parse()

	if *str == "" && *file == "" {
		return errors.New("Need eighter string or file to execute")
	}

	if *str != "" && *file != "" {
		return errors.New("Need eighter string or file to execute, not both")
	}

	return nil
}

func main() {
	C.Py_Initialize()
	defer C.Py_Finalize()

	err := argparse()
	if err != nil {
		log.Fatal(err)
		flag.PrintDefaults()
		return
	}

	if *file != "" {
		cstr := C.CString(*file)
		C.Run(cstr)
	}

	if *str != "" {
		C.PyRun_SimpleString(C.CString(*str))
	}
}
