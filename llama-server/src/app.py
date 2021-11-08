import json, flask, threading, logging, configargparse

from flask import request, jsonify, render_template, send_file
from flask_expects_json import expects_json
from datetime import datetime
from pympler import asizeof
from time import sleep

# Load configuration file and settings
p = configargparse.ArgParser(default_config_files=['.config.yml', '~/.config.yml'])
p.add('-c', '--config', required=False, is_config_file=True, help='config file path', env_var='APP_CONFIG')
p.add('-g', '--group', required=False, help='default group name', env_var='APP_GROUP')
p.add('-i', '--host', required=False, help='listening web ip', env_var='APP_HOST')
p.add('-k', '--keepalive', required=False, type=int, help='default keepalive value in seconds', env_var='APP_KEEPALIVE')
p.add('-p', '--port', required=False, help='listening web port', env_var='APP_PORT')
p.add('--interval', required=False, type=int, help='llama collection interval in seconds', env_var='LLAMA_INTERVAL')
p.add('-v', '--verbose', help='verbose logging', action='store_true', env_var='APP_VERBOSE')

config = p.parse_args()

app = flask.Flask(__name__)

# Set defaults for webserver settings
if not config.host:
    config.host = "127.0.0.1"
if not config.port:
    config.port = "5000"
if not config.interval:
    config.interval = 60

# Set logging levels
if config.verbose:
    logging.basicConfig(format="%(asctime)s %(levelname)s %(threadName)s: %(message)s", encoding='utf-8', level=logging.DEBUG)
else:
    logging.basicConfig(format="%(asctime)s %(levelname)s %(threadName)s: %(message)s", encoding='utf-8', level=logging.INFO)

# Set keepalive values, 3600 seconds if none is set
if config.keepalive:
    # How many seconds before kicking probes from service discovery
    default_keepalive = config.keepalive
else:
    # 86400 seconds = 1 day
    default_keepalive = 86400

# Set a default registration group if one is not provided
if config.group:
    # Load the default group from configuration variables.
    default_group = str(config.group)
else:
    default_group = "none"

# Debug logging for settings
logging.debug(p.format_values())
logging.debug(config)
logging.debug("Default keepalive is set to %i seconds" % default_keepalive)
logging.debug("Default group is set to '%s'" % default_group)

# Global variable to lock threads as needed
thread_lock = threading.Lock()

# Initialize global dictionaries
database = {}
metrics = {}

# Expected JSON schema for registering a probe
schema = {
    'type': 'object',
    'properties': {
        'port': {'type': 'number'},
        'keepalive': {'type': 'number', "default": default_keepalive },
        'group': {'type': 'string', "default": default_group },
        'tags': {
            'type': 'object',
            'properties': {
                'version': {'type': 'string'},
                'probe_shortname': {'type': 'string'},
                'probe_name': {'type': 'string'}
            },
            'required': ['version', 'probe_shortname', 'probe_name']
        },
    },
    'required': ['port']
}


@app.route('/', methods=['GET'])
def home():
    return "<h1>Welcome Home!</h1><p>Generic HomePage</p>"


# TODO: Prometheus or JSON metrics and stats go here
@app.route('/metrics', methods=['GET'])
def get_metrics():
    return jsonify(metrics), 200


# Returns IP address of requester, for finding NAT/Public address in the future
@app.route("/api/v1/my_ip_address", methods=['GET'])
def my_ip_address():
    return jsonify({'ip': request.remote_addr}), 200


# Interval metric
@app.route("/api/v1/interval", methods=['GET'])
def interval():
    return str(config.interval), 200

# Registration endpoint
@app.route("/api/v1/register", methods=['POST'])
@expects_json(schema, fill_defaults=True)
def add_entry():
    request_json = request.get_json()

    # Add create date to the json data
    request_json.update(create_date())

    # If IP address was not set, try to figure out what it should be. "If key NOT IN dict"
    if not request_json.get('ip'):
        # Add requestor IP address to the json data
        request_json.update({'ip': '%s' % request.remote_addr})

    # Formulate probe ID by "IP:Port", ex: "192.168.1.12:8100"
    request_json.update({'id': '%s:%s' % (request_json['ip'], request_json['port'])})

    logging.debug("Registration Update: '%s'" % request_json['id'])

    # Wait for thread lock in the event cleanup is running
    thread_lock.acquire()

    # If the group does not exsist, create it as a source key in the database
    if request_json['group'] not in database:
        database[request_json['group']] = {}

    # Turn "id" into a key to organize hosts, incert to database variable
    database[request_json['group']][request_json['id']] = request_json
    # Release thread lock
    thread_lock.release()

    return database


# A route to return all of the available entries
@app.route('/api/v1/list', methods=['GET'])
def api_list_all():
    return jsonify(database), 200


# A route to return a list of hosts that the scraper will collect from
@app.route('/api/v1/scraper', methods=['GET'])
def api_scraper():
    hosts = []
    # Cycle through all groups to formulate a list
    for group in database:
        for host in database[group]:
            hosts.append(database[group][host]['ip'])
    logging.debug("Scraper Host List: %s" % hosts)


    # If list is empty, no hosts have joined, return 127.0.0.1
    if len(hosts) == 0:
        hosts.append("127.0.0.1")

    # Turn the host LIST into a comma separated string
    joined_string = ",".join(hosts)

    return render_template("scraper.j2", hosts=joined_string)


# A route to return LLAMA Collector config file via template
@app.route('/api/v1/config/<group>', methods=['GET'])
def api_config(group):
    # Create a temporary database for data manipulation, we dont want this perm
    #database_tmp = database.copy()
    # Apparently thread safe way to perform a full database copy (not shadow)
    database_tmp = json.loads(json.dumps(database))

    if group in database_tmp:
        # If requesting probe is in the group list, change target IP to 127.0.0.1
        port = request.args.get('llamaport', None)
        reported_source_ip = request.args.get('srcip', None)
        remote_ip_address = request.remote_addr

        # If a Source IP was provided, do things (hack around certain NAT scenario)
        if reported_source_ip:
            logging.info("CONFIG: '%s' says its IP is '%s'" % (request.remote_addr, reported_source_ip))
            # TODO: Verify reported IP address is a valid ipv4 address
            remote_ip_address = reported_source_ip

        # Log if a client is not sending what port it has assigned
        if not port:
            logging.error("No port was given from probe '%s' when generating configuration" % remote_ip_address) 
            port = "null"
        
        # Store probe ID "IP_ADDRESS:PORT"
        requesting_probe_id = "%s:%s" % (remote_ip_address, port)
        logging.debug("Config request from '%s'" % requesting_probe_id)

        # Check if key not in dict python
        if requesting_probe_id not in database_tmp[group]:
            logging.error("Requesting probe '%s' has not registered, no config will be given" % requesting_probe_id)
            return jsonify({'error': "unknown probe '%s', please register" % requesting_probe_id}), 404

        # Setup source + destination pairs per probe
        for remote_id in database_tmp[group]:

            # Rewrite self to 127.0.0.1
            if remote_id == requesting_probe_id:
                database_tmp[group][requesting_probe_id]["ip"] = "127.0.0.1"
                logging.debug("Local probe translation to 127.0.0.1 - %s" % requesting_probe_id)

                database_tmp[group][requesting_probe_id]["tags"]["dst_name"] = database_tmp[group][requesting_probe_id]["tags"]["probe_name"]
                database_tmp[group][requesting_probe_id]["tags"]["dst_shortname"] = database_tmp[group][requesting_probe_id]["tags"]["probe_shortname"]
                database_tmp[group][requesting_probe_id]["tags"]["src_name"] = database_tmp[group][requesting_probe_id]["tags"]["probe_name"]
                database_tmp[group][requesting_probe_id]["tags"]["src_shortname"] = database_tmp[group][requesting_probe_id]["tags"]["probe_shortname"]
                database_tmp[group][requesting_probe_id]["tags"]["group"] = group
                pass

            database_tmp[group][remote_id]["tags"]["dst_name"] = database_tmp[group][remote_id]["tags"]["probe_name"]
            database_tmp[group][remote_id]["tags"]["dst_shortname"] = database_tmp[group][remote_id]["tags"]["probe_shortname"]
            database_tmp[group][remote_id]["tags"]["src_name"] = database_tmp[group][requesting_probe_id]["tags"]["probe_name"]
            database_tmp[group][remote_id]["tags"]["src_shortname"] = database_tmp[group][requesting_probe_id]["tags"]["probe_shortname"] 
            database_tmp[group][remote_id]["tags"]["group"] = group

            # Remove the "probe" name entries, we dont need to send those        
            #database_tmp[group][remote_id]["tags"].pop("probe_name", None)   
            #database_tmp[group][remote_id]["tags"].pop("probe_shortname", None)   

        logging.debug(database_tmp[group])
        return render_template("config.yaml.j2", template_data=database_tmp[group], template_interval=config.interval)

    # If group is not located, error
    logging.error("'/api/v1/config/%s' - Unknown group" % group)
    return jsonify({'error': "unknown group '%s'" % group}), 404


# A route to return a certain group of the available entries
@app.route('/api/v1/list/<group>', methods=['GET'])
def api_list_group(group):
    if group in database:
        return jsonify(database[group]), 200

    # If group is not located, error
    logging.error("'/api/v1/list/%s' - Unknown group" % group)
    return jsonify({'error': "unknown group '%s'" % group}), 404


# Grab current date and time, reply in json format
def create_date():
    return {'create_date': datetime.now().strftime('%Y-%m-%dT%H:%M:%S.%f')}


# Process that checks if there are duplicate probes in a group
def is_probe_dup(group, probe):
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


# Background process that removes stale entries
def clean_stale_probes():
    # Run every 60 seconds
    while(not sleep(60)):
        # Get start time for runtime metrics
        start_time = datetime.now().timestamp()

        # Aquire thread lock for variable work
        with thread_lock:
            logging.debug("Thread Locked!")

            # Initialize list 
            remove_probe_list = []
            remove_group_list = []

            # Initialize metric
            remove_probe_count = 0

            # Go through all groups for stale probes
            for group in database:
                # Scann all probes in the inventory, remove those that have aged to long
                for probe in database[group]:
                    # Caclulate current time and creation date to seconds passed
                    age = int((datetime.now() - datetime.strptime(database[group][probe]['create_date'], '%Y-%m-%dT%H:%M:%S.%f')).total_seconds())

                    logging.debug("Probe '%s' in group '%s' checked in %i seconds ago" % (probe, group, age))
                    if age > database[group][probe]['keepalive']:
                        logging.debug("Probe '%s' in group '%s' should be removed!" % (probe, group))
                        remove_probe_list.append(probe)

                    # If there is a duplicate entry mark for deletetion
                    if is_probe_dup(group, probe):
                        remove_probe_list.append(probe)
            
                # Remove old probed from global database
                for item in remove_probe_list:
                    database[group].pop(item, None)

                # Add to metric counter
                remove_probe_count = len(remove_probe_list) + remove_probe_count

                # Clear list
                remove_probe_list = []

            # Warning log on removed probes
            if remove_probe_count > 0:
                logging.warning("Removed %i probe(s) due to aging" % remove_probe_count)

            # If a group is empty add to removal list
            for group in database:
                if not database[group]:
                    logging.warning("Group '%s' is empty, removing it." % group)
                    remove_group_list.append(group)

            # Remove empty groups from global database
            for item in remove_group_list:
                database.pop(item, None)

            # Lets collect and crunch some metrics here
            global metrics

            # Calculate the number of active nodes 
            node_count = 0
            for group in database:
                node_count = node_count + len(database[group])
            logging.info("%i active probe(s) are registered" % node_count)
            
            # Write metrics
            metrics["probe_count_removed"] = remove_probe_count
            metrics["probe_count_active"] = node_count
            metrics["group_count_active"] = len(database)
            metrics["group_count_removed"] = len(remove_group_list)
            metrics["database_size_bytes"] = asizeof.asizeof(database)
            metrics["clean_runtime"] = datetime.now().timestamp() - start_time
            metrics["uptime"] = datetime.now().timestamp() - metrics["start_time"].timestamp()
            metrics["metrics_timestamp"] = datetime.now()

        logging.debug("Thread Unlocked!")


if __name__ == "__main__":
    # Gather application start time for metrics and data validation
    metrics["start_time"] = datetime.now()

    # Start background threaded process to clean stale probes
    inline_thread_cleanup = threading.Thread(target=clean_stale_probes, name="CleanThread")
    inline_thread_cleanup.start()

    logging.info("Flask server started on '%s:%s'" % (config.host, config.port))

    # Te get flask out of development mode
    # https://stackoverflow.com/questions/51025893/flask-at-first-run-do-not-use-the-development-server-in-a-production-environmen
    from waitress import serve
    serve(app, host=config.host, port=config.port)
