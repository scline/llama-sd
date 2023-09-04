'''
Contains common functions
'''

from datetime import datetime

def create_date():
    ''' Grab current date and time, reply in json format '''
    return {'create_date': datetime.now().strftime('%Y-%m-%dT%H:%M:%S.%f')}
