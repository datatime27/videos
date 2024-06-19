
from captions import Parser
from pprint import pprint
from collections import defaultdict
from nltk.sentiment.vader import SentimentIntensityAnalyzer
from nltk.corpus import stopwords
from nltk.tokenize import word_tokenize
from nltk.stem import WordNetLemmatizer
# First time only
# import nltk
# nltk.download('all')

def getSentiment(text):
    analyzer = SentimentIntensityAnalyzer()
    scores = analyzer.polarity_scores(text)
    return scores

p = Parser()
p.parse(input('Channel Name:'),cut_off_date='2005-02-14')

text_by_year = defaultdict(list)

print("# Raw Scatter Plot:")
print("videoId, words-per-min, viewCount, publishedAt, negative, neutral, positive, compound")
for videoId, first_ref in p.first_word_in_video.items():
    ref = first_ref
    year = int(ref.publishedAt.split('-')[0])
    word_counter = 0
    end = 0
    text = []
    while ref:
        word_counter += 1
        end = ref.start + ref.duration
        text_by_year[year].append(ref.word)
        text.append(ref.word)

        ref = ref.next_ref
        
    rate = word_counter/(end/60.0)
    stats = p.video_stats[videoId]
    sentiment = getSentiment('\n'.join(text))
    print ('%s %.2f %s %s %g %g %g %g' % (videoId, rate, stats['viewCount'], first_ref.publishedAt, sentiment['neg'], sentiment['neu'], sentiment['pos'], sentiment['compound']))
    #print '%s: %d words %.2f min %.2f words/min views: %s' % (videoId, word_counter, end/60.0, rate, stats)

print()
print("# Year over Year sentiment:")
print("year, negative, neutral, positive, compound")
for year,text in sorted(text_by_year.items()):
    sentiment = getSentiment('\n'.join(text))
    print ('%s %g %g %g %g' % (year, sentiment['neg'], sentiment['neu'], sentiment['pos'], sentiment['compound']) )
