package main

import (
	"fmt"
	"image/color"
	"math"
	"os"
	"slices"
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

//region customstructs
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
//endregion

var (
	_ fyne.Theme = (*appTheme)(nil)

	currentTime time.Time = time.Now()

	mainApp fyne.App
	mainWindow fyne.Window

	sideBarOffset float32 = 200
	bottomBarHeight float32 = 150

	//IMPORTANT FOR MAP
	packetInterval = 1
	defaultCircleSize float32 = 5
	selectedCircleSizeIncrease float32 = 6.7

	//region containers/labels
	mainContainer *fyne.Container

	mapContainer *fyne.Container
	mapBackgroundP fyne.CanvasObject

	mapActiveSessionLaps []sessionLap
	mapTelemetryPoints map[sessionLap][]*canvas.Circle

	mapLastSelectedPacketId int64 = -404
	mapLastSelectedSessionLap sessionLap = sessionLap{sessionId: 0, lapId: 0}

	mapWidth float32 = resolutionWidth/2 + sideBarOffset
	mapHeight float32 = resolutionHeight - bottomBarHeight

	lapColors = []color.Color{color.NRGBA{R: 255, G: 255, B: 255, A: 255}, color.NRGBA{R: 0, G: 0, B: 255, A: 255}, color.NRGBA{R: 255, G: 0, B: 0, A: 255}, color.NRGBA{R: 0, G: 255, B: 0, A: 255}, color.NRGBA{R: 255, G: 255, B: 0, A: 255}, color.NRGBA{R: 0, G: 255, B: 255, A: 255}, color.NRGBA{R: 255, G: 0, B: 255, A: 255}}

	legendContainer *fyne.Container

	legendPoints []*canvas.Circle
	legendLabels []*widget.Label

	legendPointSize float32 = 20

	graphContainer *fyne.Container

	currentInfoLabel *widget.Label

	tireTempsGraphContainer *fyne.Container
	tirePressuresGraphContainer *fyne.Container
	brakeGraphContainer *fyne.Container
	throttleGraphContainer *fyne.Container
	steeringGraphContainer *fyne.Container
	accelerationGraphContainer *fyne.Container
	speedGraphContainer *fyne.Container

	tireColors map[int]color.Color = map[int]color.Color{0: color.RGBA{R: 255, G: 255, B: 255, A: 255}, 1: color.RGBA{R: 0, G: 255, B: 0, A: 255}, 2: color.RGBA{R: 0, G: 0, B: 255, A: 255}, 3: color.RGBA{R: 255, G: 0, B: 0, A: 255}}

	graphSelectionContainer *fyne.Container

	activeGraphs []*fyne.Container
	allGraphs []*fyne.Container
	graphLines map[*fyne.Container][]*canvas.Line

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
	selectedSessionLaps []sessionLap = []sessionLap{}
	currentLapId int64 = 0
	currentPacketId int64 = 0

	playing bool = false
	forwardStep int64 = 10
	backwardStep int64 = 10
	defaultPlaybackSpeed int64 = 100
	playbackSpeed int64 = defaultPlaybackSpeed

	defaultPacketStepInterval time.Duration = time.Millisecond * 1
	nextPacketTimestamp time.Time = currentTime.Add(defaultPacketStepInterval)

	mapCirclePointers map[*canvas.Circle]*fyne.CanvasObject = make(map[*canvas.Circle]*fyne.CanvasObject)
	lastCircleOffset []float32 = []float32{} 

	screenSpaceWidth = mapWidth
	screenSpaceHeight = mapHeight

	viewSpaceMoveSpeed float32 = 10
	viewSpaceXOffset float32 = 0
	viewSpaceYOffset float32 = 0
	lastXOffset float32 = 0
	lastYOffset float32 = 0

	viewSpaceZoomFactor float32 = 1
	lastViewSpaceZoomFactor float32 = 0

	lockOn bool = false
	currentCarOffset []float32 = []float32{0, 0}
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
	graphLines = make(map[*fyne.Container][]*canvas.Line)

	objectOffsets = make(map[*fyne.CanvasObject][2]float32)
	containerPositions = make(map[*fyne.Container][2]float32)
	containerOffsets = make(map[*fyne.Container][2]float32)
	containerItems = make(map[*fyne.Container][]*fyne.CanvasObject)
	containerContainers = make(map[*fyne.Container][]*fyne.Container)

	mainContainer = container.NewWithoutLayout()
	containerPositions[mainContainer] = [2]float32{0, 0}
	containerOffsets[mainContainer] = [2]float32{0, 0}
	containerContainers[mainContainer] = []*fyne.Container{}

	mapContainer = container.NewWithoutLayout()
	addContainer(mainContainer, mapContainer, -windowOffset, -windowOffset)
	containerItems[mapContainer] = []*fyne.CanvasObject{}
	containerContainers[mapContainer] = []*fyne.Container{}

	var mapBackgroundColor = color.NRGBA{R: 33, G: 33, B: 33, A: 0}

	mapBackgroundP = canvas.NewRectangle(mapBackgroundColor)
	resizeObject(mapBackgroundP, mapWidth, mapHeight)
	addObject(mapContainer, &mapBackgroundP, 0, 0)

	mapTelemetryPoints = make(map[sessionLap][]*canvas.Circle)

	legendContainer = container.NewWithoutLayout()
	addContainer(mainContainer, legendContainer, resolutionWidth/2 + sideBarOffset - 145 -windowOffset, 80 -windowOffset)
	containerItems[legendContainer] = []*fyne.CanvasObject{}
	containerContainers[legendContainer] = []*fyne.Container{}

	var legendBackgroundColor = color.NRGBA{R: 25, G: 25, B: 25, A: 100}

	var legendBackgroundP fyne.CanvasObject = canvas.NewRectangle(legendBackgroundColor)
	resizeObject(legendBackgroundP, 150, 250)
	addObject(legendContainer, &legendBackgroundP, 0, 0)

	legendPoints = []*canvas.Circle{}
	legendLabels = []*widget.Label{}

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

	var playbackControlsBackgroundColor = color.NRGBA{R: 25, G: 25, B: 25, A: 255}

	playbackControlsContainer = container.NewWithoutLayout()
	addContainer(mainContainer, playbackControlsContainer, resolutionWidth/2 + sideBarOffset - windowOffset, resolutionHeight - bottomBarHeight - windowOffset)
	containerItems[playbackControlsContainer] = []*fyne.CanvasObject{}
	containerContainers[playbackControlsContainer] = []*fyne.Container{}

	var playbackControlsBackgroundP fyne.CanvasObject = canvas.NewRectangle(playbackControlsBackgroundColor)
	resizeObject(playbackControlsBackgroundP, resolutionWidth/2 - sideBarOffset, bottomBarHeight)
	addObject(playbackControlsContainer, &playbackControlsBackgroundP, 0, 0)

	var graphSelectionBorderColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	var graphSelectionBorderWidth float32 = 2
	var graphSelectionHeight float32 = 80
	var graphSelectionSizeOffset float32 = 40

	graphSelectionContainer = container.NewWithoutLayout()
	containerItems[graphSelectionContainer] = []*fyne.CanvasObject{}
	containerContainers[graphSelectionContainer] = []*fyne.Container{}
	addContainer(mainContainer, graphSelectionContainer, resolutionWidth/2 + sideBarOffset + graphSelectionSizeOffset/2 - windowOffset, resolutionHeight - bottomBarHeight + graphSelectionHeight/2 - windowOffset)

	var graphSelectionBackgroundColor = color.NRGBA{R: 25, G: 25, B: 25, A: 255}

	var graphSelectionBorder *canvas.Rectangle = canvas.NewRectangle(graphSelectionBackgroundColor)
	var graphSelectionBorderP fyne.CanvasObject = graphSelectionBorder
	graphSelectionBorder.StrokeColor = graphSelectionBorderColor
	graphSelectionBorder.StrokeWidth = graphSelectionBorderWidth
	resizeObject(graphSelectionBorderP, resolutionWidth/2 - sideBarOffset - graphSelectionSizeOffset, graphSelectionHeight)
	addObject(graphSelectionContainer, &graphSelectionBorderP, 0, 0)

	var checkWidthInitalOffset float32 = 40
	var checkHeightIniitalOffset float32 = 10
	var checkHeightSpacing float32 = 30

	var buttonWidth float32 = 200
	var buttonHeight float32 = 20

	var tireTempSelect *widget.Check = widget.NewCheck("Tire Temps", TireTempsCheckBoxChanged)
	var tireTempSelectP fyne.CanvasObject = tireTempSelect
	resizeObject(tireTempSelectP, buttonWidth, buttonHeight)
	addObject(graphSelectionContainer, &tireTempSelectP, checkWidthInitalOffset, checkHeightIniitalOffset)

	var tirePressureSelect *widget.Check = widget.NewCheck("Tire Pressures", TirePressuresCheckBoxChanged)
	var tirePressureSelectP fyne.CanvasObject = tirePressureSelect
	resizeObject(tirePressureSelectP, buttonWidth, buttonHeight)
	addObject(graphSelectionContainer, &tirePressureSelectP, checkWidthInitalOffset + 150, checkHeightIniitalOffset)

	var brakeSelect *widget.Check = widget.NewCheck("Brake", BrakeCheckBoxChanged)
	var brakeSelectP fyne.CanvasObject = brakeSelect
	resizeObject(brakeSelectP, buttonWidth, buttonHeight)
	addObject(graphSelectionContainer, &brakeSelectP, checkWidthInitalOffset + 300, checkHeightIniitalOffset)

	var throttleSelect *widget.Check = widget.NewCheck("Throttle", ThrottleCheckBoxChanged)
	var throttleSelectP fyne.CanvasObject = throttleSelect
	resizeObject(throttleSelectP, buttonWidth, buttonHeight)
	addObject(graphSelectionContainer, &throttleSelectP, checkWidthInitalOffset + 400, checkHeightIniitalOffset)

	var steeringSelect *widget.Check = widget.NewCheck("Steering", SteeringCheckBoxChanged)
	var steeringSelectP fyne.CanvasObject = steeringSelect
	resizeObject(steeringSelectP, buttonWidth, buttonHeight)
	addObject(graphSelectionContainer, &steeringSelectP, checkWidthInitalOffset + 520, checkHeightIniitalOffset)

	var accelerationSelect *widget.Check = widget.NewCheck("Acceleration", AccelerationCheckBoxChanged)
	var accelerationSelectP fyne.CanvasObject = accelerationSelect
	resizeObject(accelerationSelectP, buttonWidth, buttonHeight)
	addObject(graphSelectionContainer, &accelerationSelectP, checkWidthInitalOffset, checkHeightIniitalOffset + checkHeightSpacing)

	var speedSelect *widget.Check = widget.NewCheck("Speed", SpeedCheckBoxChanged)
	var speedSelectP fyne.CanvasObject = speedSelect
	resizeObject(speedSelectP, buttonWidth, buttonHeight)
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

	mainWindow.Canvas().SetContent(mainContainer)
	//endregion

	//region hotkeys
	mainWindow.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		var moveSpeed float32 = 25

		var containerHead = mainContainer
		var containerToMove = graphContainer

		var expandedGridWidth = screenSpaceWidth * viewSpaceZoomFactor
		var expandedGridHeight = screenSpaceHeight * viewSpaceZoomFactor

		var xOffset = viewSpaceXOffset * viewSpaceZoomFactor
		var yOffset = viewSpaceYOffset * viewSpaceZoomFactor

		switch key.Name {
			case fyne.KeyF:
				currentSessionLap := sessionLap{sessionId: currentSessionId, lapId: currentLapId}
				packetIds := packetIdsBySessionLap[currentSessionLap]
				maxPacketId := packetIds[1] 
				if currentPacketId + 3 > maxPacketId {
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
							currentPacketId -= 3
						}
					}
				}
				currentPacketId += 3
			case fyne.KeyK:
				lockOn = !lockOn
			case fyne.KeyUp:
				if yOffset - viewSpaceMoveSpeed * viewSpaceZoomFactor >= 0 {
					viewSpaceYOffset -= viewSpaceMoveSpeed
				}
			case fyne.KeyDown:
				if yOffset + viewSpaceMoveSpeed * viewSpaceZoomFactor <= expandedGridHeight - screenSpaceHeight {
					viewSpaceYOffset += viewSpaceMoveSpeed
				}
			case fyne.KeyLeft:
				if xOffset - viewSpaceMoveSpeed * viewSpaceZoomFactor >= 0 {
					viewSpaceXOffset -= viewSpaceMoveSpeed
				}
			case fyne.KeyRight:
				if xOffset + viewSpaceMoveSpeed * viewSpaceZoomFactor <= expandedGridWidth - screenSpaceWidth {
					viewSpaceXOffset += viewSpaceMoveSpeed
				}
			case fyne.KeyI:
				viewSpaceZoomFactor += 0.1
			case fyne.KeyO:
				if viewSpaceZoomFactor - 0.1 >= 1 {
					viewSpaceZoomFactor -= 0.1
				}
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
		nextPacketTimestamp = currentTime.Add(defaultPacketStepInterval * time.Duration(playbackSpeed))
		
		maxPacketId := packetIdsBySessionLap[sessionLap{sessionId: currentSessionId, lapId: currentLapId}][1]
		if currentPacketId + 1 <= maxPacketId {
			currentPacketId += 1
		}
	}

	checkLockOn()
	checkSelectedLaps()
	refreshLabels()
	refreshMaps()
	refreshGraphs()
	updateGraphToViewSpace()
}

func checkLockOn() {
	if lockOn {
		viewSpaceXOffset = (currentCarOffset[0] - screenSpaceWidth/2) / viewSpaceZoomFactor
		viewSpaceYOffset = (currentCarOffset[1] - screenSpaceHeight/2) / viewSpaceZoomFactor
	}
}

func checkSelectedLaps() {
	selectedSessionId, sessionErr := strconv.ParseInt(lapSelectSessionEntry.Text, 10, 64)
	if sessionErr != nil {
		return
	} 

	var lapInSelected bool = false
	
	sessionLapsSelected := []sessionLap{}
	for _, splitString := range strings.Split(lapSelectLapEntry.Text, ",") {
		lap, err := strconv.ParseInt(strings.ReplaceAll(strings.ReplaceAll(splitString, " ", ""), "	", ""), 10, 64)
		if err != nil {
			return
		}
		lapInSelected = lapInSelected || lap == currentLapId
		sessionLapsSelected = append(sessionLapsSelected, sessionLap{sessionId: selectedSessionId, lapId: lap})
	}

	if (currentSessionId != selectedSessionId || !lapInSelected) && len(sessionLapsSelected) > 0 {
		currentSessionId = selectedSessionId
		currentLapId = sessionLapsSelected[0].lapId
		currentPacketId = packetIdsBySessionLap[sessionLapsSelected[0]][0]
	}

	selectedSessionLaps = sessionLapsSelected
}

func refreshLabels() {
	currentInfoLabel.SetText(fmt.Sprintf("Session %d, Lap %d, Packet %d", currentSessionId, currentLapId, currentPacketId))
	lapSelectTotalSessions.SetText(fmt.Sprintf("%d sessions avaliable", totalSessions))
	lapSelectTotalLaps.SetText(fmt.Sprintf("%d laps avaliable", len(lapsBySession[currentSessionId])))
}

func refreshMaps() {
	lapsNeededToGenerate := []sessionLap{}
	lapsNeededToDestroy := []sessionLap{}

	if len(mapActiveSessionLaps) > 0 && mapActiveSessionLaps[0].sessionId == currentSessionId {
		for _, activeSessionLap := range mapActiveSessionLaps {
			if !slices.Contains(selectedSessionLaps, activeSessionLap) {
				lapsNeededToDestroy = append(lapsNeededToDestroy, activeSessionLap)
			}
		}
		for _, x := range selectedSessionLaps {
			if !slices.Contains(mapActiveSessionLaps, x) {
				bounds, inDict := packetIdsBySessionLap[x]
				if inDict && (bounds[0] - bounds[1]) != 0 {
					lapsNeededToGenerate = append(lapsNeededToGenerate, x)
				}
			}
		}
	} else {
		for _, x := range selectedSessionLaps {
			bounds, inDict := packetIdsBySessionLap[x]
			if inDict && (bounds[0] - bounds[1]) != 0 {
				lapsNeededToGenerate = append(lapsNeededToGenerate, x)
			}
		}
		lapsNeededToDestroy = mapActiveSessionLaps
	}

	for _, sessionLapsToDestroy := range lapsNeededToDestroy {
		circlesToDestroy, inDict := mapTelemetryPoints[sessionLapsToDestroy]
		if !inDict {
			fmt.Println("Circles Not found for given sessionLap")
			os.Exit(1)
		}

		for _, circle := range circlesToDestroy {
			circle.Hidden = true
			graphContainer.Remove(circle)
			delete(mapCirclePointers, circle)
		}

		indexInActiveSessionLaps := 0
		for index, x := range mapActiveSessionLaps {
			if x == sessionLapsToDestroy {
				indexInActiveSessionLaps = index
			}
		} 
		mapActiveSessionLaps[indexInActiveSessionLaps] = mapActiveSessionLaps[len(mapActiveSessionLaps) -1]
		mapActiveSessionLaps = mapActiveSessionLaps[:len(mapActiveSessionLaps) -1]
		
		delete(mapTelemetryPoints, sessionLapsToDestroy)
	}

	colorPicker := make(map[sessionLap]color.Color)

	for index, x := range selectedSessionLaps {
		var lapColor color.Color
		if index > len(lapColors) {
			lapColor = color.White
		} else {
			lapColor = lapColors[index]

		colorPicker[x] = lapColor
		}
	}

	if len(lapsNeededToGenerate) > 0 || len(lapsNeededToDestroy) > 0 {
		for _, x := range legendPoints {
			legendContainer.Remove(x)
		}

		for _, x := range legendLabels {
			legendContainer.Remove(x)
		}

		for index, x := range selectedSessionLaps {
			var newlegendPoint *canvas.Circle = canvas.NewCircle(colorPicker[x])
			var newlegendPointP fyne.CanvasObject = newlegendPoint
			resizeObject(newlegendPointP, legendPointSize, legendPointSize)
			addObject(legendContainer, &newlegendPointP, 20, float32(index) * 40 + 40)
			legendPoints = append(legendPoints, newlegendPoint)

			var newLegendLabel *widget.Label = widget.NewLabel(fmt.Sprintf("Lap %d", x.lapId))
			var newLegendLabelP fyne.CanvasObject = newLegendLabel
			resizeObject(newLegendLabelP, 80, 40)
			addObject(legendContainer, &newLegendLabelP, 40, float32(index) * 40 + 30)
			legendLabels = append(legendLabels, newLegendLabel)
		}
	}

	for _, sessionLapToGenerate := range lapsNeededToGenerate {
		packetsBounds := packetIdsBySessionLap[sessionLapToGenerate]
		packets := databaseAPI.QueryBetweenIdsFromPool(dbConnection, packetsBounds[0], packetsBounds[1])
		 
		lapCircles := []*canvas.Circle{}

		var minX float64 = math.MaxFloat32
		var minY float64 = math.MaxFloat32
		var maxX float64 = -math.MaxFloat32
		var maxY float64 = -math.MaxFloat32

		for _, packet := range *packets {
			location := packet.Location
			if location[0] < minX { minX = location[0] }
			if location[0] > maxX { maxX = location[0] }
			if location[1] < minY { minY = location[1] }
			if location[1] > maxY { maxY = location[1] }
		}

		var mapOffset float32 = 50
		var mapIncrease float32 = 1.75

		var nextPacketToDraw = 0

		for index, packet := range *packets {
			if index == nextPacketToDraw {
				nextPacketToDraw += packetInterval
				var newCircle *canvas.Circle = canvas.NewCircle(colorPicker[sessionLapToGenerate])
				newCircle.Resize(fyne.NewSize(1, 1))
				newCircle.StrokeWidth = 1
				var newCircleP fyne.CanvasObject = newCircle
				mapCirclePointers[newCircle] = &newCircleP
				resizeObject(newCircleP, defaultCircleSize, defaultCircleSize)
				var packetX float32 = float32((packet.Location[0] - minX) * (float64(mapWidth/2))/(maxX - minX))
				var packetY float32 = float32((packet.Location[1] - minY) * (float64(mapHeight/2))/(maxY - minY))
				addObject(mapContainer, &newCircleP, packetX * mapIncrease + mapOffset, packetY * mapIncrease + mapOffset)

				lapCircles = append(lapCircles, newCircle)
			}
		}

		mapTelemetryPoints[sessionLapToGenerate] = lapCircles
		if !slices.Contains(mapActiveSessionLaps, sessionLapToGenerate) {
			mapActiveSessionLaps = append(mapActiveSessionLaps, sessionLapToGenerate)
		}
	}
	
	if currentPacketId != mapLastSelectedPacketId {
		packetBounds := packetIdsBySessionLap[mapLastSelectedSessionLap]

		selectedIndex := int(mapLastSelectedPacketId - packetBounds[0])/packetInterval
		for index, circle := range mapTelemetryPoints[mapLastSelectedSessionLap] {
			if index == selectedIndex {
				circleP := mapCirclePointers[circle]
				resizeObject(circle, defaultCircleSize, defaultCircleSize)
				moveObject(mapContainer, circleP, lastCircleOffset[0], lastCircleOffset[1])
			}
		}

		currentSessionLap := sessionLap{sessionId: currentSessionId, lapId: currentLapId}

		packetBounds = packetIdsBySessionLap[currentSessionLap]

		selectedIndex = int(currentPacketId - packetBounds[0]/int64(packetInterval))
		for index, circle := range mapTelemetryPoints[currentSessionLap] {
			if index == selectedIndex {
				circleP := mapCirclePointers[circle]
				lastCircleOffset = []float32{objectOffsets[circleP][0], objectOffsets[circleP][1]}
				currentCarOffset = lastCircleOffset
				resizeObject(circle, defaultCircleSize * selectedCircleSizeIncrease, defaultCircleSize * selectedCircleSizeIncrease)
				moveObject(mapContainer, circleP, lastCircleOffset[0] - defaultCircleSize * selectedCircleSizeIncrease/2, lastCircleOffset[1] - defaultCircleSize * selectedCircleSizeIncrease/2)
			}
		}

		mapLastSelectedPacketId = currentPacketId
		mapLastSelectedSessionLap = sessionLap{sessionId: currentSessionId, lapId: currentLapId}
	}
}

func refreshGraphs() {
	packetBounds := packetIdsBySessionLap[sessionLap{sessionId: currentSessionId, lapId: currentLapId}]
	telemetryPackets := *databaseAPI.QueryBetweenIdsFromPool(dbConnection, packetBounds[0], packetBounds[1])

	var packetSpan int = 10
	var currentPacketIndex int = int(currentPacketId - packetBounds[0])

	var graphWidth float32 = 720
	var graphHeight float32 = 180

	var graphXSpacing float32 = graphWidth/(float32(packetSpan) * 2)

	for _, graph := range activeGraphs {
		var dataToDraw []float32
		var dataBounds [2]float32

		var tires bool = false
		var tireData [][4]float64

		switch graph{
			case tireTempsGraphContainer:
				tires = true
				for index, packet := range telemetryPackets {
					if index <= currentPacketIndex + packetSpan && index >= currentPacketIndex - packetSpan {
						tireData = append(tireData, packet.Tire_temps)
					}
				}
				dataBounds = [2]float32{20, 100}
			case tirePressuresGraphContainer:
				tires = true
				for index, packet := range telemetryPackets {
					if index <= currentPacketIndex + packetSpan && index >= currentPacketIndex - packetSpan {
						tireData = append(tireData, packet.Tire_pressures)
					}
				}
				dataBounds = [2]float32{20, 40}
			case brakeGraphContainer:
				for index, packet := range telemetryPackets {
					if index <= currentPacketIndex + packetSpan && index >= currentPacketIndex - packetSpan {
						dataToDraw = append(dataToDraw, float32(packet.Brake_input))
					}
				}
				dataBounds = [2]float32{0, 1}
			case throttleGraphContainer:
				for index, packet := range telemetryPackets {
					if index <= currentPacketIndex + packetSpan && index >= currentPacketIndex - packetSpan {
						dataToDraw = append(dataToDraw, float32(packet.Accelerator_input))
					}
				}
				dataBounds = [2]float32{0, 1}
			case steeringGraphContainer:
				for index, packet := range telemetryPackets {
					if index <= currentPacketIndex + packetSpan && index >= currentPacketIndex - packetSpan {
						dataToDraw = append(dataToDraw, float32(packet.Steering_angle))
					}
				}
				dataBounds = [2]float32{-450, 450}
			case accelerationGraphContainer:
				for index, packet := range telemetryPackets {
					if index <= currentPacketIndex + packetSpan && index >= currentPacketIndex - packetSpan {
						dataToDraw = append(dataToDraw, float32(packet.X_acceleration + packet.Z_acceleration))
					}
				}
				dataBounds = [2]float32{-0.1, 0.1}
			case speedGraphContainer:
				for index, packet := range telemetryPackets {
					if index <= currentPacketIndex + packetSpan && index >= currentPacketIndex - packetSpan {
						dataToDraw = append(dataToDraw, float32(packet.Velocity))
					}
				}
				dataBounds = [2]float32{0, 180}
		}		

		if !tires {
			linesNeeded := len(dataToDraw) -1
			existingLineCount := 0

			existingLines, inDict := graphLines[graph]
			if !inDict {
				graphLines[graph] = []*canvas.Line{}
			} else {
				existingLineCount = len(existingLines)
			}

			for i := 0; i < linesNeeded - existingLineCount; i++ {
				newLine := canvas.NewLine(color.RGBA{R: 255, G: 0, B: 0, A: 255})
				newLine.StrokeWidth = 2
				var newLineP fyne.CanvasObject = newLine
				addObject(graph, &newLineP, 0, 0)

				graphLines[graph] = append(graphLines[graph], newLine)
			}

			for i, line := range graphLines[graph] {
				if len(dataToDraw) > i + 1 {
					xI := graphXSpacing * float32(i)
					yI := graphHeight - (dataToDraw[i] - dataBounds[0]) * graphHeight/(dataBounds[1] - dataBounds[0])
					xF := graphXSpacing * float32(i + 1)
					yF := graphHeight - (dataToDraw[i + 1] - dataBounds[0]) * graphHeight/(dataBounds[1] - dataBounds[0])

					line.Position1.X = xI + containerOffsets[graph][0] + containerOffsets[graphContainer][0]
					line.Position1.Y = yI + containerOffsets[graph][1] + containerOffsets[graphContainer][1]
					line.Position2.X = xF + containerOffsets[graph][0] + containerOffsets[graphContainer][0]
					line.Position2.Y = yF + containerOffsets[graph][1] + containerOffsets[graphContainer][1]
				}
			}
		} else {
			linesNeeded := (len(tireData) -1) * 4
			existingLineCount := 0

			existingLines, inDict := graphLines[graph]
			if !inDict {
				graphLines[graph] = []*canvas.Line{}
			} else {
				existingLineCount = len(existingLines)
			}

			for i := 0; i < linesNeeded - existingLineCount; i++ {
				newLine := canvas.NewLine(color.RGBA{R: 255, G: 0, B: 0, A: 255})
				newLine.StrokeWidth = 2
				var newLineP fyne.CanvasObject = newLine
				addObject(graph, &newLineP, 0, 0)

				graphLines[graph] = append(graphLines[graph], newLine)
			}

			var lineIndex int = 0

			for index, _ := range tireData {
				if index != len(tireData) -1 {
					for tireIndex := 0; tireIndex < 4; tireIndex++ {
						line := graphLines[graph][lineIndex]
						line.StrokeColor = tireColors[tireIndex]
						xI := graphXSpacing * float32(index)
						yI := graphHeight - (float32(tireData[index][tireIndex]) - dataBounds[0]) * graphHeight/(dataBounds[1] - dataBounds[0])
						xF := graphXSpacing * float32(index + 1)
						yF := graphHeight - (float32(tireData[index + 1][tireIndex]) - dataBounds[0]) * graphHeight/(dataBounds[1] - dataBounds[0])

						line.Position1.X = xI + containerOffsets[graph][0] + containerOffsets[graphContainer][0]
						line.Position1.Y = yI + containerOffsets[graph][1] + containerOffsets[graphContainer][1]
						line.Position2.X = xF + containerOffsets[graph][0] + containerOffsets[graphContainer][0]
						line.Position2.Y = yF + containerOffsets[graph][1] + containerOffsets[graphContainer][1]
						lineIndex++
					} 
				}
			}
		}
	}
}

func updateGraphToViewSpace() {
	var xOffset float32 = -1 * viewSpaceXOffset * viewSpaceZoomFactor
	var yOffset float32 = -1 * viewSpaceYOffset * viewSpaceZoomFactor

	var changed = false

	if lastXOffset != viewSpaceXOffset || lastYOffset != viewSpaceYOffset {
		lastXOffset = viewSpaceXOffset
		lastYOffset = viewSpaceYOffset
		changed = true
	}

	if lastViewSpaceZoomFactor != viewSpaceZoomFactor {
		var scaleFactor float32 = viewSpaceZoomFactor/lastViewSpaceZoomFactor

		scaleGraph(scaleFactor)

		lastViewSpaceZoomFactor = viewSpaceZoomFactor
		changed = true
	}

	if changed {
		moveContainer(mainContainer, mapContainer, xOffset, yOffset)
	}
}

func scaleGraph(scaleFactor float32) {
	for _, points := range mapTelemetryPoints {
		for _, circle := range points {
			var circleP *fyne.CanvasObject = mapCirclePointers[circle]
			newCircleOffsetX := objectOffsets[circleP][0] * scaleFactor
			newCircleOffsetY := objectOffsets[circleP][1] * scaleFactor

			// resizeObject(*circleP, circle.Size().Width * scaleFactor, circle.Size().Height * scaleFactor)
			moveObject(mapContainer, circleP, newCircleOffsetX, newCircleOffsetY)
		}
	}

	mapBackgroundPOffsetX := objectOffsets[&mapBackgroundP][0] * scaleFactor
	mapBackgroundPOffsetY := objectOffsets[&mapBackgroundP][1] * scaleFactor
 
	// resizeObject(mapBackgroundP, mapBackgroundP.Size().Width * scaleFactor, mapBackgroundP.Size().Height * scaleFactor)
	moveObject(mapContainer, &mapBackgroundP, mapBackgroundPOffsetX, mapBackgroundPOffsetY)
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