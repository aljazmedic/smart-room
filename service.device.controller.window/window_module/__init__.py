#!/usr/bin/python
import markdown
from flask_restful import Resource, Api, reqparse
import os, shelve, requests
from flask import Flask
from redis import Redis
from flask import current_app, g
from flask.cli import with_appcontext


#Create instance
app = Flask(__name__)
api = Api(app)

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

class ModuleResource(Resource):
	def get(self, identifier):
		#shelf = get_db()
		device = requests.get(f"http://192.168.99.100:7001/device/{identifier}").json()
		print(device["data"]["controller_gateway"])
		gateway = device["data"]["controller_gateway"]
		data = requests.get(f"http://{gateway}").json()
		return {'message':f"GET on {identifier}", 'data':data["data"]}, 200

api.add_resource(ModuleResource, '/<string:identifier>')