---
title: HP Mini 1030NR Gigabit Ethernet Adapter and Linux
slug: hp-mini-1030nr-gigabit-ethernet-adapter-and-linux
publishedAt: 2009-08-04T21:42:18-07:00
tags: [Linux, Gadgets]
summary: |
  I purchased an HP 1030NR, because, IMHO, they are the sexiest netbooks
  currently on the market.  They are sleek, lightweight, and they don't have a
  wonkey keyboard layout like some netbooks (I'm looking at you, ASUS).
---
I purchased an HP 1030NR, because, IMHO, they are the sexiest netbooks
currently on the market.  They are sleek, lightweight, and they don't have a
wonkey keyboard layout like some netbooks (I'm looking at you, ASUS).

One problem: the onboard nic controller (Marvell 88E8040 adapter, sky2 kernel
module) doesn't quite work right on a default install of Ubuntu 9.04.
Sometimes it works, sometimes it's not even detected (you can't see it listed
when you type @lspci@).  Sometimes it's detected, and works, but locks the
machine up if you unplug your network cable.  Not cool.

After some digging around, I found that there is an easy way to fix this.  Edit
`/boot/grub/menu.list`, find the line that says something like this: `#
defoptions=quiet splash`, add `acpi_os_name=Linux` at the end of the line, save
the file, and type `update-grub`.  Reboot, and viola, it should work normally
from now on.

<div class="restored-from-archive">
  <h3>Restored from VimTips archive</h3>
  <p>
  This article was restored from the VimTips archive. There's probably
  missing images and broken links (and even some flash references), but it
  was still important to me to bring them back.
  </p>
</div>
