#! /bin/bash

#Port on 3306
#Windows Service Name: MySQL80

mysql -u root -p -e "USE TelemetryDB; 

                    DROP TABLE PacketInfo;
                    DROP TABLE LapInfo;
                    DROP TABLE TelemetryInfo;
                    DROP TABLE TireInfo;

                    CREATE TABLE PacketInfo (
                        PacketID bigint,
                        SessionID int,
                        LapID int,
                        PacketDatetime datetime
                    );

                    CREATE TABLE LapInfo (
                        SessionID int,
                        LapID int,
                        LapTime bigint,
                        DriverName varchar(255),
                        TrackName varchar(255),
                        TrackConfiguration varchar(255),
                        CarName varchar(255)
                    );
                    
                    CREATE TABLE TelemetryInfo (
                        PacketID bigint,
                        SpeedMPH float(24),
                        Gas float(24),
                        Brake float(24),
                        Steer float(24),
                        Clutch float(24),
                        Gear tinyint,
                        RPM float(24),
                        TurboBoost float(24),
                        LocalAngularVelocityX float(24),
                        LocalAngularVelocityY float(24),
                        LocalAngularVelocityZ float(24),
                        VelocityX float(24),
                        VelocityY float(24),
                        VelocityZ float(24),
                        WorldPositionX float(24),
                        WorldPositionY float(24),
                        WorldPositionZ float(24),
                        Aero_DragCoeffcient float(24),
                        Aero_LiftCoefficientFront float(24),
                        Aero_LiftCoefficientRear float(24)
                    );
                    
                    CREATE TABLE TireInfo (
                        PacketID bigint,
                        FL_CamberRad float(24),
                        FR_CamberRad float(24),
                        RL_CamberRad float(24),
                        RR_CamberRad float(24),
                        FL_SlipAngle float(24),
                        FR_SlipAngle float(24),
                        RL_SlipAngle float(24),
                        RR_SlipAngle float(24),
                        FL_SlipRatio float(24),
                        FR_SlipRatio float(24),
                        RL_SlipRatio float(24),
                        RR_SlipRatio float(24),
                        FL_SelfAligningTorque float(24),
                        FR_SelfAligningTorque float(24),
                        RL_SelfAligningTorque float(24),
                        RR_SelfAligningTorque float(24),
                        FL_Load float(24),
                        FR_Load float(24),
                        RL_Load float(24),
                        RR_Load float(24),
                        FL_TyreSlip float(24),
                        FR_TyreSlip float(24),
                        RL_TyreSlip float(24),
                        RR_TyreSlip float(24),
                        FL_ThermalState float(24),
                        FR_ThermalState float(24),
                        RL_ThermalState float(24),
                        RR_ThermalState float(24),
                        FL_DynamicPressure float(24),
                        FR_DynamicPressure float(24),
                        RL_DynamicPressure float(24),
                        RR_DynamicPressure float(24),
                        FL_TyreDirtyLevel float(24),
                        FR_TyreDirtyLevel float(24),
                        RL_TyreDirtyLevel float(24),
                        RR_TyreDirtyLevel float(24)
                    );"