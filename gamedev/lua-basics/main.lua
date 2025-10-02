print("Hello, world!")
print("")
A = "goofy gopher" -- Globals are ok.
print("a is " .. A)
A = "goofy lua"
if A ~= "goofy gopher" then
	print("re assign a as.. " .. A .. " ..new point in mem, the gc will clean up the old one")
end
print("")

local countMe = 0
for i = 1, 5 do -- The range includes both ends.
	countMe = i
	print(countMe .. " normal count")
	countMe = 0
end
print("")

for i = 1, 5 do
	countMe = countMe + i
	print(countMe .. " plus 1")
	countMe = 0
end
print("")

print("no more counting, countMe is " .. countMe)
print("")
local fredSum = 0
-- In general, the range is begin, end[, step].
for j = 5, 1, -1 do
	print("before " .. fredSum)
	fredSum = fredSum + j
	print("after  " .. fredSum)
end
fredSum = 0
print("")

-- Another loop construct: prob good for game loops maybe.
local num = 5 -- Numbers can be integer or floating point.
repeat
	print("the way of the future num -> " .. num)
	num = num - 1
until num == 0
print("")

local aBoolValue = false

-- Only nil and false are falsy; 0 and '' are true!
if not aBoolValue then
	print("it was false")
end

-- 'or' and 'and' are short-circuited.
-- This is similar to the a?b:c operator in C/js:
local ans = aBoolValue and "yes" or "no" --> 'no'
print("ans " .. ans)
print("")

local function fib(n)
	if n < 2 then
		return 1
	end
	return fib(n - 2) + fib(n - 1)
end
print("fib(10) = " .. fib(10))
print("")
-- Closures and anonymous functions are ok:
local function adder(x)
	-- The returned function is created when adder is
	-- called, and remembers the value of x:
	return function(y)
		return x + y
	end
end
local a1 = adder(9)
local a2 = adder(36)
print(a1(16)) --> 25
print(a2(64)) --> 100

-- Returns, func calls, and assignments all work
-- with lists that may be mismatched in length.
-- Unmatched receivers are nil;
-- unmatched senders are discarded.

-- NOTE:
--local x, y, z = 1, 2, 3, 4
-- NOTE: Now x = 1, y = 2, z = 3, and 4 is thrown away.
local function bar(a, b, c)
	print(a, b, c)
	return 4, 8, 15, 16, 23, 42
end
bar("zaphod") --> prints "zaphod  nil nil"
print("")

print(
	"thats the super basics..\n"
		.. "see https://learnxinyminutes.com/lua/ ..for more\n"
		.. "everything after this is play with objects, tables, meta tables, etc\n"
)

print("requiring tables...\n")
local mod = require("./tables")
local grid = mod.MkGrid(3, 4)
mod.PrettyPrintGrid(grid)
mod.PrettyPrintGrid(mod.MkGrid(2, 2))
