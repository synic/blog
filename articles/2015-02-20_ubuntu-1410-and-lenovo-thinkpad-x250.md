<!-- :metadata:

title: Ubuntu 14.10 and the Lenovo ThinkPad x250
tags: Linux
published: 2015-02-20T22:12:40-0700
summary:

I recently purchased and received a [Lenovo ThinkPad
x250](http://shop.lenovo.com/us/en/laptops/thinkpad/x-series/x250/) and
immediately proceeded to install [Ubuntu](http://ubuntu.org) over whatever
version of Windows it came with.

-->

I recently purchased and received a [Lenovo ThinkPad
x250](http://shop.lenovo.com/us/en/laptops/thinkpad/x-series/x250/) and
immediately proceeded to install [Ubuntu](http://ubuntu.org) over whatever
version of Windows it came with.

Unfortunately, not everything works out of the box (multiple monitors using the
ultradock, sound with external speakers and the ultradock, 3d video, and the
trackpoint buttons). It's not very surprising, as these machines were just
released in the US last week, and I'm sure the next version of Ubuntu will have
much better support for them. In the mean time, however, I did get some things
to work, and I'm providing some instructions for those here.

### Multiple Monitors using the UltraDock
I've always loved ThinkPads, and most of my laptops have been one, but it's
only been recently that they've been enough to replace my desktop needs. That
happened for me because of the UltraDock, which allows you to use multiple
monitors with your laptop. On my T410 (which had a different ultradock), it
just worked out of the box. With the x250 and it's dock, X would see both
monitors as one large monitor. This is because the new dock uses a DisplayPort
1.2 feature, which didn't make it into the Linux kernel until 3.17 (Ubuntu
14.10 uses 3.16). Following <a
href='http://cweiske.de/tagebuch/thinkpad-ultradock-screens.htm'>this guide</a>
fixed the problem for me, with one change. I used the 3.18.7 kernel instead of
3.17.1 as he describes in the article, because 3.18 just happens to help fix
the trackpoint button problems as well.

### Trackpoint Buttons
Doing what I describe here will make the trackpoint buttons work (including
middle click scrolling), but it will completely disable the touchpad. I'm ok
with that, touchpads aren't my favorite. First, follow [this
guide](http://cweiske.de/tagebuch/thinkpad-ultradock-screens.htm) to update
your kernel, but instead of 3.17.1, use kernel 3.18.7 instead. (If you're
feeling brave, I've got precompiled debs
[here](http://synicworld.com/media/debs/). Edit `/etc/modprobe.d/psmouse.conf`
and add the following line:

`options psmouse proto=imps`

**Update**:  Putting options psmouse=imps in /etc/modprobe.d/psmouse.conf seems
to do nothign. Instead, add `psmouse.proto=imps` to your kernel cmdline, as
directed below in the "Volume Control and Backlight Brightness Keys".

Then, create a file at `/usr/share/X11/xorg.conf.d/90-evdev.conf`

```
Section "InputClass"
  Identifier "Touchpad/TrackPoint"
  MatchProduct "PS/2 Synaptics TouchPad"
  MatchDriver "evdev"
  Option "EmulateWheel" "1"
  Option "EmulateWheelButton" "2"
  Option "Emulate3Buttons" "0"
  Option "XAxisMapping" "6 7"
  Option "YAcisMapping" "4 5"
EndSection
```

Restart your X server, and your trackpoint buttons should work as expected.

### Volume Control and Backlight Brightness Keys

To get the volume keys to work, I had to pass `acpi_osi=Linux` to the kernel
boot options. Also, to get backlight control to work, add
`acpi_backlight=vendor` as well. In `/etc/default/grub`, my
`GRUB_CMDLINE_LINUX_DEFAULT` looks like this:

```
GRUB_CMDLINE_LINUX_DEFAULT="quiet splash acpi_osi=Linux acpi_backlight=vendor psmouse.proto=imps"
```

After making those changes, run `update-grub` and reboot. This made the
multimedia keys work for me, but the brightness controls still did not work. I
installed a program called `xbacklight` and mapped the brightness controls to
`xbacklight -dec 10` and `xbacklight -inc 10`. There's some more work that
could be done here (like automatically setting it to full brightness on AC
power, etc), but this is good for now.

### In Progress...

1. The audio out jack on the Ultradock doesn't work at all, though the speakers
   on the laptop and the headphone port do work.

### Update, Feb 21, 2015

I found the Padoka PPA here:
<https://launchpad.net/~paulo-miguel-dias/+archive/ubuntu/mesa/>. Adding this
PPA, running `apt-get update && apt-get dist-upgrade` <strong>eliminates the
need to compile xf86-video-intel as directed in the "Multiple Monitors using
the UltraDock" instructions, and makes 3d video work .</strong>

Still not working: the sound port on the UltraDock.

### Update, Mar 01, 2015

**I got the audio port on the ultradock working!**

To make this work, create a file at `/lib/firmware/x250.fw` with the following
contents:

```
[codec]
0x10ec0292 0x17aa2226 0

[pincfg]
0x16 0x21211010
0x19 0x21a11010
```

... then, create a file called `/etc/modprobe.d/hda-intel.conf` with the
following contents:

```
options snd-hda-intel patch=x250.fw,x250.fw,x250.fw
```

Reboot, and the audio port on your dock should work. I imagine the process is
the same for any of 2015 thinkpads. With each thinkpad, the only number that
changes is the second one under the `[codec]` heading (in this case,
`0x17aa2226`). That information can be found by downloading this script:
<http://www.alsa-project.org/alsa-info.sh>, and running it via `bash
alsa-info.sh` as root. Search fo **0x17aa** in the output, and you'll find the
other half of the number you need.


### Update Mar 2, 2015 - Fingerprint scanner works

Following the [instructions found in this
gist](https://gist.github.com/foosel/3abd45bc1b6ae121965b), I was able to get
the fingerprint scanner to work. It doesn't work very well. I often have to try
more than once to get it to recognize my fingerprint, and there doesn't seem to
be any way to scan multiple images of your finger to improve reliability.

Coupled with the fact that the things are easy to fool, I will be disabling it.

<div class="restored-from-archive">
  <h3>Restored from VimTips archive</h3>
  <p>
  This article was restored from the VimTips archive. There's probably
  missing images and broken links (and even some flash references), but it
  was still important to me to bring them back.
  </p>
</div>
