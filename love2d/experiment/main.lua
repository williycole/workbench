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
  GREET = "🔥 Hello Charizard!!"
  GREET = "💧 Hello Squirtle!!"
  C = GREET
  print("SAME: math.random " .. math.random(0, 20))
  print("DIFF: love.math.random " .. love.math.random(0, 20))

  X_POS = 0
  Y_POS = 0
  T = buildXRandomCoords(2)

  X = 100
end

function love.draw()
  love.graphics.print(GREET, 0, 0)

  for _, pos in ipairs(T) do
    love.graphics.circle("fill", pos.x + X_POS, pos.y + Y_POS + 1, 10)
  end

  love.graphics.rectangle("line", X, 50, 200, 150)
end

-- on update x and y are new numbers this means each rectangle is drawn in a new position on iteration
function love.update(dt)
  X = X + 100 * dt
  X_POS = love.math.random(0, 100)
  Y_POS = love.math.random(0, 100)
  -- print("x:" .. X_POS .. " y:" .. Y_POS)
  -- print("dt: " .. dt)
end
