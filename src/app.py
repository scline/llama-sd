import flask, threading, logging, configargparse

from flask import request, jsonify
from flask_expects_json import expects_json
from datetime import datetime
from time import sleep

# Load configuration file and settings
p = configargparse.ArgParser(default_config_files=['.config.yml', '~/.config.yml'])
p.add('-c', '--config', required=False, is_config_file=True, help='config file path')
p.add('-p', '--port', required=False, help='listening web port', env_var='APP_PORT')
p.add('-i', '--host', required=False, help='listening web ip', env_var='APP_HOST')
p.add('-k', '--keepalive', help='default keepalive value in seconds', env_var='APP_KEEPALIVE')
p.add('-v', '--verbose', help='verbose logging', action='store_true', env_var='APP_VERBOSE')

config = p.parse_args()

app = flask.Flask(__name__)

# Set defaults for webserver settings
if not config.host:
    config.host = "127.0.0.1"
if not config.port:
    config.port = "5000"

# Set logging levels
if config.verbose:
    logging.basicConfig(format="%(asctime)s %(levelname)s %(threadName)s: %(message)s", encoding='utf-8', level=logging.DEBUG)
else:
    logging.basicConfig(format="%(asctime)s %(levelname)s %(threadName)s: %(message)s", encoding='utf-8', level=logging.INFO)

# Set keepalive values, 3600 seconds if none is set
if config.keepalive:
    # How many seconds before kicking probes from service discovery
    default_keepalive = int(config.keepalive)
else:
    # 86400 seconds = 1 day
    default_keepalive = 86400

# Debug logging for settings
logging.debug(p.format_values())
logging.debug(config)
logging.debug("Default keepalive is set to %i seconds" % default_keepalive)

# Global variable to lock threads as needed
thread_lock = threading.Lock()

# Initialize global dictionaries
database = {}
metrics = {}

# Expected JSON schema for registering a probe
schema = {
    'type': 'object',
    'properties': {
        'address': {'type': 'string'},
        'port': {'type': 'number'},
        'keepalive': {'type': 'number', "default": default_keepalive },
        'meta': {
            'type': 'object',
            'properties': {
                'version': {'type': 'string'}
            },
            'required': ['version']
        },
    },
    'required': ['id', 'address', 'port']
}


@app.route('/', methods=['GET'])
def home():
    return "<h1>Welcome Home!</h1><p>Generic HomePage</p>"


# TODO: Prometheus metrics and stats go here
@app.route('/metrics', methods=['GET'])
def get_metrics():
    # Registration Count
    # Registration Deletes
    # Database Size
    # Uptime
    # Memory Usage (is that easy?)
    return "Metrics page for prometheus in the future"


# Registration endpoint
@app.route("/api/v1/register", methods=['POST'])
@expects_json(schema, fill_defaults=True)
def add_entry():
    request_json = request.get_json()

    # Add create date to the json data
    request_json.update(create_date())

    # Turn "id" into a key to organize hosts, incert to database variable
    database[request_json['id']] = request_json

    return database


# A route to return all of the available entries
@app.route('/api/v1/list', methods=['GET'])
def api_list_all():
    return jsonify(database)


# Grab current date and time, reply in json format
def create_date():
    return {'create_date': datetime.now().strftime('%Y-%m-%dT%H:%M:%S.%f')}


# Background process that removes stale entries
def clean_stale_probes():
    global thread_lock
    # Run every 60 seconds
    while(not sleep(60)):
        # Aquire thread lock for variable work
        with thread_lock:
            logging.debug("Thread Locked!")
            logging.debug("%i probe entries are registered" % len(database))

            # Initialize list 
            remove_probe_list = []

            # Scann all probes in the inventory, remove those that have aged to long
            for probe in database:
                # Caclulate current time and creation date to seconds passed
                age = int((datetime.now() - datetime.strptime(database[probe]['create_date'], '%Y-%m-%dT%H:%M:%S.%f')).total_seconds())

                logging.debug("Probe '%s' checked in %i seconds ago" % (probe, age))
                if age > database[probe]['keepalive']:
                    logging.debug("Probe '%s' should be removed!" % probe)
                    remove_probe_list.append(probe)
            
            # Remove old probed from global database
            for item in remove_probe_list:
                database.pop(item, None)
            
            # Warning log on removed probes
            if remove_probe_list:
                logging.warning("Removed %i probe(s) due to aging - %s" % (len(remove_probe_list), str(remove_probe_list)))

        logging.debug("Thread Unlocked!")

# Gather application start time for metrics and data validation
metrics = { "start_time",  datetime.now() }

# Start background threaded process to clean stale probes
inline_thread_cleanup = threading.Thread(target=clean_stale_probes, name="CleanThread")
inline_thread_cleanup.start()

if __name__ == "__main__":
    # Te get flask out of development mode
    # https://stackoverflow.com/questions/51025893/flask-at-first-run-do-not-use-the-development-server-in-a-production-environmen
    from waitress import serve
    serve(app, host=config.host, port=config.port)
