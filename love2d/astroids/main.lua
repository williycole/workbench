function love.load()
  GREET = "Hello Astroids!!"
  Astroid_Radius = love.math.random(10, 100)
  Astroid_X = love.math.random(10, 100)
  Astroid_Y = love.math.random(10, 100)
  Astroid_Base_Count = 5
end

function love.draw()
  love.graphics.print(GREET, 0, 0)
  SpawnAstroids(Astroid_Base_Count)
end

function love.update(dt)
  local Astroid_Speed = 25
  Astroid_X = Astroid_X + Astroid_Speed * dt
  Astroid_Y = Astroid_Y + Astroid_Speed * dt
end

function SpawnAstroids(num_astroids)
  local i = 0
  while i <= num_astroids do
    i = i + 1
    love.graphics.circle("line", Astroid_X, Astroid_Y, Astroid_Radius)
    print("astroid " .. i .. " spawned at " .. Astroid_X .. ", " .. Astroid_Y .. "with radius " .. Astroid_Radius)
  end
end
