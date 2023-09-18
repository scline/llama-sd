'''
Contains flask JSON schema for Flask calls.
'''

from helpers.config import load_conf


# Load configuration variable for defaults within scema's
config = load_conf()

# JSON schema for registering a remote llama probe
registration_schema = {
    'type': 'object',
    'properties': {
        'port': {'type': 'number'},
        'keepalive': {'type': 'number', "default": config.keepalive },
        'group': {'type': 'string', "default": config.group },
        'tags': {
            'type': 'object',
            'properties': {
                'version': {'type': 'string'},
                'probe_shortname': {'type': 'string'},
                'probe_name': {'type': 'string'}
            },
            'required': ['version', 'probe_shortname', 'probe_name'],
            'additionalProperties': True
        },
    },
    'required': ['port'],
}
