package main

import (
	"fmt"
	"image/color"
	"os"
	"strconv"
	"strings"
	"time"

	"databaseAPI"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/exp/maps"
)

type sessionLap struct {
	sessionId int64
	lapId int64
}

func sessionLapCmp (a  sessionLap, other sessionLap) int {
	if a.sessionId != other.sessionId { return int(a.sessionId - other.sessionId) }
	return int(a.lapId - other.lapId)
}

type appTheme struct {}
func (m appTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
		case theme.ColorNameInputBackground:
			return color.Black
		case theme.ColorNameBackground:
			return color.Black
		case theme.ColorNameButton:
			return color.White
		case theme.ColorNameHover:
			return color.NRGBA{R: 150, G: 150, B: 150, A: 255}
		case theme.ColorNameForeground:
			return color.White		
	}

	return theme.DefaultTheme().Color(name, variant)
}
func (m appTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}
func (m appTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}
func (m appTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}


var (
	_ fyne.Theme = (*appTheme)(nil)

	currentTime time.Time = time.Now()

	mainApp fyne.App
	mainWindow fyne.Window

	sideBarOffset float32 = 200
	bottomBarHeight float32 = 300

	//region containers/labels
	mapContainer *fyne.Container

	mapActiveSessionLaps []sessionLap
	mapTelemetryPoints []*widget.Button

	graphContainer *fyne.Container

	currentInfoLabel *widget.Label

	tireTempsGraphContainer *fyne.Container
	tirePressuresGraphContainer *fyne.Container
	brakeGraphContainer *fyne.Container
	throttleGraphContainer *fyne.Container
	steeringGraphContainer *fyne.Container
	accelerationGraphContainer *fyne.Container
	speedGraphContainer *fyne.Container

	graphSelectionContainer *fyne.Container

	activeGraphs []*fyne.Container
	allGraphs []*fyne.Container


	lapSelectContainer *fyne.Container

	lapSelectTotalSessions *widget.Label
	lapSelectTotalLaps *widget.Label

	lapSelectSessionEntry *widget.Entry
	lapSelectLapEntry *widget.Entry

	playbackControlsContainer *fyne.Container

	
	objectOffsets map[*fyne.CanvasObject][2]float32
	containerPositions map[*fyne.Container][2]float32
	containerOffsets map[*fyne.Container][2]float32
	containerItems map[*fyne.Container][]*fyne.CanvasObject
	containerContainers map[*fyne.Container][]*fyne.Container

	//endregion

	dbConnection *pgxpool.Pool

	allMiniPackets []databaseAPI.MiniTelemetryPacket
	lapsBySession map[int64][]int64
	//Min and max packetIds for a sessionLap: [inital, final]
	packetIdsBySessionLap map[sessionLap][2]int64
	totalSessions int64

	currentSessionId int64 = 0
	selectedLapIds []int64 = []int64{}
	currentLapId int64 = 0
	currentPacketId int64 = 0

	playing bool = false
	forwardStep int64 = 10
	backwardStep int64 = 10
	defaultPlaybackSpeed int64 = 500
	playbackSpeed int64 = 500

	defaultPacketStepInterval time.Duration = time.Millisecond * 1
	nextPacketTimestamp time.Time = currentTime.Add(defaultPacketStepInterval)
)

const (
	refreshInterval time.Duration = time.Millisecond * 5
	resolutionWidth float32 = 1920
	resolutionHeight float32 = 1080

	windowOffset float32 = 4
)

func onStart() {
	dbConnection = databaseAPI.NewConnection()
	allMiniPackets = *databaseAPI.QueryMiniPacketsFromPool(dbConnection)
	fmt.Printf("%d Packets\n", len(allMiniPackets))

	mainApp = app.New()
	mainWindow = mainApp.NewWindow("Vehicle Telemetry Monitor")
	mainWindow.Resize(fyne.NewSize(resolutionWidth, resolutionHeight))

	mainApp.Settings().SetTheme(&appTheme{})

	//region getting sessionLap data
	lapsBySession = make(map[int64][]int64)
	packetIdsBySessionLap = make(map[sessionLap][2]int64)

	var lastLapId int64 = 0
	var currentSession int64 = 0
	for _, miniPacket := range allMiniPackets {
		if lastLapId > miniPacket.LapId {
			currentSession += 1
			lapsBySession[currentSession] = []int64{miniPacket.LapId}

			lastSessionLap := sessionLap{sessionId: currentSession -1, lapId: lastLapId}
			packetIds := packetIdsBySessionLap[lastSessionLap]
			packetIds[1] = miniPacket.Id -1
			packetIdsBySessionLap[lastSessionLap] = packetIds

			lastLapId = miniPacket.LapId

			newSessionLap := sessionLap{sessionId: currentSession, lapId: miniPacket.LapId}
			packetIdsBySessionLap[newSessionLap] = [2]int64{miniPacket.Id, 0}
		} else if miniPacket.LapId != lastLapId {
			lapsBySession[currentSession] = append(lapsBySession[currentSession], miniPacket.LapId)

			lastSessionLap := sessionLap{sessionId: currentSession, lapId: lastLapId}
			packetIds := packetIdsBySessionLap[lastSessionLap]
			packetIds[1] = miniPacket.Id -1
			packetIdsBySessionLap[lastSessionLap] = packetIds

			lastLapId = miniPacket.LapId

			newSessionLap := sessionLap{sessionId: currentSession, lapId: miniPacket.LapId}
			packetIdsBySessionLap[newSessionLap] = [2]int64{miniPacket.Id, 0}
		}
	}
	veryLastSessionLap := sessionLap{sessionId: currentSession, lapId: lastLapId}
	packetIds := packetIdsBySessionLap[veryLastSessionLap]
	packetIds[1] = allMiniPackets[len(allMiniPackets) -1].Id
	packetIdsBySessionLap[veryLastSessionLap] = packetIds

	totalSessions = currentSession + 1
	// keys := maps.Keys(packetIdsBySessionLap)
	// slices.SortFunc(keys, func(a, b sessionLap) int {return sessionLapCmp(a, b)})

	// for _, sessionLapKey := range keys {
	// 	packetIds := packetIdsBySessionLap[sessionLapKey]
	// 	fmt.Printf("Session %d.%d has packets %d - %d\n", sessionLapKey.sessionId, sessionLapKey.lapId, packetIds[0], packetIds[1])
	// }
	//endregion

	//region container initalization
	activeGraphs = []*fyne.Container{}
	allGraphs = []*fyne.Container{}

	objectOffsets = make(map[*fyne.CanvasObject][2]float32)
	containerPositions = make(map[*fyne.Container][2]float32)
	containerOffsets = make(map[*fyne.Container][2]float32)
	containerItems = make(map[*fyne.Container][]*fyne.CanvasObject)
	containerContainers = make(map[*fyne.Container][]*fyne.Container)

	mainContainer := container.NewWithoutLayout()
	containerPositions[mainContainer] = [2]float32{0, 0}
	containerOffsets[mainContainer] = [2]float32{0, 0}
	containerContainers[mainContainer] = []*fyne.Container{}

	mapContainer = container.NewWithoutLayout()
	addContainer(mainContainer, mapContainer, -windowOffset, -windowOffset)
	containerItems[mapContainer] = []*fyne.CanvasObject{}
	containerContainers[mapContainer] = []*fyne.Container{}

	var mapBackgroundColor = color.NRGBA{R: 33, G: 33, B: 33, A: 255}

	var mapBackgroundP fyne.CanvasObject = canvas.NewRectangle(mapBackgroundColor)
	resizeObject(mapBackgroundP, resolutionWidth/2 + sideBarOffset, resolutionHeight - bottomBarHeight)
	addObject(mapContainer, &mapBackgroundP, 0, 0)

	graphContainer = container.NewWithoutLayout()
	addContainer(mainContainer, graphContainer, resolutionWidth/2 + sideBarOffset - windowOffset, -windowOffset)
	containerItems[graphContainer] = []*fyne.CanvasObject{}
	containerContainers[graphContainer] = []*fyne.Container{}

	var graphBackgroundColor = color.NRGBA{R: 33, G: 33, B: 33, A: 255}

	var graphBackgroundP fyne.CanvasObject = canvas.NewRectangle(graphBackgroundColor)
	resizeObject(graphBackgroundP, resolutionWidth/2 - sideBarOffset, resolutionHeight * 2)
	addObject(graphContainer, &graphBackgroundP, 0, 0)

	tireTempsGraphContainer	= container.NewWithoutLayout()
	addContainer(graphContainer, tireTempsGraphContainer, 0, 0)
	containerItems[tireTempsGraphContainer] = []*fyne.CanvasObject{}
	containerContainers[tireTempsGraphContainer] = []*fyne.Container{}
	tireTempsGraphContainer.Hidden = true
	allGraphs = append(allGraphs, tireTempsGraphContainer)

	var graphsBackgroundColor = color.NRGBA{R: 40, G: 40, B: 40, A: 255}
	var graphsWidth float32 = resolutionWidth/2 - sideBarOffset - 40
	var graphsHeight float32 = 180

	var tireTempsGraphBackgroundP fyne.CanvasObject = canvas.NewRectangle(graphsBackgroundColor)
	resizeObject(tireTempsGraphBackgroundP, graphsWidth, graphsHeight)
	addObject(tireTempsGraphContainer, &tireTempsGraphBackgroundP, 0, 0)

	var tireTempsGraphTitle *widget.Label = widget.NewLabel("tireTemps")
	var tireTempsGraphTitleP fyne.CanvasObject = tireTempsGraphTitle
	resizeObject(tireTempsGraphTitleP, 100, 40)
	addObject(tireTempsGraphContainer, &tireTempsGraphTitleP, 20, 10)

	tirePressuresGraphContainer = container.NewWithoutLayout()
	addContainer(graphContainer, tirePressuresGraphContainer, 0, 0)
	containerItems[tirePressuresGraphContainer] = []*fyne.CanvasObject{}
	containerContainers[tirePressuresGraphContainer] = []*fyne.Container{}
	tirePressuresGraphContainer.Hidden = true
	allGraphs = append(allGraphs, tirePressuresGraphContainer)

	var tirePressuresGraphBackgroundP fyne.CanvasObject = canvas.NewRectangle(graphsBackgroundColor)
	resizeObject(tirePressuresGraphBackgroundP, graphsWidth, graphsHeight)
	addObject(tirePressuresGraphContainer, &tirePressuresGraphBackgroundP, 0, 0)

	var tirePressuresGraphTitle *widget.Label = widget.NewLabel("tirePressures")
	var tirePressuresGraphTitleP fyne.CanvasObject = tirePressuresGraphTitle
	resizeObject(tirePressuresGraphTitleP, 100, 40)
	addObject(tirePressuresGraphContainer, &tirePressuresGraphTitleP, 20, 10)

	brakeGraphContainer = container.NewWithoutLayout()
	addContainer(graphContainer, brakeGraphContainer, 0, 0)
	containerItems[brakeGraphContainer] = []*fyne.CanvasObject{}
	containerContainers[brakeGraphContainer] = []*fyne.Container{}
	brakeGraphContainer.Hidden = true
	allGraphs = append(allGraphs, brakeGraphContainer)

	var brakeGraphBackgroundP fyne.CanvasObject = canvas.NewRectangle(graphsBackgroundColor)
	resizeObject(brakeGraphBackgroundP, graphsWidth, graphsHeight)
	addObject(brakeGraphContainer, &brakeGraphBackgroundP, 0, 0)

	var brakeGraphTitle *widget.Label = widget.NewLabel("brake")
	var brakeGraphTitleP fyne.CanvasObject = brakeGraphTitle
	resizeObject(brakeGraphTitleP, 100, 40)
	addObject(brakeGraphContainer, &brakeGraphTitleP, 20, 10)

	throttleGraphContainer = container.NewWithoutLayout()
	addContainer(graphContainer, throttleGraphContainer, 0, 0)
	containerItems[throttleGraphContainer] = []*fyne.CanvasObject{}
	containerContainers[throttleGraphContainer] = []*fyne.Container{}
	throttleGraphContainer.Hidden = true
	allGraphs = append(allGraphs, throttleGraphContainer)

	var throttleGraphBackgroundP fyne.CanvasObject = canvas.NewRectangle(graphsBackgroundColor)
	resizeObject(throttleGraphBackgroundP, graphsWidth, graphsHeight)
	addObject(throttleGraphContainer, &throttleGraphBackgroundP, 0, 0)

	var throttleGraphTitle *widget.Label = widget.NewLabel("throttle")
	var throttleGraphTitleP fyne.CanvasObject = throttleGraphTitle
	resizeObject(throttleGraphTitleP, 100, 40)
	addObject(throttleGraphContainer, &throttleGraphTitleP, 20, 10)

	steeringGraphContainer = container.NewWithoutLayout()
	addContainer(graphContainer, steeringGraphContainer, 0, 0)
	containerItems[steeringGraphContainer] = []*fyne.CanvasObject{}
	containerContainers[steeringGraphContainer] = []*fyne.Container{}
	steeringGraphContainer.Hidden = true
	allGraphs = append(allGraphs, steeringGraphContainer)

	var steeringGraphBackgroundP fyne.CanvasObject = canvas.NewRectangle(graphsBackgroundColor)
	resizeObject(steeringGraphBackgroundP, graphsWidth, graphsHeight)
	addObject(steeringGraphContainer, &steeringGraphBackgroundP, 0, 0)

	var steeringGraphTitle *widget.Label = widget.NewLabel("steering")
	var steeringGraphTitleP fyne.CanvasObject = steeringGraphTitle
	resizeObject(steeringGraphTitleP, 100, 40)
	addObject(steeringGraphContainer, &steeringGraphTitleP, 20, 10)

	accelerationGraphContainer = container.NewWithoutLayout()
	addContainer(graphContainer, accelerationGraphContainer, 0, 0)
	containerItems[accelerationGraphContainer] = []*fyne.CanvasObject{}
	containerContainers[accelerationGraphContainer] = []*fyne.Container{}
	accelerationGraphContainer.Hidden = true
	allGraphs = append(allGraphs, accelerationGraphContainer)

	var accelerationGraphBackgroundP fyne.CanvasObject = canvas.NewRectangle(graphsBackgroundColor)
	resizeObject(accelerationGraphBackgroundP, graphsWidth, graphsHeight)
	addObject(accelerationGraphContainer, &accelerationGraphBackgroundP, 0, 0)

	var accelerationGraphTitle *widget.Label = widget.NewLabel("acceleration")
	var accelerationGraphTitleP fyne.CanvasObject = accelerationGraphTitle
	resizeObject(accelerationGraphTitleP, 100, 40)
	addObject(accelerationGraphContainer, &accelerationGraphTitleP, 20, 10)

	speedGraphContainer = container.NewWithoutLayout()
	addContainer(graphContainer, speedGraphContainer, 0, 0)
	containerItems[speedGraphContainer] = []*fyne.CanvasObject{}
	containerContainers[speedGraphContainer] = []*fyne.Container{}
	speedGraphContainer.Hidden = true
	allGraphs = append(allGraphs, speedGraphContainer)

	var speedGraphBackgroundP fyne.CanvasObject = canvas.NewRectangle(graphsBackgroundColor)
	resizeObject(speedGraphBackgroundP, graphsWidth, graphsHeight)
	addObject(speedGraphContainer, &speedGraphBackgroundP, 0, 0)

	var speedGraphTitle *widget.Label = widget.NewLabel("speed")
	var speedGraphTitleP fyne.CanvasObject = speedGraphTitle
	resizeObject(speedGraphTitleP, 100, 40)
	addObject(speedGraphContainer, &speedGraphTitleP, 20, 10)

	var graphSelectionBorderColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	var graphSelectionBorderWidth float32 = 2
	var graphSelectionHeight float32 = 80
	var graphSelectionSizeOffset float32 = 40

	graphSelectionContainer = container.NewWithoutLayout()
	containerItems[graphSelectionContainer] = []*fyne.CanvasObject{}
	containerContainers[graphSelectionContainer] = []*fyne.Container{}
	addContainer(mainContainer, graphSelectionContainer, resolutionWidth/2 + sideBarOffset + graphSelectionSizeOffset/2 - windowOffset, resolutionHeight - bottomBarHeight - graphSelectionHeight - graphSelectionSizeOffset/2 - windowOffset)

	var emptyColor = color.NRGBA{R: 0, G: 0, B: 0, A: 0}

	var graphSelectionBorder *canvas.Rectangle = canvas.NewRectangle(emptyColor)
	var graphSelectionBorderP fyne.CanvasObject = graphSelectionBorder
	graphSelectionBorder.StrokeColor = graphSelectionBorderColor
	graphSelectionBorder.StrokeWidth = graphSelectionBorderWidth
	resizeObject(graphSelectionBorderP, resolutionWidth/2 - sideBarOffset - graphSelectionSizeOffset, graphSelectionHeight)
	addObject(graphSelectionContainer, &graphSelectionBorderP, 0, 0)

	var checkWidthInitalOffset float32 = 40
	var checkHeightIniitalOffset float32 = 10
	var checkHeightSpacing float32 = 30

	var tireTempSelect *widget.Check = widget.NewCheck("Tire Temps", TireTempsCheckBoxChanged)
	var tireTempSelectP fyne.CanvasObject = tireTempSelect
	resizeObject(tireTempSelectP, 20, 20)
	addObject(graphSelectionContainer, &tireTempSelectP, checkWidthInitalOffset, checkHeightIniitalOffset)

	var tirePressureSelect *widget.Check = widget.NewCheck("Tire Pressures", TirePressuresCheckBoxChanged)
	var tirePressureSelectP fyne.CanvasObject = tirePressureSelect
	resizeObject(tirePressureSelectP, 20, 20)
	addObject(graphSelectionContainer, &tirePressureSelectP, checkWidthInitalOffset + 150, checkHeightIniitalOffset)

	var brakeSelect *widget.Check = widget.NewCheck("Brake", BrakeCheckBoxChanged)
	var brakeSelectP fyne.CanvasObject = brakeSelect
	resizeObject(brakeSelectP, 20, 20)
	addObject(graphSelectionContainer, &brakeSelectP, checkWidthInitalOffset + 300, checkHeightIniitalOffset)

	var throttleSelect *widget.Check = widget.NewCheck("Throttle", ThrottleCheckBoxChanged)
	var throttleSelectP fyne.CanvasObject = throttleSelect
	resizeObject(throttleSelectP, 20, 20)
	addObject(graphSelectionContainer, &throttleSelectP, checkWidthInitalOffset + 400, checkHeightIniitalOffset)

	var steeringSelect *widget.Check = widget.NewCheck("Steering", SteeringCheckBoxChanged)
	var steeringSelectP fyne.CanvasObject = steeringSelect
	resizeObject(steeringSelectP, 20, 20)
	addObject(graphSelectionContainer, &steeringSelectP, checkWidthInitalOffset + 520, checkHeightIniitalOffset)

	var accelerationSelect *widget.Check = widget.NewCheck("Acceleration", AccelerationCheckBoxChanged)
	var accelerationSelectP fyne.CanvasObject = accelerationSelect
	resizeObject(accelerationSelectP, 20, 20)
	addObject(graphSelectionContainer, &accelerationSelectP, checkWidthInitalOffset, checkHeightIniitalOffset + checkHeightSpacing)

	var speedSelect *widget.Check = widget.NewCheck("Speed", SpeedCheckBoxChanged)
	var speedSelectP fyne.CanvasObject = speedSelect
	resizeObject(speedSelectP, 20, 20)
	addObject(graphSelectionContainer, &speedSelectP, checkWidthInitalOffset + 150, checkHeightIniitalOffset + checkHeightSpacing)

	currentInfoLabel = widget.NewLabel("Session , Lap , Packet ")
	currentInfoLabel.TextStyle = fyne.TextStyle{Bold: true}
	var currentInfoLabelP fyne.CanvasObject = currentInfoLabel
	resizeObject(currentInfoLabelP, 200, 40)
	addObject(graphContainer, &currentInfoLabelP, (resolutionWidth/2 - sideBarOffset)/2 - 200/2, 20)

	lapSelectContainer = container.NewWithoutLayout()
	addContainer(mainContainer, lapSelectContainer, -windowOffset, resolutionHeight - bottomBarHeight - windowOffset)
	containerItems[lapSelectContainer] = []*fyne.CanvasObject{}
	containerContainers[lapSelectContainer] = []*fyne.Container{}

	var lapSelectBackgroundColor = color.NRGBA{R: 25, G: 25, B: 25, A: 255}

	var lapSelectBackgroundP fyne.CanvasObject = canvas.NewRectangle(lapSelectBackgroundColor)
	resizeObject(lapSelectBackgroundP, resolutionWidth/2 + sideBarOffset, bottomBarHeight)
	addObject(lapSelectContainer, &lapSelectBackgroundP, 0, 0)

	var lapSelectItemWidthOffset float32 = -50

	lapSelectTotalSessions = widget.NewLabel(fmt.Sprintf("%d sessions", totalSessions))
	var lapSelectTotalSessionsP fyne.CanvasObject = lapSelectTotalSessions
	resizeObject(lapSelectTotalSessionsP, 300, 40)
	addObject(lapSelectContainer, &lapSelectTotalSessionsP, 325 + lapSelectItemWidthOffset, bottomBarHeight/4)

	var lapSelectSessionText *widget.Label = widget.NewLabel("Session")
	lapSelectSessionText.TextStyle = fyne.TextStyle{Bold: true}
	var lapSelectSessionTextP fyne.CanvasObject = lapSelectSessionText
	resizeObject(lapSelectSessionTextP, 200, 80)
	addObject(lapSelectContainer, &lapSelectSessionTextP, 40 + 125 + lapSelectItemWidthOffset, bottomBarHeight/2)

	lapSelectSessionEntry = widget.NewEntry()
	lapSelectSessionEntry.SetPlaceHolder("Enter Session ID")
	var lapSelectSessionEntryP fyne.CanvasObject = lapSelectSessionEntry
	resizeObject(lapSelectSessionEntryP, 300, 40)
	addObject(lapSelectContainer, &lapSelectSessionEntryP, 40 + 200 + lapSelectItemWidthOffset, bottomBarHeight/2)

	lapSelectTotalLaps = widget.NewLabel(fmt.Sprintf("%d laps", len(lapsBySession[currentSession])))
	var lapSelectTotalLapsP fyne.CanvasObject = lapSelectTotalLaps
	resizeObject(lapSelectTotalLapsP, 300, 40)
	addObject(lapSelectContainer, &lapSelectTotalLapsP, (resolutionWidth + sideBarOffset)/2 - 125 + lapSelectItemWidthOffset, bottomBarHeight/4)

	var lapSelectLapText *widget.Label = widget.NewLabel("Lap")
	lapSelectLapText.TextStyle = fyne.TextStyle{Bold: true}
	var lapSelectLapTextP fyne.CanvasObject = lapSelectLapText
	resizeObject(lapSelectLapTextP, 200, 80)
	addObject(lapSelectContainer, &lapSelectLapTextP, (resolutionWidth + sideBarOffset)/2 - 300 + lapSelectItemWidthOffset, bottomBarHeight/2)

	lapSelectLapEntry = widget.NewEntry()
	lapSelectLapEntry.SetPlaceHolder("Enter Lap IDs")
	var lapSelectLapEntryP fyne.CanvasObject = lapSelectLapEntry
	resizeObject(lapSelectLapEntryP, 300, 40)
	addObject(lapSelectContainer, &lapSelectLapEntryP, (resolutionWidth + sideBarOffset)/2 - 300*1.5 + 200 + lapSelectItemWidthOffset, bottomBarHeight/2)

	var playbackControlsBackgroundColor = color.NRGBA{R: 25, G: 25, B: 25, A: 255}

	playbackControlsContainer = container.NewWithoutLayout()
	addContainer(mainContainer, playbackControlsContainer, resolutionWidth/2 + sideBarOffset - windowOffset, resolutionHeight - bottomBarHeight - windowOffset)
	containerItems[playbackControlsContainer] = []*fyne.CanvasObject{}
	containerContainers[playbackControlsContainer] = []*fyne.Container{}

	var playbackControlsBackgroundP fyne.CanvasObject = canvas.NewRectangle(playbackControlsBackgroundColor)
	resizeObject(playbackControlsBackgroundP, resolutionWidth/2 - sideBarOffset, bottomBarHeight)
	addObject(playbackControlsContainer, &playbackControlsBackgroundP, 0, 0)

	mainWindow.Canvas().SetContent(mainContainer)
	//endregion

	//region hotkeys
	mainWindow.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		var moveSpeed float32 = 25

		var containerHead = mainContainer
		var containerToMove = graphContainer

		switch key.Name {
			case fyne.KeyS:
				if containerOffsets[containerToMove][1] + moveSpeed > resolutionHeight * -1 {
					moveContainer(containerHead, containerToMove, containerOffsets[containerToMove][0], containerOffsets[containerToMove][1] - moveSpeed)
				}
			case fyne.KeyW:
				if containerOffsets[containerToMove][1] - moveSpeed < -29 {
					moveContainer(containerHead, containerToMove, containerOffsets[containerToMove][0], containerOffsets[containerToMove][1] + moveSpeed)
				}
			case fyne.KeyA:
				currentSessionLap := sessionLap{sessionId: currentSessionId, lapId: currentLapId}
				packetIds := packetIdsBySessionLap[currentSessionLap]
				minPacketId := packetIds[0]
				if currentPacketId - 1 < minPacketId {
					var minLap int64 = 9223372036854775807
					for _, lap := range lapsBySession[currentSessionId] {
						if lap < minLap {
							minLap = lap
						}
					}
					if currentLapId - 1 >= minLap {
						currentLapId -= 1
					} else {
						var minSession int64 = 9223372036854775807
						for _, session := range maps.Keys(lapsBySession) {
							if session < minSession {
								minSession = session
							}
						}

						if currentSessionId - 1 >= minSession {
							currentSessionId -= 1
							
							var maxLap int64 = -9223372036854775808
							for _, lap := range lapsBySession[currentSessionId] {
								if lap > maxLap {
									maxLap = lap
								}
							}	
							currentLapId = maxLap
						} else {
							currentPacketId += 1
						}
					}
				}
				currentPacketId -= 1
			case fyne.KeyD:
				currentSessionLap := sessionLap{sessionId: currentSessionId, lapId: currentLapId}
				packetIds := packetIdsBySessionLap[currentSessionLap]
				maxPacketId := packetIds[1] 
				if currentPacketId + 1 > maxPacketId {
					var maxLap int64 = -9223372036854775808
					for _, lap := range lapsBySession[currentSessionId] {
						if lap > maxLap {
							maxLap = lap
						}
					}
					if currentLapId + 1 <= maxLap {
						currentLapId += 1
					} else {
						var maxSession int64 = -9223372036854775808
						for _, session := range maps.Keys(lapsBySession) {
							if session > maxSession {
								maxSession = session
							}
						}

						if currentSessionId + 1 <= maxSession {
							currentSessionId += 1

							var minLap int64 = 9223372036854775807
							for _, lap := range lapsBySession[currentSessionId] {
								if lap < minLap {
									minLap = lap
								}
							}
							currentLapId = minLap
						} else {
							currentPacketId -= 1
						}
					}
				}
				currentPacketId += 1
			case fyne.KeySpace:
				playing = !playing
		}
	})
	//endregion

}

func onUpdate(deltaTime time.Duration) {
	currentTime = currentTime.Add(deltaTime)

	if playing && currentTime.Compare(nextPacketTimestamp) >= 0 {
		nextPacketTimestamp = nextPacketTimestamp.Add(defaultPacketStepInterval * time.Duration(playbackSpeed))
		
		maxPacketId := packetIdsBySessionLap[sessionLap{sessionId: currentSessionId, lapId: currentLapId}][1]
		if currentPacketId + 1 <= maxPacketId {
			currentPacketId += 1
		}
	}

	checkSelectedSessionLap()
	refreshLabels()
	refreshMaps()
	refreshGraphs()
}

func checkSelectedSessionLap() {
	selectedSessionId, sessionErr := strconv.ParseInt(lapSelectSessionEntry.Text, 10, 64)
	if sessionErr != nil {
		return
	} 
	
	lapsSelected := []int64{}
	for _, splitString := range strings.Split(lapSelectLapEntry.Text, ",") {
		lap, err := strconv.ParseInt(strings.ReplaceAll(strings.ReplaceAll(splitString, " ", ""), "	", ""), 10, 64)
		if err != nil {
			return
		}
		lapsSelected = append(lapsSelected, lap)
	}

	if selectedSessionId != currentSessionId {
		currentSessionId = selectedSessionId
		var minLap int64 = 9223372036854775807
		for _, lap := range lapsSelected {
			if lap < minLap {
				minLap = lap
			}
		}
		currentLapId = minLap
		currentPacketId = packetIdsBySessionLap[sessionLap{sessionId: currentSessionId, lapId: currentLapId}][0]
	} else {
		selectedLapIds = append(selectedLapIds, currentLapId)
	}
	selectedLapIds = lapsSelected
}

func refreshLabels() {
	currentInfoLabel.SetText(fmt.Sprintf("Session %d, Lap %d, Packet %d", currentSessionId, currentLapId, currentPacketId))
	lapSelectTotalSessions.SetText(fmt.Sprintf("%d sessions", totalSessions))
	lapSelectTotalLaps.SetText(fmt.Sprintf("%d laps", len(lapsBySession[currentSessionId])))
}


func refreshMaps() {
	
}

func refreshGraphs() {

}

func refreshGraphVisiblity() {
	for _, graphContainer := range allGraphs {
		graphContainer.Hidden = true
	}

	var graphContainerWidthOffset float32 = 20
	var graphContainerHeightOffset float32 = 80
	var graphContainerHeightSpacing float32 = 200

	for index, container := range activeGraphs {
		moveContainer(graphContainer, container, graphContainerWidthOffset, graphContainerHeightOffset + float32(index) * graphContainerHeightSpacing)
		container.Hidden = false
	}
}

func removeFromArray(arr []*fyne.Container, item *fyne.Container) []*fyne.Container {
	var indexOfItem int = 0
	var found bool = false
	for index, arrItem := range arr {
		if !found && item == arrItem {
			indexOfItem = index
			found = true
		}
	}

	if !found {
		fmt.Println("called removeFromArray when item did not exist in array")
		return arr
	} else {
		arr[indexOfItem] = arr[len(arr) -1]
		return arr[:len(arr) -1]
	}
}

//region checkBoxFunctions
func TireTempsCheckBoxChanged(status bool) {if status {activeGraphs = append(activeGraphs, tireTempsGraphContainer) } else {activeGraphs = removeFromArray(activeGraphs, tireTempsGraphContainer)}; refreshGraphVisiblity()}
func TirePressuresCheckBoxChanged(status bool) {if status {activeGraphs = append(activeGraphs, tirePressuresGraphContainer) } else {activeGraphs = removeFromArray(activeGraphs, tirePressuresGraphContainer)}; refreshGraphVisiblity()}
func BrakeCheckBoxChanged(status bool) {if status {activeGraphs = append(activeGraphs, brakeGraphContainer) } else {activeGraphs = removeFromArray(activeGraphs, brakeGraphContainer)}; refreshGraphVisiblity()}
func ThrottleCheckBoxChanged(status bool) {if status {activeGraphs = append(activeGraphs, throttleGraphContainer) } else {activeGraphs = removeFromArray(activeGraphs, throttleGraphContainer)}; refreshGraphVisiblity()}
func SteeringCheckBoxChanged(status bool) {if status {activeGraphs = append(activeGraphs, steeringGraphContainer) } else {activeGraphs = removeFromArray(activeGraphs, steeringGraphContainer)}; refreshGraphVisiblity()}
func AccelerationCheckBoxChanged(status bool) {if status {activeGraphs = append(activeGraphs, accelerationGraphContainer) } else {activeGraphs = removeFromArray(activeGraphs, accelerationGraphContainer)}; refreshGraphVisiblity()}
func SpeedCheckBoxChanged(status bool) {if status {activeGraphs = append(activeGraphs, speedGraphContainer) } else {activeGraphs = removeFromArray(activeGraphs, speedGraphContainer)}; refreshGraphVisiblity()}
//endregion

func moveContainer(headContainer *fyne.Container, childContainer *fyne.Container, xPos float32, yPos float32) {
	containerOffsets[childContainer] = [2]float32{xPos, yPos}

	headContainerPosition, inDict := containerPositions[headContainer]
	if !inDict {
		fmt.Println("Container Offset Not Found")
		os.Exit(1)
	}
	
	containerPositions[childContainer] = [2]float32{headContainerPosition[0] + xPos, headContainerPosition[1] + yPos}
	for _, object := range containerItems[childContainer] {
		objectOffset, inDict := objectOffsets[object]
		if !inDict {
			fmt.Println("Object Offset Not Found (please use moveObject() atleast once for any created canvas object assigned to a container)")
			os.Exit(1)
		}
		moveObject(childContainer, object, objectOffset[0], objectOffset[1])
	}
	for _, container := range containerContainers[childContainer] {
		containerOffset, inDict := containerOffsets[container]
		if !inDict {
			fmt.Println("Container Offset Not Found (please add container to containerOffsets any created containers)")
			os.Exit(1)
		}
		moveContainer(childContainer, container, containerOffset[0], containerOffset[1])
	}
}

func removeContainer(container *fyne.Container) {
	for headContainer, childContainers := range containerContainers{
		var inContainer bool = false
		var indexInArray int = 0
		for index, childContainer := range childContainers {
			if childContainer == container {
				inContainer = true
				indexInArray = index
			}
		}

		if inContainer {
			headContainer.Remove(container)
			childContainers[indexInArray] = childContainers[len(childContainers) -1]
			containerContainers[headContainer] = childContainers[:len(childContainers) -1]
		}
	}

	for _, object := range containerItems[container] {
		delete(objectOffsets, object)
		container.Remove(*object)
	}

	delete(containerItems, container)
	delete(containerContainers, container)
	delete(containerOffsets, container)
	delete(containerPositions, container)

	container.RemoveAll()
}

func addContainer(headContainer *fyne.Container, childContainer *fyne.Container, xPos float32, yPos float32) {
	containerContainers[headContainer] = append(containerContainers[headContainer], childContainer)
	headContainer.Add(childContainer)
	moveContainer(headContainer, childContainer, xPos, yPos)
}

func addObject(container *fyne.Container, object *fyne.CanvasObject, xPos float32, yPos float32) {
	containerItems[container] = append(containerItems[container], object)
	container.Add(*object)
	moveObject(container, object, xPos, yPos)
}

func moveObject(container *fyne.Container, object *fyne.CanvasObject, xPos float32, yPos float32) {
	objectOffsets[object] = [2]float32{xPos, yPos}

	containerPosition, inDict := containerPositions[container]
	if !inDict {
		fmt.Println("Container Offset Not Found")
		os.Exit(1)
	}

	(*object).Move(fyne.NewPos(containerPosition[0] + xPos, containerPosition[1] + yPos))
}

func resizeObject(object fyne.CanvasObject, width float32, height float32) {
	object.Resize((fyne.NewSize(width, height)))
}

func main() {
	onStart()

	go func() {
		for range time.Tick(refreshInterval) {
			onUpdate(refreshInterval)
		}
	}()

	mainWindow.ShowAndRun()
}