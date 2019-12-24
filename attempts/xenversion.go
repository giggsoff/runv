package main

/*
#cgo LDFLAGS: -lxenlight -lxentoollog
#include <stdio.h>
#include <stdlib.h>
#include <libxl.h>
static inline int testNil(const char ***test_input)
{
   int i = 0;
   if (*test_input)
        while ((*test_input)[i]){
			printf("%s\n", (*test_input)[i]);
			printf("%d\n", i);
            i++;
		}
   return i;
}
*/
import "C"
import (
	"fmt"
	"log"
	"strconv"
	"unsafe"
)

/*
 * Other flags that may be needed at some point:
 *  -lnl-route-3 -lnl-3
 *
 * To get back to static linking:
 * #cgo LDFLAGS: -lxenlight -lyajl_s -lxengnttab -lxenstore -lxenguest -lxentoollog -lxenevtchn -lxenctrl -lblktapctl -lxenforeignmemory -lxencall -lz -luuid -lutil
 */

type Error int

const (
	ErrorNonspecific                  = Error(-C.ERROR_NONSPECIFIC)
	ErrorVersion                      = Error(-C.ERROR_VERSION)
	ErrorFail                         = Error(-C.ERROR_FAIL)
	ErrorNi                           = Error(-C.ERROR_NI)
	ErrorNomem                        = Error(-C.ERROR_NOMEM)
	ErrorInval                        = Error(-C.ERROR_INVAL)
	ErrorBadfail                      = Error(-C.ERROR_BADFAIL)
	ErrorGuestTimedout                = Error(-C.ERROR_GUEST_TIMEDOUT)
	ErrorTimedout                     = Error(-C.ERROR_TIMEDOUT)
	ErrorNoparavirt                   = Error(-C.ERROR_NOPARAVIRT)
	ErrorNotReady                     = Error(-C.ERROR_NOT_READY)
	ErrorOseventRegFail               = Error(-C.ERROR_OSEVENT_REG_FAIL)
	ErrorBufferfull                   = Error(-C.ERROR_BUFFERFULL)
	ErrorUnknownChild                 = Error(-C.ERROR_UNKNOWN_CHILD)
	ErrorLockFail                     = Error(-C.ERROR_LOCK_FAIL)
	ErrorJsonConfigEmpty              = Error(-C.ERROR_JSON_CONFIG_EMPTY)
	ErrorDeviceExists                 = Error(-C.ERROR_DEVICE_EXISTS)
	ErrorCheckpointDevopsDoesNotMatch = Error(-C.ERROR_CHECKPOINT_DEVOPS_DOES_NOT_MATCH)
	ErrorCheckpointDeviceNotSupported = Error(-C.ERROR_CHECKPOINT_DEVICE_NOT_SUPPORTED)
	ErrorVnumaConfigInvalid           = Error(-C.ERROR_VNUMA_CONFIG_INVALID)
	ErrorDomainNotfound               = Error(-C.ERROR_DOMAIN_NOTFOUND)
	ErrorAborted                      = Error(-C.ERROR_ABORTED)
	ErrorNotfound                     = Error(-C.ERROR_NOTFOUND)
	ErrorDomainDestroyed              = Error(-C.ERROR_DOMAIN_DESTROYED)
	ErrorFeatureRemoved               = Error(-C.ERROR_FEATURE_REMOVED)
)

var errors = [...]string{
	ErrorNonspecific:                  "Non-specific error",
	ErrorVersion:                      "Wrong version",
	ErrorFail:                         "Failed",
	ErrorNi:                           "Not Implemented",
	ErrorNomem:                        "No memory",
	ErrorInval:                        "Invalid argument",
	ErrorBadfail:                      "Bad Fail",
	ErrorGuestTimedout:                "Guest timed out",
	ErrorTimedout:                     "Timed out",
	ErrorNoparavirt:                   "No Paravirtualization",
	ErrorNotReady:                     "Not ready",
	ErrorOseventRegFail:               "OS event registration failed",
	ErrorBufferfull:                   "Buffer full",
	ErrorUnknownChild:                 "Unknown child",
	ErrorLockFail:                     "Lock failed",
	ErrorJsonConfigEmpty:              "JSON config empty",
	ErrorDeviceExists:                 "Device exists",
	ErrorCheckpointDevopsDoesNotMatch: "Checkpoint devops does not match",
	ErrorCheckpointDeviceNotSupported: "Checkpoint device not supported",
	ErrorVnumaConfigInvalid:           "VNUMA config invalid",
	ErrorDomainNotfound:               "Domain not found",
	ErrorAborted:                      "Aborted",
	ErrorNotfound:                     "Not found",
	ErrorDomainDestroyed:              "Domain destroyed",
	ErrorFeatureRemoved:               "Feature removed",
}

func (e Error) Error() string {
	if 0 < int(e) && int(e) < len(errors) {
		s := errors[e]
		if s != "" {
			return s
		}
	}
	return fmt.Sprintf("libxl error: %d", -e)

}

type Context struct {
	ctx    *C.libxl_ctx
	logger *C.xentoollog_logger_stdiostream
}

type VersionInfo struct {
	XenVersionMajor int
	XenVersionMinor int
	XenVersionExtra string
	Compiler        string
	CompileBy       string
	CompileDomain   string
	CompileDate     string
	Capabilities    string
	Changeset       string
	VirtStart       uint64
	Pagesize        int
	Commandline     string
	BuildId         string
}

func (cinfo *C.libxl_version_info) toGo() (info *VersionInfo) {
	info = &VersionInfo{}
	info.XenVersionMajor = int(cinfo.xen_version_major)
	info.XenVersionMinor = int(cinfo.xen_version_minor)
	info.XenVersionExtra = C.GoString(cinfo.xen_version_extra)
	info.Compiler = C.GoString(cinfo.compiler)
	info.CompileBy = C.GoString(cinfo.compile_by)
	info.CompileDomain = C.GoString(cinfo.compile_domain)
	info.CompileDate = C.GoString(cinfo.compile_date)
	info.Capabilities = C.GoString(cinfo.capabilities)
	info.Changeset = C.GoString(cinfo.changeset)
	info.VirtStart = uint64(cinfo.virt_start)
	info.Pagesize = int(cinfo.pagesize)
	info.Commandline = C.GoString(cinfo.commandline)
	info.BuildId = C.GoString(cinfo.build_id)

	return
}

func (Ctx *Context) CheckOpen() (err error) {
	if Ctx.ctx == nil {
		err = fmt.Errorf("Context not opened")
	}
	return
}

func (Ctx *Context) Open() (err error) {
	if Ctx.ctx != nil {
		return
	}

	Ctx.logger = C.xtl_createlogger_stdiostream(C.stderr, C.XTL_ERROR, 0)
	if Ctx.logger == nil {
		err = fmt.Errorf("Cannot open stdiostream")
		return
	}

	ret := C.libxl_ctx_alloc(&Ctx.ctx, C.LIBXL_VERSION,
		0, (*C.xentoollog_logger)(unsafe.Pointer(Ctx.logger)))

	if ret != 0 {
		err = Error(-ret)
	}
	return
}

func (Ctx *Context) GetVersionInfo() (info *VersionInfo, err error) {
	err = Ctx.CheckOpen()
	if err != nil {
		return
	}

	var cinfo *C.libxl_version_info

	cinfo = C.libxl_get_version_info(Ctx.ctx)

	info = cinfo.toGo()

	return
}
func (Ctx *Context) Close() (err error) {
	ret := C.libxl_ctx_free(Ctx.ctx)
	Ctx.ctx = nil

	if ret != 0 {
		err = Error(-ret)
	}
	C.xtl_logger_destroy((*C.xentoollog_logger)(unsafe.Pointer(Ctx.logger)))
	return
}

func main() {
	var Ctx Context
	err := Ctx.Open()
	if err != nil {
		log.Fatal(err)
	}
	v, err := Ctx.GetVersionInfo()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Version: %d %d %s", v.XenVersionMajor, v.XenVersionMinor, v.XenVersionExtra)
	err = Ctx.Close()
	if err != nil {
		log.Fatal(err)
	}
	l := 15
	cExtra := (**C.char)(C.malloc(C.ulong(l+1) * C.ulong(unsafe.Sizeof(uintptr(0)))))
	defer C.free(unsafe.Pointer(cExtra))
	arrayPtr := (*[1 << 30]*C.char)(unsafe.Pointer(cExtra))[0 : l+1 : l+1]
	for i := 0; i < l; i++ {
		cstr := C.CString("testing" + strconv.Itoa(i))
		defer C.free(unsafe.Pointer(cstr))
		arrayPtr[i] = (*C.char)(cstr)
	}
	arrayPtr[l] = nil
	res := C.testNil((***C.char)(unsafe.Pointer(&cExtra)))
	fmt.Printf("Test: %d", int(res))
}
