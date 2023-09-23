'''
Contains models and dataclasses for Constants
'''

from dataclasses import dataclass

@dataclass
class ApiDefaults:
    ''' Class for modeling default values for API registration '''
    keepalive: int
    interval: int
    group: str

@dataclass
class FlaskDefaults:
    ''' Class for modeling Flask web server default values '''
    host: str
    port: int

@dataclass
class InfluxdbDefaults:
    ''' Class for modeling InfluxDB default values '''
    name: str
    port: int
