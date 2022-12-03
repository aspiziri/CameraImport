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
* <a href="https://questionary.readthedocs.io/en/stable/pages/installation.html">Questionary</a> - Needed to accept user input for how to run the script

### Installation

Getting ths script up and running after ensuring you have the pre-requiusites installed is quite simple:

1. Download both the `CameraImport.py` and `config.ini` files (or clone/pull the repo) and place them in the same directory somewhere on your machine
2. Open the `config.ini` file
3. Make sure the file formats to scan for are updated to include your cameras supported file types (this should largely be the same for every camera but the `rawFormat` value may change
```
imgFormat = jpg,JPG,png,PNG,gif,GIF,tif,TIF,jpeg,JPEG
videoFormat = mp4,MP4,avi,AVI,mov,MOV
rawFormat = arw,ARW
```

4. Make sure the source and destination locations are set to where your SD card/camera are connected and where you want to output to
```
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
<img src="https://i.imgur.com/SKVfksz.gif">

2. **Selective Import**<br />
In this mode a you will be presented with a list of all of the dates for which you have photo/videos from and you can select which specific dates you want to import from. You will then be asked for the directory name you want to save everything to and all images/videos from that specific date will be migrated into the same folder.<br /><br />
*Use case: You don't always clear your SD card after each shoot and you have already copied over most of it's contents in the past. Today you just want to import all the photos taken from your most recent shoot.*
<img src="https://i.imgur.com/3BU1w2i.gif">

<p align="right">(<a href="#readme-top">back to top</a>)</p>


## Sub-folder structure

In addition to importing your photos and videos based on date this script can also ensure each file type goes to a specific sub-folder based on the heirchy you want. 

As an example, with my camera I capture both JPG and RAW copies of every photo and I also record a lot of videos. When I import I want all the RAW assets to go into a "Capture" folder as that is the structure thet my editing software, Capture One, used. The JPGs I keep around for quick reference but they can live elsewhere. Similarly I want my videos all dumped into a single spot. Given that here is the my desired output structure:

```
2022-10-01
├── Capture
│   ├── 2022-10-01 14.24.26.ARW
│   ├── 2022-10-01 14.31.16.ARW
│   ├── 2022-10-01 14.32.36.ARW
├── Videos
│   ├── 2018-09-11 10.51.46 - 3840w - 23.976 - 12sec
│   ├── 2018-09-11 22.06.19 - 3840w - 23.976 - 83sec.MP4
├── 2022-10-01 14.24.26.jpg
├── 2022-10-01 14.31.16.jpg
├── 2022-10-01 14.32.36.jpg
```

To get this structure, I would use the following configuration:
```
imgRelativePath = /
videoRelativePath = /Videos/
rawRelativePath = /Capture/
```

This tells:
* "images" (jpg) to go to `\2022-10-01\`
* "videos" (mp4) to go to `\2022-10-01\Videos\`
* "raw" (ARW) to go to `\2022-10-01\Capture\`

## License

Distributed under the MIT License. See `LICENSE.txt` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>