-- dynamical build a grid
-- global functions in lua are accessed/exported like this
local M = {}
function M.MkGrid(numRow, numCol)
	local grid = {}
	local row = {}

	local num = numCol -- Numbers can be integer or floating point.
	repeat
		table.insert(row, "[]")
		num = num - 1
	until num == 0

	for i = 1, numRow do
		table.insert(grid, table.concat(row))
		i = i + 1
	end
	return grid
end

function M.PrettyPrintGrid(grid)
	-- NOTE: ipairs is a Lua iterator function for array-like tables.
	for _, row in ipairs(grid) do
		print(row)
	end
end

return M
