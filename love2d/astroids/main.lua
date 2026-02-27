function love.load()
  GREET = "Hello Astroids!!"
  Astroids = SpawnAstroids(10)
  Ship_x, Ship_y, Ship_angle = 200, 100, math.pi / 4
end

function love.draw()
  DrawAstroids()
  -- drawShip
  DrawShip("fill", Ship_x, Ship_y, 20, 10, Ship_angle)
end

function love.update(delta)
  UpdateAstroidsPosition(delta)
end

function SpawnAstroids(num_astroids)
  Astroids = {}
  for i = 1, num_astroids do
    Astroids[i] = {
      x = math.random(0, love.graphics.getWidth()),
      y = math.random(0, love.graphics.getHeight()),
      radius = math.random(10, 50),
      speed = math.random(1, 50)
    }

    print("Astroid spawned at X: " .. Astroids[i].x
      .. " Y: " .. Astroids[i].y
      .. " Radius: " .. Astroids[i].radius
      .. " Speed: " .. Astroids[i].speed)
  end

  return Astroids
end

function DrawAstroids()
  love.graphics.print(GREET, 0, 0)
  for i = 1, #Astroids do
    love.graphics.circle("line", Astroids[i].x, Astroids[i].y, Astroids[i].radius)
  end
end

function UpdateAstroidsPosition(delta)
  for i = 1, #Astroids do
    Astroids[i].x = Astroids[i].x + Astroids[i].speed * delta
    Astroids[i].y = Astroids[i].y + Astroids[i].speed * delta
  end
end

function DrawShip(mode, x, y, length, width, angle) -- position, length, width and angle
  love.graphics.push()
  love.graphics.translate(x, y)
  love.graphics.rotate(angle)
  love.graphics.polygon(mode, -length / 2, -width / 2, -length / 2, width / 2, length / 2, 0)
  love.graphics.pop()
end
