
import json
from captions import Parser
from pprint import pprint
from collections import defaultdict
from statistics import mean,stdev

p = Parser()
p.parse(input('Channel Name:'),cut_off_date='2005-02-14') # Date youtube was founded
words_per_year = defaultdict(int)
duration_per_year = defaultdict(int)
videos_per_year = defaultdict(int)
rates_per_year = defaultdict(list)
best_videos_per_year = {}
years = set()
lexicon = set()
longest_words = defaultdict(list)
for videoId, first_ref in p.first_word_in_video.items():
    ref = first_ref
    word_counter = 0
    end = 0
    year = int(ref.publishedAt.split('-')[0])
    years.add(year)
    videos_per_year[year] += 1
    viewCount = int(p.video_stats[videoId]['viewCount'])
    if year not in best_videos_per_year or best_videos_per_year[year]['viewCount'] < viewCount:
        best_videos_per_year[year] = {'videoId' : videoId, 'viewCount':  viewCount}
    while ref:
        word_counter += 1
        words_per_year[year] += 1
        if ref.is_in_dictionary:
            lexicon.add(ref.word)
            if len(ref.word) > 6:
                longest_words[len(ref.word)].append(ref.word+' '+ref.link)
        end = ref.start + ref.duration
        ref = ref.next_ref

    rate = word_counter/(end/60.0)
    rates_per_year[year].append(rate)
    #print '%s: %d words %.2f min %.2f words/min' % (videoId, word_counter, end/60.0, rate)
    duration_per_year[year] += end

total_words = 0
total_duration = 0
total_videos = 0

for year in sorted(years):
    total_words += words_per_year[year]
    total_duration += duration_per_year[year]
    total_videos += videos_per_year[year]
    
    rate = words_per_year[year]/(duration_per_year[year]/60.0)
    print ('%d: %d videos %d words %.2f hrs %.2f words/min (%.2f stdev) best video: %s %d' % (
        year, videos_per_year[year], words_per_year[year], duration_per_year[year]/3600.0, mean(rates_per_year[year]), 
        stdev(rates_per_year[year]), best_videos_per_year[year]['videoId'], best_videos_per_year[year]['viewCount']))
print()
print('Total words:',total_words)
print('Total duration: %.2f hrs' % (total_duration/3600.0) )
print('Total videos:',total_videos)

print ('lexicon_count:',len(lexicon))
longest_length = sorted(longest_words.keys())[-1]
print ('longest_words:',longest_length)
pprint(longest_words[longest_length])
