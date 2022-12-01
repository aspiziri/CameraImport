<a name="readme-top"></a>


## CameraImport

This is a utility that makes sorting photos and videos into a structured format simple. Given a source folder it will rename and organize all of your media based on your chosen configuration.


## About The Project

I created this script to help me automate the tedious and repetative job of importing all of my photos and videos off my my camera after a shoot.

My goal was to have this accomplish the following:
* Automatically find and move all photos and videos from an SD card onto my machine
* Rename all of the files to date/time filenames for easy sorting in the future
* Put each of the files into specific sub-folders by type based on my archiving and workflow needs

<p align="right">(<a href="#readme-top">back to top</a>)</p>


## Getting Started

See below for details on how to get this script up and running and what options you have to customize it to your needs.

### Prerequisites

* <a href="https://www.python.org/downloads/">Python</a> - Needed to run the script
* <a href="https://ffmpeg.org/download.html">ffmpeg</a> - Needed for determining video attibutes used in the renaming process
* <a href="https://ffmpeg.org/download.html">ffprobe</a> - Needed for determining video attibutes used in the renaming process

### Installation

Getting ths script up and running after ensuring you have the pre-requiusites installed is quite simple:

1. Download both the `CameraImport.py` and `config.ini` files (or clone/pull the repo) and place them in the same directory somewhere on your machine
2. Open the `config.ini` file
3. Make sure the file formats to scan for are updated to include your cameras supported file types (this should largely be the same for every camera but the `rawFormat` value may change
```python
imgFormat = jpg,JPG,png,PNG,gif,GIF,tif,TIF,jpeg,JPEG
videoFormat = mp4,MP4,avi,AVI,mov,MOV
rawFormat = arw,ARW
```

4. Make sure the source and destination locations are set to where your SDc ar/camera are connected and where you want to output to
```python
source = E:
destination = D:\Media\Camera Import
``` 
5. Run the CameraImport script

<p align="right">(<a href="#readme-top">back to top</a>)</p>


## Operating modes

This script has two different modes that it can operate in. Simply run the script and it will ask you which mode you want to execute in.

1. **Full Import**<br />
This is the "one-click" mode which will grab every single photo/video on the device and automatically import them into date seperated folders.<br /><br />
*Use case: You want to dump all content from your camera in one go*
<img src="http://pulseblend.com/other/CameraImport/FullImport.gif">

2. **Selective Import**<br />
In this mode a you will be presented with a list of all of the dates for which you have photo/videos from and you can select which specific dates you want to import from. You will then be asked for the directory name you want to save everything to and all images/videos from that specific date will be migrated into the same folder.<br /><br />
*Use case: You don't always clear your SD card after each shoot and you have already copied over most of it's contents in the past. Today you just want to import all the photos taken from your most recent shoot.*
<img src="http://pulseblend.com/other/CameraImport/SelectiveImport.gif">

<p align="right">(<a href="#readme-top">back to top</a>)</p>


## License

Distributed under the MIT License. See `LICENSE.txt` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>