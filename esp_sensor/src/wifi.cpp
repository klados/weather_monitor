#include "wifi.h"
#include "esp_err.h"
#include "esp_event.h"
#include "esp_http_client.h"
#include "esp_log.h"
#include "esp_wifi.h"
#include "freertos/FreeRTOS.h"
#include "freertos/event_groups.h"
#include "freertos/task.h"
#include "nvs_flash.h"
#include <cstdio>
#include <cstring>
#include <ctime>

#include "mbedtls/md.h"
#include "sdkconfig.h"

static constexpr const char *TAG = "wifi";

#define WIFI_CONNECTED_BIT BIT0
#define WIFI_FAIL_BIT      BIT1

static EventGroupHandle_t s_wifi_event_group;
static int s_retry_num = 0;

static bool compute_hmac_sha256_hex(const char *key,
                                    const char *device_id,
                                    const char *timestamp,
                                    const uint8_t *body,
                                    size_t body_len,
                                    char *out_hex,
                                    size_t out_hex_size) {

    if (!key || !device_id || !timestamp || !body || !out_hex || out_hex_size < 65) {
        return false;
    }

    const mbedtls_md_info_t *md_info = mbedtls_md_info_from_type(MBEDTLS_MD_SHA256);
    if (!md_info) {
        return false;
    }

    unsigned char hmac[32];
    mbedtls_md_context_t ctx;
    mbedtls_md_init(&ctx);

    const int rc = mbedtls_md_setup(&ctx, md_info, 1);
    if (rc != 0) {
        mbedtls_md_free(&ctx);
        return false;
    }

    if (mbedtls_md_hmac_starts(&ctx,reinterpret_cast<const unsigned char *>(key),std::strlen(key)) != 0 ||
        mbedtls_md_hmac_update(&ctx,reinterpret_cast<const unsigned char *>(device_id),std::strlen(device_id)) != 0 ||
        mbedtls_md_hmac_update(&ctx,reinterpret_cast<const unsigned char *>(timestamp),std::strlen(timestamp)) != 0 ||
        mbedtls_md_hmac_update(&ctx, body, body_len) != 0 ||
        mbedtls_md_hmac_finish(&ctx, hmac) != 0) {

            mbedtls_md_free(&ctx);
            return false;
        }


    mbedtls_md_free(&ctx);

    for (size_t i = 0; i < sizeof(hmac); ++i) {
        std::snprintf(out_hex + (i * 2), out_hex_size - (i * 2), "%02x", hmac[i]);
    }
    out_hex[64] = '\0';
    return true;
}

static void wifi_event_handler(void *arg, esp_event_base_t event_base,
                               int32_t event_id, void *event_data) {
    if (event_base == WIFI_EVENT && event_id == WIFI_EVENT_STA_START) {
        esp_wifi_connect();
    } else if (event_base == WIFI_EVENT && event_id == WIFI_EVENT_STA_DISCONNECTED) {
        if (s_retry_num < CONFIG_WIFI_MAX_RETRIES) {
            esp_wifi_connect();
            s_retry_num++;
            ESP_LOGW(TAG, "Retry connection to AP (%d/%d)", s_retry_num, CONFIG_WIFI_MAX_RETRIES);
        } else {
            xEventGroupSetBits(s_wifi_event_group, WIFI_FAIL_BIT);
        }
        ESP_LOGW(TAG, "Connection attempt failed");
    } else if (event_base == IP_EVENT && event_id == IP_EVENT_STA_GOT_IP) {
        const auto *event = static_cast<ip_event_got_ip_t *>(event_data);
        ESP_LOGI(TAG, "Got IP: " IPSTR, IP2STR(&event->ip_info.ip));
        s_retry_num = 0;
        xEventGroupSetBits(s_wifi_event_group, WIFI_CONNECTED_BIT);
    }
}

static esp_err_t http_event_handler(esp_http_client_event_t *evt) {
    (void)evt;
    return ESP_OK;
}

bool wifi_connect(void) {
    esp_err_t ret = nvs_flash_init();
    if (ret == ESP_ERR_NVS_NO_FREE_PAGES || ret == ESP_ERR_NVS_NEW_VERSION_FOUND) {
        ESP_ERROR_CHECK(nvs_flash_erase());
        ESP_ERROR_CHECK(nvs_flash_init());
    } else {
        ESP_ERROR_CHECK(ret);
    }

    s_wifi_event_group = xEventGroupCreate();

    ESP_ERROR_CHECK(esp_netif_init());
    ESP_ERROR_CHECK(esp_event_loop_create_default());
    esp_netif_create_default_wifi_sta();

    wifi_init_config_t cfg = WIFI_INIT_CONFIG_DEFAULT();
    ESP_ERROR_CHECK(esp_wifi_init(&cfg));

    esp_event_handler_instance_t instance_any_id;
    esp_event_handler_instance_t instance_got_ip;
    ESP_ERROR_CHECK(esp_event_handler_instance_register(WIFI_EVENT,
                                                        ESP_EVENT_ANY_ID,
                                                        &wifi_event_handler,
                                                        nullptr, &instance_any_id));
    ESP_ERROR_CHECK(esp_event_handler_instance_register(IP_EVENT,
                                                        IP_EVENT_STA_GOT_IP,
                                                        &wifi_event_handler,
                                                        nullptr, &instance_got_ip));

    wifi_config_t wifi_config = {};
    strncpy(reinterpret_cast<char *>(wifi_config.sta.ssid), CONFIG_WIFI_SSID, sizeof(wifi_config.sta.ssid) - 1);
    strncpy(reinterpret_cast<char *>(wifi_config.sta.password), CONFIG_WIFI_PASSWORD, sizeof(wifi_config.sta.password) - 1);
    wifi_config.sta.threshold.authmode = WIFI_AUTH_WPA2_PSK;

    ESP_ERROR_CHECK(esp_wifi_set_mode(WIFI_MODE_STA));
    ESP_ERROR_CHECK(esp_wifi_set_config(WIFI_IF_STA, &wifi_config));
    ESP_ERROR_CHECK(esp_wifi_start());

    ESP_LOGI(TAG, "Connecting to SSID: %s", CONFIG_WIFI_SSID);

    const EventBits_t bits = xEventGroupWaitBits(
        s_wifi_event_group,
        WIFI_CONNECTED_BIT | WIFI_FAIL_BIT,
        pdFALSE, pdFALSE,
        pdMS_TO_TICKS(CONFIG_WIFI_CONNECT_TIMEOUT_MS)
    );

    esp_event_handler_instance_unregister(WIFI_EVENT, ESP_EVENT_ANY_ID, instance_any_id);
    esp_event_handler_instance_unregister(IP_EVENT, IP_EVENT_STA_GOT_IP, instance_got_ip);
    vEventGroupDelete(s_wifi_event_group);

    if (bits & WIFI_CONNECTED_BIT) {
        ESP_LOGI(TAG, "Connected to WiFi");
        return true;
    }
    ESP_LOGE(TAG, "Failed to connect after %d ms", CONFIG_WIFI_CONNECT_TIMEOUT_MS);
    return false;
}

void wifi_disconnect(void) {
    ESP_ERROR_CHECK(esp_wifi_stop());
    ESP_ERROR_CHECK(esp_wifi_deinit());
}

bool wifi_post_sensor_data(float temperature_c, float humidity_percent) {
    char post_data[128];
    const int len = snprintf(post_data, sizeof(post_data),
                             "{\"temperature\":%.1f,\"humidity\":%.1f}",
                             temperature_c, humidity_percent);
    if (len < 0 || static_cast<size_t>(len) >= sizeof(post_data)) {
        ESP_LOGE(TAG, "Failed to serialize JSON");
        return false;
    }

    char device_id[33];
    std::snprintf(device_id, sizeof(device_id), "%s", CONFIG_DEVICE_ID);

    char timestamp_hex[17];
    const auto now = static_cast<long long>(std::time(nullptr));
    std::snprintf(timestamp_hex, sizeof(timestamp_hex), "%lld", now);

    char signature_hex[65];
    if (!compute_hmac_sha256_hex(CONFIG_SERVER_HMAC_SECRET,
                                 device_id,
                                 timestamp_hex,
                                 reinterpret_cast<const uint8_t *>(post_data),
                                 static_cast<size_t>(len),
                                 signature_hex,
                                 sizeof(signature_hex))) {
        ESP_LOGE(TAG, "Failed to compute HMAC");
        return false;
                                 }

    char url[128];
    const int url_len = snprintf(url, sizeof(url), "http://%s:%d%s",
                                 CONFIG_SERVER_HOST, CONFIG_SERVER_PORT, CONFIG_SERVER_PATH);
    if (url_len < 0 || static_cast<size_t>(url_len) >= sizeof(url)) {
        ESP_LOGE(TAG, "Server URL is too long");
        return false;
    }

    esp_http_client_config_t config = {};
    config.url = url;
    config.method = HTTP_METHOD_POST;
    config.timeout_ms = CONFIG_HTTP_TIMEOUT_MS;
    config.event_handler = http_event_handler;

    esp_http_client_handle_t client = esp_http_client_init(&config);
    if (!client) {
        ESP_LOGE(TAG, "Failed to initialize HTTP client");
        return false;
    }

    esp_http_client_set_header(client, "Content-Type", "application/json");
    esp_http_client_set_header(client, "X-Signature", signature_hex);
    esp_http_client_set_header(client, "X-Timestamp", timestamp_hex);
    esp_http_client_set_header(client, "X-Device-ID", device_id);

    esp_http_client_set_post_field(client, post_data, len);

    const esp_err_t err = esp_http_client_perform(client);
    const int status = esp_http_client_get_status_code(client);
    esp_http_client_cleanup(client);

    if (err != ESP_OK) {
        ESP_LOGE(TAG, "HTTP POST failed: %s", esp_err_to_name(err));
        return false;
    }
    if (status < 200 || status >= 300) {
        ESP_LOGE(TAG, "HTTP status %d", status);
        return false;
    }
    ESP_LOGI(TAG, "POST success, status %d", status);
    return true;
}
