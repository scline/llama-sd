'''
Contains functions related to llama probes
'''
import logging


def is_probe_dup(group, probe, database) -> bool:
    ''' Process that checks if there are duplicate probes in a group '''
    for probe_compaire in database[group]:

        # Ignore self entry
        if probe == probe_compaire:
            continue

        # If two different probes habe the same shortname, return True
        if database[group][probe]["tags"]["probe_shortname"] == database[group][probe_compaire]["tags"]["probe_shortname"]:
            logging.error("Duplicate probe entry found in Group: '%s', ID: '%s', probe_shortname: '%s'" % (group, probe, database[group][probe]["tags"]["probe_shortname"]))
            return True
    
    # No duplicates found
    return False
