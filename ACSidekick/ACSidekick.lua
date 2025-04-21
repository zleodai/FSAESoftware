UIElements = require(".src/UIElements")
TelemetryCollector = require(".src/TelemetryCollector")

alreadyStart = false
windowActive = false
uiState = ac.getUI()

---@type UIElement[]
activeElements = {}
buttonElements = {}

direction = {
  n = vec2(0, 1),
  ne = vec2(math.sqrt(0.5), math.sqrt(0.5)),
  e = vec2(1, 0),
  se = vec2(math.sqrt(0.5), -math.sqrt(0.5)),
  s = vec2(0, -1),
  sw = vec2(-math.sqrt(0.5), -math.sqrt(0.5)),
  w = vec2(-1, 0),
  nw = vec2(-math.sqrt(0.5), math.sqrt(0.5))
}

function script.onShowWindow()
  -- Window Opened
  windowActive = true

  if not alreadyStart then
    -- Load configs/defaultValues

    ini = ac.getFolder(ac.FolderID.ContentTracks) .. '/' .. ac.getTrackFullID('/') .. '/data/map.ini'
    config = ac.INIConfig.load(ini):mapSection('PARAMETERS', { SCALE_FACTOR = 1, Z_OFFSET = 1, X_OFFSET = 1, WIDTH=500, HEIGHT=500, MARGIN=20, DRAWING_SIZE=10, MAX_SIZE=1000})
    config.OFFSETS = vec2(config.X_OFFSET, config.Z_OFFSET)

    displaySize = uiState.windowSize

    -- UI Element Setup

    ---@type UIElement[]
    activeElements = {}

    ---@type UIElement
    dataCollectionDisplay = UIElement:new()
    dataCollectionDisplay:setBackground(rgbm(0.1, 0.1, 0.1, 1))
    table.insert(activeElements, dataCollectionDisplay)

    ---@type UIElement
    collectDataButton = UIElement:new()
    table.insert(activeElements, collectDataButton)
    table.insert(buttonElements, collectDataButton)
    collectingData = false
    packetsSent = 0

    ---@type UIElement
    mapDisplay = UIElement:new()
    mapDisplay:setBackground(rgbm(0.1, 0.1, 0.1, 1))
    table.insert(activeElements, mapDisplay)

    local map_mini = ac.getFolder(ac.FolderID.ContentTracks) .. '\\' .. ac.getTrackFullID('\\') .. '\\map_mini.png'
    local map = ac.getFolder(ac.FolderID.ContentTracks) .. '\\' .. ac.getTrackFullID('\\') .. '\\map.png'
    mapImageSrc = io.exists(map_mini) and map_mini or map
    mapImageOffset = vec2(0, 50)
    mapImagePadding = vec2(20, 20)
    mapImageSize = ui.imageSize(mapImageSrc)
    mapScale = 1

    playerCar = ac.getCar(0)

    -- TelemetryCollection

    ---@type DBLink
    dbLink = DBLink:new("http://localhost:8090")

    alreadyStart = true
    guiInitalized = false
  end
end

function script.onHideWindow()
  -- Window Close
  windowActive = false
end

function script.onWindowUpdate(dt)
  -- On Window Update

  if not alreadyStart then return end
  local windows = ac.getAppWindows()
  for _, window in pairs(windows) do
    if window.title == "ACSidekick" then
      AppInfo = window
    end
  end
  ac.debug("windowInfo", AppInfo)

  windowSize = AppInfo.size
  defaultTitleFontSize = (windowSize.x * 0.1 + windowSize.y * 0.9)/32
  
  -- Handle onClicks
  if collectDataButton.onClick then
    collectingData = not collectingData
  end

  -- Update GUI
  if (dataCollectionDisplay ~= nil and mapDisplay ~= nil) then
    dataCollectionDisplay.pos = AppInfo.position
    dataCollectionDisplay.size = vec2(windowSize.x, windowSize.y * 0.5)
    dataCollectionDisplay:clearElements()

    dataCollectionDisplay:addText("Timing", defaultTitleFontSize, vec2(0, 0), vec2(windowSize.x, 50))
    dataCollectionDisplay:addText("Packets Sent: " .. packetsSent, defaultTitleFontSize/1.5, vec2(0, dataCollectionDisplay.size.y * 0.9), dataCollectionDisplay.size)

    collectDataButton.size = vec2(dataCollectionDisplay.size.x/3, dataCollectionDisplay.size.y/8)
    collectDataButton.pos = AppInfo.position + dataCollectionDisplay.size/2 - collectDataButton.size/2 + vec2(0, windowSize.y*0.175)
    collectDataButton:clearElements()

    if collectingData then 
      collectDataButton:addText("Unlink DB", defaultTitleFontSize/1.5, vec2(0, 0), collectDataButton.size)
      collectDataButton:setBackground(rgbm(0.3, 0, 0, 1))
    else
      collectDataButton:addText("Link DB", defaultTitleFontSize/1.5, vec2(0, 0), collectDataButton.size)
      collectDataButton:setBackground(rgbm(0, 0.3, 0, 1))
    end

    mapDisplay.pos = AppInfo.position + vec2(0, dataCollectionDisplay.size.y)
    mapDisplay.size = vec2(windowSize.x, windowSize.y - dataCollectionDisplay.size.y)
    mapDisplay:clearElements()

    -- Map positioning
    local mapSizeFactor = math.pow(1.15, mapScale)
    local mapOffset = (vec2(playerCar.position.x, playerCar.position.z) + config.OFFSETS) * (mapSizeFactor/config.SCALE_FACTOR) - mapDisplay.size/2
    rotationangle = 270 - math.deg(math.atan2(playerCar.look.x, playerCar.look.z))
    local p1 = -mapOffset
    local p2 = -mapOffset + (mapSizeFactor * mapImageSize)
    mapDisplay:addImage(mapImageSrc, p1, p2, rotationangle, mapDisplay.size/2)
    mapDisplay:addImage("./data/Arrow.png", mapDisplay.size/2 - vec2(25, 25), mapDisplay.size/2 + vec2(25, 25), 90, mapDisplay.size/2)

    guiInitalized = true
  end
end

function script.update(dt)
  if collectingData then

  end
end

function script.scenePreRenderUpdate()
  -- Called before a scene has started rendering
end

function script.postGeometryRenderUpdate()
  -- Called when opaque geometry has finished rendering
end

function script.preRenderUIUpdate()
  -- Called before rendering ImGui apps to draw things on screen

  if windowActive then
      uiState = ac.getUI()

      if guiInitalized then
      -- map scaling update
      if uiState.mouseWheel ~= 0 then
        if (uiState.mousePos.x > AppInfo.position.x and uiState.mousePos.y > AppInfo.position.y and uiState.mousePos.x < (AppInfo.position + AppInfo.size).x and uiState.mousePos.y < (AppInfo.position + AppInfo.size).y) then
          mapScale = mapScale + uiState.mouseWheel
        end
      end
      
      ui.pushDWriteFont(ui.DWriteFont("Chakra Petch;Weight=Light;", "./data"))

      for _, element in pairs(activeElements) do
        element:draw()
      end

      for _, element in pairs(buttonElements) do
        if (uiState.isMouseLeftKeyClicked) then
          if (uiState.mousePos.x > element.pos.x and uiState.mousePos.y > element.pos.y and uiState.mousePos.x < (element.pos + element.size).x and uiState.mousePos.y < (element.pos + element.size).y) then
            element.onClick = true
          else
            element.onClick = false
          end
        else
          element.onClick = false
        end
      end
    end
  end
end

function script.simulationUpdate()
  -- Called after a whole simulation update
end