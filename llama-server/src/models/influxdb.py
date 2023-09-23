'''
Contains models and dataclasses for InfluxDB
'''

from dataclasses import dataclass

@dataclass
class InfluxDataPoint:
    ''' Class for modeling InfluxDB metric point '''
    measurement: str
    tags: dict
    fields: dict
