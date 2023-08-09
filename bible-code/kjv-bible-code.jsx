
function addColorText(layer, index, color) {
    var animator = layer.property('Text').property('Animators').addProperty('ADBE Text Animator');
    var colorProp = animator.property('Properties').addProperty('ADBE Text Fill Color');
    colorProp.setValueAtTime(0,[0,0,0]);
    colorProp.setValueAtTime(1,color);
    var selector = animator.property('Selectors').addProperty('ADBE Text Selector');
    var advanced = selector.property('Advanced');
    advanced.property('Based On').setValue(2); // Characters
    advanced.property('Units').setValue(2); // Index
    var start = selector.property('ADBE Text Index Start').setValue(index)
    var end = selector.property('ADBE Text Index End').setValue(index+1)
}


var font_size = 30;
var tracking = 700;

var positionX = 960;
var positionY = 540;

var width = 1800;
var height = 900;

var letterWidth = 22.2
var letterHeight = 36.0
var wrapLength = 22

var originX = positionX - width/2
var originY = positionY - height/2

var text = "GEMYCOMMANDMENTSMYSTATUTESANDMYLAWSANDISAACDWELTINGERARANDTHEMENOFTHEPLACEASKEDHIMOFHISWIFEANDHESAIDSHEISMYSISTERFORHEFEAREDTOSAYSHEISMYWIFELESTSAIDHETHEMENOFTHEPLACESHOULDKILLMEFORREBEKAHBECAUSESHEWASFAIRTOLOOKUPONANDITCAMETOPASSWHENHEHADBEENTHEREALONGTIMETHATABIMELECHKINGOFTHEPHILISTINESLOOKEDOUTATAWINDOWANDSAWANDBEHOLDISAACWASSPORTINGWITHREBEKAHHISWIFEANDABIMELECHCALLEDISAACANDSAIDBEHOLDOFASURETYSHEISTHYWIFEANDHOWSAIDSTTHOUSHEISMYSISTERANDISAACSAIDUNTOHIMBECAUSEISAIDLESTIDIEFORHERANDABIMELECHSAIDWHATISTHISTHOUHASTDO";

app.beginUndoGroup('bible-code');
var comp = app.project.activeItem;

label = comp.layers.addBoxText([width,height], text);
label.threeDLayer = true

var textProp = label.property("Source Text");
label.position.setValue([positionX, positionY]);
var textDocument = textProp.value;
textDocument.fontSize = font_size;
textDocument.fillColor = [0, 0, 0];
textDocument.tracking = tracking;
textProp.setValue(textDocument);

var animator = label.property('Text').property('Animators').addProperty('ADBE Text Animator');
var colorProp = animator.property('Properties').addProperty('ADBE Text Fill Color');
colorProp.setValueAtTime(0,[0,0,0]);
colorProp.setValueAtTime(1,[0.8, 0.8, 0.8]);


// "ELBIB": 62-462 skip:100 
addColorText(label, 62, [0.6046603, 0.47025454, 0.6645601])  //E
addColorText(label, 162, [0.6046603, 0.47025454, 0.6645601])  //L
addColorText(label, 262, [0.6046603, 0.47025454, 0.6645601])  //B
addColorText(label, 362, [0.6046603, 0.47025454, 0.6645601])  //I
addColorText(label, 462, [0.6046603, 0.47025454, 0.6645601])  //B

// "EDOC": 356-464 skip:36 
addColorText(label, 356, [0.4377142, 0.21231875, 0.68682307])  //E
addColorText(label, 392, [0.4377142, 0.21231875, 0.68682307])  //D
addColorText(label, 428, [0.4377142, 0.21231875, 0.68682307])  //O
addColorText(label, 464, [0.4377142, 0.21231875, 0.68682307])  //C

app.endUndoGroup();
