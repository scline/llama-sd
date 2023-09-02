'''
Contains flask JSON schema for Flask calls.
'''

# TODO this needs to leave after var refactor
from common.constants import default_keepalive, default_group

# JSON schema for registering a remote llama probe
registration_schema = {
    'type': 'object',
    'properties': {
        'port':      {'type': 'number'},
        'keepalive': {'type': 'number', "default": default_keepalive },
        'group':     {'type': 'string', "default": default_group },
        'tags': {
            'type': 'object',
            'properties': {
                'version':         {'type': 'string'},
                'probe_shortname': {'type': 'string'},
                'probe_name':      {'type': 'string'}
            },
            'required': ['version', 'probe_shortname', 'probe_name']
        },
    },
    'required': ['port']
}
