#include <esp_log.h>
#include <nvs.h>
#include <stdlib.h>
#include <string.h>

#include "naos.h"
#include "utils.h"

typedef struct {
  naos_type_t type;
  const char *param;
  void *pointer;
} naos_params_sync_item_t;

static naos_params_sync_item_t naos_params_sync_registry[CONFIG_NAOS_SYNC_REGISTRY_SIZE];

static size_t naos_params_sync_registry_count = 0;

static nvs_handle naos_params_nvs_handle;

static bool naos_params_add_sync(const char *param, naos_params_sync_item_t item) {
  // check param length
  if (strlen(param) == 0) {
    return false;
  }

  // check registry count
  if (naos_params_sync_registry_count >= CONFIG_NAOS_SYNC_REGISTRY_SIZE) {
    ESP_LOGE(NAOS_LOG_TAG, "naos_params_add_sync: registry full");
    return false;
  }

  // add entry to registry
  naos_params_sync_registry[naos_params_sync_registry_count] = item;

  // increment counter
  naos_params_sync_registry_count++;

  return true;
}

static void naos_params_update_sync(const char *param) {
  // update synchronized variables
  for (size_t i = 0; i < naos_params_sync_registry_count; i++) {
    // get item
    naos_params_sync_item_t item = naos_params_sync_registry[i];

    // check param
    if (strcmp(item.param, param) != 0) {
      continue;
    }

    // check type
    switch (item.type) {
      case NAOS_STRING: {
        // get pointer
        char **pointer = item.pointer;

        // free existing value if pointer is set
        if (*pointer != NULL) {
          free(*pointer);
        }

        // set new value
        *pointer = strdup(naos_get(param));

        break;
      }
      case NAOS_BOOL: {
        // get pointer
        bool *pointer = item.pointer;

        // set new value
        *pointer = naos_get_b(param);

        break;
      }
      case NAOS_LONG: {
        // get pointer
        int32_t *pointer = item.pointer;

        // set new value
        *pointer = naos_get_l(param);

        break;
      }
      case NAOS_DOUBLE: {
        // get pointer
        double *pointer = item.pointer;

        // set new value
        *pointer = naos_get_d(param);

        break;
      }
    }
  }
}

// TODO: Rename nvs namespace.

void naos_params_init() {
  // open nvs namespace
  ESP_ERROR_CHECK(nvs_open("naos-manager", NVS_READWRITE, &naos_params_nvs_handle));

  // initialize params
  for (int i = 0; i < naos_config()->num_parameters; i++) {
    // get param
    naos_param_t param = naos_config()->parameters[i];

    // check_type
    switch (param.type) {
      case NAOS_STRING:
        naos_ensure(param.name, param.default_s);
        break;
      case NAOS_BOOL:
        naos_ensure_b(param.name, param.default_b);
        break;
      case NAOS_LONG:
        naos_ensure_l(param.name, param.default_l);
        break;
      case NAOS_DOUBLE:
        naos_ensure_d(param.name, param.default_d);
        break;
    }
  }

  // setup synchronization
  for (int i = 0; i < naos_config()->num_parameters; i++) {
    // get param
    naos_param_t param = naos_config()->parameters[i];

    // check_type
    switch (param.type) {
      case NAOS_STRING:
        if (param.sync_s != NULL) {
          naos_sync(param.name, param.sync_s);
        }
        break;
      case NAOS_BOOL:
        if (param.sync_b != NULL) {
          naos_sync_b(param.name, param.sync_b);
        }
        break;
      case NAOS_LONG:
        if (param.sync_l != NULL) {
          naos_sync_l(param.name, param.sync_l);
        }
        break;
      case NAOS_DOUBLE:
        if (param.sync_d != NULL) {
          naos_sync_d(param.name, param.sync_d);
        }
        break;
    }
  }
}

char *naos_params_list() {
  // return empty string if there are no params
  if (naos_config()->num_parameters == 0) {
    return strdup("");
  }

  // determine list length
  size_t length = 0;
  for (int i = 0; i < naos_config()->num_parameters; i++) {
    // get param
    naos_param_t param = naos_config()->parameters[i];

    // add length
    length += strlen(param.name) + 3;
  }

  // allocate buffer
  char *buf = malloc(length);

  // write names
  size_t pos = 0;
  for (int i = 0; i < naos_config()->num_parameters; i++) {
    // get param
    naos_param_t param = naos_config()->parameters[i];

    // copy name
    strcpy(buf + pos, param.name);
    pos += strlen(param.name);

    // write separator
    buf[pos] = ':';
    pos++;

    // write type
    switch (param.type) {
      case NAOS_STRING:
        buf[pos] = 's';
        break;
      case NAOS_BOOL:
        buf[pos] = 'b';
        break;
      case NAOS_LONG:
        buf[pos] = 'l';
        break;
      case NAOS_DOUBLE:
        buf[pos] = 'd';
        break;
    }
    pos++;

    // write comma or zero
    buf[pos] = (char)((i == naos_config()->num_parameters - 1) ? '\0' : ',');
    pos++;
  }

  return buf;
}

char *naos_get(const char *param) {
  // static reference to buffer
  static char *buf;

  // free last param
  if (buf != NULL) {
    free(buf);
    buf = NULL;
  }

  // get param size
  size_t required_size;
  esp_err_t err = nvs_get_str(naos_params_nvs_handle, param, NULL, &required_size);
  if (err == ESP_ERR_NVS_NOT_FOUND) {
    buf = strdup("");
    return buf;
  } else {
    ESP_ERROR_CHECK(err);
  }

  // allocate size
  buf = malloc(required_size);
  ESP_ERROR_CHECK(nvs_get_str(naos_params_nvs_handle, param, buf, &required_size));

  return buf;
}

bool naos_get_b(const char *param) { return strtol(naos_get(param), NULL, 10) == 1; }

int32_t naos_get_l(const char *param) { return (int32_t)strtol(naos_get(param), NULL, 10); }

double naos_get_d(const char *param) { return strtod(naos_get(param), NULL); }

void naos_set(const char *param, const char *value) {
  // set parameter
  ESP_ERROR_CHECK(nvs_set_str(naos_params_nvs_handle, param, value));

  // sync param
  naos_params_update_sync(param);
}

void naos_set_b(const char *param, bool value) { naos_set(param, naos_i2str(value)); }

void naos_set_l(const char *param, int32_t value) { naos_set(param, naos_i2str(value)); }

void naos_set_d(const char *param, double value) { naos_set(param, naos_d2str(value)); }

bool naos_ensure(const char *param, const char *value) {
  // check parameter
  size_t required_size;
  esp_err_t err = nvs_get_str(naos_params_nvs_handle, param, NULL, &required_size);
  if (err == ESP_ERR_NVS_NOT_FOUND) {
    naos_set(param, value);
    return true;
  } else {
    ESP_ERROR_CHECK(err);
  }

  return false;
}

bool naos_ensure_b(const char *param, bool value) { return naos_ensure(param, naos_i2str(value)); }

bool naos_ensure_l(const char *param, int32_t value) { return naos_ensure(param, naos_i2str(value)); }

bool naos_ensure_d(const char *param, double value) { return naos_ensure(param, naos_d2str(value)); }

bool naos_unset(const char *param) {
  // erase parameter
  esp_err_t err = nvs_erase_key(naos_params_nvs_handle, param);
  if (err == ESP_ERR_NVS_NOT_FOUND) {
    return false;
  } else {
    ESP_ERROR_CHECK(err);
  }

  // sync param
  naos_params_update_sync(param);

  return true;
}

bool naos_sync(const char *param, char **pointer) {
  // prepare item
  naos_params_sync_item_t item = {
      .type = NAOS_STRING,
      .param = param,
      .pointer = pointer,
  };

  // add sync item
  bool ret = naos_params_add_sync(param, item);

  // read current value
  *pointer = strdup(naos_get(param));

  return ret;
}

bool naos_sync_b(const char *param, bool *pointer) {
  // prepare item
  naos_params_sync_item_t item = {
      .type = NAOS_BOOL,
      .param = param,
      .pointer = pointer,
  };

  // add sync item
  bool ret = naos_params_add_sync(param, item);

  // read current value
  *pointer = naos_get_b(param);

  return ret;
}

bool naos_sync_l(const char *param, int32_t *pointer) {
  // prepare item
  naos_params_sync_item_t item = {
      .type = NAOS_LONG,
      .param = param,
      .pointer = pointer,
  };

  // add sync item
  bool ret = naos_params_add_sync(param, item);

  // read current value
  *pointer = naos_get_l(param);

  return ret;
}

bool naos_sync_d(const char *param, double *pointer) {
  // prepare item
  naos_params_sync_item_t item = {
      .type = NAOS_DOUBLE,
      .param = param,
      .pointer = pointer,
  };

  // add sync item
  bool ret = naos_params_add_sync(param, item);

  // read current value
  *pointer = naos_get_d(param);

  return ret;
}
