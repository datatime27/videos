
import sys
import os
import re
import json
from collections import defaultdict
from pprint import pprint

PUNCTUATION_RE = re.compile('[-;,.\'\"\\!\\?,]')
DICTIONARY_FILE = 'words_dictionary.json'

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
    def parse(self, directory, cut_off_date=None):
        self.word_index = defaultdict(list)
        self.first_word_in_video = {}
        self.video_lengths = []
        self.video_counter = 0
        self.video_stats = {}
        self.total_caption_time_secs = 0.0
        with open(DICTIONARY_FILE) as f:
            self.words_dictionary = json.load(f)
        
        filenames = os.listdir(directory)
        for filename in filenames:
            filepath = os.path.join(directory, filename)
            with open(filepath) as f:
                obj = json.load(f)
                title = obj['title'].replace('&#39;',"'")
                publishedAt = obj['publishedAt']
                if cut_off_date and publishedAt < cut_off_date:
                    continue
                self.video_counter += 1

                #print (self.video_counter,publishedAt, title.encode('utf-8'), filepath)

                videoId = obj['id']
                captions = obj['captions']
                self.video_stats[videoId] = obj['stats']
                self.buildIndex(videoId, title.encode('utf-8'), publishedAt, captions)
                
                if len(captions) > 10:
                    length = captions[-1]['start'] + captions[-1]['duration']
                    self.video_lengths.append((length,videoId))
            
        self.video_lengths.sort()
        #print self.video_lengths
        #print self.video_counter
        #print self.total_caption_time_secs

        
    def buildIndex(self, videoId, title, publishedAt, captions):
        prev_ref = None
        for caption in captions:
            words = caption['text'].split()
            for word in words:
                last_in_sentence = False
                if re.search(PUNCTUATION_RE, word):
                    last_in_sentence = True
                    
                word = re.sub("[\\W]",'',word.lower())
                if not word:
                    continue
                ref = WordRef()
                ref.word = word
                ref.start = caption['start']
                ref.duration = caption['duration']
                ref.text = caption['text']
                ref.videoId = videoId
                ref.title = title
                ref.publishedAt = publishedAt
                ref.prev_ref = prev_ref
                ref.link = 'https://www.youtube.com/watch?v=%s&t=%ds' % (videoId, caption['start'])
                ref.last_in_sentence = last_in_sentence
                ref.is_in_dictionary = word in self.words_dictionary

                if prev_ref:
                    prev_ref.next_ref = ref
                else: # first ref for this video
                    self.first_word_in_video[videoId] = ref
                self.word_index[word].append(ref)
                prev_ref = ref       
        
    # Search for an exact word and return all references
    def findWord(self, word):
        if word not in self.word_index:
            print ("'%s' not found" % (word))
            return []
            
        return sorted(self.word_index[word], key=lambda x: x.link)
    
    # Expand regex pattern into all of the words that match it
    def reWord(self, pattern):
        results = {}
        p = re.compile(pattern)
        for word, refs in self.word_index.items():
            if re.search(p,word):
                results[word] = len(refs)
        return results

    # Search for an exact phrase and return all references
    # only_last_in_sentence only works with transcripts that have punctuation 
    def findWords(self, words, only_last_in_sentence=False):
        first_word = words[0].lower()
        if first_word not in self.word_index:
            print ("'%s' not found" % (first_word))
            return []
            
        matches = {}
        refs = self.word_index[first_word]
        for ref in refs:
            if self.matchNextWords(ref, words[1:], only_last_in_sentence):
                id = (ref.videoId,ref.start)
                matches[id] = ref
                
        return list(matches.values())
        
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
        