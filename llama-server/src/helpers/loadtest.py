'''
Contains functions related to self load-testing
'''

import random
import requests
import logging


def generate_random_ip_address() -> str:
    ''' Generates a random IP address '''

    local_random = random.Random()
    return f"{local_random.randrange(1,255)}.{local_random.randrange(1,255)}.{local_random.randrange(1,255)}.{local_random.randrange(1,255)}"


def random_group() -> str:
    ''' Generates a random group name, this is out of a list of 8 '''

    local_random = random.Random()
    group_list = [
        "LoadTest_01",
        "LoadTest_02",
        "LoadTest_03",
        "LoadTest_04",
        "LoadTest_05",
        "LoadTest_06",
        "LoadTest_07",
        "LoadTest_08",]
    
    return f"{local_random.choices(group_list)}"


def loadtest_register_probe(config, keepalive: int) -> None:
    ''' Register a loadtest probe for X secconds, returns status code '''
    ip = generate_random_ip_address()

    payload = {
        'ip': ip,
        'group': random_group(),
        'port': 8100,
        'keepalive': keepalive,
        'tags': {
            'version': 'loadtest',
            'probe_shortname': ip,
            'probe_name': ip
        }
    }

    try:
        r = requests.post(f'http://127.0.0.1:{config.port}/api/v1/register', json=payload)
    except Exception as e:
        logging.error(e)
