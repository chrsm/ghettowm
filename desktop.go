package main

import (
	"log"

	"github.com/chrsm/winapi"
	"github.com/chrsm/winapi/user"
)

type desktop struct {
	head, tail *window
}

func (d *desktop) push(hwnd winapi.HWND) {
	w := &window{
		hwnd: hwnd,
	}

	if d.head == nil {
		d.head = w
	} else {
		d.tail.next = w
		w.prev = d.tail
	}

	d.tail = w
}

func (d *desktop) find(hwnd winapi.HWND) *window {
	for w := d.head; w != nil; w = w.next {
		if w.hwnd == hwnd {
			return w
		}
	}

	return nil
}

func (d *desktop) remove(hwnd winapi.HWND) bool {
	w := d.find(hwnd)
	if w == nil {
		return false
	}

	w.prev.next = w.next
	w.next.prev = w.prev

	w.next, w.prev = nil, nil

	return true
}

// This obviously doesn't work ;-)
func (d *desktop) Tile() {
	// hard-coded max size
	maxw := 2560
	maxh := 1080

	log.Printf("attempting to tile windows on desktop")

	var nw int
	for w := d.head; w != nil; w = w.next {
		nw++
	}

	if nw == 0 {
		log.Printf("no windows on this desktop, man")
		return
	}

	// how much can we give each window?
	maxiw := maxw / nw
	maxih := maxh / nw
	log.Printf("max window size: %dx%d", maxiw, maxih)
	for w := d.head; w != nil; w = w.next {
		user.SetWindowPos(w.hwnd, user.HwndTop, 0, 0, maxiw, maxih, user.SwpNoActivate)
	}
}
