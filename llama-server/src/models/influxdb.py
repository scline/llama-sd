'''
Contains models and dataclasses for InfluxDB
'''

from dataclasses import dataclass

@dataclass
class InfluxServerSettings:
    ''' Class for modeling InfluxDB server settings '''
    host: str
    port: int
    database: str
    verify_ssl: bool

@dataclass
class InfluxDataPoint:
    ''' Class for modeling InfluxDB metric point '''
    measurement: str
    tags: dict
    fields: dict
