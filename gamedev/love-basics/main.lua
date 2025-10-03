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
		local position = { math.random(0, 20), math.random(0, 20) }

		table.insert(coords, position)
		i = i + 1
	end

	return coords
end

function love.load()
	GREET = "󰐝 Hello Charizard!!"
	GREET = "󰐝 Hello Pikachu!!"
	C = GREET

	X_POS = 100
	Y_POS = 20

	COORDS = buildXRandomCoords(1)
end

function love.draw()
	-- love.graphics.print(C, 400, 300)

	for i, pos in ipairs(COORDS) do
		print("coords:" .. i .. " x:" .. pos[1] .. "y:" .. pos[2])
		love.graphics.rectangle("fill", pos[1], pos[2] + 1, 30, 30)
	end
end

-- on update x and y are new numbers this means each rectangle is drawn in a new position on iteration
function love.update()
	X_POS = X_POS + love.math.random(0, 20)
	Y_POS = Y_POS + love.math.random(0, 20)
end
