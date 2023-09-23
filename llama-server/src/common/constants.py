'''
Contains application constants that should not change during a run
'''

from models.constants import ApiDefaults, FlaskDefaults, InfluxdbDefaults

# API registration defaults, 86400 seconds = 1 day
ApiDefaults = ApiDefaults(
    keepalive= 86400,
    interval=10,
    group="default")

# Flask webserver defaults
FlaskDefaults = FlaskDefaults(
    host="127.0.0.1",
    port=5000)

# Flask webserver defaults
InfluxdbDefaults = InfluxdbDefaults(
    name="llama",
    port=8086)
