package main

type keybinds struct {
	set map[int]*keybind
	len int
}

type keybind struct {
	id int
	cb func(int)
}
