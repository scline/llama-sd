'''
Contains common functions
'''

from datetime import datetime

# Grab current date and time, reply in json format
def create_date():
    return {'create_date': datetime.now().strftime('%Y-%m-%dT%H:%M:%S.%f')}