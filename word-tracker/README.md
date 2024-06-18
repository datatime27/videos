
PLEASE NOTE THIS CODE IS PRETTY HACKY AND MAY NOT WORK OUT OF THE BOX.

# Word Tracker

Requirements to download transcripts:

* Set up Account with YouTube API Client 
https://developers.google.com/youtube/v3/quickstart/python

* Install YouTubeTranscriptApi 
https://pypi.org/project/youtube-transcript-api/


Requirements to build sentiment:

* Install nltk with Vader
https://www.nltk.org/install.html


Open up run-download.py and edit `CAPTIONS_DIR` and `YOUTUBE_CHANNEL_ID` to be the ID of the YouTube channel you want to download the captions for.
Then run:
```
> python run-download.py
```

This will open up a website and walk you through the process to give access to the script.
At some point it might say this is unsafe, but you'll have to accept it in order for the script to work.

Once it runs, it should generate your new directory with a .json file for every video in the channel named by its video id.


Even if you don't download anything, I have provided a few Mr Beast Transcripts in the mrbeast folder.
You can run these scripts to see data from the provided mrbeast files:
```
> python build-overview.py
> python build-scatter-plot.py
> python search.py
```
