
import sys
import os
import re
import captions
from pprint import pprint
from optparse import OptionParser
   

if __name__ == '__main__':
    parser = OptionParser()
    (options, args) = parser.parse_args()

    p = captions.Parser()
    p.parse(input('Channel Name:'))
    for ref in p.findWords(input('Enter String to search for:').lower().split()):
        print(ref.link,ref.text)
    print()
    
