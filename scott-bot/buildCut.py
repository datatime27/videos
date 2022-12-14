
import sys
import os
import re
import captions
from pprint import pprint
from optparse import OptionParser

CUTS_DIR = 'cuts'
PRPROJ_BIN_NAME = 'tom-scott-go'

aeTemplate = '''
app.beginUndoGroup('import shark');

function load(filepath, inPoint, duration, startTime) {
    try {
        var importOptions = new ImportOptions(new File(filepath));
    } catch(e) {
        alert('Missing '+filepath);
        return startTime+duration;
    }

    footageItem = project.importFile(importOptions);
    var footageLayer = comp.layers.add(footageItem);
    footageLayer.inPoint = inPoint;
    footageLayer.startTime = startTime -footageLayer.inPoint;
    footageLayer.outPoint = footageLayer.inPoint + duration;
    return startTime+duration
}
// Creating project
var project = app.project;

// Creating comp
var comp = project.items.addComp(
    'data',
    1920,
    1080,
    1,
    %g,
    30
    );

// Open the comp
comp.openInViewer();

%s

app.endUndoGroup();
'''

pprojTemplate = '''
var binName = "%s";
var videoDir = %s;
var allMediaTypes = 4;

insertClips();

function insertClips() {
    var seq =  app.project.activeSequence;
    videos = getChildrenMap(getItem(binName));
    var failed_videos = 0;
%s
if (failed_videos > 0) {
    alert('Failed to import ' + failed_videos + ' video files');
    }
}

function overwriteClip(seq, videos, videoName, inPoint, duration) {
    video = videos[videoName]
    if (video == null){ // Load video into bin
        bin = getItem(binName)
        path = videoDir+'\\\\'+videoName
        app.project.importFiles([path], false, bin);
        videos = getChildrenMap(bin);
        video = videos[videoName]
        if (video == null){
            return 1;
        }
    }
    var videoTrack = seq.videoTracks[0];
    
    video.setInPoint(inPoint,allMediaTypes);
    video.setOutPoint(inPoint+duration,allMediaTypes);
    video.setScaleToFrameSize();
    videoTrack.insertClip(video,seq.getPlayerPosition())
    var myTime = new Time();
    myTime.seconds = seq.getPlayerPosition().seconds+duration;
    seq.setPlayerPosition(myTime.ticks);
    return 0
}

function getItem(name) {
    items = app.project.rootItem.children;
    for(var i = 0; i < items.numItems; i++) {
        if (items[i].name == name) {
            return items[i]
        }
    }
    throw Exception(binName+" not found")
}

function getChildrenMap(bin) {
    m = {}
    items = bin.children
    for(var i = 0; i < items.numItems; i++) {
        m[items[i].name] = items[i];
    }
    return m;
}
'''

def fileSafe(s):
    return os.path.join(CUTS_DIR,re.sub("[^a-zA-Z0-9]+","_",str(s)))
    
def ascii(s):
    return ''.join([i if ord(i) < 128 else ' ' for i in s])

def buildAEScript(refs):
    downloads = {}
    duration = 0.0
    lines = ['var startTime = 0;']
    for ref in refs:
        print
        print ref.publishedAt, ref.videoId, ref.title, ref.link
        print ref.start, ref.duration, repr(ref.text)
        lines.append('startTime = load("%s/%s.mp4", %g, %g, startTime);' % (captions.VIDEOS_DIR, ref.videoId, ref.start/1000.0, (ref.duration+0.5)/1000 ))
        downloads[ref.videoId] = 'youtube-dl.exe -o "'+captions.VIDEOS_DIR+'/%(id)s.%(ext)s" https://www.youtube.com/watch?v='+ref.videoId
        duration += (ref.duration+0.5)/1000 
    with open(word+'.jsx','w') as f:
        f.write(aeTemplate % (duration, '\n'.join(lines)))
        
    for i in downloads.values():
        print i
   
def buildPPROJcript(file_name, refs):
    if not refs:
        'No words found'
        return
    lines = []
    for ref in refs:
        #print
        print ref.publishedAt, ref.videoId, ref.title, ref.link
        #print ref.start, ref.duration, repr(ref.text)
        lines.append('failed_videos += overwriteClip(seq, videos, "%s.mp4", %g, %g); // %s' % (ref.videoId, ref.start/1000.0, (ref.duration+0.5)/1000, ascii(ref.text) ))
    with open(file_name+'.jsx','w') as f:
        f.write(pprojTemplate % (PRPROJ_BIN_NAME, repr(os.path.abspath(captions.VIDEOS_DIR)), '\n'.join(lines)))
        
    print 'Wrote', file_name+'.jsx'
            

if __name__ == '__main__':
    if not os.path.exists(CUTS_DIR):
        os.mkdir(CUTS_DIR)
    parser = OptionParser()
    parser.add_option("-e", "--end_only",
                      action="store_true", dest="end_only", default=False,
                      help="Only return results at the end of a sentence. These candidates are easier to edit")

    (options, args) = parser.parse_args()

    # python buildCut.py 2022-01-01 never going to let you down
    cut_off_date = args[0] # 2022-01-01
    words = args[1:] # never going to let you down
    parser = captions.Parser()
    parser.parse(cut_off_date)
    refs = parser.findWords(words, only_last_in_sentence=options.end_only) 
    buildPPROJcript(fileSafe(words), refs)
    parser.printVideosToDownload()
