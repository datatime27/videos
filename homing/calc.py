import math

earthRadiusKm = 6371
earthRadiusMi = 3958.8

def degreesToRadians(degrees) :
    return degrees * math.pi / 180.0


def distanceBetweenEarthCoordinates(lat1, lon1, lat2, lon2):
    dLat = degreesToRadians(lat2-lat1)
    dLon = degreesToRadians(lon2-lon1)

    lat1 = degreesToRadians(lat1)
    lat2 = degreesToRadians(lat2)

    a = math.sin(dLat/2) * math.sin(dLat/2) + math.sin(dLon/2) * math.sin(dLon/2) * math.cos(lat1) * math.cos(lat2)
    c = 2 * math.atan2(math.sqrt(a), math.sqrt(1-a))
    mi = earthRadiusMi * c
    return mi
    