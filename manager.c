#include <esp_log.h>
#include <esp_system.h>
#include <freertos/FreeRTOS.h>
#include <freertos/semphr.h>
#include <freertos/task.h>
#include <string.h>

#include <nadk.h>

#include "ble.h"
#include "device.h"
#include "general.h"
#include "led.h"
#include "mqtt.h"
#include "wifi.h"

SemaphoreHandle_t nadk_manager_mutex;

typedef enum nadk_manager_state_t {
  NADK_MANAGER_STATE_DISCONNECTED,
  NADK_MANAGER_STATE_CONNECTED,
  NADK_MANAGER_STATE_NETWORKED
} nadk_manager_state_t;

static nadk_manager_state_t nadk_manager_current_state;

static void nadk_manager_set_state(nadk_manager_state_t new_state) {
  // default state name
  const char *name = "Unknown";

  // triage state
  switch (new_state) {
    // handle disconnected state
    case NADK_MANAGER_STATE_DISCONNECTED: {
      name = "Disconnected";
      nadk_led_set(false, false);
      break;
    }

    // handle connected state
    case NADK_MANAGER_STATE_CONNECTED: {
      name = "Connected";
      nadk_led_set(true, false);
      break;
    }

    // handle networked state
    case NADK_MANAGER_STATE_NETWORKED: {
      name = "Networked";
      nadk_led_set(false, true);
      break;
    }
  }

  // change state
  nadk_manager_current_state = new_state;

  // update connection status
  nadk_ble_set_string(NADK_BLE_ID_CONNECTION_STATUS, (char *)name);

  ESP_LOGI(NADK_LOG_TAG, "nadk_manager_set_state: %s", name)
}

static void nadk_manager_configure_wifi() {
  // get ssid
  char wifi_ssid[NADK_BLE_STRING_SIZE];
  nadk_ble_get_string(NADK_BLE_ID_WIFI_SSID, wifi_ssid);

  // get password
  char wifi_password[NADK_BLE_STRING_SIZE];
  nadk_ble_get_string(NADK_BLE_ID_WIFI_PASSWORD, wifi_password);

  // configure wifi
  nadk_wifi_configure(wifi_ssid, wifi_password);
}

static void nadk_manager_start_mqtt() {
  // get host
  char mqtt_host[NADK_BLE_STRING_SIZE];
  nadk_ble_get_string(NADK_BLE_ID_MQTT_HOST, mqtt_host);

  // get client id
  char mqtt_client_id[NADK_BLE_STRING_SIZE];
  nadk_ble_get_string(NADK_BLE_ID_MQTT_CLIENT_ID, mqtt_client_id);

  // get username
  char mqtt_username[NADK_BLE_STRING_SIZE];
  nadk_ble_get_string(NADK_BLE_ID_MQTT_USERNAME, mqtt_username);

  // get password
  char mqtt_password[NADK_BLE_STRING_SIZE];
  nadk_ble_get_string(NADK_BLE_ID_MQTT_PASSWORD, mqtt_password);

  // get base topic
  char base_topic[NADK_BLE_STRING_SIZE];
  nadk_ble_get_string(NADK_BLE_ID_BASE_TOPIC, base_topic);

  // start mqtt
  nadk_mqtt_start(mqtt_host, 1883, mqtt_client_id, mqtt_username, mqtt_password, base_topic);
}

static void nadk_manager_ble_callback(nadk_ble_id_t id) {
  // dismiss any other changed characteristic
  if (id != NADK_BLE_ID_COMMAND) {
    return;
  }

  // acquire mutex
  NADK_LOCK(nadk_manager_mutex);

  // get value
  char value[NADK_BLE_STRING_SIZE];
  nadk_ble_get_string(NADK_BLE_ID_COMMAND, value);

  // detect command
  bool restart_mqtt = strcmp(value, "restart-mqtt") == 0;
  bool restart_wifi = strcmp(value, "restart-wifi") == 0;

  // handle wifi restart
  if (restart_wifi) {
    ESP_LOGI(NADK_LOG_TAG, "nadk_manager_ble_callback: restart wifi");

    switch (nadk_manager_current_state) {
      case NADK_MANAGER_STATE_NETWORKED: {
        // stop device
        nadk_device_stop();

        // fallthrough
      }

      case NADK_MANAGER_STATE_CONNECTED: {
        // stop mqtt client
        nadk_mqtt_stop();

        // change state
        nadk_manager_set_state(NADK_MANAGER_STATE_DISCONNECTED);

        // fallthrough
      }

      case NADK_MANAGER_STATE_DISCONNECTED: {
        // restart wifi
        nadk_manager_configure_wifi();
      }
    }
  }

  // handle mqtt restart
  if (restart_mqtt) {
    ESP_LOGI(NADK_LOG_TAG, "nadk_manager_ble_callback: restart mqtt");

    switch (nadk_manager_current_state) {
      case NADK_MANAGER_STATE_NETWORKED: {
        // stop device
        nadk_device_stop();

        // change state
        nadk_manager_set_state(NADK_MANAGER_STATE_CONNECTED);

        // fallthrough
      }

      case NADK_MANAGER_STATE_CONNECTED: {
        // stop mqtt client
        nadk_mqtt_stop();

        // restart mqtt
        nadk_manager_start_mqtt();

        // fallthrough
      }

      case NADK_MANAGER_STATE_DISCONNECTED: {
        // do nothing if not yet connected
      }
    }
  }

  // release mutex
  NADK_UNLOCK(nadk_manager_mutex);
}

static void nadk_manager_wifi_callback(nadk_wifi_status_t status) {
  // acquire mutex
  NADK_LOCK(nadk_manager_mutex);

  switch (status) {
    case NADK_WIFI_STATUS_CONNECTED: {
      ESP_LOGI(NADK_LOG_TAG, "nadk_manager_wifi_callback: connected");

      // check if connection is new
      if (nadk_manager_current_state == NADK_MANAGER_STATE_DISCONNECTED) {
        // change sate
        nadk_manager_set_state(NADK_MANAGER_STATE_CONNECTED);

        // start wifi
        nadk_manager_start_mqtt();
      }

      break;
    }

    case NADK_WIFI_STATUS_DISCONNECTED: {
      ESP_LOGI(NADK_LOG_TAG, "nadk_manager_wifi_callback: disconnected");

      // check if disconnection is new
      if (nadk_manager_current_state >= NADK_MANAGER_STATE_CONNECTED) {
        // stop mqtt
        nadk_mqtt_stop();

        // change state
        nadk_manager_set_state(NADK_MANAGER_STATE_DISCONNECTED);
      }

      break;
    }
  }

  // release mutex
  NADK_UNLOCK(nadk_manager_mutex);
}

static void nadk_manager_mqtt_callback(esp_mqtt_status_t status) {
  // acquire mutex
  NADK_LOCK(nadk_manager_mutex);

  switch (status) {
    case ESP_MQTT_STATUS_CONNECTED: {
      ESP_LOGI(NADK_LOG_TAG, "nadk_manager_mqtt_callback: connected");

      // check if connection is new
      if (nadk_manager_current_state == NADK_MANAGER_STATE_CONNECTED) {
        // change state
        nadk_manager_set_state(NADK_MANAGER_STATE_NETWORKED);

        // start device
        nadk_device_start();
      }

      break;
    }

    case ESP_MQTT_STATUS_DISCONNECTED: {
      ESP_LOGI(NADK_LOG_TAG, "nadk_manager_mqtt_callback: disconnected");

      // change state
      nadk_manager_set_state(NADK_MANAGER_STATE_CONNECTED);

      // stop device
      nadk_device_stop();

      // restart mqtt
      nadk_manager_start_mqtt();

      break;
    }
  }

  // release mutex
  NADK_UNLOCK(nadk_manager_mutex);
}

void nadk_manager_init(nadk_device_t *device) {
  // delay startup by max 5000ms
  int delay = esp_random() / 858994;
  ESP_LOGI(NADK_LOG_TAG, "nadk_manager_init: delay startup by %dms", delay);
  vTaskDelay(delay / portTICK_PERIOD_MS);

  // create mutex
  nadk_manager_mutex = xSemaphoreCreateMutex();

  // save device
  nadk_device_init(device);

  // initialize LEDs
  nadk_led_init();

  // initialize bluetooth stack
  nadk_ble_init(nadk_manager_ble_callback, device->name);

  // initialize wifi stack
  nadk_wifi_init(nadk_manager_wifi_callback);

  // initialize mqtt client
  nadk_mqtt_init(nadk_manager_mqtt_callback, nadk_device_forward);

  // set initial state
  nadk_manager_set_state(NADK_MANAGER_STATE_DISCONNECTED);

  // initially configure wifi
  nadk_manager_configure_wifi();
}
