---
title: "NeoVim: Normalize and Copy Code Block"
slug: neovim-copy-without-leading-indent
tags: [NeoVim, Programming]
publishedAt: 2025-01-29T03:35:36-07:00
summary: |
  When copying code blocks to share on GitHub or other platforms, preserving
  proper indentation can be tricky. I created a NeoVim function that
  automatically normalizes indentation when copying, making your shared code
  snippets clean and consistent. Here's how to set it up!
---

Let's say I wanted to copy this block of code so that I can paste it into a
github comment:

![Copy block](https://github.com/user-attachments/assets/4ce11648-30fd-447d-b10c-74a5552a9740)
If I just press `y` to yank the block and paste it into a `lua` block on
github, it will look like this:

![Bad copy block](https://gist.github.com/user-attachments/assets/17770b00-46dd-4795-a165-a3eb32a871be)

If I use my function (see below), it looks like this when pasted:

![Good copy block](https://gist.github.com/user-attachments/assets/12faa947-af21-4f0f-9ff5-8bb18b88f451)

Much better!

Here's the function:

```lua
local function copy_normalized_block()
  local mode = vim.fn.mode()
  if mode ~= "v" and mode ~= "V" then
    return
  end

  vim.cmd([[silent normal! "xy]])
  local text = vim.fn.getreg("x")
  local lines = vim.split(text, "\n", { plain = true })

  local converted = {}
  for _, line in ipairs(lines) do
    local l = line:gsub("\t", "  ")
    table.insert(converted, l)
  end

  local min_indent = math.huge
  for _, line in ipairs(converted) do
    if line:match("[^%s]") then
      local indent = #(line:match("^%s*"))
      min_indent = math.min(min_indent, indent)
    end
  end
  min_indent = min_indent == math.huge and 0 or min_indent

  local result = {}
  for _, line in ipairs(converted) do
    if line:match("^%s*$") then
      table.insert(result, "")
    else
      local processed = line:sub(min_indent + 1)
      processed = processed:gsub("^%s+", function(spaces)
        return string.rep("  ", math.floor(#spaces / 2))
      end)
      table.insert(result, processed)
    end
  end

  local normalized = table.concat(result, "\n")
  vim.fn.setreg("+", normalized)
  vim.notify("Copied normalized text to clipboard")
end
```

And then I bind it to a key like so:

```lua
vim.keymap.set(
  "v",
  "<leader>y",
  copy_normalized_block,
  { desc = "Copy normalized" },
)
```
