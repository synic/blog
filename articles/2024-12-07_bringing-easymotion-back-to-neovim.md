<!-- :metadata:

title: Bringing EasyMotion back to NeoVim!
tags: VIM, NeoVim
publishedAt: 2024-12-07T08:56:23-07:00
summary:

I'm pretty excited about this one. I've been using
[EasyMotion](https://github.com/easymotion/vim-easymotion) for a long time -
before [NeoVim](https://github.com/neovim/neovim) existed, and after, even
though it seemed to go out of style (see https://github.com/neovim/neovim for
more information about why). There are lots of alternatives, like
[leap.nvim](https://github.com/ggandor/leap.nvim),
[flash.nvim](https://github.com/folke/flash.nvim),
[sneak.nvim](https://github.com/justinmk/vim-sneak) and etc, but they all
required typing more characters, or didn't work across windows, or were just
generally buggy. I stuck with EasyMotion; it _just worked_ for me, even despite
it's drawbacks. Until it didn't. Recently, it has started crashing NeoVim. I
[reported](https://github.com/easymotion/vim-easymotion/issues/507) the issue,
but no response. EasyMotion is quite old and it doesn't appear as though anyone
is working on it anymore.

I found a solution, though, and that is to write a custom
[hop.nvim](https://github.com/smoka7/hop.nvim) command that _works_ like
EasyMotion used to. Read the rest of the article if you want to see how to do
it!

-->

This is my [lazy.nvim](https://github.com/folke/lazy.nvim) hop configuration,
with the custom command:

```lua
{
  "smoka7/hop.nvim",
  version = "*",
  opts = { keys = "etovxpdygfblzhckisuran;,", quit_key = "q" },
  config = function(_, opts)
    local hop = require("hop")
    hop.setup(opts)

    -- Create custom command for hopping to character/word matches
    vim.api.nvim_create_user_command("HopEasyMotion", function()
      local char = vim.fn.getchar()
      -- Convert numeric char code to string
      char = type(char) == "number" and vim.fn.nr2char(char) or char

      -- Create pattern based on input character type
      local pattern
      if char:match("%a") then
        -- For letters: match words starting with that letter (case insensitive)
        pattern = "\\c\\<" .. char
      elseif char:match("[%(%)]") then
        -- For parentheses: match them literally
        pattern = char
      elseif char == "." then
        -- For period: match literal period
        pattern = "\\."
      else
        -- For other non-letters: escape special characters
        pattern = char:gsub("([%^%$%(%)%%%.%[%]%*%+%-%?])", "%%%1")
      end

      ---@diagnostic disable-next-line: missing-fields
      hop.hint_patterns({
        current_line_only = false,
        multi_windows = true,
        hint_position = require("hop.hint").HintPosition.BEGIN,
      }, pattern)
    end, { desc = "Hop to words starting with input character" })
  end,
  keys = {
    { "<leader><leader>", "<cmd>HopEasyMotion<cr>", desc = "Hop to word" },
    { ";b", "<cmd>HopWordBC<cr>", desc = "Hop to word before cursor" },
    { ";w", "<cmd>HopWord<cr>", desc = "Hop to word in current buffer" },
    { ";a", "<cmd>HopWordAC<cr>", desc = "Hop to word after cursor" },
    { ";d", "<cmd>HopLineMW<cr>", desc = "Hop to line" },
    { ";f", "<cmd>HopNodes<cr>", desc = "Hop to node" },
    { ";s", "<cmd>HopPatternMW<cr>", desc = "Hop to pattern" },
    { ";j", "<cmd>HopVertical<cr>", desc = "Hop to location vertically" },
  },
}
```

With this, the command `HopEasyMotion` does what EasyMotion does (at least, the
part I used most, anyway). I have it bound to the `<leader><leader>`, (I have
`<leader>` bound to the space key).

So, I can type `<space><space>` followed by either a character that is the
start of a word in any visible window, OR a single non-letter character (like
parentheses, period, braces, etc) and hop will label each match in all windows.
I then type the label, and hop will take the cursor right to it.

This means that, without using the mouse, I can nearly instantly jump to any
location that I can see, even across windows. Of course, hop has all kinds of
other ways to move around, you can see I've bound most of them to chords that
start with `;`. Hop is pretty cool, [check it
out](https://github.com/smoka7/hop.nvim)!
