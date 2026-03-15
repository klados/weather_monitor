#include "dht.h"
#include "esp_err.h"
#include "esp_log.h"
#include "esp_sleep.h"
#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include "wifi.h"
#include <cmath>

static constexpr const char *TAG = "weather_station";
static constexpr gpio_num_t DHT_GPIO = GPIO_NUM_18;
static constexpr dht_sensor_type_t DHT_SENSOR_TYPE = DHT_TYPE_AM2301;
static constexpr int MAX_READ_ATTEMPTS = 3;
static constexpr uint64_t SLEEP_TIME_US = 10ULL * 60ULL * 1000000ULL;

extern "C" void app_main(void) {
    const esp_sleep_wakeup_cause_t wakeup_cause = esp_sleep_get_wakeup_cause();
    if (wakeup_cause == ESP_SLEEP_WAKEUP_TIMER) {
        esp_log_write(ESP_LOG_INFO, TAG, "Wakeup from deep sleep timer");
    } else {
        esp_log_write(ESP_LOG_INFO, TAG, "Cold boot / reset");
    }

    // DHT22 needs a brief warm-up after power-up/wakeup.
    vTaskDelay(pdMS_TO_TICKS(2000));

    float temperature_c = 0.0f;
    float humidity_percent = 0.0f;
    bool valid_reading = false;

    for (int attempt = 1; attempt <= MAX_READ_ATTEMPTS; ++attempt) {
        const esp_err_t err = dht_read_float_data(
            DHT_SENSOR_TYPE,
            static_cast<gpio_num_t>(DHT_GPIO),
            &humidity_percent,
            &temperature_c
        );

        if (err == ESP_OK && !std::isnan(temperature_c) && !std::isnan(humidity_percent)) {
            valid_reading = true;
            break;
        }

        esp_log_write(
        ESP_LOG_WARN,
            TAG,
            "DHT read attempt %d/%d failed: %s",
            attempt,
            MAX_READ_ATTEMPTS,
            esp_err_to_name(err)
        );
        vTaskDelay(pdMS_TO_TICKS(1000));
    }

    if (valid_reading) {
        ESP_LOGI(TAG, "Temperature: %.1f C, Humidity: %.1f %%", temperature_c, humidity_percent);

        if (wifi_connect()) {
            if (wifi_post_sensor_data(temperature_c, humidity_percent)) {
                ESP_LOGI(TAG, "Sensor data sent successfully");
            } else {
                ESP_LOGW(TAG, "Failed to send sensor data to server");
            }
            wifi_disconnect();
        } else {
            ESP_LOGW(TAG, "WiFi connection failed, skipping upload");
        }
    } else {
        ESP_LOGW(TAG, "Failed to read DHT sensor after %d attempts", MAX_READ_ATTEMPTS);
    }

    ESP_LOGI(TAG, "Sleeping for 10 minutes");
    ESP_ERROR_CHECK(esp_sleep_enable_timer_wakeup(SLEEP_TIME_US));
    vTaskDelay(pdMS_TO_TICKS(100));
    esp_deep_sleep_start();
}