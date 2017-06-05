#ifndef NADK_H
#define NADK_H

#include <stdarg.h>
#include <stdbool.h>
#include <stdint.h>

/**
 * The messages scopes.
 *
 * The 'local' scope denotes messages that are transferred under the configured base topic of the device while the
 * 'global' scope denotes messages that are transferred directly below the root.
 */
typedef enum { NADK_LOCAL, NADK_GLOBAL } nadk_scope_t;

/**
 * Get the string representation of the specified scope.
 *
 * @param scope - The scope.
 * @return The string value.
 */
const char *nadk_scope_str(nadk_scope_t scope);

/**
 * The system statuses.
 */
typedef enum { NADK_DISCONNECTED, NADK_CONNECTED, NADK_NETWORKED } nadk_status_t;

/**
 * Get the string representation of the specified status.
 *
 * @param scope - The status.
 * @return The string value.
 */
const char *nadk_status_str(nadk_status_t status);

/**
 * The main configuration object.
 */
typedef struct {
  /**
   * The device type.
   */
  const char *device_type;

  /**
   * The firmware version.
   */
  const char *firmware_version;

  /**
   * The callback that is called once the device comes online.
   */
  void (*online_callback)();

  /**
   * The callback that is called when a parameter has been updated.
   *
   * @param param - The parameter.
   * @param value - The value.
   */
  void (*update_callback)(const char *param, const char *value);

  /**
   * The message callback is called with incoming messages.
   *
   * Note: The base topic has already been removed from the topic and should not start with a '/'.
   *
   * @param topic
   * @param payload
   * @param len
   * @param scope
   */
  void (*message_callback)(const char *topic, const char *payload, unsigned int len, nadk_scope_t scope);

  /**
   * The loop callback is called in over and over as long as the device is online.
   */
  void (*loop_callback)();

  /**
   * The interval of the loop callback in milliseconds.
   */
  int loop_interval;

  /**
   * The offline callback is called once the device becomes offline.
   */
  void (*offline_callback)();

  /**
   * The callback is called once the device has changed its status.
   */
  void (*status_callback)(nadk_status_t status);

  /**
   * If set, the device will randomly (up to 5s) delay startup to overcome WiFi and MQTT congestion issues if many
   * devices restart at the same time.
   */
  bool delay_startup;
} nadk_config_t;

/**
 * Write a log message.
 *
 * The message will be printed in the console and published to the broker if the device has logging enabled.
 *
 * @param fmt - The message format.
 * @param ... - The used arguments.
 */
void nadk_log(const char *fmt, ...);

/**
 * Initialize the NADK.
 *
 * Note: Should only be called once on boot.
 *
 * @param config - The configuration object.
 */
void nadk_init(nadk_config_t *config);

/**
 * Subscribe to specified topic.
 *
 * The topic is automatically prefixed with the configured base topic if the scope is local.
 *
 * @param topic
 * @param scope
 * @return
 */
bool nadk_subscribe(const char *topic, int qos, nadk_scope_t scope);

/**
 * Unsubscribe from specified topic.
 *
 * The topic is automatically prefixed with the configured base topic if the scope is local.
 *
 * @param topic
 * @param scope
 * @return
 */
bool nadk_unsubscribe(const char *topic, nadk_scope_t scope);

/**
 * Publish bytes payload to specified topic.
 *
 * The topic is automatically prefixed with the configured base topic if the scope is local.
 *
 * @param topic
 * @param payload
 * @param len
 * @param qos
 * @param retained
 * @param scope
 * @return
 */
bool nadk_publish(const char *topic, void *payload, uint16_t len, int qos, bool retained, nadk_scope_t scope);

/**
 * Publish string to specified topic.
 *
 * The topic is automatically prefixed with the configured base topic if the scope is local.
 *
 * @param topic
 * @param str
 * @param qos
 * @param retained
 * @param scope
 * @return
 */
bool nadk_publish_str(const char *topic, const char *str, int qos, bool retained, nadk_scope_t scope);

/**
 * Publish integer to specified topic.
 *
 * The topic is automatically prefixed with the configured base topic if the scope is local.
 *
 * @param topic
 * @param num
 * @param qos
 * @param retained
 * @param scope
 * @return
 */
bool nadk_publish_int(const char *topic, int num, int qos, bool retained, nadk_scope_t scope);

/**
 * Will return the value of the requested parameter.
 *
 * Note: The returned pointer is valid until the next call to nadk_get().
 *
 * @param param - The parameter.
 * @return Pointer to value.
 */
char *nadk_get(const char *param);

#endif  // NADK_H
