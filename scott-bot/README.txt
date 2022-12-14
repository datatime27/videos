
PLEASE NOTE THIS CODE IS PRETTY HACKY AND PROBABLY WON'T WORK OUT OF THE BOX.
PLEASE USE THIS AS REFERENCE TO BUILD YOUR OWN CODE.

Requirements to build "The Scott Bot":

* Make sure you have Adobe Premiere and ExtendScript Toolkit

* Set up Account with YouTube API Client

* Download youtube-dl.exe for downloading the actual YouTube video files

Create a root directory to run all commands inside of.

Open up download.py and edit YOUTUBE_CHANNEL_ID to be the ID of the YouTube channel you want to download the captions for.
RUN:
> python download.py

This should generate a 'captions' directory with a .json file for every video in the channel named by its video id.

Open buildCut.py and modify PRPROJ_BIN_NAME to be the name you want the project bin to be in premiere (make sure you create your project bin folder in Premiere first).
buildCut.py will parse all of the .json each time you run it in order to build the database and find your desired words.
buildCut.py needs the oldest start time you want to look back to.

Here are some examples:
> python buildCut.py 2022-01-01 going to
> python buildCut.py 2022-01-01 data

Add the -e option to only find phrase and the end of sentences (they tend to be better for editing)

buildCut.py prints the youtube-dl.exe commands that you will need to run in order to download the actual video files from YouTube.
YOU NEED TO DOWNLOAD THE VIDEO FILES. THIS CODE DOES NOT DO THAT AUTOMATICALLY.

buildCut.py should also write a .jsx file in your 'cuts' directory which looks something like _going_to_.jsx

Once the .jsx is written, and you have downloaded the video files, you can load the .jsx file into Adobe ExtendScript Toolkit and run it against your Premiere timeline.


Project Bin Folder:
Make sure you have created a Project bin folder in Premiere that matches the name you set in PRPROJ_BIN_NAME.
The .jsx files should ask you to load in the video file of any video that's not already in your project bin folder.

I have also provided extendscript-loadbin.jsx
This is a convenience script that you can run in ExtendScript 
to load all of the files in your 'videos' directory into a Premiere project bin folder.
If you don't pre load the video files into Premiere with extendscript-loadbin.jsx that is fine, they will get loaded into Premiere when you run the .jsx files.
