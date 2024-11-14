# FSAESoftware
Repo for the software found in Formula SAE car for LMU. Senior Project


# Telemetry App User Manual

## Table of Contents
- [Overview](#overview)
- [Getting Started](#getting-started)
  - [System Requirements](#system-requirements)
  - [Installation](#installation)
- [Using the Telemetry App](#using-the-telemetry-app)
  - [Main Interface](#main-interface)
  - [Real-Time Data View](#real-time-data-view)
  - [Historical Data Access](#historical-data-access)
  - [Alerts and Notifications](#alerts-and-notifications)
  - [Customizing the Dashboard](#customizing-the-dashboard)
- [Troubleshooting](#troubleshooting)
- [Support](#support)

---

## Overview
The Telemetry App is designed for use with the Formula SAE vehicle, allowing users to view real-time and historical telemetry data. This data includes critical metrics such as speed, suspension status, and GPS location, all aimed at optimizing vehicle performance.

## Getting Started

### System Requirements
- **Operating System:** Windows 10 or higher
- **Software Dependencies:**
  - PostgreSQL (for database)
  - Go programming language
  - Fyne.io framework (included with app installation)

### Installation
1. **Download the Telemetry App** Pull the code or download the zip file from the Github Repo.
2. **Launch the App:** Double-click the Telemetry App executable to open the application.

---

## Using the Telemetry App

### Main Interface
Once opened, the app displays the main dashboard, which provides real-time access to all telemetry data. Navigation options are available on the sidebar, including:
- **Dashboard:** Default view displaying core telemetry metrics.
- **Historical Data:** Access past telemetry logs.
- **Settings:** Customize the data display options.

### Real-Time Data View
- The main dashboard offers real-time data updates every second.
- Core metrics shown:
  - **Speed**
  - **Throttle and Brake Position**
  - **Steering Angle**
  - **Suspension Status**
  - **GPS Location**

Each data point updates as new telemetry data is received from the CCU.

### Historical Data Access
- Navigate to the **Historical Data** tab.
- Use the date and time filters to specify the time range for the data you want to view.
- Data displays as graphs and tables, showing trends and allowing you to analyze past performance.

### Alerts and Notifications
- **Real-Time Alerts:** The app will generate alerts if telemetry data crosses certain thresholds (e.g., high motor temperature, high tire temperature).
- Visual alerts appear on the dashboard, while audio notifications are triggered for critical alerts.
- **Alert History:** Navigate to the Alerts section to review past alerts and their timestamps.

### Customizing the Dashboard
- Open **Settings** from the sidebar to choose which data metrics to display on the dashboard.
- **Custom Layouts:** Arrange widgets to prioritize specific data, such as speed or suspension.
- **Themes and Units:** Customize display colors and units of measurement for each metric.

---

## Troubleshooting

- **Data Not Displaying:**
  - Ensure the CCU is connected and transmitting data.
  - Check the database connection in the Settings section.
- **Slow Performance:**
  - Reduce the number of metrics displayed on the dashboard.
  - Ensure the computer meets the minimum system requirements.
- **Connection Issues:**
  - Verify that WiFi connectivity is active between the app and CCU.
  - Restart the app and reconnect to the database.

---
