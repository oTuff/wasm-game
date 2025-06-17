SCREEN_WIDTH = 640
SCREEN_HEIGHT = 480
BUNNY_SCALE = 0.2
UPPER_BOUND = 40

---@class Bunny
---@field posX number
---@field posY number
---@field scaleX number
---@field scaleY number
---@field speedX number
---@field speedY number

function love.load()
	love.window.setMode(SCREEN_WIDTH, SCREEN_HEIGHT)
	BunnySprite = love.graphics.newImage("assets/bunny.png")
	---@type Bunny[]
	BunnyList = {}
	Gravity = 0.75
end

---@param amount integer
local function newBunny(amount)
	for _ = 1, amount do
		---@type Bunny
		local bunny = {
			posX = math.random() * 5,
			posY = math.random() * 5,
			scaleX = BUNNY_SCALE,
			scaleY = BUNNY_SCALE,
			speedX = (math.random() * 2) + 2,
			speedY = (math.random() * 2) + 2,
		}
		table.insert(BunnyList, bunny)
	end
end

-- Trigger once
function love.mousepressed(_, _, button)
	if button == 2 then
		local toAdd = 100 - (#BunnyList % 100)
		if toAdd == 0 then
			toAdd = 100
		end
		newBunny(toAdd)
	end

	if button == 3 then
		local toAdd = 1000 - (#BunnyList % 1000)
		if toAdd == 0 then
			toAdd = 1000
		end
		newBunny(toAdd)
	end
end

local accumulator = 0
local tickRate = 1 / 60

-- Debug tickrate info
local fixedUpdates = 0
local fixedUpdatesLastSecond = 0
local timeCounter = 0

function love.update(dt)
	accumulator = accumulator + dt
	-- debug
	timeCounter = timeCounter + dt

	while accumulator >= tickRate do
		-- Rapid trigger
		if love.mouse.isDown(1) then
			newBunny(10)
		end

		for _, bunny in pairs(BunnyList) do
			bunny.posX = bunny.posX + bunny.speedX
			bunny.posY = bunny.posY + bunny.speedY
			bunny.speedY = bunny.speedY + Gravity

			local scaledWidth = BunnySprite:getWidth() * bunny.scaleX
			local scaledHeight = BunnySprite:getHeight() * bunny.scaleY

			if bunny.posX < 0 or bunny.posX > SCREEN_WIDTH - scaledWidth then
				bunny.speedX = -bunny.speedX
			end

			if bunny.posY > SCREEN_HEIGHT - scaledHeight then
				bunny.posY = SCREEN_HEIGHT - scaledHeight
				bunny.speedY = -bunny.speedY
			elseif bunny.posY < UPPER_BOUND and bunny.speedY < 0 then
				bunny.speedY = bunny.speedY * 0.7
			end
		end

		-- debug
		fixedUpdates = fixedUpdates + 1

		accumulator = accumulator - tickRate
	end

	-- debug
	if timeCounter >= 1 then
		fixedUpdatesLastSecond = fixedUpdates
		fixedUpdates = 0
		timeCounter = timeCounter - 1
	end
end

function love.draw()
	for _, bunny in pairs(BunnyList) do
		love.graphics.draw(BunnySprite, bunny.posX, bunny.posY, 0, 0.2, 0.2)
	end
	-- love.graphics.print(string.format("Current FPS: %d\nBunnies: %d", love.timer.getFPS(), #BunnyList))
	love.graphics.print(
		string.format("FPS: %d\nTPS: %d\nBunnies: %d", love.timer.getFPS(), fixedUpdatesLastSecond, #BunnyList)
	)
end
