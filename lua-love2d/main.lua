function love.load()
	love.window.setMode(640, 480)
	BunnySprite = love.graphics.newImage("assets/bunny.png")

	Player = {}
	Player.x = 400
	Player.y = 200
	Player.speed = 1
end

function love.update(dt)
	if love.mouse.isDown(1) then
	end
end

function love.draw()
	love.graphics.print("Hello World", 400, 300)

	love.graphics.draw(BunnySprite)

	love.graphics.setColor(255, 0, 0) -- color of rectangle
	love.graphics.rectangle("fill", Player.x, Player.y, 10, 10)
	love.graphics.setColor(255, 255, 255) -- reset color

	love.graphics.print("Current FPS: " .. tostring(love.timer.getFPS()), 10, 10)
end
