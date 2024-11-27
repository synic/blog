<!-- :metadata:

title: My Gnome Configuration Settings
tags: Linux
publishedAt: 2009-08-02T23:26:52-07:00
summary:

I've installed Ubuntu on my netbook twice today (don't ask).  Probably like the
rest of you, part of installing Ubuntu involves going through and changing all
of the moronic settings that Gnome comes with as default.  This includes
changing the default terminal background from white to black, changing font
sizes, setting toolbars to "icon only", setting up hotkeys, etc.

-->

<p>I've installed Ubuntu on my netbook twice today (don't ask).  Probably like
the rest of you, part of installing Ubuntu involves going through and changing
all of the moronic settings that Gnome comes with as default.  This includes
changing the default terminal background from white to black, changing font
sizes, setting toolbars to "icon only", setting up hotkeys, etc.  </p>

<p>It was getting a little annoying to do all the time, so I finally went
through gconf-editor and found the settings I change, and wrote up a script to
do all of them at once.  Figured I'd post it here, not only so I can use it
    later, but so you can derive from it if you wish.  It's pretty simple, just
    using the command line utility "gconftool".  Here it is:</p>

```bash
$ gconftool -s --type bool /apps/update-notifier/auto_launch false
$ gconftool -s --type string /apps/metacity/general/action_double_click_titlebar toggle_shade
$ gconftool -s --type string /desktop/gnome/interface/toolbar_style icons
$ gconftool -s --type string /apps/nautilus/preferences/desktop_font "Bitstream Vera Sans 9"
$ gconftool -s --type string /apps/metacity/general/titlebar_font "Bitstream Vera Sans Bold 9"
$ gconftool -s --type string /desktop/gnome/interface/monospace_font_name "Bitstream Vera Sans Mono 9"
$ gconftool -s --type string /desktop/gnome/interface/document_font_name "Bitstream Vera Sans 9"
$ gconftool -s --type string /desktop/gnome/interface/font_name "Bitstream Vera Sans 9"
$ gconftool -s --type bool /apps/gnome-terminal/profiles/Default/allow_bold false
$ gconftool -s --type bool /apps/gnome-terminal/profiles/Default/use_theme_colors false
$ gconftool -s --type bool /apps/gnome-terminal/global/confirm_window_close false
$ gconftool -s --type string /apps/gnome-terminal/profiles/Default/foreground_color '#FFFFFFFFFFFF'
$ gconftool -s --type string /apps/gnome-terminal/profiles/Default/background_color '#000000000000'
$ gconftool -s --type string /apps/gnome-terminal/profiles/Default/palette "#000000000000:#AAAA00000000:#0000AAAA0000:#AAAA55550000:#00000000AAAA:#AAAA0000AAAA:#0000AAAAAAAA:#AAAAAAAAAAAA:#555555555555:#FFFF55555555:#5555FFFF5555:#FFFFFFFF5555:#55555555FFFF:#FFFF5555FFFF:#5555FFFFFFFF:#FFFFFFFFFFFF"
$ gconftool -s --type string /apps/gnome-terminal/profiles/Default/scrollbar_position hidden
$ gconftool -s --type bool /apps/gnome-terminal/profiles/Default/default_show_menubar false
$ gconftool -s --type int /apps/metacity/general/num_workspaces 4
$ gconftool -s --type string /apps/metacity/global_keybindings/switch_to_workspace_1 "<Control>F1"
$ gconftool -s --type string /apps/metacity/global_keybindings/switch_to_workspace_2 "<Control>F2"
$ gconftool -s --type string /apps/metacity/global_keybindings/switch_to_workspace_3 "<Control>F3"
$ gconftool -s --type string /apps/metacity/global_keybindings/switch_to_workspace_4 "<Control>F4"
$ gconftool -s --type string /apps/nautilus/icon_view/default_zoom_level small
$ gconftool -s --type bool /apps/nautilus/desktop/computer_icon_visible true
$ gconftool -s --type bool /apps/nautilus/desktop/home_icon_visible true
$ gconftool -s --type bool /apps/nautilus/desktop/network_icon_visible true
$ gconftool -s --type bool /apps/nautilus/desktop/trash_icon_visible true
```

<p>I'll probably keep editing this article and adding new things as I find
them.</p>

<div class="restored-from-archive">
  <h3>Restored from VimTips archive</h3>
  <p>
  This article was restored from the VimTips archive. There's probably
  missing images and broken links (and even some flash references), but it
  was still important to me to bring them back.
  </p>
</div>
