ghettowm
========

[![CircleCI](https://circleci.com/gh/chrsm/ghettowm.svg?style=svg)](https://circleci.com/gh/chrsm/ghettowm)

This is a work-in-progress ghetto window manager for Windows. I miss i3.

ghettowm plugs in to Windows 10's native virtual desktop support by way of
VirtualDesktopAccessor, rather than keeping track of per-desktop windows
and hiding/unhiding them as it did initially.

The benefit with using the IVirtualDesktopManager is that in the event
of an issue with ghettowm, your setup is still running 100% fine.

This may present some issues in the future as IVirtualDesktopManager has
almost no functionality - the useful stuff is in IVirtualDesktopManagerInternal,
which is undocumented and likely to change. As such, ghettowm is completely
tied to VirtualDesktopAccessor.

Building
========

`GOOS=windows go build -ldflags -H=windowsgui`

Alternatively, you can download one of the [releases](/chrsm/ghettowm/releases)

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

- The Windows key (L or R) don't seem to trigger registered hotkeys. It is not
recommended to use it until this is fixed.
- Likely various race conditions

Example Configuration
=====================

```
-- ghettowm config
local bit32 = require('bit32')
local windows = require('windows')
local ghetto_util = require('ghetto_util')

local modkey = bit32.bor(ghetto_util.get_mod('Alt'), ghetto_util.get_mod('NoRepeat'))

ghettowm:RegisterHotkey(modkey, ghetto_util.get_key('LeftArrow'), function()
  ghettowm:SwitchDesktopPrev()
end)

ghettowm:RegisterHotkey(modkey, ghetto_util.get_key('RightArrow'), function()
  ghettowm:SwitchDesktopNext()
end)

ghettowm:RegisterHotkey(bit32.bor(modkey, ghetto_util.get_mod('Control')), ghetto_util.get_key('Q'), function()
  ghettowm:Quit()
end)

ghettowm:RegisterHotkey(modkey, ghetto_util.get_key('One'), function()
  ghettowm:SwitchDesktop(0)
end)

ghettowm:RegisterHotkey(modkey, ghetto_util.get_key('Two'), function()
  ghettowm:SwitchDesktop(1)
end)

ghettowm:RegisterHotkey(modkey, ghetto_util.get_key('H'), function()
  ghettowm:NextWindow()
end)

ghettowm:RegisterHotkey(modkey, ghetto_util.get_key('L'), function()
  ghettowm:PrevWindow()
end)
```

Credits
=======

[Ciantic](https://github.com/Ciantic) - [VirtualDesktopAccessor](https://github.com/Ciantic/VirtualDesktopAccessor)
