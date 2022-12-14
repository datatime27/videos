# -*- coding: utf-8 -*-

import io
import os
import os.path
import json
import requests
import json
import datetime
import time
from base64 import b64encode
from pprint import pprint
import google_auth_oauthlib.flow
from google.oauth2.credentials import Credentials
from googleapiclient.discovery import build
from googleapiclient.errors import HttpError
from google_auth_oauthlib.flow import InstalledAppFlow
from google.auth.transport.requests import Request

SCOPES = [
    'https://www.googleapis.com/auth/youtube.readonly',
    'https://www.googleapis.com/auth/youtube.force-ssl',
    'https://www.googleapis.com/auth/youtubepartner-channel-audit',
    ]

API_SERVICE_NAME = 'youtube'
API_VERSION = 'v3'
CLIENT_SECRETS_FILE = 'client_secret.json' # Add your youtube api client login file here.

# This is the channel ID for TomScottGo - replace with your youtube channel id.
# Don't know your id? Use this tool https://commentpicker.com/youtube-channel-id.php
YOUTUBE_CHANNEL_ID = "UCBa659QWEk1AI4Tg--mrJ2A"  #TomScottGo
    
captions_dir = 'captions'

class YoutubeConnection:
    def __init__(self):
        self.api = self.get_service()
        
    def get_service(self):
        creds = None
        # The file token.json stores the user's access and refresh tokens, and is
        # created automatically when the authorization flow completes for the first
        # time.
        if os.path.exists('token.json'):
            creds = Credentials.from_authorized_user_file('token.json', SCOPES)
        # If there are no (valid) credentials available, let the user log in.
        if not creds or not creds.valid:
            if creds and creds.expired and creds.refresh_token:
                creds.refresh(Request())
            else:
                flow = InstalledAppFlow.from_client_secrets_file(
                    CLIENT_SECRETS_FILE, SCOPES)
                creds = flow.run_local_server(port=0)
            # Save the credentials for the next run
            with open('token.json', 'w') as token:
                token.write(creds.to_json())

        try:
            return build(API_SERVICE_NAME, API_VERSION, credentials = creds)

        except HttpError as err:
            print(err)

    def download_channel_captions(self, channelId):
        publishedBefore = datetime.datetime.now() + datetime.timedelta(days=1)
        publishedBefore = publishedBefore.isoformat()+'Z'
        counter = 0
        while True:
            search_results = self.api.search().list(
                part="snippet",
                type="video",
                channelId=channelId,
                order='date',
                maxResults=50,
                publishedBefore=publishedBefore).execute()
            for search_result in search_results['items']:
                publishedBefore=search_result['snippet']['publishedAt']
                publishedBefore = datetime.datetime.strptime(publishedBefore,'%Y-%m-%dT%H:%M:%SZ') - datetime.timedelta(seconds=1)
                publishedBefore = publishedBefore.isoformat()+'Z'
                print counter, search_result['id']['videoId'], search_result['snippet']['publishedAt']
                counter+=1
                #download_captions(search_result)
            if not search_results['items']:
                break

    # Youtube pagination is broken - don't use this function. Use above function instead.
    '''
    def download_channel_captions_pagination(self, channelId):
        pageToken = None
        while True:
            search_results = self.api.search().list(
                part="snippet",
                type="video",
                channelId=channelId,
                order='date',
                maxResults=50,
                publishedBefore='2017-01-06T13:07:03Z',
                pageToken=pageToken).execute()
            for search_result in search_results['items']:
               print search_result['id']['videoId'], search_result['snippet']['publishedAt']
               #download_captions(search_result)
            
            if 'nextPageToken' not in search_results:
                break
            pageToken = search_results['nextPageToken']
            print pageToken
    '''

# Sample curl command
# curl -s 'https://www.youtube.com/youtubei/v1/get_transcript?key=AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8' -H 'Content-Type: application/json' --data-raw 
# "{\"context\":{\"client\":{\"clientName\":\"WEB\",\"clientVersion\":\"2.2021111\"}},\"params\":\"$(printf '\n\x0bEi6LKHNFpeE' | base64)\"}"
def download_captions(item):
    videoId = str(item['id']['videoId'])
    params = b64encode('\n\x0b'+videoId)
    url = 'https://www.youtube.com/youtubei/v1/get_transcript?key=AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8' 
    myobj = {"context":{"client":{"clientName":"WEB","clientVersion":"2.2021111"}},"params":params}
    
    headers = {'Content-Type': 'application/json'}
    r = requests.post(url, json = myobj , headers=headers)
    
    obj = json.loads(r.text)
    obj['id'] = item['id']
    obj['snippet'] = item['snippet']
    filepath = os.path.join(captions_dir, videoId+'.json')
    with open(filepath, 'wb') as f:
        json.dump(obj,f,indent=4)

    
if __name__ == '__main__':
    if not os.path.exists(captions_dir):
        os.mkdir(captions_dir)
    client = YoutubeConnection()
    client.download_channel_captions(YOUTUBE_CHANNEL_ID)
