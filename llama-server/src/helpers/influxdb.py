'''
Contains functions related to Influx database
'''

import socket
import logging
from influxdb import InfluxDBClient

from models.influxdb import InfluxDataPoint


def write_influx(client: InfluxDBClient, points) -> None:
    ''' Function to write server metrics to influxDB '''

    # Attempt to write metrics to InfluxDB
    try:
        client.write_points(points)
    except Exception as e:
        logging.error("Error writing to InfluxDB")
        logging.error(e)
        return None

    # Log how many metrics we wrote
    logging.info("Wrote %i metrics to influxDB" % len(points))


# Create the InfluxDB if one does not already exsist
def setup_influx(config) -> InfluxDBClient or None:
    ''' Setup influxDB database '''
    # Setup InfluxDB client
    client = InfluxDBClient(
        host=config.influxdb_host,
        port=config.influxdb_port,
        database=config.influxdb_name,
        verify_ssl=False)

    # Create database if it does not exsist
    try:
        client.create_database(config.influxdb_name)
        return client
    except Exception as e:
        logging.error("Error creating InfluxDB Database, please verify one exsists")
        logging.error(e)
        return None


def metrics_log_point(metrics) -> InfluxDataPoint:
    ''' Format metrics into something InfluxDB can use '''
    try:
        hostname = socket.gethostname()
        ipaddress = socket.gethostbyname(hostname)
    # We dont really care at the moment if this does not return
    except:
        hostname = "localhost"
        ipaddress = "127.0.0.1"

    return [{
        "measurement": "llama_server",
        "tags": {
            "hostname": hostname,
            "ipaddress": ipaddress
        },
        "fields": {
            "probe_count_removed": int(metrics["probe_count_removed"]),
            "probe_count_active": int(metrics["probe_count_active"]),
            "group_count_active": int(metrics["group_count_active"]),
            "group_count_removed": int(metrics["group_count_removed"]),
            "database_size_bytes": int(metrics["database_size_bytes"]),
            "clean_runtime": float(metrics["clean_runtime"]),
            "uptime": int(metrics["uptime"]),
        }
    }]