
# coding=latin1

import math
xml = '''<?xml version="1.0"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.0//EN" "http://www.w3.org/TR/2001/REC-SVG-20010904/DTD/svg10.dtd">
<svg viewBox="0 0 8.5in 11in" xmlns="http://www.w3.org/2000/svg">
    <style>
        .ruler {
          font: 11px sans-serif;
          fill: red;
        }
        .degrees {
          font: 11px sans-serif;
          fill: red;
        }
        .altitude {
          font: 11px sans-serif;
          fill: green;
        }
    </style>
    <rect x="0" y="0" width="8.5in" height="11in" fill="none" stroke="black"/>

%s
</svg>
'''

height = 11.811 # inches (30 cm)
pages = 4
total_degrees = 80
degree_increments = 2
page_offset = 10
page_length = 10.75
degree_symbol = '\xb0'
degree_symbol = '°'

def buildLogo():
    return '<image href="logo-sm-v1.png" x="6.2in" y="9.5in" width="1.5in" height="1.5in"/>'

def buildRuler(page):
    lines = []
    for i in range(50):
        inches = page_length + page_offset*page - i
        d = {
            'y': inches,
            'quant': i}
        lines.append('''
            <text text-anchor="middle" x="0.4in" y="%(y)sin" class="ruler">%(quant)d in</text>
            <line x1="0in" y1="%(y)sin" x2="0.4in" y2="%(y)sin" stroke="black" />''' % d)
           
    for i in range(110):
        cm = page_length*2.54 + page_offset*2.54*page - i
        d = {
            'y': cm,
            'quant': i}
        lines.append('''
            <text text-anchor="middle" x="8.0in" y="%(y)scm" class="ruler">%(quant)d cm</text>
            <line x1="8.5in" y1="%(y)scm" x2="8.0in" y2="%(y)scm" stroke="black" />''' % d)
    return '\n'.join(lines)
    
def buildSVG(page):
    lines = []
    for i in range(1,total_degrees* degree_increments):
        centerx = 4.25
        centery = page_offset * page
        degrees = i/2.0
        
        r = height * math.tan(math.radians(degrees))
        
        stroke = "1" if i%2 == 0 else "0.2"

        d = {
            'cx': centerx,
            'cy': page_length+centery,
            'r': r,
            'texty': page_length+centery-r,
            'degrees': int(degrees),
            'degreesx': centerx-0.2,
            'altitude': 90-int(degrees),
            'altitudex': centerx+0.2,
            'stroke': stroke,
            'degree_symbol': degree_symbol,
        }
        lines.append('<circle cx="%(cx)sin" cy="%(cy)sin" r="%(r)sin" fill="none" stroke="black" stroke-width="%(stroke)s"/>' %d)
        if i%2 == 0:
            lines.append('<text text-anchor="middle" x="%(degreesx)sin" y="%(texty)sin" class="degrees">%(degrees)s%(degree_symbol)s</text>' %d)
            lines.append('<text text-anchor="middle" x="%(altitudex)sin" y="%(texty)sin" class="altitude">%(altitude)s%(degree_symbol)s</text>' %d)
    
    return '\n'.join(lines) + buildRuler(page)
def main():
    for index,page in enumerate(range(pages)):
        svg = buildSVG(page)
        if index == 0:
            svg += buildLogo()
        filepath = 'page%d.svg' % (page) 
        with open(filepath, 'w') as f:
            f.write(xml % (svg))


main()
