<!-- :metadata:

title: Smart character movement (f,F,t,T) with hop.nvim
tags: NeoVim, Programming
publishedAt: 2024-12-19T18:09:22-07:00
summary:

I've enhanced Neovim's character navigation motions (`f`, `F`, `t`, `T`) by
integrating them with [hop.nvim](https://github.com/smoka7/hop.nvim). The
default behavior remains unchanged when there's only one occurrence of the
target character, or when using a count prefix (like `3f`). However, when
multiple matches exist, hop.nvim automatically labels each occurrence, making
it much easier to jump directly to your desired location.

-->

## Goal

Enhance Neovim's character navigation motions (`f`, `F`, `t`, `T`) with smart,
context-aware behavior. This enhancement integrates
[hop.nvim](https://github.com/smoka7/hop.nvim) to provide visual labels when
multiple target characters exist in the desired direction. The implementation
preserves Neovim's default behavior for single-character matches and
count-prefixed commands (like `3f`), while adding efficient multi-character
navigation through hop.nvim's labeling system.

```lua
-- Smart hop on `f`, `F`, `t`, and `T`
--
-- If there's only one of the target char in the direction specified, just go
-- there (default behavior). Otherwise, use hop to label the duplicates with
-- target labels
local function smart_hop(opts)
  -- If there's a count, use default vim behavior
  if vim.v.count > 0 then
    local char = vim.fn.getchar()
    char = type(char) == "number" and vim.fn.nr2char(char) or char
    vim.cmd("normal! " .. vim.v.count .. opts.motion .. char)
    return
  end

  -- Store if we're in operator-pending mode
  local is_operator = vim.fn.mode(1):match("[vo]")

  local hop = require("hop")
  local hint = require("hop.hint")
  local default_opts = setmetatable({}, { __index = require("hop.defaults") })
  local jump_regex = require("hop.jump_regex")

  local function check_opts(o)
    if not o then
      return
    end

    if vim.version.cmp({ 0, 10, 0 }, vim.version()) < 0 then
      o.hint_type = hint.HintType.OVERLAY
    end
  end

  local function override_opts(o)
    check_opts(o)
    return setmetatable(o or {}, { __index = default_opts })
  end
  local char = vim.fn.getchar()

  char = type(char) == "number" and vim.fn.nr2char(char) or char

  -- Get current line and cursor position
  local line = vim.api.nvim_get_current_line()
  local col = vim.api.nvim_win_get_cursor(0)[2]

  -- Count occurrences based on direction
  local count
  if opts.direction == hint.HintDirection.AFTER_CURSOR then
    local after_cursor = line:sub(col + 2)
    count = select(2, after_cursor:gsub(vim.pesc(char), ""))
  else
    local before_cursor = line:sub(1, col + 1)
    count = select(2, before_cursor:gsub(vim.pesc(char), ""))
  end

  if count <= 1 then
    -- Use native motion for 0 or 1 occurrence
    if is_operator then
      -- In operator-pending mode, execute the operator with the motion
      local op = vim.v.operator
      local reg = vim.v.register ~= '"' and '"' .. vim.v.register or ""
      vim.cmd("normal! " .. reg .. op .. opts.motion .. char)
    else
      -- In normal mode, just use the motion directly
      vim.cmd("normal! " .. opts.motion .. char)
    end
  else
    opts = override_opts({
      direction = opts.direction,
      current_line_only = true,
      hint_offset = opts.hint_offset,
    })
    if is_operator then
      -- For operators, we need to set the register and operator type
      local reg = vim.v.register
      local op = vim.v.operator
      hop.hint_with_regex_opts = {
        callback = function(pos)
          -- Apply the operator from current position to target
          local target_line = pos.line + 1
          local target_col = pos.column + 1
          local cmd = string.format(
            "normal! %s%s%dl%dh",
            reg ~= "" and '"' .. reg or "",
            op,
            target_line - vim.fn.line("."),
            target_col - vim.fn.col(".")
          )
          vim.cmd(cmd)
        end,
      }
    end
    -- Use hop for multiple occurrences
    hop.hint_with_regex(
      jump_regex.regex_by_case_searching(char, true, opts), opts)
  end
end

return {
  "smoka7/hop.nvim",
  version = "*",
  config = true,
  keys = {
    {
      "f",
      function()
        smart_hop({
          direction = require("hop.hint").HintDirection.AFTER_CURSOR,
          motion = "f",
          hint_offset = 0,
        })
      end,
      desc = "Smart hop char after cursor",
      mode = { "n", "v", "o" },
    },
    {
      "F",
      function()
        smart_hop({
          direction = require("hop.hint").HintDirection.BEFORE_CURSOR,
          motion = "F",
          hint_offset = 0,
        })
      end,
      desc = "Smart hop char before cursor",
      mode = { "n", "v", "o" },
    },
    {
      "t",
      function()
        smart_hop({
          direction = require("hop.hint").HintDirection.AFTER_CURSOR,
          motion = "t",
          hint_offset = -1,
        })
      end,
      desc = "Smart hop before char after cursor",
      mode = { "n", "v", "o" },
    },
    {
      "T",
      function()
        smart_hop({
          direction = require("hop.hint").HintDirection.BEFORE_CURSOR,
          motion = "T",
          hint_offset = 1,
        })
      end,
      desc = "Smart hop before char before cursor",
      mode = { "n", "v", "o" },
    },
  },
}
```
