#!/bin/bash

# Globals
#   BGPS_CONFIG - map of config values
#   BGPS_MAX    - length of largest key in config map

declare -A BGPS_CONFIG
VALID_CONFIG_KEYS=(
  "COLOR_PROMPT"
  # TODO add more valid keys
)

#
# Add configuration to config map
#
_set_config_value() {
  local key="${1^^}"
  local value="$2"
  for i in "${VALID_CONFIG_KEYS[@]}"; do 
    if [[ $key == $i ]] ; then
      # record max key length for future formatting during print
      local key_length=${#key}
      if ((!BGPS_MAX || BGPS_MAX < key_length)) ; then
        BGPS_MAX=$key_length
      fi
      
      # add to config map
      BGPS_CONFIG["${key}"]="${value}"
    fi
  done
}

#
# Get configuration from config map 
#
_get_config_value() {
  local key="${1^^}"
  echo "${BGPS_CONFIG[$key]}"
}

#
# Read file and store in global config map
#
_read_config() {
  local config_file="$1"
  if [ -f "$config_file" ]; then
    while IFS='=' read -r key value; do
      _set_config_value $key $value
    done < "$config_file"
  fi
}

#
# Set the config map by reading the values from the provided file.
# If the config map already exists, consider it cached and do not
# read from the file.
#
_set_config() {
  if [[ ! -v BGPS_CONFIG[@] ]] ; then
    _read_config "$@" 
  fi
}

#
# Pretty print the config 
#
_print_config() {
  if [[ -v BGPS_CONFIG[@] ]] ; then
    for i in "${!BGPS_CONFIG[@]}"; do 
      printf "%-${BGPS_MAX}s = %s\n" $i ${BGPS_CONFIG[$i]} 
    done
  fi
}

#
# Unset global variables
#
_delete_config() {
  unset BGPS_CONFIG
  unset BGPS_MAX
}

#
# Unset variables and functions
#
_bgps_config_unset() {
  unset VALID_CONFIG_KEYS
  unset -f _set_config_value
  unset -f _get_config_value
  unset -f _read_config
  unset -f _set_config
  unset -f _print_config
  unset -f _delete_config
  unset -f _bgps_config_unset
}
