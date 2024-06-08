
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
    p.parse('mrbeast')
    
    # Always lowercase
    print ('"lets go":')
    for ref in p.findWords(['lets','go']):
        print(ref.link,ref.text)
    print()
    
    print ('"feast.*":')
    for words in p.reWord('feast.*'):
        print(words)
    print()
    
    print ('"feastables":')
    for ref in p.findWord('feastables'):
        print(ref.link,ref.text)
    print()
        
    