ghettowm
========

This is a work-in-progress ghetto window manager for Windows. I miss i3.

There are currently 3 "virtual desktops".
Not in the same sense as Windows Virtual Desktops, because I have no idea how those work.
Basically, you move between them with `ALT+(LEFT|RIGHT)`, left decrements, right increments.
Loops back around on both sides.

Right now, ghettowm uses VirtualDesktopAccessor to switch between Win10 virtual desktops,
rather than keeping track of per-desktop windows and hiding/unhiding them.

Hotkeys
=======

- `ALT+LEFT` - Move to desktop on the "left"
- `ALT+RIGHT` - Move to desktop on the "right"
- `ALT+SHIFT+Q` - Quit

Building
========

`GOOS=windows go build -ldflags -H=windowsgui`

Plans
=====

As of `42263c`, configuration is controlled by a lua file (`ghetto.lua`).
Keybinds and their handlers can be registered in code!

In the near future, I want to:

- Support pinning of windows
- Support moving windows between desktops
- Allow use of the WIN key
- Support defining layouts/tiling per desktop

Known Issues
============

- The Windows key (L or R) don't seem to trigger registered hotkeys.
- Likely various race conditions

Example Configuration
=====================

```
local bit32 = require('bit32')
local modkey = bit32.bor(get_mod_code('Alt'), get_mod_code('NoRepeat'))

ghettowm:RegisterHotkey(modkey, get_key_code('LeftArrow'), function()
  ghettowm:SwitchDesktopPrev()
end)

ghettowm:RegisterHotkey(modkey, get_key_code('RightArrow'), function()
  ghettowm:SwitchDesktopNext()
end)

ghettowm:RegisterHotkey(bit32.bor(modkey, get_mod_code('Control')), get_key_code('Q'), function()
  ghettowm:Quit()
end)

ghettowm:RegisterHotkey(modkey, get_key_code('One'), function()
  ghettowm:SwitchDesktop(0)
end)

ghettowm:RegisterHotkey(modkey, get_key_code('Two'), function()
  ghettowm:SwitchDesktop(1)
end)
```

Credits
=======

[Ciantic](https://github.com/Ciantic) - [VirtualDesktopAccessor](https://github.com/Ciantic/VirtualDesktopAccessor)
