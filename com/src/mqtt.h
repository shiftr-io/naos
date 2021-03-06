#ifndef _NAOS_MQTT_H
#define _NAOS_MQTT_H

#include <esp_mqtt.h>

#include <naos.h>

/**
 * The message callback.
 */
typedef void (*naos_mqtt_message_callback_t)(const char *topic, uint8_t *payload, size_t len, naos_scope_t scope);

/**
 * Initialize the MQTT subsystem.
 *
 * @note Should only be called once on boot.
 *
 * @param scb The status callback.
 * @param mcb The message callback.
 */
void naos_mqtt_init(esp_mqtt_status_callback_t scb, naos_mqtt_message_callback_t mcb);

/**
 * Start the MQTT process.
 *
 * @param host The broker host.
 * @param port The broker port.
 * @param client_id The client id.
 * @param username The client username.
 * @param password The client password.
 * @param base_topic The base topic.
 */
void naos_mqtt_start(const char *host, char *port, const char *client_id, const char *username, const char *password,
                     const char *base_topic);

/**
 * Stop the MQTT process.
 */
void naos_mqtt_stop();

#endif  // _NAOS_MQTT_H
