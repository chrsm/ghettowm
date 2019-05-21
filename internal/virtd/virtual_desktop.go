// Package virtd implements APIs from Ciantic's VirtualDesktopAccessor.
// more info about VDA can be found at: https://github.com/Ciantic/VirtualDesktopAccessor
package virtd

import (
	"syscall"

	"github.com/chrsm/winapi"
)

var (
	vdapi = syscall.NewLazyDLL("VirtualDesktopAccessor.dll")

	pGetCurrentDesktopNumber         = vdapi.NewProc("GetCurrentDesktopNumber")
	pGetDesktopCount                 = vdapi.NewProc("GetDesktopCount")
	pGetDesktopIdByNumber            = vdapi.NewProc("GetDesktopIdByNumber")
	pGetDesktopNumber                = vdapi.NewProc("GetDesktopNumber")
	pGetDesktopNumberById            = vdapi.NewProc("GetDesktopNumberById")
	pGetWindowDesktopId              = vdapi.NewProc("GetWindowDesktopId")
	pGetWindowDesktopNumber          = vdapi.NewProc("GetWindowDesktopNumber")
	pIsWindowOnCurrentVirtualDesktop = vdapi.NewProc("IsWindowOnCurrentVirtualDesktop")
	pMoveWindowToDesktopNumber       = vdapi.NewProc("MoveWindowToDesktopNumber")
	pGoToDesktopNumber               = vdapi.NewProc("GoToDesktopNumber")
	pRegisterPostMessageHook         = vdapi.NewProc("RegisterPostMessageHook")
	pUnregisterPostMessageHook       = vdapi.NewProc("UnregisterPostMessageHook")
	pIsPinnedWindow                  = vdapi.NewProc("IsPinnedWindow")
	pPinWindow                       = vdapi.NewProc("PinWindow")
	pUnPinWindow                     = vdapi.NewProc("UnPinWindow")
	pIsPinnedApp                     = vdapi.NewProc("IsPinnedApp")
	pPinApp                          = vdapi.NewProc("PinApp")
	pUnPinApp                        = vdapi.NewProc("UnPinApp")
	pIsWindowOnDesktopNumber         = vdapi.NewProc("IsWindowOnDesktopNumber")
	pRestartVirtualDesktopAccessor   = vdapi.NewProc("RestartVirtualDesktopAccessor")
)

func GetCurrentDesktopNumber() int {
	ret, _, _ := pGetCurrentDesktopNumber.Call()

	return int(ret)
}

func GetDesktopCount() int {
	ret, _, _ := pGetDesktopCount.Call()

	return int(ret)
}

func GetWindowDesktopNumber(w winapi.HWND) int {
	ret, _, _ := pGetWindowDesktopNumber.Call(uintptr(w))

	return int(int32(ret))
}

func IsWindowOnCurrentVirtualDesktop(w winapi.HWND) bool {
	ret, _, _ := pIsWindowOnCurrentVirtualDesktop.Call(uintptr(w))

	return ret == 1
}

func MoveWindowToDesktopNumber(w winapi.HWND, i int) bool {
	ret, _, _ := pMoveWindowToDesktopNumber.Call(uintptr(w), uintptr(i))

	return ret == 1
}

func GoToDesktopNumber(i int) {
	pGoToDesktopNumber.Call(uintptr(i))
}

func RegisterPostMessageHook(l winapi.HWND, offset int) {
	pRegisterPostMessageHook.Call(uintptr(l), uintptr(offset))
}

func UnregisterPostMessageHook(l winapi.HWND) {
	pUnregisterPostMessageHook.Call(uintptr(l))
}

func IsPinnedWindow(w winapi.HWND) bool {
	ret, _, _ := pIsPinnedWindow.Call(uintptr(w))

	return ret == 1
}

func PinWindow(w winapi.HWND) {
	pPinWindow.Call(uintptr(w))
}

func UnpinWindow(w winapi.HWND) {
	pUnPinWindow.Call(uintptr(w))
}

func IsPinnedApp(w winapi.HWND) bool {
	ret, _, _ := pIsPinnedApp.Call(uintptr(w))

	return ret == 1
}

func PinApp(w winapi.HWND) {
	pPinApp.Call(uintptr(w))
}

func UnpinApp(w winapi.HWND) {
	pUnPinApp.Call(uintptr(w))
}

func IsWindowOnDesktopNumber(w winapi.HWND, i int) bool {
	ret, _, _ := pIsWindowOnDesktopNumber.Call(uintptr(w), uintptr(i))

	return ret == 1
}

func RestartVirtualDesktopAccessor() {
	pRestartVirtualDesktopAccessor.Call()
}

func GetDesktopIdByNumber(i int) winapi.GUID {
	panic("not implemented: broken")
	// ret, _, _ := pGetDesktopIdByNumber.Call(uintptr(i))

	return winapi.GUID{}
}

func GetDesktopNumberById(id winapi.GUID) int {
	panic("not implemented: broken")

	return 0
	/*ret, _, _ := pGetDesktopNumberById.Call(id)

	return int(ret)*/
}
