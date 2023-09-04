'''
Contains configuration options via 'configargparse' library
'''

import configargparse
import logging



def load_conf() -> dict:
    ''' Load and injest config or env variables '''

    # Load configuration file and settings
    p = configargparse.ArgParser(default_config_files=['.config.yml', '~/.config.yml'])

    p.add('-c', '--config', required=False, is_config_file=True, help='config file path', env_var='APP_CONFIG')
    p.add('-g', '--group', required=False, help='default group name', env_var='APP_GROUP')
    p.add('-i', '--host', required=False, help='listening web ip', env_var='APP_HOST')
    p.add('-k', '--keepalive', required=False, type=int, help='default keepalive value in seconds', env_var='APP_KEEPALIVE')
    p.add('-p', '--port', required=False, help='listening web port', env_var='APP_PORT')
    p.add('--interval', required=False, type=int, help='llama collection interval in seconds', env_var='LLAMA_INTERVAL')
    p.add('--influxdb-host', required=False, help='InfluxDB Hostname', env_var='INFLUXDB_HOST')
    p.add('--influxdb-port', required=False, type=int, help='InfluxDB Port', env_var='INFLUXDB_PORT')
    p.add('--influxdb-name', required=False, help='InfluxDB Name, defaults to "enphase"', env_var='INFLUXDB_NAME')
    p.add('-v', '--verbose', help='verbose logging', action='store_true', env_var='APP_VERBOSE')

    config = p.parse_args()

    # Set defaults for webserver settings
    if not config.host:
        config.host = "127.0.0.1"
    if not config.port:
        config.port = "5000"
    if not config.interval:
        config.interval = 10

    # Set defaults for InfluxDB settings
    if config.influxdb_host:
        if not config.influxdb_port:
            config.influxdb_port = 8086
        if not config.influxdb_name:
            config.influxdb_name = "llama"

    # Set keepalive values, 3600 seconds if none is set
    if config.keepalive:
        # How many seconds before kicking probes from service discovery
        default_keepalive = config.keepalive
    else:
        # 86400 seconds = 1 day
        default_keepalive = 86400

    # Set a default registration group if one is not provided
    if config.group:
        # Load the default group from configuration variables.
        default_group = str(config.group)
    else:
        default_group = "none"

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