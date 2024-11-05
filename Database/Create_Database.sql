CREATE TABLE telemetry (
    date                            date,
    time_step                       time,
    tire_temps                      float[4][1][1],
    tire_pressures                  float[4],
    velocity                        NUMERIC(12, 6),
    location                        NUMERIC(12, 6)[2],
    accelerator_input               NUMERIC(12, 8),
    brake_input                     NUMERIC(12, 8),
    steering_angle                  NUMERIC(12, 6),
    gyro_pitch                      NUMERIC(12, 6),
    gyro_yaw                        NUMERIC(12, 6),
    gyro_roll                       NUMERIC(12, 6),
    x_acceleration                  NUMERIC(12, 6),
    y_acceleration                  NUMERIC(12, 6),
    z_acceleration                  NUMERIC(12, 6),
    total_power_draw                float,
    active_suspension_power_draw    float[4],
    motor_power_draw                float[2],
    battery_voltage                 float,
    traction_lost                   boolean,
    abs_throttle_limiting           float[2],
    limited_slip_usage              float[2],
    log                             varchar[100]
);

-- INSERT INTO telemetry (date, time_step, tire_temps, tire_pressures, velocity, location, accelerator_input, brake_input, steering_angle, gyro_pitch, gyro_yaw, gyro_roll, x_acceleration, y_acceleration, z_acceleration, total_power_draw, active_suspension_power_draw, motor_power_draw, battery_voltage, traction_lost, abs_throttle_limiting, limited_slip_usage)

-- INSERT INTO telemetry (date, time_step, x_acceleration, y_acceleration, z_acceleration)
--     VALUES (CURRENT_DATE::date, NOW()::time, 0, 0, 0);

-- SELECT date, time_step, x_acceleration, y_acceleration, z_acceleration FROM telemetry;