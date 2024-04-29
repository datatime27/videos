
import MySQLdb

readTemplate = """
    SELECT timestamp, latitude, longitude
    FROM `@@table_name@@` 
    ORDER BY timestamp DESC
    LIMIT 1 """
writeTemplate = "INSERT INTO `@@table_name@@` VALUES (%s, %s, %s)"

def readLastLocation(conn, table):
    cursor = conn.cursor()

    timestamp = None
    latitude = None
    longitude = None
    
    try:
        cursor.execute (readTemplate.replace('@@table_name@@',table))

        for (timestamp, latitude, longitude) in cursor:
            break 
        
    except MySQLdb.Error, e:
        if 'Duplicate entry' in e.args[1]:
            return
        print
        print "Error %d: %s" % (e.args[0], e.args[1])
        
    cursor.close()
    return timestamp, latitude, longitude

def writeLocation(conn, table, timestamp, latitude, longitude):
    cursor = conn.cursor()

    try:
        sql = (writeTemplate.replace('@@table_name@@',table))
        cursor.execute (sql, (timestamp, latitude, longitude))
        conn.commit()
    except MySQLdb.Error, e:
        if 'Duplicate entry' not in e.args[1]:
            print "Error %d: %s" % (e.args[0], e.args[1])
        
    cursor.close()
