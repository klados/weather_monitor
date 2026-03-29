#include <ctime>
#include "esp_sntp.h"
#include "esp_log.h"

static constexpr const char *TAG = "time";

bool sync_time_with_sntp() {
    sntp_setoperatingmode(SNTP_OPMODE_POLL);
    sntp_setservername(0, "pool.ntp.org");
    sntp_init();

    // Wait for time to be set
    std::time_t now = 0;
    std::tm timeinfo = {};
    int retry = 0;
    const int retry_count = 30;

    while (retry < retry_count) {
        std::time(&now);
        localtime_r(&now, &timeinfo);
        if (timeinfo.tm_year >= (2020 - 1900)) {
            ESP_LOGI(TAG, "Time synchronized: %lld", static_cast<long long>(now));
            return true;
        }
        vTaskDelay(pdMS_TO_TICKS(1000));
        ++retry;
    }

    ESP_LOGW(TAG, "SNTP sync timed out");
    return false;
}