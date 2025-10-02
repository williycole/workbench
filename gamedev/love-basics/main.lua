-- y
-- ↑
-- |
-- |
-- |
-- +------------→ x

function love.draw()
	local x = "cat"
	x = "dog"
	x = "pikachu"
	local c = x

	love.graphics.print(c, 400, 300)
	love.graphics.rectangle("fill", 400, 0, 100, 100)
end
