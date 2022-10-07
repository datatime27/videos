var dialog = new Window("dialog");
dialog.add("statictext", undefined, "Message:");
messageText = dialog.add("edittext", [0,0,800,20], "test");
dialog.add("statictext", undefined, "Duration (secs):");
durationText = dialog.add("edittext", [0,0,40,20], "1.0");
dialog.add( "button", undefined, "OK", { name: "ok" } );
ret = dialog.show();


var message = messageText.text;
var totalDuration = parseFloat(durationText.text);

app.beginUndoGroup('generate-subtitles');
var comp = app.project.activeItem;
var layer = comp.layers.addBoxText([1280,600],message);

var animator = layer.property('Text').property('Animators').addProperty('ADBE Text Animator');
var opacity = animator.property('Properties').addProperty('ADBE Text Opacity');
opacity.setValue(0);
var selector = animator.property('Selectors').addProperty('ADBE Text Selector');
var advanced = selector.property('Advanced');
advanced.property('Based On').setValue(3); // Words
advanced.property('Units').setValue(2); // Index
var start = selector.property('ADBE Text Index Start')


var words = message.split(" ");
var currentTime = 0;
for (var i=0; i < words.length; i++) {
    var word = words[i];
    start.setValueAtTime(currentTime,i+1);
    start.setInterpolationTypeAtKey(i+1, KeyframeInterpolationType.HOLD);
    currentTime += totalDuration * (word.length+1) / message.length;
}
layer.startTime = app.project.activeItem.time;
app.endUndoGroup();



