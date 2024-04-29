
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


def calcDistance(conn):
    transmitterTimestamp, transmitterLat, transmitterLong = db.readLastLocation(conn, 'transmitter')
    if transmitterTimestamp == None:
        print json.dumps({'distance': -1, 'delay':-1})
        return
    distance = calc.distanceBetweenEarthCoordinates(transmitterLat, transmitterLong, latitude, longitude)
    delay = now-transmitterTimestamp
    print json.dumps({'distance': distance, 'delay': delay.seconds})


print "Content-type: application/json"
print

try:
    conn = connection.getConnection()
except MySQLdb.Error, e:
    print
    print "Error %d: %s" % (e.args[0], e.args[1])
    sys.exit (1)
    
now = datetime.datetime.utcnow()

calcDistance(conn)

db.writeLocation(conn, 'receiver', now, latitude, longitude)
conn.close()
