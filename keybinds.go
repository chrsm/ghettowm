package ghettowm

import (
	"log"
	"strings"
	"sync/atomic"
	"unsafe"

	"github.com/chrsm/winapi"
	"github.com/chrsm/winapi/user"
)

type keybinds struct {
	keys map[string]*keybind
	len  int
}

type keybind struct {
	id int
	cb func(int)
}

func newKeybinds() *keybinds {
	k := &keybinds{
		keys: make(map[string]*keybind),
	}

	return k
}

func (k *keybinds) set(modifiers user.ModKey, vkey user.VirtualKey, cb func(int)) {
	k.len++
	k.keys[keyString(modifiers, vkey)] = &keybind{
		id: k.len,
		cb: cb,
	}
}

func (k *keybinds) get(modifiers user.ModKey, vkey user.VirtualKey) *keybind {
	str := keyString(modifiers, vkey)

	if v, ok := k.keys[str]; ok {
		return v
	}

	return nil
}

func (k *keybinds) getByID(id int) *keybind {
	for i := range k.keys {
		if k.keys[i].id == id {
			return k.keys[i]
		}
	}

	return nil
}

func keyString(modifiers user.ModKey, vkey user.VirtualKey) string {
	// translate modifiers + vkey into a string for comparison
	var bits []string

	mods := []user.ModKey{user.ModAlt, user.ModControl, user.ModShift, user.ModWin}
	for i := range mods {
		if modifiers&mods[i] == mods[i] {
			bits = append(bits, user.GetModKeyName(mods[i]))
		}
	}

	bits = append(bits, user.GetVirtualKeyName(vkey))

	return strings.Join(bits, "+")
}

var winkeyState atomic.Value // 0 = not pressed, 1 = down

func (gwm *WindowManager) keyboardHook(code int, wp winapi.WPARAM, lp winapi.LPARAM) uintptr {
	info := *(*user.KBDLLHOOKSTRUCT)(unsafe.Pointer(lp))

	// log.Printf("code(%d), wp(%d), lp(%d): %#v", code, wp, lp, info)

	// swallow it, but make note of the key state.
	if info.VKCode == user.VirtKeyLeftWin || info.VKCode == user.VirtKeyRightWin {
		if wp == user.WmKeyDown || wp == user.WmSysKeyDown {
			winkeyState.Store(1)
		} else {
			winkeyState.Store(0)
		}

		return 1
	}

	if v := winkeyState.Load().(int); v != 0 {
		// check if we have a keybind registered for the key being pressed
		kb := gwm.keybinds.get(user.ModWin, user.VirtualKey(info.VKCode))
		if kb != nil {
			log.Printf("received a callback for %s", keyString(user.ModWin, user.VirtualKey(info.VKCode)))
			// gwm.cbch <- kb
			user.SendMessage(gwm.hwnd, user.WmHotkey, winapi.WPARAM(kb.id), 0)

			// swallow it
			return 1
		} else {
			log.Printf("no keybind registered: %s", keyString(user.ModWin, user.VirtualKey(info.VKCode)))
		}
	}

	return user.CallNextHook(0, code, wp, lp)
}
