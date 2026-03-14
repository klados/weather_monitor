#include "dht.h"
#include "esp_err.h"
#include "esp_log.h"
#include "esp_sleep.h"
#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include <cmath>

static constexpr const char *TAG = "weather_station";
static constexpr gpio_num_t DHT_GPIO = GPIO_NUM_18;
static constexpr dht_sensor_type_t DHT_SENSOR_TYPE = DHT_TYPE_AM2301;
static constexpr int MAX_READ_ATTEMPTS = 3;
static constexpr uint64_t SLEEP_TIME_US = 10ULL * 60ULL * 1000000ULL;

extern "C" void app_main(void) {
    const esp_sleep_wakeup_cause_t wakeup_cause = esp_sleep_get_wakeup_cause();
    if (wakeup_cause == ESP_SLEEP_WAKEUP_TIMER) {
        ESP_LOGI(TAG, "Wakeup from deep sleep timer");
    } else {
        ESP_LOGI(TAG, "Cold boot / reset");
    }

    // DHT22 benefits from a brief warm-up after power-up/wakeup.
    vTaskDelay(pdMS_TO_TICKS(2000));

    float temperature_c = 0.0f;
    float humidity_percent = 0.0f;
    bool valid_reading = false;

    for (int attempt = 1; attempt <= MAX_READ_ATTEMPTS; ++attempt) {
        const esp_err_t err = dht_read_float_data(
            DHT_SENSOR_TYPE,
            DHT_GPIO,
            &humidity_percent,
            &temperature_c
        );

        if (err == ESP_OK && !std::isnan(temperature_c) && !std::isnan(humidity_percent)) {
            valid_reading = true;
            break;
        }

        ESP_LOGW(
            TAG,
            "DHT22 read attempt %d/%d failed: %s",
            attempt,
            MAX_READ_ATTEMPTS,
            esp_err_to_name(err)
        );

        if (attempt < MAX_READ_ATTEMPTS) {
            vTaskDelay(pdMS_TO_TICKS(1000));
        }
    }

    if (valid_reading) {
        ESP_LOGI(
            TAG,
            "Temperature: %.1f C, Humidity: %.1f %%",
            temperature_c,
            humidity_percent
        );
    } else {
        ESP_LOGW(TAG, "Failed to read DHT22 after %d attempts", MAX_READ_ATTEMPTS);
    }

    ESP_LOGI(TAG, "Sleeping for %llu minutes", SLEEP_TIME_US / (60ULL * 1000000ULL));
    ESP_ERROR_CHECK(esp_sleep_enable_timer_wakeup(SLEEP_TIME_US));
    vTaskDelay(pdMS_TO_TICKS(100));
    esp_deep_sleep_start();
}