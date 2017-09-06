#!/bin/bash

# Globals
#   BGPS_CONFIG - map of config values
#   BGPS_MAX    - length of largest key in config map

source ./bgps-config.sh

BGPS_CONFIG_FILE="${HOME}/.bgps_config"

# TODO implement flags better than this
if [[ $1 == "--ls-config" ]] ; then
  _print_config
elif [[ $1 == "--clear-config" ]] ; then
  _delete_config
else
  _set_config $BGPS_CONFIG_FILE
fi

# clean up
_bgps_config_unset
unset BGPS_CONFIG_FILE
