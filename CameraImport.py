#! python3
import configparser
import json
import os
import subprocess
from pathlib import Path
from shutil import copy2
from datetime import datetime

# Load in all config values
config = configparser.ConfigParser()
config.sections()
config.read('config.ini')

# Iterate through all files in the source location
p = Path(config['STRUCTURE']['source'])
for p in p.glob('**/*.*'):
	base = config['STRUCTURE']['destination']

	# Break out the file name so we can check if its format is within our list
	format = p.name.split('.')
	if format[len(format)-1] in config['CORE']['imgFormat'] or format[len(format)-1] in config['CORE']['videoFormat'] or format[len(format)-1] in config['CORE']['rawFormat']:

		# Check if the date folder is requested
		if config['STRUCTURE']['directoryPerDate']:
			taken = datetime.fromtimestamp(os.path.getmtime(p)).strftime("%Y-%m-%d")
			base = base + "/" + taken

			# Check if the date folder exists and if not create it
			if not os.path.isdir(config['STRUCTURE']['destination'] + '/' + taken):
				print("Creating directory " + config['STRUCTURE']['destination'] + '/' + taken)
				os.mkdir(config['STRUCTURE']['destination'] + '/' + taken)

	# Check if the file is an image
	if format[len(format)-1] in config['CORE']['imgFormat']:

		# Check if the image folder exists and if not create it
		if not os.path.isdir(base + config['STRUCTURE']['imgRelativePath']):
			print("Creating directory " + base + config['STRUCTURE']['imgRelativePath'])
			os.mkdir(base + config['STRUCTURE']['imgRelativePath'])

		# Copy the file to its destination
		print("Copying file: " + base + config['STRUCTURE']['imgRelativePath'] + p.name)
		copy2(p, base + config['STRUCTURE']['imgRelativePath'] + p.name)

	# Check if the file is a video
	if format[len(format)-1] in config['CORE']['videoFormat']:

		# Check if the video folder exists and if not create it
		if not os.path.isdir(base + config['STRUCTURE']['videoRelativePath']):
			print("Creating directory " + base + config['STRUCTURE']['videoRelativePath'])
			os.mkdir(base + config['STRUCTURE']['videoRelativePath'])

		# Call ffprobe to pull details about the video
		returned_data = subprocess.check_output(['ffprobe', '-v', 'quiet', '-print_format', 'json', '-show_format', '-show_streams', str(p.resolve())])

		# Loading of the json file
		data = json.loads(returned_data.decode('utf-8'))
		t = (data["streams"][0]["avg_frame_rate"])
		fps = [float(x) for x in t.split('/')]

		# Generate the file name fragments
		nameDate = datetime.fromtimestamp(os.path.getmtime(p)).strftime("%Y-%m-%d %H.%M.%S")
		nameFPS = "{:.3f}".format(fps[0] / fps[1])
		nameRes = str(data["streams"][0]["width"]) + "w"
		nameDuration = str(round(float(data["streams"][0]["duration"]))) + "sec"
		nameFormat = format[len(format)-1]

		# Copy the file to its destination
		print("Copying file: " + base + config['STRUCTURE']['videoRelativePath'] + nameDate + " - " + nameRes + " - " + nameFPS + " - " + nameDuration + "." + nameFormat)
		copy2(p, base + config['STRUCTURE']['videoRelativePath'] + nameDate + " - " + nameRes + " - " + nameFPS + " - " + nameDuration + "." + nameFormat)

	if format[len(format)-1] in config['CORE']['rawFormat']:

		# Check if the raw folder exists and if not create it
		if not os.path.isdir(base + config['STRUCTURE']['rawRelativePath']):
			print("Creating directory " + base + config['STRUCTURE']['rawRelativePath'])
			os.mkdir(base + config['STRUCTURE']['rawRelativePath'])

		# Copy the file to its destination
		print("Copying file: " + base + config['STRUCTURE']['rawRelativePath'] + p.name)
		copy2(p, base + config['STRUCTURE']['rawRelativePath'] + p.name)