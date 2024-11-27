<!-- :metadata:

title: Getting Ubuntu to work on your Eee Box
tags: Linux, Gadgets
publishedAt: 2008-09-25T01:00:44-07:00
summary:

My Eee Box (model B202) arrived today.  I actually plan on using this
machine to run Windows (I have a very short list of things that I *require* the
blasted OS to run)...

-->

<p>My Eee Box (model B202) arrived today.  I actually plan on using this
machine to run Windows (I have a very short list of things that I *require* the
blasted OS to run).</p>

<p>First impressions:  Yeah, they really are 'cute'.  About the size of a Wii,
but thinner.  The <a href='http://splashtop.com'>splashtop</a> mini OS that's
built into the BIOS is pretty groovy, but I haven't really been able to think
of a good use for it.  The VGA mounting bracket for LCD monitors is just plain
awesome.  Out of sight, out of mind.  I plan on using Windows XP via remote
desktop from my main Ubuntu machine.</p>

<p>Who cares about all that, though.  I obviously had to try and see if I could
get Ubuntu running on it, which I did, and that's what this post is all
about.</p>

<p>To do this, you'll need a linux machine with internet access and a USB
key.</b>

<p>Create a bootable USB key with <a
href='http://mirrors.xmission.com/ubuntu-cd/hardy/ubuntu-8.04.1-desktop-i386.iso'>ubuntu-8.04.1-desktop_i386.iso</a>
and the <a href='http://www.startx.ro/sugar/isotostick.sh'>isotostick.sh</a>
script mentioned on the Ubuntu <a
href='https://help.ubuntu.com/community/Installation/FromUSBStick'>FromUSBStick</a>
Wiki page.  This assumes that you have an already pre formatted (ext2 or fat32)
USB key, and that the device is /dev/sde1 (to find out what device your key is
using, type 'dmesg' a few seconds after plugging it in).
</p>

```bash
$ wget http://mirrors.xmission.com/ubuntu-cd/hardy/ubuntu-8.04.1-desktop-i386.iso
$ wget http://www.startx.ro/sugar/isotostick.sh
$ chmod +x isotostick.sh
$ sudo ./isotostick.sh ./ubuntu-8.04.1-desktop-i386.iso /dev/sde1
Not verifying image...(no checkisomd5 in Ubuntu so skipping)!
Copying live image to USB stick
cp: cannot create symbolic link `/media/usbdev.Dn7319/dists/stable': Operation not permitted
cp: cannot create symbolic link `/media/usbdev.Dn7319/dists/unstable': Operation not permitted
Installing boot loader
USB stick set up as live image!
```

<p>
The 'cannot create symbolic link' errors are normal if your key is formatted
with fat32.  Next, you'll want to get the initial custom eeepc kernel packages
from <a href='http://array.org/ubuntu'>array.org</a> and copy them on to the
usb key.</p>

```bash
$ mkdir eeepc
$ cd eeepc
$ wget http://www.array.org/ubuntu/dists/hardy/eeepc/binary-i386/linux-image-2.6.24-21-eeepc_2.6.24-21.39eeepc1_i386.deb
$ wget http://www.array.org/ubuntu/dists/hardy/eeepc/binary-i386/linux-ubuntu-modules-2.6.24-21-eeepc_2.6.24-21.30eeepc5_i386.deb
$ sudo mkdir /mnt/disk
$ sudo mount /dev/sde1 /mnt/disk
$ sudo cp -rv ../eeepc /mnt/disk
$ sudo umount /mnt/disk
```

<p>Next, you'll want to configure your Eee Box to boot from your USB key.
Unlike the manual (and the interwebs) would like you to believe, pressing F8 at
the splashtop boot screen will not allow you to choose your USB key as the boot
device.  Plug the key in, and enter the BIOS setup screen.  Press the right
arrow key to "Boot", down arrow to "Hard Disk Drives" (if you don't see this
option, make sure you plug your USB key in before you power on the Eee Box),
select the "1rst Drive" option and choose your USB device.  Press F10 and
enter.  The Box should now boot the Ubuntu install ISO.</p>

<p>If you are using an LCD monitor connected via the VGA->DVI converter, once
you get past the "Ubuntu" boot logo screen with the progress bar, your monitor
might just go into powersaving mode.  If this happens to you, press CTRL+ALT+F2
to go to a virtual terminal.  Your monitor should come back on.  Type 'sudo
bash'. Edit /etc/X11/xorg.conf and change it so that it looks like the
xorg.conf file located <a
href='http://www.vimtips.org/media/files/xorgeeebox.conf.txt'>here</a>.  Once
you've saved it, type 'killall -9 X'. If X doesn't start back up automatically
in a few seconds, just type 'startx'.  If the installation doesn't start
automatically after X is running, hit ALT+F2, type 'ubiquity' in the box, and
press enter.  The familiar Ubuntu installation screen should pop up.</p>

<p>I'm not going to go over the installation process, I assume you already know
how to install Ubuntu.  The Eee Box 80GB hdd is separated nicely into two
partitions already.  One for Windows, and one for "data".  I chose to dual boot
the machine, so I just erased the "data" partition, created a 2GB swap and used
the rest as the root ext3 filesystem.  Obviously you can partition the system
however you want.</p>

<p>Once you get Ubuntu installed, and you've booted into the new installation,
you're going to want to install those two custom kernel .deb files you copied
onto the USB key.  Because you installed Ubuntu from the USB device, your
/etc/fstab file will have an entry telling the OS that your USB key is a cdrom
drive, and it won't mount correctly when you plug it in.  Edit /etc/fstab and
remove the last line.  Plug your USB key in if you haven't already, and if you
have, remove it and plug it back in.  It should mount automatically.  Open a
terminal, see where it mounted using the 'mount' command, and install the .deb
files:</p>

```bash
$ cd /media/disk/eeepc
$ sudo dpkg -i *.deb
```

<p>Reboot.  Everything should be working at this point.  WiFi, the wired
network device, sound, compiz, etc.  Woot!</p>

<div class="restored-from-archive">
  <h3>Restored from VimTips archive</h3>
  <p>
  This article was restored from the VimTips archive. There's probably
  missing images and broken links (and even some flash references), but it
  was still important to me to bring them back.
  </p>
</div>
