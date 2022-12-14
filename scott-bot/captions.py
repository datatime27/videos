
import sys
import os
import re
import json
from collections import defaultdict
from pprint import pprint

CAPTIONS_DIR = "captions"
VIDEOS_DIR = "videos"
PUNCTUATION = '-;,.\'\"\!\?,'
PUNCTUATION_RE = re.compile('['+PUNCTUATION+']')

class WordRef:
    def __init__(self):
        self.word = None
        self.start = None
        self.duration = None
        self.text = None
        self.videoId = None
        self.title = None
        self.publishedAt = None
        self.prev_ref = None
        self.next_ref = None
        self.link = None
        self.last_in_sentence = False
    
    def __str__(self):
        return self.__repr__()
    def __repr__(self):
        return str({
            'word': self.word,
            'start': self.start,
            'duration': self.duration,
            'text' : self.text,
            'videoId' : self.videoId,
            'title' : self.title,
            'publishedAt' : self.publishedAt,
            'link': self.link,
            'last_in_sentence': self.last_in_sentence,
            })
            
class Parser:        
    def parse(self, cut_off_date):
        self.word_index = defaultdict(list)
        self.first_word_in_video = {}
        self.videos_to_download = set()
        self.video_counter = 0
        self.total_caption_time_secs = 0.0
        
        for filename in os.listdir(CAPTIONS_DIR):
            filepath = os.path.join(CAPTIONS_DIR, filename)
            with open(filepath) as f:
                captions = json.load(f)
                title = captions['snippet']['title'].replace('&#39;',"'")
                publishedAt = captions['snippet']['publishedAt']
                if publishedAt < cut_off_date:
                    continue
                self.video_counter += 1

                #print publishedAt, title.encode('utf-8'), filepath

                videoId = captions['id']['videoId']

            cues = self.parseCaptions(captions)
            self.buildIndex(videoId, title.encode('utf-8'), publishedAt, cues)
            
        #print self.video_counter
        #print self.total_caption_time_secs

    def parseCaptions(self, captions):
        cues = []
        for action in captions["actions"]:
            for cueGroup in action["updateEngagementPanelAction"]["content"]["transcriptRenderer"]["body"]["transcriptBodyRenderer"]["cueGroups"]:
                for cue in cueGroup["transcriptCueGroupRenderer"]["cues"]:
                    duration = int(cue["transcriptCueRenderer"]["durationMs"])
                    cues.append({
                    'start' : int(cue["transcriptCueRenderer"]["startOffsetMs"]),
                    'duration' : duration,
                    'text' : cue["transcriptCueRenderer"]["cue"]["simpleText"].replace('&#39;',"'").replace('\n',' ').encode('utf-8'),
                    })
                self.total_caption_time_secs += duration/1000.0
                    
        return cues
        
    def buildIndex(self, videoId, title, publishedAt, cues):
        prev_ref = None
        for cue in cues:
            words = cue['text'].split()
            for word in words:
                last_in_sentence = False
                if re.search(PUNCTUATION_RE, word):
                    last_in_sentence = True
                    
                word = word.lower().strip(PUNCTUATION)
                if not word:
                    continue
                ref = WordRef()
                ref.word = word
                ref.start = cue['start']
                ref.duration = cue['duration']
                ref.text = cue['text']
                ref.videoId = videoId
                ref.title = title
                ref.publishedAt = publishedAt
                ref.prev_ref = prev_ref
                ref.link = 'https://www.youtube.com/watch?v=%s&t=%ds' % (videoId, cue['start']/1000)
                ref.last_in_sentence = last_in_sentence

                if prev_ref:
                    prev_ref.next_ref = ref
                else: # first ref for this video
                    self.first_word_in_video[videoId] = ref
                self.word_index[word].append(ref)
                prev_ref = ref            
            
    def findWord(self, word):
        if word not in self.word_index:
            print "'%s' not found" % (word)
            return []
            
        return sorted(self.word_index[word])

    def findWords(self, words, only_last_in_sentence=False):
        self.videos_to_download = set()
        
        first_word = words[0].lower()
        if first_word not in self.word_index:
            print "'%s' not found" % (first_word)
            return []
            
        matches = {}
        refs = self.word_index[first_word]
        for ref in refs:
            if self.matchNextWords(ref, words[1:], only_last_in_sentence):
                id = (ref.videoId,ref.start)
                matches[id] = ref
                
                video_file = os.path.join(VIDEOS_DIR, ref.videoId + '.mp4')
                if not os.path.exists(video_file):
                    self.videos_to_download.add(ref.videoId)
                        
        return sorted(matches.values())
        
    def matchNextWords(self, ref, words, only_last_in_sentence):
        for word in words:
            word = word.lower()
            if ref.next_ref == None or ref.next_ref.word != word:
                return False
            ref = ref.next_ref
            
        # If we don't care about only finding the phrase at the end of a sentence then return true.
        if not only_last_in_sentence: 
            return True
        
        # We only want the phrase if it occurs at the end of a sentence.
        return ref == None or ref.last_in_sentence
        
    def printVideosToDownload(self):
        for videoId in self.videos_to_download:
            # You need to download youtube-dl.exe and then you can copy and paste this command into the command line.
            print 'youtube-dl.exe -o "'+VIDEOS_DIR+'/%(id)s.%(ext)s" https://www.youtube.com/watch?v='+videoId
