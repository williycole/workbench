-- TODO: do my natural way then play around with fat structs for optimization
function love.load()
  GREET = "Hello Astroids!!"
  Astroids = SpawnAstroids(10)
  Ship_x, Ship_y, Ship_angle = 400, 300, -math.pi / 2
end

function love.draw()
  DrawAstroids()
  DrawShip("fill", Ship_x, Ship_y, 20, 10, Ship_angle)
end

function love.update(delta)
  UpdateAstroidsPosition(delta)
  -- TODO: fix here - pick back up
  MoveShip(Ship_x, Ship_y, Ship_angle, delta)
end

function SpawnAstroids(num_astroids)
  Astroids = {}
  for i = 1, num_astroids do
    local speed = math.random(1, 50)
    local angle = math.random() * 2 * math.pi -- 2pi ~ 6.28, i.e. full circle in radians

    Astroids[i] = {
      x = math.random(0, love.graphics.getWidth()),
      y = math.random(0, love.graphics.getHeight()),
      radius = math.random(10, 50),
      x_v = math.sin(angle) * speed,
      y_v = math.cos(angle) * speed
    }
    -- math.sin(angle) * speed
    print("Astroid spawned at"
      .. " X: " .. Astroids[i].x
      .. " Y: " .. Astroids[i].y
      .. " Radius: " .. Astroids[i].radius
      .. " Velocity X: " .. Astroids[i].x_v
      .. " Velocity Y: " .. Astroids[i].y_v)
  end

  return Astroids
end

function DrawAstroids()
  love.graphics.print(GREET, 0, 0)
  for i = 1, #Astroids do
    love.graphics.circle("line", Astroids[i].x, Astroids[i].y, Astroids[i].radius)
  end
end

function UpdateAstroidsPosition(dt)
  for i = 1, #Astroids do
    Astroids[i].x = Astroids[i].x + Astroids[i].x_v * dt
    Astroids[i].y = Astroids[i].y + Astroids[i].y_v * dt
  end
end

function DrawShip(mode, x, y, length, width, angle) -- position, length, width and angle
  love.graphics.push()
  love.graphics.translate(x, y)
  love.graphics.rotate(angle)
  love.graphics.polygon(mode, -length / 2, -width / 2, -length / 2, width / 2, length / 2, 0)
  love.graphics.pop()
end

function GetShipPosition()

end

function MoveShip(x, y, angle, dt)
  if love.keyboard.isDown("right") then
    x = x + 100 * dt
  end

  if love.keyboard.isDown("left") then
    x = x - 100 * dt
  end

  if love.keyboard.isDown("up") then
    y = y - 100 * dt
  end

  if love.keyboard.isDown("down") then
    y = y + 100 * dt
  end
end
