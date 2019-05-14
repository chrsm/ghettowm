ghettowm
========

This is a work-in-progress ghetto window manager for Windows. I miss i3.

There are currently 3 "virtual desktops".
Not in the same sense as Windows Virtual Desktops, because I have no idea how those work.
Basically, you move between them with `ALT+(LEFT|RIGHT)`, left decrements, right increments.
Loops back around on both sides.

I'm sure WVD uses more complicated methods, but ghettowm just sets `SW_HIDE` on windows not
in the current "desktop".

Hotkeys
=======

`ALT+LEFT` - Move to desktop on the "left"
`ALT+RIGHT` - Move to desktop on the "right"
`ALT+SHIFT+Q` - Quit

Building
========

`GOOS=windows go build -ldflags -H=windowsgui`

Plans
=====

I plan on expanding this in the future, namely to support configuration outside of recompiling the binary..

