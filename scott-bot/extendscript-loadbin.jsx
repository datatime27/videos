// This is a convenience script that you can run in ExtendScript 
// to load all of the files in your videos directory into a Premiere project bin folder.

var binName = "tom-scott-go";
var allMediaTypes = 4;

// main function where everything is checked and ran
function loadVideos() {
    var importFolder = new Folder;
    importFolder = Folder.selectDialog("Open a folder");

    // if folder is not selected, tell them to run script again
    if(importFolder == null) {
            alert("No folder selected", "Please run again");
            return false;
        }
    // if folder is selected, get all the files inside of it
    var paths = getPaths(importFolder.getFiles());

    var project = app.project;
    var projectItem = project.rootItem;

    var bin = projectItem.createBin(binName);
    project.importFiles(paths,false,bin);
}

function getPaths(files) {
    var thisName;
    var paths = [];
    for(var i = 0; i < files.length; i++) {
        thisName = files[i].name;
        paths.push(files[i].fsName);
    }
    return paths;
}

