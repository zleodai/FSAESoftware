CREATE TABLE telemetry (
    id                              int8,
<<<<<<< Updated upstream:Telemetry App/databaseAPI/SQLCreateDB.sql
    lap_id                          int8,
=======
>>>>>>> Stashed changes:Database/SQLCreateDB.sql
    date_entry                      date,
    time_step                       time,
    tire_temps                      float[4],
    tire_pressures                  float[4],
    velocity                        float,
    location                        float[2],
    accelerator_input               float,
    brake_input                     float,
    steering_angle                  float,
    gyro_pitch                      float,
    gyro_yaw                        float,
    gyro_roll                       float,
    x_acceleration                  float,
    y_acceleration                  float,
    z_acceleration                  float,
    total_power_draw                float,
    active_suspension_power_draw    float[4],
    motor_power_draw                float[2],
    battery_voltage                 float,
    traction_loss                   float,
    abs_throttle_limiting           float[2],
    limited_slip_usage              float[2]
);