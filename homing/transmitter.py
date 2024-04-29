
import sys
import datetime
import cgi, cgitb
import calc
import json
import connection
import db

cgitb.enable()
form = cgi.FieldStorage()
latitude = float(form['latitude'].value)
longitude = float(form['longitude'].value)


def returnData(conn):
    receiverTimestamp, receiverLat, receiverLong = db.readLastLocation(conn, 'receiver')
    if receiverTimestamp == None:
        print json.dumps({'distance': -1, 'delay':-1})
        return
    #receiverLat = 34.0
    #receiverLong = -118.6

    distance = calc.distanceBetweenEarthCoordinates(receiverLat, receiverLong, latitude, longitude)
    delay = now-receiverTimestamp    
    print json.dumps({'latitude':receiverLat, 'longitude': receiverLong, 'distance': distance, 'delay': delay.seconds})

print "Content-type: application/json"
print

try:
    conn = connection.getConnection()
except MySQLdb.Error, e:
    print
    print "Error %d: %s" % (e.args[0], e.args[1])
    sys.exit (1)
    
now = datetime.datetime.utcnow()

returnData(conn)

db.writeLocation(conn, 'transmitter', now, latitude, longitude)
conn.close()
