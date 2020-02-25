INACTIVE
========

I am no longer working on this project. I only use Windows via KVM now and do
not have a need for this. I may revisit this project in the future but I do
not currently have the time.

For those looking for a tiling desktop manager for Windows, I recommend:

- [HashTWM](https://github.com/ZaneA/HashTWM)
- [bug.n](https://github.com/fuhsjr00/bug.n)
- [math0ne's tiling wm](https://github.com/math0ne/windows-tiling-window-manager)

HashTWM is native (written in C), while bug.n and math0ne's are both AutoHotKey.
I don't think any of them are designed to work with virtual desktops,
unfortunately.


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

Alternatively, you can download one of the [releases](https://github.com/chrsm/ghettowm/releases)

Plans
=====

As of `42263c`, configuration is controlled by a lua file (`ghetto.lua`).
Keybinds and their handlers can be registered in code!

In the near future, I want to:

- Support pinning of windows
- Support moving windows between desktops
- Support defining layouts/tiling per desktop

As of `a009caa`, WINKEY as a meta key works.
However, this breaks ALT/SHIFT or key combos with more than one other key.

I plan to reimplement hotkeys asap, but frankly everything here needs a bit of
cleanup so perhaps I will spend some time refactoring.

Known Issues
============

- Likely various race conditions
- WINKEY meta works, but now nothing else does ;)

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
