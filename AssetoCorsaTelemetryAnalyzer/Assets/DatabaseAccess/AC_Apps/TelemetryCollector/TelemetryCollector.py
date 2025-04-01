import urllib.parse
import urllib.request
import ac
import acsys
import os
import requests
import hashlib
from third_party.sim_info import *

#AC Docs https://docs.google.com/document/d/13trBp6K1TjWbToUQs_nfFsB291-zVJzRZCNaTYt4Dzc/pub

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
nextInsertTimestamp = insertCooldown
#endregion

#region Telemetry Globals
lastVelocity = [0, 0, 0]
acceleration = [0, 0, 0]
currentSession = 0
currentLap = 0
lapTimes = {}
key = "945"

lapStartLocation = [0, 0, 0]
lapStartSet = False
lapEndLocation = [0, 0, 0]
lapEndSet = False

lapStarted = False
distanceForNewLap = 30
timeForNewLap = 10

insertQueries = []
#endregion Global Vars

session = requests.Session()

simInfo = SimInfo()

def acMain(ac_version):
    global appWindow # <- you'll need to update your window in other functions.

    appWindow = ac.newApp(appName)
    ac.setTitle(appWindow, appName)
    ac.setSize(appWindow, width, height)

    ac.addRenderCallback(appWindow, appGL) # -> links this app's window to an OpenGL render function

    assignLabels()
    assignButtons()

    hash = hashlib.sha256()
    hash.update(str.encode(key))
    publicKey = hash.hexdigest() 

    global currentSession
    request = urllib.request.urlopen("http://localhost:5432/getCurrentSession.php?publicKey={0}".format(publicKey)).read()
    currentSession = int(bytes.decode(request))

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

def sendSessionLapPacket(carID):
    insertType = 2
    hash = hashlib.sha256()
    hash.update(str.encode(str(insertType * int(key) * (currentSession + currentLap))))
    publicKey = hash.hexdigest()
    trackConfig = ac.getTrackConfiguration(carID) if ac.getTrackConfiguration(carID) != "" else "Default"
    urllib.request.urlopen("http://localhost:5432/insertIntoDatabase.php?publicKey={0}&insertType={1}&SessionID={2}&LapID={3}&LapTime={4}&DriverName={5}&TrackName={6}&TrackConfiguration={7}&CarName={8}".format(publicKey, insertType, currentSession, currentLap, int(lapTimes[currentLap]*1000), ac.getDriverName(carID), ac.getTrackName(carID), trackConfig, ac.getCarName(carID)))

def recordTelemetryPacket(carID):
    insertType = 1
    hash = hashlib.sha256()
    hash.update(str.encode(str(insertType * int(key) * (currentSession + currentLap))))
    publicKey = hash.hexdigest()
    
    localAngularVelocity = ac.getCarState(carID, acsys.CS.LocalAngularVelocity)
    velocity = ac.getCarState(carID, acsys.CS.Velocity)
    worldPosition = ac.getCarState(carID, acsys.CS.WorldPosition)
    camberRad = ac.getCarState(carID, acsys.CS.CamberRad)
    slipAngle = ac.getCarState(carID, acsys.CS.SlipAngle)
    slipRatio = ac.getCarState(carID, acsys.CS.SlipRatio)
    selfAligningTorque = ac.getCarState(carID, acsys.CS.Mz)
    load = ac.getCarState(carID, acsys.CS.Load)
    tyreSlip = ac.getCarState(carID, acsys.CS.TyreSlip)
    thermalState = ac.getCarState(carID, acsys.CS.CurrentTyresCoreTemp)
    dynamicPressure = ac.getCarState(carID, acsys.CS.DynamicPressure)
    tyreDirtyLevel = ac.getCarState(carID, acsys.CS.TyreDirtyLevel)
    
    insertString = "http://localhost:5432/insertIntoDatabase.php?publicKey={0}&insertType={1}&SessionID={2}&LapID={3}&SpeedMPH={4}&Gas={5}&Brake={6}&Steer={7}&Clutch={8}&Gear={9}&RPM={10}&TurboBoost={11}&LocalAngularVelocityX={12}&LocalAngularVelocityY={13}&LocalAngularVelocityZ={14}&VelocityX={15}&VelocityY={16}&VelocityZ={17}&WorldPositionX={18}&WorldPositionY={19}&WorldPositionZ={20}&Aero_DragCoeffcient={21}&Aero_LiftCoefficientFront={22}&Aero_LiftCoefficientRear={23}&FL_CamberRad={24}&FR_CamberRad={25}&RL_CamberRad={26}&RR_CamberRad={27}&FL_SlipAngle={28}&FR_SlipAngle={29}&RL_SlipAngle={30}&RR_SlipAngle={31}&FL_SlipRatio={32}&FR_SlipRatio={33}&RL_SlipRatio={34}&RR_SlipRatio={35}&FL_SelfAligningTorque={36}&FR_SelfAligningTorque={37}&RL_SelfAligningTorque={38}&RR_SelfAligningTorque={39}&FL_Load={40}&FR_Load={41}&RL_Load={42}&RR_Load={43}&FL_TyreSlip={44}&FR_TyreSlip={45}&RL_TyreSlip={46}&RR_TyreSlip={47}&FL_ThermalState={48}&FR_ThermalState={49}&RL_ThermalState={50}&RR_ThermalState={51}&FL_DynamicPressure={52}&FR_DynamicPressure={53}&RL_DynamicPressure={54}&RR_DynamicPressure={55}&FL_TyreDirtyLevel={56}&FR_TyreDirtyLevel={57}&RL_TyreDirtyLevel={58}&RR_TyreDirtyLevel={59}".format(publicKey, insertType, currentSession, currentLap,ac.getCarState(carID, acsys.CS.SpeedMPH),ac.getCarState(carID, acsys.CS.Gas),ac.getCarState(carID, acsys.CS.Brake),ac.getCarState(carID, acsys.CS.Steer),ac.getCarState(carID, acsys.CS.Clutch),ac.getCarState(carID, acsys.CS.Gear),ac.getCarState(carID, acsys.CS.RPM),ac.getCarState(carID, acsys.CS.TurboBoost),localAngularVelocity[0],localAngularVelocity[1],localAngularVelocity[2],velocity[0],velocity[1],velocity[2],worldPosition[0],worldPosition[1],worldPosition[2],ac.getCarState(carID, acsys.CS.Aero),ac.getCarState(carID, acsys.AERO.CL_Front),ac.getCarState(carID, acsys.AERO.CL_Rear),camberRad[0],camberRad[1],camberRad[2],camberRad[3],slipAngle[0],slipAngle[1],slipAngle[2],slipAngle[3],slipRatio[0],slipRatio[1],slipRatio[2],slipRatio[3],selfAligningTorque[0],selfAligningTorque[1],selfAligningTorque[2],selfAligningTorque[3],load[0],load[1],load[2],load[3],tyreSlip[0],tyreSlip[1],tyreSlip[2],tyreSlip[3],thermalState[0],thermalState[1],thermalState[2],thermalState[3],dynamicPressure[0],dynamicPressure[1],dynamicPressure[2],dynamicPressure[3],tyreDirtyLevel[0],tyreDirtyLevel[1],tyreDirtyLevel[2],tyreDirtyLevel[3])
    session.put(insertString)
    # global insertQueries
    # insertQueries.append(insertString)
    # ac.log("insertQueries has {0} items".format(len(insertQueries)))

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

def startLaps(x, y):
    ac.setPosition(StartLapsButton, buttonWidthOffset+(buttonWidth/2)+(buttonWidthSpacing/2) - (buttonWidth * buttonSizeIncrease - buttonWidth)/2, buttonOffset + 100 - (40 * buttonSizeIncrease - 40)/2)
    ac.setSize(StartLapsButton, buttonWidth * buttonSizeIncrease, 40 * buttonSizeIncrease)

    global startLapsButtonPressed
    global startLapsButtonTimestamp
    startLapsButtonTimestamp = currentTime + buttonSizeIncreaseDuration
    startLapsButtonPressed = True
    
    global lapStarted
    global currentLap
    lapStarted = True
    currentLap += 1

def exportToDB(x, y):
    ac.setPosition(ExportToDBButton, buttonWidthOffset+(buttonWidth/2)+(buttonWidthSpacing/2) - (buttonWidth * buttonSizeIncrease - buttonWidth)/2, buttonOffset + 100 - (25 * buttonSizeIncrease - 25)/2)
    ac.setSize(ExportToDBButton, buttonWidth * buttonSizeIncrease, 25 * buttonSizeIncrease)

    global exportButtonTimestamp
    global exportButtonPressed
    exportButtonTimestamp = currentTime + buttonSizeIncreaseDuration
    exportButtonPressed = True

    global insertQueries
    for insertQuery in insertQueries:
        urllib.request.urlopen(insertQuery)
    
    insertQueries = []

def updateLabels():
    ac.setText(LapNumberLabel, "Lap {0}".format(currentLap))
    ac.setText(LapTimeLabel, "{0:.3f}".format(lapTimes[currentLap]))

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

    if lapStarted and currentTime > nextInsertTimestamp:
        recordTelemetryPacket(0)
        nextInsertTimestamp = currentTime + insertCooldown

    if lapStarted and lapStartSet and lapEndSet:
        location = ac.getCarState(0, acsys.CS.WorldPosition)
        distance = abs(lapEndLocation[0] - location[0]) + abs(lapEndLocation[1] - location[1]) + abs(lapEndLocation[2] - location[2])
        if distance < distanceForNewLap and lapTimes[currentLap] > timeForNewLap:
            sendSessionLapPacket(0)
            currentLap += 1