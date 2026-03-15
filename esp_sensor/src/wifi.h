#pragma once

#include <cstdbool>

#ifdef __cplusplus
extern "C" {
#endif

/**
 * Initialize WiFi and connect with retries.
 * Blocks until connected or max retries exceeded.
 *
 * @return true if connected, false if failed
 */
bool wifi_connect(void);

/**
 * Disconnect WiFi and deinitialize.
 * Call before deep sleep to reduce power.
 */
void wifi_disconnect(void);

/**
 * POST sensor data (temperature, humidity) to the configured server.
 *
 * @param temperature_c Temperature in Celsius
 * @param humidity_percent Humidity in percent
 * @return true if POST succeeded (2xx response), false otherwise
 */
bool wifi_post_sensor_data(float temperature_c, float humidity_percent);

#ifdef __cplusplus
}
#endif
