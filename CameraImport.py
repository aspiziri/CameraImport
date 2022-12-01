import configparser
import json
import os
import sys
import subprocess
import questionary
from pathlib import Path
from shutil import copy2
from datetime import datetime
from os import system

# Load in all config values
config = configparser.ConfigParser()
config.sections()
config.read('config.ini')

def printDynamic(output=""):
    #system("cls")
    sys.stdout.write(output + '%\r')
    sys.stdout.flush()

def getFileCount(importDates):
    i = 0
    p = Path(config['STRUCTURE']['source'])
    for p in p.glob('**/*.*'):
        format = p.name.split('.')
        if format[len(format)-1] in config['CORE']['imgFormat'] or format[len(format)-1] in config['CORE']['videoFormat'] or format[len(format)-1] in config['CORE']['rawFormat']:
            taken = datetime.fromtimestamp(os.path.getmtime(p)).strftime("%Y-%m-%d")
            if taken in importDates:
                i+=1
    return i

def runImport(importDates, totalFiles, folderName=None):
    i=0
    p = Path(config['STRUCTURE']['source'])
    for p in p.glob('**/*.*'):
        format = p.name.split('.')
        if format[len(format)-1] in config['CORE']['imgFormat'] or format[len(format)-1] in config['CORE']['videoFormat'] or format[len(format)-1] in config['CORE']['rawFormat']:
            base = config['STRUCTURE']['destination']
            taken = datetime.fromtimestamp(os.path.getmtime(p)).strftime("%Y-%m-%d")
            takenFull = datetime.fromtimestamp(os.path.getmtime(p)).strftime("%Y-%m-%d %H.%M.%S")

            if folderName is not None:
                base = base + "/" + folderName
            else:
                base = base + "/" + taken

            if folderName is None or taken in importDates:
                i+=1

                # Check if the folder exists and create it if not
                if format[len(format)-1] in config['CORE']['imgFormat'] or format[len(format)-1] in config['CORE']['videoFormat'] or format[len(format)-1] in config['CORE']['rawFormat']:
                    if not os.path.isdir(base):
                        printDynamic("(" + str(round((i/totalFiles)*100)) + "% - " + str(i) + "/" + str(totalFiles) + ") Creating directory " + config['STRUCTURE']['imgRelativePath'])
                        os.mkdir(base)

                # Check if the file is a JPG image
                if format[len(format)-1] in config['CORE']['imgFormat']:

                    # Check if the image folder exists and if not create it
                    if not os.path.isdir(base + config['STRUCTURE']['imgRelativePath']):
                        printDynamic("(" + str(round((i/totalFiles)*100)) + "% - " + str(i) + "/" + str(totalFiles) + ") Creating directory: " + config['STRUCTURE']['imgRelativePath'])
                        os.mkdir(base + config['STRUCTURE']['imgRelativePath'])

                    # Copy the file to its destination
                    o = 1
                    if(os.path.exists(base + config['STRUCTURE']['imgRelativePath'] + takenFull + "." + format[len(format)-1])):
                        while os.path.exists(base + config['STRUCTURE']['imgRelativePath'] + takenFull + "_" + str(o) + "." + format[len(format)-1]):
                            o+=1
                        printDynamic("(" + str(round((i/totalFiles)*100)) + "% - " + str(i) + "/" + str(totalFiles) + ") Copying file: " + config['STRUCTURE']['imgRelativePath'] + takenFull + "_" + str(o) + "." + format[len(format)-1])
                        copy2(p, base + config['STRUCTURE']['imgRelativePath'] + takenFull + "_" + str(o) + "." + format[len(format)-1])
                    else:
                        printDynamic("(" + str(round((i/totalFiles)*100)) + "% - " + str(i) + "/" + str(totalFiles) + ") Copying file: " + config['STRUCTURE']['imgRelativePath'] + takenFull + "." + format[len(format)-1])
                        copy2(p, base + config['STRUCTURE']['imgRelativePath'] + takenFull + "." + format[len(format)-1])

                # Check if the file is a RAW image
                if format[len(format)-1] in config['CORE']['rawFormat']:

                    # Check if the raw folder exists and if not create it
                    if not os.path.isdir(base + config['STRUCTURE']['rawRelativePath']):
                        printDynamic("(" + str(round((i/totalFiles)*100)) + "% - " + str(i) + "/" + str(totalFiles) + ") Creating directory: " + config['STRUCTURE']['rawRelativePath'])
                        os.mkdir(base + config['STRUCTURE']['rawRelativePath'])

                    # Copy the file to its destination
                    o = 1
                    if(os.path.exists(base + config['STRUCTURE']['rawRelativePath'] + takenFull + "." + format[len(format)-1])):
                        while os.path.exists(base + config['STRUCTURE']['rawRelativePath'] + takenFull + "_" + str(o) + "." + format[len(format)-1]):
                            o+=1
                        printDynamic("(" + str(round((i/totalFiles)*100)) + "% - " + str(i) + "/" + str(totalFiles) + ") Copying file: " + config['STRUCTURE']['rawRelativePath']  + takenFull + "_" + str(o) + "." + format[len(format)-1])
                        copy2(p, base + config['STRUCTURE']['rawRelativePath'] + takenFull + "_" + str(o) + "." + format[len(format)-1])
                    else:
                        printDynamic("(" + str(round((i/totalFiles)*100)) + "% - " + str(i) + "/" + str(totalFiles) + ") Copying file: " + config['STRUCTURE']['rawRelativePath']  + takenFull + "." + format[len(format)-1])
                        copy2(p, base + config['STRUCTURE']['rawRelativePath'] + takenFull + "." + format[len(format)-1])

                # Check if the file is a video
                if format[len(format)-1] in config['CORE']['videoFormat']:

                    # Check if the video folder exists and if not create it
                    if not os.path.isdir(base + config['STRUCTURE']['videoRelativePath']):
                        printDynamic("(" + str(round((i/totalFiles)*100)) + "% - " + str(i) + "/" + str(totalFiles) + ") Creating directory: " + config['STRUCTURE']['videoRelativePath'])
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
                    printDynamic("(" + str(round((i/totalFiles)*100)) + "% - " + str(i) + "/" + str(totalFiles) + ") Copying file: " + base + config['STRUCTURE']['videoRelativePath'] + nameDate + " - " + nameRes + " - " + nameFPS + " - " + nameDuration + "." + nameFormat)
                    copy2(p, base + config['STRUCTURE']['videoRelativePath'] + nameDate + " - " + nameRes + " - " + nameFPS + " - " + nameDuration + "." + nameFormat)

# Print intro
print("CameraImport v2.0")
print()
print("== Settings =====================================")
print()

# Ask what type of import to run
importType = questionary.select("What type of import would you like to perform?", choices=['Full import', 'Selective import']).ask()

# Iterate through all files in the source location
datesFound = []

p = Path(config['STRUCTURE']['source'])
for p in p.glob('**/*.*'):
    base = config['STRUCTURE']['destination']

    # Break out the file name so we can check if its format is within our list
    format = p.name.split('.')
    if format[len(format)-1] in config['CORE']['imgFormat'] or format[len(format)-1] in config['CORE']['videoFormat'] or format[len(format)-1] in config['CORE']['rawFormat']:
        fileDate = datetime.fromtimestamp(os.path.getmtime(p)).strftime("%Y-%m-%d")
        if fileDate not in datesFound:
            datesFound.append(fileDate)

# Sort the dates that we found
datesFound = sorted(datesFound)

if(importType == "Full import"):
    print()
    print("== Import =======================================")
    print()

    count = getFileCount(datesFound)
    runImport(datesFound, count)

    print("Complete                                                                                                                                  ")

elif(importType == "Selective import"):
    datesSelected = questionary.checkbox("Select dates to import", choices=datesFound).ask()
    folderName = questionary.text("What would you like to name the folder to import to").ask()

    print()
    print("== Import =======================================")
    print()

    count = getFileCount(datesSelected)
    runImport(datesSelected, count, folderName)

    print("Complete                                                                                                                                  ")