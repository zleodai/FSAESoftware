import ac
import acsys
import os
from third_party.sim_info import *

#region GUI Globals
appName = "TelemetryCollector"
width, height = 335 , 400

buttonOffset = 150
buttonHeight = 85
buttonWidthOffset = 30
buttonWidth = 125
buttonWidthSpacing = 25
buttonFontSize = 20
buttonSizeIncrease = 1.1

labelOffset = 40
labelWidthOffset = 30
labelSpacing = 45
labelFontSize = 25
#endregion

#region Timed Globals
buttonSizeIncreaseDuration = 0.1

startButtonPressed = False
startButtonTimestamp = 1
endButtonPressed = False
endButtonTimestamp = 1
exportButtonPressed = False
exportButtonTimestamp = 1
startLapsButtonPressed = False
startLapsButtonTimestamp = 1

currentTime = 0
lapEndCooldownTimestamp = 0

insertCooldown = 0.1
nextInsertTimestamp = 0.15
#endregion

#region Telemetry Globals
lastVelocity = [0, 0, 0]
acceleration = [0, 0, 0]
currentLap = 0
lapTimes = {}

lapStartLocation = [0, 0, 0]
lapStartSet = False
lapEndLocation = [0, 0, 0]
lapEndSet = False

lapStarted = False
distanceForNewLap = 30
timeForNewLap = 10

exportString = ""
#endregion Global Vars

simInfo = SimInfo()

def acMain(ac_version):
    global appWindow # <- you'll need to update your window in other functions.

    appWindow = ac.newApp(appName)
    ac.setTitle(appWindow, appName)
    ac.setSize(appWindow, width, height)

    ac.addRenderCallback(appWindow, appGL) # -> links this app's window to an OpenGL render function

    assignLabels()
    assignButtons()

    return appName

def assignLabels():
    global LapNumberLabel
    LapNumberLabel = ac.addLabel(appWindow, "Lap ")
    ac.setPosition(LapNumberLabel, labelWidthOffset, labelOffset)
    ac.setFontSize(LapNumberLabel, labelFontSize)

    global LapTimeLabel
    LapTimeLabel = ac.addLabel(appWindow, "")
    ac.setPosition(LapTimeLabel, labelWidthOffset, labelOffset + labelSpacing)
    ac.setFontSize(LapTimeLabel, labelFontSize)

def assignButtons():
    global SetLapStartButton
    SetLapStartButton = ac.addButton(appWindow, "\nSet Start")
    ac.setPosition(SetLapStartButton, buttonWidthOffset, buttonOffset)
    ac.setSize(SetLapStartButton, buttonWidth, buttonHeight)
    ac.setFontSize(SetLapStartButton, buttonFontSize)
    ac.setFontAlignment(SetLapStartButton, "center")
    ac.setFontAlignment
    ac.addOnClickedListener(SetLapStartButton, setLapStartButtonPress)

    global SetLapEndButton
    SetLapEndButton = ac.addButton(appWindow, "\nSet End")
    ac.setPosition(SetLapEndButton, buttonWidthOffset + buttonWidth + buttonWidthSpacing, buttonOffset)
    ac.setSize(SetLapEndButton, buttonWidth, buttonHeight)
    ac.setFontSize(SetLapEndButton, buttonFontSize)
    ac.setFontAlignment(SetLapEndButton, "center")
    ac.addOnClickedListener(SetLapEndButton, setLapEndButtonPress)

    global ExportToDBButton
    ExportToDBButton = ac.addButton(appWindow, "Export To DB")
    ac.setPosition(ExportToDBButton, buttonWidthOffset+(buttonWidth/2)+(buttonWidthSpacing/2), buttonOffset + 100)
    ac.setSize(ExportToDBButton, buttonWidth, 25)
    ac.setFontSize(ExportToDBButton, 16)
    ac.setFontAlignment(ExportToDBButton, "center")
    ac.addOnClickedListener(ExportToDBButton, exportToDB)

    global StartLapsButton
    StartLapsButton = ac.addButton(appWindow, "Start")
    ac.setPosition(StartLapsButton, buttonWidthOffset+(buttonWidth/2)+(buttonWidthSpacing/2), buttonOffset + 175)
    ac.setSize(StartLapsButton, buttonWidth, 40)
    ac.setFontSize(StartLapsButton, 24)
    ac.setFontAlignment(StartLapsButton, "center")
    ac.addOnClickedListener(StartLapsButton, startLaps)

def appGL(deltaT):#-------------------------------- OpenGL UPDATE
    """
    This is where you redraw your openGL graphics
    if you need to use them .
    """
    pass # -> Delete this line if you do something here !

def updateLabels():
    ac.setText(LapNumberLabel, "Lap {0}".format(currentLap))
    ac.setText(LapTimeLabel, "{0:.3f}".format(lapTimes[currentLap]))

def updatePacketString():
    global exportString

    tireTemps = ac.getCarState(0, acsys.CS.CurrentTyresCoreTemp)
    tirePressure = ac.getCarState(0, acsys.CS.DynamicPressure)
    location = ac.getCarState(0, acsys.CS.WorldPosition)
    
    exportString += "\n{0} {1} {2} {3} {4} {5} {6} {7} {8} {9} {10} {11} {12} {13} {14} {15} {16} {17} {18} {19} {20}".format(
        currentLap, #0
        tireTemps[0], #1-4
        tireTemps[1],
        tireTemps[2],
        tireTemps[3],
        tirePressure[0], #5-8
        tirePressure[1],
        tirePressure[2],
        tirePressure[3],
        ac.getCarState(0, acsys.CS.SpeedMPH), #9
        location[0], #10-11
        location[2],
        ac.getCarState(0, acsys.CS.Gas), #12
        ac.getCarState(0, acsys.CS.Brake), 
        ac.getCarState(0, acsys.CS.Steer), 
        0, #15
        0, 
        0, 
        acceleration[0], #18
        acceleration[1], 
        acceleration[2],
    )

def setLapStartButtonPress(x, y):
    ac.setPosition(SetLapStartButton, buttonWidthOffset - (buttonWidth * buttonSizeIncrease - buttonWidth)/2, buttonOffset - (buttonHeight * buttonSizeIncrease - buttonHeight)/2)
    ac.setSize(SetLapStartButton, buttonWidth * buttonSizeIncrease, buttonHeight * buttonSizeIncrease)

    global startButtonTimestamp
    global startButtonPressed
    startButtonTimestamp = currentTime + buttonSizeIncreaseDuration
    startButtonPressed = True

    global lapStartLocation
    global lapStartSet

    lapStartLocation = ac.getCarState(0, acsys.CS.WorldPosition)
    lapStartSet = True

def setLapEndButtonPress(x, y):
    ac.setPosition(SetLapEndButton, buttonWidthOffset - (buttonWidth * buttonSizeIncrease - buttonWidth)/2 + buttonWidth + buttonWidthSpacing, buttonOffset - (buttonHeight * buttonSizeIncrease - buttonHeight)/2)
    ac.setSize(SetLapEndButton, buttonWidth * buttonSizeIncrease, buttonHeight * buttonSizeIncrease)

    global endButtonTimestamp
    global endButtonPressed
    endButtonTimestamp = currentTime + buttonSizeIncreaseDuration
    endButtonPressed = True

    global lapEndLocation
    global lapEndSet

    lapEndLocation = ac.getCarState(0, acsys.CS.WorldPosition)
    lapEndSet = True

def exportToDB(x, y):
    ac.setPosition(ExportToDBButton, buttonWidthOffset+(buttonWidth/2)+(buttonWidthSpacing/2) - (buttonWidth * buttonSizeIncrease - buttonWidth)/2, buttonOffset + 100 - (25 * buttonSizeIncrease - 25)/2)
    ac.setSize(ExportToDBButton, buttonWidth * buttonSizeIncrease, 25 * buttonSizeIncrease)

    global exportButtonTimestamp
    global exportButtonPressed
    exportButtonTimestamp = currentTime + buttonSizeIncreaseDuration
    exportButtonPressed = True

    global exportString

    f = open("telemetryData.txt", "w")
    f.write(exportString)
    f.close()

    os.system("insertIntoDB.exe")
    exportString = ""

def startLaps(x, y):
    ac.setPosition(StartLapsButton, buttonWidthOffset+(buttonWidth/2)+(buttonWidthSpacing/2) - (buttonWidth * buttonSizeIncrease - buttonWidth)/2, buttonOffset + 175 - (40 * buttonSizeIncrease - 40)/2)
    ac.setSize(StartLapsButton, buttonWidth * buttonSizeIncrease, 40 * buttonSizeIncrease)

    global startLapsButtonPressed
    global startLapsButtonTimestamp
    startLapsButtonTimestamp = currentTime + buttonSizeIncreaseDuration
    startLapsButtonPressed = True
    
    global lapStarted
    global currentLap
    lapStarted = True
    currentLap += 1

def updateButtons():
    global startButtonPressed
    global endButtonPressed
    global exportButtonPressed
    global startLapsButtonPressed

    if startButtonPressed and currentTime > startButtonTimestamp:
        ac.setPosition(SetLapStartButton, buttonWidthOffset, buttonOffset)
        ac.setSize(SetLapStartButton, buttonWidth, buttonHeight)
        startButtonPressed = False
    if endButtonPressed and currentTime > endButtonTimestamp:
        ac.setPosition(SetLapEndButton, buttonWidthOffset + buttonWidth + buttonWidthSpacing, buttonOffset)
        ac.setSize(SetLapEndButton, buttonWidth, buttonHeight)
        endButtonPressed = False
    if exportButtonPressed and currentTime > exportButtonTimestamp:
        ac.setPosition(ExportToDBButton, buttonWidthOffset+(buttonWidth/2)+(buttonWidthSpacing/2), buttonOffset + 100)
        ac.setSize(ExportToDBButton, buttonWidth, 25)
        exportButtonPressed = False
    if startLapsButtonPressed and currentTime > startLapsButtonTimestamp:
        ac.setPosition(StartLapsButton, buttonWidthOffset+(buttonWidth/2)+(buttonWidthSpacing/2), buttonOffset + 175)
        ac.setSize(StartLapsButton, buttonWidth, 40)
        startLapsButtonPressed = False

def acUpdate(deltaT):#-------------------------------- AC UPDATE
    """
    This is where you update your app window ( != OpenGL graphics )
    such as : labels , listener , ect ...
    """
    global currentTime
    global nextInsertTimestamp
    global lastVelocity
    global acceleration
    global currentLap

    if currentLap in lapTimes.keys():
        lapTimes[currentLap] += deltaT
    else:
        lapTimes[currentLap] = deltaT

    currentTime += deltaT

    updateLabels()
    updateButtons()

    currentVelocity = ac.getCarState(0, acsys.CS.LocalVelocity)
    acceleration = [(currentVelocity[0] - lastVelocity[0])/deltaT/1000, (currentVelocity[1] - lastVelocity[1])/deltaT/1000, (currentVelocity[2] - lastVelocity[2])/deltaT/1000]
    lastVelocity = currentVelocity

    if currentTime > nextInsertTimestamp:
        updatePacketString()
        nextInsertTimestamp = currentTime + insertCooldown

    if lapStarted and lapStartSet and lapEndSet:
        location = ac.getCarState(0, acsys.CS.WorldPosition)
        distance = abs(lapEndLocation[0] - location[0]) + abs(lapEndLocation[1] - location[1]) + abs(lapEndLocation[2] - location[2])
        if distance < distanceForNewLap and lapTimes[currentLap] > timeForNewLap:
            currentLap += 1