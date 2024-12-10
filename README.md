# Telemetry App User Guide

## Table of Contents
- [Introduction](#introduction)
- [System Requirements](#system-requirements)
- [Installation](#installation)
- [Getting Started](#getting-started)
- [Using the Telemetry App](#using-the-telemetry-app)
  - [Main Dashboard](#main-dashboard)
  - [Real-Time Data Monitoring](#real-time-data-monitoring)
  - [Historical Data Analysis](#historical-data-analysis)
  - [Alerts and Notifications](#alerts-and-notifications)
  - [Customizing the Interface](#customizing-the-interface)
- [Troubleshooting](#troubleshooting)
- [Support and Feedback](#support-and-feedback)

---

## Introduction
The Telemetry App is designed to provide real-time and historical data monitoring for the Formula SAE vehicle. It offers insights into various vehicle parameters, aiding in performance optimization and diagnostics.

## System Requirements
- **Operating System:** Windows 10 or higher
- **Hardware:**
  - Minimum 4 GB RAM
- **Software:**
  - PostgreSQL (for database management)
  - Go programming language runtime
  - Fyne.io framework (included with the app)
  - Bash/Shell (included with git bash just move the bin folder to PATH in system variables)

## Installation/Setup
1. **Clone the repo** Clone the repo into a file location of your choice.
2. **Run the Installer:** Enter the directory Revised_Telemetry_App/databaseAPI and run StartDB.sh. This can be done with bash or sh via using ```sh Setup.sh``` or ```bash Setup.sh``` depending on which one is binded in your path variables
3. **Start the Database** Start the database by running ```sh StartDB.sh``` or ```bash StartDB.sh```
4. **Change Resolution** For the program ensure that your resolution is 1920x1080 or higher. Also if it is at 1920x1080 ensure that the scale setting in display settings is set to 100%.
5. **Run the Program** The program should be precompiled to windows and is located in Revised_Telemetry_App/main.exe.

## Using the Telemetry App

## App Controls
1. **Session/Lap Control** You will see options to enter in SessionID or LapID. It will also show you the Sessions avaliable and Laps avaliable for each entered Session ID. Here the IDs start from 0 so if it says 2 Sessions avaliable the sessionIDS you can enter are 0 and 1.
2. **Graph Viewer** You will also see checkboxes on different telemetry settings you can view. Simply check or uncheck the boxes to display them.
3. **Graph Scroll** To scroll through the graphs selected simply press "W" or "S" to move up and down respectively
4. **Map Controls** Use arrow keys to move the map (NOTE at the start the map will not move as you can only move when zoomed in) Use the keys "I" and "O" to zoom in and out respectively. Use the "K" key to lock in on the car (Will lock in when you start playing).
5. **Pause/Play** Simple press the Space key to play/pause. (NOTE DO NOT ZOOM IN OR OUT WHEN PLAYING. MAKE SURE TO PAUSE WHEN ZOOMING IN AND OUT). Press "A" or "D" to rewind or move forward in time respectively. Press "F" to super fast-forward.
6. **Switching To New Session/Lap** Make sure to zoom all the way out when switching to a new session. When changing the session/lap simply enter in a new sessionID or lapID. (NOTE Currently the view for the map can get thrown off when making a new map so just press K to lock in on the car when this happens)
   
## Troubleshooting
- **No Data Displayed:**
  - Enter in FSAE_AC_TelemetryCollector Folder and simply run the insertIntoDB.exe when the database is started.
- **Shell Scripts not working**
  - If the setup.sh failed to install the postgresql binaries -please ensure the lastest version of postgresql is installed and set to your path variables.
