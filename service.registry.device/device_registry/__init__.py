#!/usr/bin/python
import markdown
from flask_restful import Resource, Api, reqparse
import os, shelve
from flask import Flask
from redis import Redis
from flask import current_app, g
from flask.cli import with_appcontext

#Create instance
app = Flask(__name__)
api = Api(app)
redis = Redis(host='redis', port=6379)

def get_db():
	db = getattr(g, '_database', None)
	if db is None:
		db = g._database = shelve.open("devices.db")
	return db

@app.teardown_appcontext
def teardown_db(exception=None):
	db = g.pop('_database', None)
	if db is not None:
		db.close()


@app.route("/")
def index():
	"""Documentation"""
	with open(os.path.dirname(app.root_path) + '/README.md', 'r') as markdown_file:
		content = markdown_file.read()

		return markdown.markdown(content)

class DeviceList(Resource):
	def get(self):
		shelf = get_db()
		keys = list(shelf.keys())
		devices = []

		for key in keys:
			devices.append(shelf[key])

		return {
			"message":'Success',
			"data":devices
		}, 200
	def post(self):
		parser = reqparse.RequestParser()

		parser.add_argument('identifier', type=str, required=True)
		parser.add_argument('name', type=str, required=True)
		parser.add_argument('device_type', type=int, required=True)
		parser.add_argument('controller_name', type=str, required=True)
		parser.add_argument('controller_gateway', type=str, required=True)

		args = parser.parse_args()

		shelf = get_db()
		shelf[args['identifier']] = args

		msg = 'Device registered!' if args['identifier'] not in shelf else 'Device overwritten!'

		return {'message':msg, 'data':args}, 201

class Device(Resource):
	def get(self, identifier):
		shelf = get_db()

		if not (identifier in shelf):
			return {'message':'Device no found.', 'data':{}}, 404
		return {'message':'Device found', 'data':shelf[identifier]}, 200

	def delete(self, identifier):
		shelf = get_db()

		if not (identifier in shelf):
			return {'message':'Device no found.', 'data':{}}, 404
		del shelf[identifier]
		return '', 204

api.add_resource(DeviceList, '/devices')
api.add_resource(Device, '/device/<string:identifier>')