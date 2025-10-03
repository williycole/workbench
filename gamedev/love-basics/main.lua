--zip : ❯    zip -9 -r love-basics.love .
--run : ❯    love love-basics.love
-- y
-- ↑
-- |
-- |
-- |
-- +------------→ x

-- build x number of random coordinates
local function buildXRandomCoords(numOfRects)
	local coords = {}

	for i = 1, numOfRects do
		-- TODO: look up love.graphics.getHeight()
		-- and the like to not go out of bounds, also
		-- look up why the window is small on launch
		--
		local width = love.math.random(0, love.graphics.getWidth())
		local height = love.math.random(0, love.graphics.getHeight())

		table.insert(coords, { x = width, y = height })
		print("inserted coords:" .. i .. " x:" .. coords[i].x .. "y:" .. coords[i].y)

		i = i + 1
	end

	return coords
end

function love.load()
	GREET = "󰐝 Hello Charizard!!"
	GREET = "󰐝 Hello Pikachu!!"
	C = GREET
	print("SAME: math.random " .. math.random(0, 20))
	print("DIFF: love.math.random " .. love.math.random(0, 20))

	X_POS = 0
	Y_POS = 0
	T = buildXRandomCoords(2)
end

function love.draw()
	love.graphics.print(C, 400, 300)

	for _, pos in ipairs(T) do
		love.graphics.rectangle("fill", pos.x + X_POS, pos.y + Y_POS + 1, 30, 30)
	end
end

-- on update x and y are new numbers this means each rectangle is drawn in a new position on iteration
function love.update()
	X_POS = love.math.random(0, 100)
	Y_POS = love.math.random(0, 100)
	-- print("x:" .. X_POS .. " y:" .. Y_POS)
end
