'''
Contains configuration options via 'configargparse' library
'''

import configargparse
import logging

from common.constants import ApiDefaults, FlaskDefaults, InfluxdbDefaults

def load_conf():
    ''' Load and injest config or env variables '''

    # Load configuration file and settings
    p = configargparse.ArgParser()

    p.add('--loadtest', required=False, action='store_true', help='enable loadtest mode', env_var='APP_LOADTEST')
    p.add('--group', required=False, help='default group name', env_var='APP_GROUP')
    p.add('--host', required=False, help='listening web ip', env_var='APP_HOST')
    p.add('--keepalive', required=False, type=int, help='default keepalive value in seconds', env_var='APP_KEEPALIVE')
    p.add('--port', required=False, help='listening web port', env_var='APP_PORT')
    p.add('--interval', required=False, type=int, help='llama collection interval in seconds', env_var='LLAMA_INTERVAL')
    p.add('--influxdb-host', required=False, help='InfluxDB Hostname', env_var='INFLUXDB_HOST')
    p.add('--influxdb-port', required=False, type=int, help='InfluxDB Port', env_var='INFLUXDB_PORT')
    p.add('--influxdb-name', required=False, help='InfluxDB Name, defaults to "enphase"', env_var='INFLUXDB_NAME')
    p.add('-v', '--verbose', help='verbose logging', action='store_true', env_var='APP_VERBOSE')

    config = p.parse_args()

    # Set defaults for webserver settings
    if not config.host:
        config.host = FlaskDefaults.host
    if not config.port:
        config.port = FlaskDefaults.port
    if not config.interval:
        config.interval = ApiDefaults.interval

    # Set defaults for InfluxDB settings
    if config.influxdb_host:
        if not config.influxdb_port:
            config.influxdb_port = InfluxdbDefaults.port
        if not config.influxdb_name:
            config.influxdb_name = InfluxdbDefaults.name

    # Set keepalive values, 86400 seconds if none is set
    if config.keepalive:
        # How many seconds before kicking probes from service discovery
        default_keepalive = config.keepalive
    else:
        default_keepalive = ApiDefaults.keepalive

    # Set a default registration group if one is not provided
    if config.group:
        # Load the default group from configuration variables.
        default_group = str(config.group)
    else:
        default_group = ApiDefaults.group

    # Set logging levels
    if config.verbose:
        logging.basicConfig(format="%(asctime)s %(levelname)s %(threadName)s: %(message)s", encoding='utf-8', level=logging.DEBUG)
    else:
        logging.basicConfig(format="%(asctime)s %(levelname)s %(threadName)s: %(message)s", encoding='utf-8', level=logging.INFO)

    # Debug logging for settings
    logging.debug(p.format_values())
    logging.debug(config)
    logging.debug("Default keepalive is set to %i seconds" % default_keepalive)
    logging.debug("Default group is set to '%s'" % default_group)

    return config