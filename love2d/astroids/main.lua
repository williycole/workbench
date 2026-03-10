-- TODO: do my natural way then play around with fat structs for optimization
function love.load()
  GREET = "Hello Astroids!!"
  Astroids = SpawnAstroids(10)

  Ship_x = 400
  Ship_y = 300
  Ship_y_velocity = 0
  Ship_x_velocity = 0
  Ship_angle = -math.pi / 2
end

function love.draw()
  DrawAstroids()
  DrawShip("fill", Ship_x, Ship_y, 20, 10, Ship_angle)
end

function love.update(dt)
  UpdateAstroidsPosition(dt)
  -- TODO: fix here - pick back up
  UpdateShipPosition(dt)
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
      x_velocity = math.cos(angle) * speed,
      y_velocity = math.sin(angle) * speed
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

function UpdateShipPosition(dt)
  if love.keyboard.isDown("right") then
    Ship_x = Ship_x + 100 * dt * Ship_angle
  end

  if love.keyboard.isDown("left") then
    Ship_x = Ship_x - 100 * dt
  end

  if love.keyboard.isDown("up") then
    Ship_y = Ship_y - 100 * dt
  end

  if love.keyboard.isDown("down") then
    Ship_y = Ship_y + 100 * dt
  end
end
