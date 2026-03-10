-- TODO: do my natural way then play around with fat structs for optimization
function love.load()
  GREET = "Hello Astroids!!"

  Astroids = SpawnAstroids(10)

  Ship = {
    length = 20,
    width = 10,
    x = 400,
    y = 300,
    y_velocity = 0,
    x_velocity = 0,
    angle = -math.pi / 2,
  }
end

function love.draw()
  DrawAstroids()
  DrawShip()
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
      x_position = math.random(0, love.graphics.getWidth()),
      y_position = math.random(0, love.graphics.getHeight()),
      radius = math.random(10, 50),
      x_velocity = math.cos(angle) * speed,
      y_velocity = math.sin(angle) * speed
    }
    -- math.sin(angle) * speed
    print("Astroid spawned at"
      .. " X Postion: " .. Astroids[i].x_position
      .. " Y Position: " .. Astroids[i].y_position
      .. " Radius: " .. Astroids[i].radius
      .. " X Velocity: " .. Astroids[i].x_velocity
      .. " Y Velocity: " .. Astroids[i].y_velocity)
  end

  return Astroids
end

function DrawAstroids()
  love.graphics.print(GREET, 0, 0)
  for i = 1, #Astroids do
    love.graphics.circle("line", Astroids[i].x_position, Astroids[i].y_position, Astroids[i].radius)
  end
end

function UpdateAstroidsPosition(dt)
  for i = 1, #Astroids do
    Astroids[i].x_position = Astroids[i].x_position + Astroids[i].x_velocity * dt
    Astroids[i].y_position = Astroids[i].y_position + Astroids[i].y_velocity * dt
  end
end

function DrawShip() -- position, length, width and angle
  local mode = "fill"

  love.graphics.push()
  love.graphics.translate(Ship.x, Ship.y)
  love.graphics.rotate(Ship.angle)
  love.graphics.polygon(mode, -Ship.length / 2, -Ship.width / 2, -Ship.length / 2, Ship.width / 2, Ship.length / 2, 0)
  love.graphics.pop()
end

function GetShipPosition()

end

function UpdateShipPosition(dt)
  if love.keyboard.isDown("right") then
    Ship.x = Ship.x + 100 * dt * Ship.angle
  end

  if love.keyboard.isDown("left") then
    Ship.x = Ship.x - 100 * dt
  end

  if love.keyboard.isDown("up") then
    Ship.y = Ship.y - 100 * dt
  end

  if love.keyboard.isDown("down") then
    Ship.y = Ship.y - 100 * dt
  end
end
