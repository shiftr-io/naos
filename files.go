package naos

// SdkconfigContent holds the default content of the 'sdkconfig' file.
const SdkconfigContent = `#
# Automatically generated file; DO NOT EDIT.
# Espressif IoT Development Framework Configuration
#

#
# SDK tool configuration
#
CONFIG_TOOLPREFIX="xtensa-esp32-elf-"
CONFIG_PYTHON="python"

#
# Bootloader config
#
# CONFIG_LOG_BOOTLOADER_LEVEL_NONE is not set
# CONFIG_LOG_BOOTLOADER_LEVEL_ERROR is not set
# CONFIG_LOG_BOOTLOADER_LEVEL_WARN is not set
CONFIG_LOG_BOOTLOADER_LEVEL_INFO=y
# CONFIG_LOG_BOOTLOADER_LEVEL_DEBUG is not set
# CONFIG_LOG_BOOTLOADER_LEVEL_VERBOSE is not set
CONFIG_LOG_BOOTLOADER_LEVEL=3

#
# Security features
#
# CONFIG_SECURE_BOOT_ENABLED is not set
# CONFIG_FLASH_ENCRYPTION_ENABLED is not set

#
# Serial flasher config
#
CONFIG_ESPTOOLPY_PORT="/dev/cu.SLAB_USBtoUART"
# CONFIG_ESPTOOLPY_BAUD_115200B is not set
# CONFIG_ESPTOOLPY_BAUD_230400B is not set
CONFIG_ESPTOOLPY_BAUD_921600B=y
# CONFIG_ESPTOOLPY_BAUD_2MB is not set
# CONFIG_ESPTOOLPY_BAUD_OTHER is not set
CONFIG_ESPTOOLPY_BAUD_OTHER_VAL=115200
CONFIG_ESPTOOLPY_BAUD=921600
CONFIG_ESPTOOLPY_COMPRESSED=y
# CONFIG_FLASHMODE_QIO is not set
# CONFIG_FLASHMODE_QOUT is not set
CONFIG_FLASHMODE_DIO=y
# CONFIG_FLASHMODE_DOUT is not set
CONFIG_ESPTOOLPY_FLASHMODE="dio"
# CONFIG_ESPTOOLPY_FLASHFREQ_80M is not set
CONFIG_ESPTOOLPY_FLASHFREQ_40M=y
# CONFIG_ESPTOOLPY_FLASHFREQ_26M is not set
# CONFIG_ESPTOOLPY_FLASHFREQ_20M is not set
CONFIG_ESPTOOLPY_FLASHFREQ="40m"
# CONFIG_ESPTOOLPY_FLASHSIZE_1MB is not set
CONFIG_ESPTOOLPY_FLASHSIZE_2MB=y
# CONFIG_ESPTOOLPY_FLASHSIZE_4MB is not set
# CONFIG_ESPTOOLPY_FLASHSIZE_8MB is not set
# CONFIG_ESPTOOLPY_FLASHSIZE_16MB is not set
CONFIG_ESPTOOLPY_FLASHSIZE="2MB"
CONFIG_ESPTOOLPY_FLASHSIZE_DETECT=y
CONFIG_ESPTOOLPY_BEFORE_RESET=y
# CONFIG_ESPTOOLPY_BEFORE_NORESET is not set
CONFIG_ESPTOOLPY_BEFORE="default_reset"
CONFIG_ESPTOOLPY_AFTER_RESET=y
# CONFIG_ESPTOOLPY_AFTER_NORESET is not set
CONFIG_ESPTOOLPY_AFTER="hard_reset"
# CONFIG_MONITOR_BAUD_9600B is not set
# CONFIG_MONITOR_BAUD_57600B is not set
CONFIG_MONITOR_BAUD_115200B=y
# CONFIG_MONITOR_BAUD_230400B is not set
# CONFIG_MONITOR_BAUD_921600B is not set
# CONFIG_MONITOR_BAUD_2MB is not set
# CONFIG_MONITOR_BAUD_OTHER is not set
CONFIG_MONITOR_BAUD_OTHER_VAL=115200
CONFIG_MONITOR_BAUD=115200

#
# Partition Table
#
# CONFIG_PARTITION_TABLE_SINGLE_APP is not set
CONFIG_PARTITION_TABLE_TWO_OTA=y
# CONFIG_PARTITION_TABLE_CUSTOM is not set
CONFIG_PARTITION_TABLE_CUSTOM_FILENAME="partitions.csv"
CONFIG_PARTITION_TABLE_CUSTOM_APP_BIN_OFFSET=0x10000
CONFIG_PARTITION_TABLE_FILENAME="partitions_two_ota.csv"
CONFIG_APP_OFFSET=0x10000
CONFIG_OPTIMIZATION_LEVEL_DEBUG=y
# CONFIG_OPTIMIZATION_LEVEL_RELEASE is not set

#
# Component config
#
# CONFIG_AWS_IOT_SDK is not set
CONFIG_BT_ENABLED=y
CONFIG_BLUEDROID_ENABLED=y
CONFIG_BTC_TASK_STACK_SIZE=3072
# CONFIG_BLUEDROID_MEM_DEBUG is not set
# CONFIG_CLASSIC_BT_ENABLED is not set
# CONFIG_BT_DRAM_RELEASE is not set
CONFIG_GATTS_ENABLE=y
CONFIG_GATTC_ENABLE=y
CONFIG_SMP_ENABLE=y
# CONFIG_BT_STACK_NO_LOG is not set
CONFIG_BT_ACL_CONNECTIONS=4
CONFIG_BTDM_CONTROLLER_RUN_CPU=0
CONFIG_BT_RESERVE_DRAM=0x10000

#
# esp-mqtt
#
CONFIG_ESP_MQTT_ENABLED=y
CONFIG_ESP_MQTT_TASK_STACK_SIZE=4096
CONFIG_ESP_MQTT_TASK_STACK_PRIORITY=5
CONFIG_ESP_MQTT_EVENT_QUEUE_SIZE=64

#
# ESP32-specific
#
# CONFIG_ESP32_DEFAULT_CPU_FREQ_80 is not set
# CONFIG_ESP32_DEFAULT_CPU_FREQ_160 is not set
CONFIG_ESP32_DEFAULT_CPU_FREQ_240=y
CONFIG_ESP32_DEFAULT_CPU_FREQ_MHZ=240
CONFIG_MEMMAP_SMP=y
# CONFIG_MEMMAP_TRACEMEM is not set
# CONFIG_MEMMAP_TRACEMEM_TWOBANKS is not set
# CONFIG_ESP32_TRAX is not set
CONFIG_TRACEMEM_RESERVE_DRAM=0x0
# CONFIG_ESP32_ENABLE_COREDUMP_TO_FLASH is not set
# CONFIG_ESP32_ENABLE_COREDUMP_TO_UART is not set
CONFIG_ESP32_ENABLE_COREDUMP_TO_NONE=y
# CONFIG_ESP32_ENABLE_COREDUMP is not set
# CONFIG_ESP32_APPTRACE_DEST_TRAX is not set
# CONFIG_ESP32_APPTRACE_DEST_UART is not set
CONFIG_ESP32_APPTRACE_DEST_NONE=y
# CONFIG_ESP32_APPTRACE_ENABLE is not set
# CONFIG_TWO_UNIVERSAL_MAC_ADDRESS is not set
CONFIG_FOUR_UNIVERSAL_MAC_ADDRESS=y
CONFIG_NUMBER_OF_UNIVERSAL_MAC_ADDRESS=4
CONFIG_SYSTEM_EVENT_QUEUE_SIZE=32
CONFIG_SYSTEM_EVENT_TASK_STACK_SIZE=4096
CONFIG_MAIN_TASK_STACK_SIZE=4096
CONFIG_NEWLIB_STDOUT_ADDCR=y
# CONFIG_NEWLIB_NANO_FORMAT is not set
CONFIG_CONSOLE_UART_DEFAULT=y
# CONFIG_CONSOLE_UART_CUSTOM is not set
# CONFIG_CONSOLE_UART_NONE is not set
CONFIG_CONSOLE_UART_NUM=0
CONFIG_CONSOLE_UART_BAUDRATE=115200
# CONFIG_ULP_COPROC_ENABLED is not set
CONFIG_ULP_COPROC_RESERVE_MEM=0
# CONFIG_ESP32_PANIC_PRINT_HALT is not set
CONFIG_ESP32_PANIC_PRINT_REBOOT=y
# CONFIG_ESP32_PANIC_SILENT_REBOOT is not set
# CONFIG_ESP32_PANIC_GDBSTUB is not set
CONFIG_ESP32_DEBUG_OCDAWARE=y
CONFIG_INT_WDT=y
CONFIG_INT_WDT_TIMEOUT_MS=300
CONFIG_INT_WDT_CHECK_CPU1=y
CONFIG_TASK_WDT=y
# CONFIG_TASK_WDT_PANIC is not set
CONFIG_TASK_WDT_TIMEOUT_S=5
CONFIG_TASK_WDT_CHECK_IDLE_TASK=y
CONFIG_TASK_WDT_CHECK_IDLE_TASK_CPU1=y
# CONFIG_ESP32_TIME_SYSCALL_USE_RTC is not set
CONFIG_ESP32_TIME_SYSCALL_USE_RTC_FRC1=y
# CONFIG_ESP32_TIME_SYSCALL_USE_FRC1 is not set
# CONFIG_ESP32_TIME_SYSCALL_USE_NONE is not set
CONFIG_ESP32_RTC_CLOCK_SOURCE_INTERNAL_RC=y
# CONFIG_ESP32_RTC_CLOCK_SOURCE_EXTERNAL_CRYSTAL is not set
CONFIG_ESP32_RTC_CLK_CAL_CYCLES=1024
CONFIG_ESP32_DEEP_SLEEP_WAKEUP_DELAY=0
# CONFIG_ESP32_XTAL_FREQ_40 is not set
# CONFIG_ESP32_XTAL_FREQ_26 is not set
CONFIG_ESP32_XTAL_FREQ_AUTO=y
CONFIG_ESP32_XTAL_FREQ=0
CONFIG_WIFI_ENABLED=y
# CONFIG_SW_COEXIST_ENABLE is not set
CONFIG_ESP32_WIFI_STATIC_RX_BUFFER_NUM=10
CONFIG_ESP32_WIFI_DYNAMIC_RX_BUFFER_NUM=32
# CONFIG_ESP32_WIFI_STATIC_TX_BUFFER is not set
CONFIG_ESP32_WIFI_DYNAMIC_TX_BUFFER=y
CONFIG_ESP32_WIFI_TX_BUFFER_TYPE=1
CONFIG_ESP32_WIFI_DYNAMIC_TX_BUFFER_NUM=32
CONFIG_ESP32_WIFI_AMPDU_ENABLED=y
CONFIG_ESP32_WIFI_NVS_ENABLED=y
CONFIG_PHY_ENABLED=y

#
# PHY
#
CONFIG_ESP32_PHY_CALIBRATION_AND_DATA_STORAGE=y
# CONFIG_ESP32_PHY_INIT_DATA_IN_PARTITION is not set
CONFIG_ESP32_PHY_MAX_WIFI_TX_POWER=20
CONFIG_ESP32_PHY_MAX_TX_POWER=20
# CONFIG_ETHERNET is not set

#
# FAT Filesystem support
#
CONFIG_FATFS_CODEPAGE_ASCII=y
# CONFIG_FATFS_CODEPAGE_437 is not set
# CONFIG_FATFS_CODEPAGE_720 is not set
# CONFIG_FATFS_CODEPAGE_737 is not set
# CONFIG_FATFS_CODEPAGE_771 is not set
# CONFIG_FATFS_CODEPAGE_775 is not set
# CONFIG_FATFS_CODEPAGE_850 is not set
# CONFIG_FATFS_CODEPAGE_852 is not set
# CONFIG_FATFS_CODEPAGE_855 is not set
# CONFIG_FATFS_CODEPAGE_857 is not set
# CONFIG_FATFS_CODEPAGE_860 is not set
# CONFIG_FATFS_CODEPAGE_861 is not set
# CONFIG_FATFS_CODEPAGE_862 is not set
# CONFIG_FATFS_CODEPAGE_863 is not set
# CONFIG_FATFS_CODEPAGE_864 is not set
# CONFIG_FATFS_CODEPAGE_865 is not set
# CONFIG_FATFS_CODEPAGE_866 is not set
# CONFIG_FATFS_CODEPAGE_869 is not set
# CONFIG_FATFS_CODEPAGE_932 is not set
# CONFIG_FATFS_CODEPAGE_936 is not set
# CONFIG_FATFS_CODEPAGE_949 is not set
# CONFIG_FATFS_CODEPAGE_950 is not set
CONFIG_FATFS_CODEPAGE=1
CONFIG_FATFS_MAX_LFN=255

#
# FreeRTOS
#
# CONFIG_FREERTOS_UNICORE is not set
CONFIG_FREERTOS_CORETIMER_0=y
# CONFIG_FREERTOS_CORETIMER_1 is not set
CONFIG_FREERTOS_HZ=1000
CONFIG_FREERTOS_ASSERT_ON_UNTESTED_FUNCTION=y
# CONFIG_FREERTOS_CHECK_STACKOVERFLOW_NONE is not set
# CONFIG_FREERTOS_CHECK_STACKOVERFLOW_PTRVAL is not set
CONFIG_FREERTOS_CHECK_STACKOVERFLOW_CANARY=y
CONFIG_FREERTOS_WATCHPOINT_END_OF_STACK=y
CONFIG_FREERTOS_THREAD_LOCAL_STORAGE_POINTERS=1
CONFIG_FREERTOS_ASSERT_FAIL_ABORT=y
# CONFIG_FREERTOS_ASSERT_FAIL_PRINT_CONTINUE is not set
# CONFIG_FREERTOS_ASSERT_DISABLE is not set
CONFIG_FREERTOS_BREAK_ON_SCHEDULER_START_JTAG=y
# CONFIG_ENABLE_MEMORY_DEBUG is not set
CONFIG_FREERTOS_ISR_STACKSIZE=1536
# CONFIG_FREERTOS_LEGACY_HOOKS is not set
CONFIG_FREERTOS_MAX_TASK_NAME_LEN=16
# CONFIG_SUPPORT_STATIC_ALLOCATION is not set
CONFIG_TIMER_TASK_PRIORITY=1
CONFIG_TIMER_TASK_STACK_DEPTH=2048
CONFIG_TIMER_QUEUE_LENGTH=10
# CONFIG_FREERTOS_DEBUG_INTERNALS is not set

#
# Log output
#
# CONFIG_LOG_DEFAULT_LEVEL_NONE is not set
# CONFIG_LOG_DEFAULT_LEVEL_ERROR is not set
# CONFIG_LOG_DEFAULT_LEVEL_WARN is not set
CONFIG_LOG_DEFAULT_LEVEL_INFO=y
# CONFIG_LOG_DEFAULT_LEVEL_DEBUG is not set
# CONFIG_LOG_DEFAULT_LEVEL_VERBOSE is not set
CONFIG_LOG_DEFAULT_LEVEL=3
CONFIG_LOG_COLORS=y

#
# LWIP
#
# CONFIG_L2_TO_L3_COPY is not set
CONFIG_LWIP_MAX_SOCKETS=10
CONFIG_LWIP_THREAD_LOCAL_STORAGE_INDEX=0
# CONFIG_LWIP_SO_REUSE is not set
CONFIG_LWIP_SO_RCVBUF=y
CONFIG_LWIP_DHCP_MAX_NTP_SERVERS=1
# CONFIG_LWIP_IP_FRAG is not set
# CONFIG_LWIP_IP_REASSEMBLY is not set
CONFIG_TCP_MAXRTX=12
CONFIG_TCP_SYNMAXRTX=6
CONFIG_LWIP_DHCP_DOES_ARP_CHECK=y
CONFIG_TCPIP_TASK_STACK_SIZE=2560
# CONFIG_PPP_SUPPORT is not set

#
# mbedTLS
#
CONFIG_MBEDTLS_SSL_MAX_CONTENT_LEN=16384
# CONFIG_MBEDTLS_DEBUG is not set
CONFIG_MBEDTLS_HARDWARE_AES=y
CONFIG_MBEDTLS_HARDWARE_MPI=y
CONFIG_MBEDTLS_MPI_USE_INTERRUPT=y
CONFIG_MBEDTLS_HARDWARE_SHA=y
CONFIG_MBEDTLS_HAVE_TIME=y
# CONFIG_MBEDTLS_HAVE_TIME_DATE is not set

#
# naos
#
CONFIG_NAOS_MQTT_BUFFER_SIZE=10000
CONFIG_NAOS_MQTT_COMMAND_TIMEOUT=2000
CONFIG_NAOS_HEARTBEAT_INTERVAL=5000
CONFIG_NAOS_UPDATE_MAX_CHUNK_SIZE=4096

#
# OpenSSL
#
# CONFIG_OPENSSL_DEBUG is not set
CONFIG_OPENSSL_ASSERT_DO_NOTHING=y
# CONFIG_OPENSSL_ASSERT_EXIT is not set

#
# SPI Flash driver
#
# CONFIG_SPI_FLASH_ENABLE_COUNTERS is not set
CONFIG_SPI_FLASH_ROM_DRIVER_PATCH=y
`

// MakefileContent holds the default content of the 'Makefile' file.
const MakefileContent = `PROJECT_NAME := naos-project

export IDF_PATH := ../esp-idf
include $(IDF_PATH)/make/project.mk
`
