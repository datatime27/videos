
var powerUp = new Audio('./power-up-1.mp3');
var powerDown = new Audio('./power-down-1.mp3');
const deadzone = 0.001;
const ftInMi = 5280;


function handleError(err) {
  output.innerHTML += `ERROR(${err.code}): ${err.message}`;
}

const options = {
    enableHighAccuracy: true,
    timeout: 5000,
    maximumAge: 0,
}

let wakelock;
const canWakeLock = () => 'wakeLock' in navigator;


async function lockWakeState() {
  if(!canWakeLock()) return;
  try {
    wakelock = await navigator.wakeLock.request();
    wakelock.addEventListener('release', () => {
      console.log('Screen Wake State Locked:', !wakelock.released);
    });
    console.log('Screen Wake State Locked:', !wakelock.released);
  } catch(e) {
    console.log('Failed to lock wake state with reason:', e.message);
  }
}

function formatDistance(distance) {
  if (distance > 1)
    return distance.toFixed(4) + ' mi'
  ft = distance * ftInMi
  return ft.toFixed(4) + ' ft'
}
function mapsLink(latitude, longitude) {
    return '<a href="https://maps.google.com/?q='+latitude+','+longitude+'" target="_blank">' + latitude + ' ' + longitude + '</a>'
}

function warmerColder(distance) {
    if (lastDistance < 0) {
        lastDistance = distance
    }
    if (distance > lastDistance + deadzone) {
        if (enablesounds.checked) {
            powerDown.play();
        }
        document.body.style.background = '#ffaaaa'
        lastDistance = distance
    }
    else if (distance < lastDistance - deadzone) {
        if (enablesounds.checked) {
            powerUp.play();
        }
        document.body.style.background = '#aaffaa'
        lastDistance = distance
    } else {
        document.body.style.background = 'white'
    }
}

function load() {
    refreshPosition();
}

function refreshPosition() {
    lockWakeState();
    refresh.disabled=true;
    
    if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition(updatePosition, handleError, options);
    } else {
        output.innerHTML = "Geolocation is not supported by this browser.";
        refresh.disabled=false;
    }
}


var intervalId = setInterval(refreshPosition, 10000);
function updateInterval(checkbox) {
    if (checkbox.checked) {
        intervalId = setInterval(refreshPosition, 10000);
        refreshPosition();
    } else {
        clearInterval(intervalId);
    }
}