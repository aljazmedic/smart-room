#!/usr/bin/python
from device_registry import app as device_registry_app

device_registry_app.run(host='0.0.0.0', port=80, use_reloader=True, threaded=True, debug=True)