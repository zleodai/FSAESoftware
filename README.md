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
  - Minimum 8 GB RAM
  - 20 GB free disk space
- **Software:**
  - PostgreSQL (for database management)
  - Go programming language runtime
  - Fyne.io framework (included with the app)

## Installation
1. **Download the Installer:** Obtain the latest version of the Telemetry App installer from the official repository.
2. **Run the Installer:** Double-click the downloaded file and follow the on-screen instructions to complete the installation.
3. **Database Setup:** Ensure that PostgreSQL is installed and configured.

## Getting Started
1. **Launch the App:** Open the Telemetry App from the Start Menu or desktop shortcut.
2. **Initial Sync:** The app will synchronize with the vehicle's Central Control Unit (CCU) to fetch real-time data. Ensure the vehicle's telemetry system is active.

## Using the Telemetry App

### Main Dashboard
The main dashboard provides an overview of critical vehicle parameters:
- **Speed:** Current vehicle speed in km/h.
- **Throttle and Brake Position:** Percentage values indicating pedal positions.
- **Steering Angle:** Degree of steering input.
- **Tire Temperatures:** Real-time temperature readings for each tire.
- **Suspension Status:** Indicators showing suspension system performance.
- **GPS Location:** Current coordinates and track position.

### Real-Time Data Monitoring
- **Live Updates:** Data fields update every second to reflect current vehicle status.
- **Graphical Displays:** Visual graphs for parameters like speed and tire temperatures provide trend analysis.

### Historical Data Analysis
- **Access Logs:** Navigate to the 'Historical Data' section to view past telemetry records.
- **Filter Options:** Use date and time filters to narrow down specific sessions or events.
- **Export Data:** Export data logs in CSV format for external analysis.

### Alerts and Notifications
- **Real-Time Alerts:** The app generates alerts for parameters exceeding predefined thresholds (e.g., high tire temperature).
- **Notification Center:** Access a log of all alerts with timestamps and details.
- **Customization:** Set custom alert thresholds in the 'Settings' section.

### Customizing the Interface
- **Dashboard Layout:** Rearrange widgets to prioritize specific data points.
- **Themes:** Choose between light and dark modes for optimal visibility.
- **Units of Measurement:** Select preferred units (e.g., km/h or mph) in the 'Settings' menu.

## Troubleshooting
- **No Data Displayed:**
  - Verify that the vehicle's telemetry system is active and transmitting data.
  - Check the database connection settings in the app.
- **App Crashes or Freezes:**
  - Ensure your system meets the minimum requirements.
  - Restart the app and try again.
- **Incorrect Data Readings:**
  - Confirm that all sensors are calibrated and functioning correctly.
  - Check for any firmware updates for the CCU.
