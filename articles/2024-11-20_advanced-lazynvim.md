---
title: Advanced lazy.nvim
slug: advanced-lazynvim
publishedAt: 2024-11-20T18:59:00-07:00
tags: [NeoVim, Programming]
openGraph:
  description: Does lazy.nvim configuration confuse you? This article is for you!
  image: /static/img/screenshots/lazynvim.webp
summary: |
  <img src="/static/img/screenshots/lazynvim.webp" width="592" height="299" alt="lazy.nvim screenshot" />

  There have been many plugin managers in the Vim ecosystem over the years. I've
  used quite a number of them, but
  [lazy.nvim](https://github.com/folke/lazy.nvim) is the one I've enjoyed the
  most, by far. It's easy, it's concise, there's a ton of "hidden" little tricks,
  and it allows me to easily separate my plugin configuration into separate files
  (which was sort of a pain in previous plugin managers).

  Most of these this information can be found in the lazy.nvim README, but
  without just getting your hands dirty, it can be difficult to understand how it
  all comes together. That's where this article comes in.
---
<img src="/static/img/screenshots/lazynvim.webp" width="592" height="299" alt="lazy.nvim screenshot" />

There have been many plugin managers in the Vim ecosystem over the years. I've
used quite a number of them, but
[lazy.nvim](https://github.com/folke/lazy.nvim) is the one I've enjoyed the
most, by far. It's easy, it's concise, there's a ton of "hidden" little tricks,
and it allows me to easily separate my plugin configuration into separate files
(which was sort of a pain in previous plugin managers).

Most of these this information can be found in the lazy.nvim README, but
without just getting your hands dirty, it can be difficult to understand how it
all comes together. That's where this article comes in.

Let's jump right in with an explanation of how modern NeoVim plugins are set
up, and how that works with lazy.nvim.

# The Ideal NeoVim Plugin

For this example, we will invent a fake NeoVim plugin called "noun". Noun will
have a structure like this:

```
~/Projects/noun.nvim ❯ tree
.
├── LICENSE
├── README.md
└── lua
    └── noun
        ├── init.lua
        └── main.lua

3 directories, 4 files
```

What's important here is that there is a module that matches the repository
name ("noun" in this case, more on that later) that exports a setup function
with the following signature: setup(opts), where opts is a table that contains
the configuration options for this plugin. This function is meant to called
when the plugin is loaded (with lazy.nvim, this may or may not be when vim
starts up; more on that later).


Let's say that the repository for noun is https://github.com/jbz/noun.nvim.
lazy.nvim will try to automatically discover where the setup function is by
removing everything except the last path element and then stripping `.vim`
or `.nvim` from the end of the path.


If the plugin you are trying to use follows this format, the following 4 plugin
specs would be equivalent:

```lua
-- this calls `require("noun").setup({})`
{ "jbz/noun.nvim", config = true }

-- this also calls `require("noun").setup({})
{ "jbz/noun.nvim", opts = {} }

-- samesies
{ "jbz/noun.nvim", opts = function() return {} end }

-- doing it manually
{ "jbz/noun.nvim", config = function() require("noun.nvim").setup({}) end }
```

If you don't pass opts or config, or if you pass config = false, the setup
function will not be called automatically. 

What if the plugin's module name is not the same as the repository name? If
noun's module was named something like noun_nvim, you can tell lazy.nvim what
the name of the module is by setting the name option in the plugin spec:

```lua
{ "jbz/noun.nvim", name="noun_nvim", config = true }
```

Of course, if it's an older vim plugin, or something that doesn't follow this
format, you can override the config function and do whatever you want:

```lua
{
  "rcarriga/nvim-notify",
   event = "VeryLazy",
   opts = {
     render = "minimal",
     stages = "fade",
   },
   config = function(_, opts)
     local notify = require("notify")
     notify.setup(opts)
     vim.notify = notify
   end,
}
```

This plugin contains two reasons to pass a custom config function: First,
because the main module is "notify" and not "nvim-notify", and we needed to set
`vim.notify = notify`. If the module name matches but you still need to do
something other than call setup, you need a custom function.

# "Classic" and Other Plugins

Some plugins (older usually, but not always) do not use the setup pattern, and
instead want you to set configuration options using global variables (such as
let g:EasyMotion_smartcase = true). In order to set these, you can use the
config function, and the settings will be set when the plugin loads. If you
want these settings to be set before the plugin loads, you can use the init
function, like this:

```lua
{
  "easymotion/vim-easymotion",
  init = function()
    vim.g.EasyMotion_smartcase = true
    vim.g.EasyMotion_use_upper = true
  end,
}
```

The init function is called when lazy.nvim itself loads, before any other
plugins are loaded, even if the plugin itself is configured to lazy load.

With the config and init functions, you can load and configure almost any
plugin that is supported by NeoVim.

# What about the Vim's Plugin directories?

If you've been using Vim for a while, you may be used to setting up plugins
with files in these directories. While you can still do that, it will be harder
to use the lazy loading features of lazy.nvim.

# Lazy Loading

One of the cool features of lazy.nvim is that you can delay loading a specific
plugin until it is needed. This can improve startup time and memory usage if
you use a lot of plugins.

If you specify one of the following keys in your plugin spec, your plugin will
be lazy loaded depending on the events you chose:

* `ft` &nbsp; - &nbsp; "filetype": can be a single string filetype (like
`ft="html"`) or it can be a table of multiple filetypes, (like `ft={"html",
"css"}`). This will cause the plugin to load any time a buffer with the given
filetype is encountered.<br>
* `event` &nbsp; - &nbsp; this will cause the plugin to load when the autocmd event is
encountered for the first time, such as `BufReadPre`, `LspAttach`, `OptionSet`.
User events are supported too, in those cases, you pass the pattern of the
event instead of User. One such user autocmd is defined by lazy.nvim itself:
VeryLazy. This event is fired once lazy.nvim has completed loading itself and
all of the non-lazy loaded plugins in your configuration. Like ft, this can be
a single event or a table of multiple events. <br>
* `cmd` &nbsp; - &nbsp; "command": this will cause the plugin to load if the
command specified is executed. If you set `cmd = 'Telescope'` it will load any
time you run a telescope command. Can also be a single command or a table of
multiple commands.

You can also specify `lazy = false`, and then load the plugin manually using
`:Lazy load [name]`.

## Implications of Lazy Loading

At this point, you may be tempted to lazy load everything. Various sources
(including the LazyVim distribution) will load plugins (such as nvim-lspconfig)
using the `event = {"BufReadPost", "BufWritePost", "BufNewFile"}` (which will
lazy load the plugin after any file is read), but I've found that this can
cause the plugin to load after other events, like `FileType`. The result is
that if you open files from the command line, like `nvim file.go`, LSP won't
work until you open a second go file to trigger the `FileType` event.
[Folke](https://github.com/folke) (lazy.nvim's author) seems to be trying to
solve this with the `LazyFile` event in the LazyVim distribution, but hasn't
exposed it to lazy.nvim itself:
https://github.com/LazyVim/LazyVim/discussions/1583)

# Dependencies

Lazy.nvim allows you to specify dependencies for plugins, which, of course,
allows you to say what plugins need to be loaded for a given plugin to work,
such as telescope depending on plenary:

```lua
{
  "nvim-telescope/telescope.nvim",
   dependencies = { "nvim-lua/plenary.nvim" }
}
```

However, they can be used for a second purpose. In the lazy.nvim documentation
it says the following: "Dependencies are always lazy-loaded unless specified
otherwise"

This means that any dependencies you add to a plugin, whether or not that
plugin actually depends on them, will be lazy loaded with the dependent plugin.

```lua
{
  "nvim-telescope/telescope.nvim",
   cmd = "Telescope",
   dependencies = {
     "nvim-lua/plenary.nvim",
     "nvim-telescope/telescope-ui-select.nvim",
   }
}
```

In this case, `telescope` doesn't actually require `telescope-ui-select` to be
loaded, but it does mean that `telescope-ui-select` will be lazy loaded along
with telescope itself when the `:Telescope` command is used. Pretty cool!

# Multiple Plugins, Same Event

You can load multiple plugins for using the same event. I use it to load all
themes before showing the Telescope theme picker. Note the use of
`vim.deepcopy`, you can't use the exact same table, so you can either define
one and deepcopy it for every plugins, or just manually list the keys for each
one pointing to the same target function.

```lua
-- themes.lua

local function colorscheme_picker()
	local target = vim.fn.getcompletion

	-- only show themes that were installed via lazy (and habamax)
	vim.fn.getcompletion = function()
		return vim.tbl_filter(function(color)
			return color == "habamax" or not vim.tbl_contains(config.options.default_colorschemes, color)
			---@diagnostic disable-next-line: redundant-parameter
		end, target("", "color"))
	end

	vim.cmd.Telescope("colorscheme")
	vim.fn.getcompletion = target
end

local keys = {
	{ "<leader>st", colorscheme_picker, desc = "List themes" },
}

return {
	{ "sainnhe/gruvbox-material", keys = vim.deepcopy(keys) },
	{ "catppuccin/nvim", name = "catppuccin", keys = vim.deepcopy(keys) },
	{ "neanias/everforest-nvim", name = "everforest", keys = vim.deepcopy(keys) },
	{ "rebelot/kanagawa.nvim", keys = vim.deepcopy(keys) },
	{ "Mofiqul/dracula.nvim", keys = vim.deepcopy(keys) },
	{ "EdenEast/nightfox.nvim", keys = vim.deepcopy(keys) },
	{ "oxfist/night-owl.nvim", keys = vim.deepcopy(keys) },
	{ "AlexvZyl/nordic.nvim", keys = vim.deepcopy(keys),
	{ "ribru17/bamboo.nvim", keys = vim.deepcopy(keys) },
}
```

Because all the themes use the same event, they will load and show up in the
theme picker. Otherwise, they wouldn't appear until they were loaded manually
(alternatively, you could set `lazy = false` on all of them).

Happy Vimming!
