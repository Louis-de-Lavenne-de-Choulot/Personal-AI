//create a list of keywords
var keywords = ["to language", "from language", "activation word", "speech"];
var started = false;

var ssu = new SpeechSynthesisUtterance();
ssu.lang = 'en-US';

// function GetSortedVoices() {
// 	// speechSynthesis.getVoices(); but in a variable
// 	const voices = speechSynthesis.getVoices();
// 	// put it as list for select element
// 	let sel = "";
// 	voices.forEach((voice) => {
// 		// values are the voice name and the language
// 		sel += `<option value="${voice.lang}, ${voice.name}">${voice.lang}, ${voice.name}</option>`;
// 	});

// 	// sort the list by language
// 	sel = sel.split('</option>').sort().join('</option>');
// // add </option> to the end of the list and remove it from the beginning
// 	sel = sel.substring(9)+ '</option>';;

// 	console.log(sel);
// 	// add it to the select element with id voices
// 	document.querySelector('#voices').innerHTML = sel;
// }

try {
    var SpeechRecognition = window.SpeechRecognition || window.webkitSpeechRecognition;
    var recognition = new SpeechRecognition();
    recognition.continuous = true;
    //recognition lang to chinese
    recognition.lang = "en-UK";
} catch (e) {
    console.error(e);
}

// if enter is pressed then update recognition language
document.onkeyup = function (e) {
    if (e.keyCode == 13) {
        recognition.lang = (document.querySelector('.speech_recognition').value);
        ssu.lang = recognition.lang;
    }
}

// eventlistener on <select> with id voices
document.querySelector('#voices').addEventListener('change', function () {
    // get the value of the selected option
    var voice = this.value;
    // get the language from the value
    var lang = voice.substring(0, voice.indexOf(","));
    // get the name from the value
    var name = voice.substring(voice.indexOf(",") + 2);
    recognition.stop();
    // set the language of the speechSynthesisUtterance
    ssu.lang = lang;
    // set the name of the speechSynthesisUtterance
    ssu.voice = name;
});


recognition.onstart = function () {
    started = true;
}

recognition.onend = function (event) {
    started = false;
}

recognition.onerror = function (event) {
    started = false;
    if (event.error == 'no-speech') {
        console.log('No speech was detected. Try again.');
    }else{
        console.log(event.error);
    }
}

let trans = '';
var noteContent = '';
var waitOnce = false;
recognition.onresult = async function (event) {
    var current = event.resultIndex;
    var transcript = event.results[current][0].transcript;
    trans = transcript;
    var mobileRepeatBug = (current == 1 && transcript == event.results[0][0].transcript);
    if (!mobileRepeatBug) {
        var doNotWait = false;
        var actw = (document.querySelector('.activation_word').value).toLowerCase();
        // if  activation word in transcript and if there something after it
        if ((transcript.toLowerCase()).includes(actw) && transcript[transcript.toLowerCase().indexOf(actw) +
                actw.length-1].length > 0) {
            doNotWait = true;
        }
        noteContent += transcript;
        transcript = transcript.toLowerCase();
        //write noteContent to the component with class .ia
        ia.innerText = (noteContent);
        //log doNotWait and waitOnce
        console.log(doNotWait, waitOnce);
        // if transcript contains the word set in the input field .activation_word then call Discuss() and Send() and await for the response
        if (doNotWait || waitOnce) {
            // var indexesAnd = getAllIndexes(transcript, 'and');
            var minIndex = [999, "none"];
            //foreach keyword in the list of keywords
            keywords.forEach(async keyword => {
                if (transcript.includes(keyword)) {
                    var index = transcript.indexOf(keyword);
                    if (index < minIndex[0]) {
                        minIndex = [index, keyword];
                    }
                    // }
                }
            });
            //if minIndex is not 999 then call the function that is in minIndex[1]
            if (minIndex[0] != 999) {
                if (minIndex[1] == "to language") {
                    // take word after to language
                    let word = transcript.substring(minIndex[0] + minIndex[1].length);
                    // if word contains  spaces then take the first word
                    while (word.includes(" ")) {
                        word = word.substring(0, word.indexOf(" "));
                    }

                    lang = transcript.substring(minIndex[0] + minIndex[1].length);

                } else if (minIndex[1] == "activation word") {
                    activation_word.value = transcript.substring(minIndex[0] + minIndex[1].length);
                } else if (minIndex[1] == "speech language") {
                    speech_recognition.value = transcript.substring(minIndex[0] + minIndex[1].length);

                    recognition.lang = speech_recognition.innerText;
                }
            } else {
                // remove the activation word from the transcript
                transcript = transcript.substring(transcript.indexOf(actw) + actw.length);
                await SpeechSend(transcript);
                waitOnce = false;
            }
        }

        if (transcript.includes(document.querySelector('.activation_word').value.toLowerCase()) && !
            doNotWait) {
            waitOnce = true;
        }

    }
}


async function soundProcessing(stream) {
    // TODO: add sound processing
    const audioContext = new AudioContext();
    const analyser = audioContext.createAnalyser();
    const microphone = audioContext.createMediaStreamSource(stream);
    const scriptProcessor = audioContext.createScriptProcessor(2048, 1, 1);

    analyser.smoothingTimeConstant = 0.8;
    analyser.fftSize = 1024;

    microphone.connect(analyser);
    analyser.connect(scriptProcessor);
    scriptProcessor.connect(audioContext.destination);
    scriptProcessor.onaudioprocess = function () {
        const array = new Uint8Array(analyser.frequencyBinCount);
        analyser.getByteFrequencyData(array);
        const arraySum = array.reduce((a, value) => a + value, 0);
        const average = arraySum / array.length;
        if (Math.round(average) > 44 && !started) {
            //if recognition is not running then start it
            console.log("Starting recognition..., average: " + Math.round(average));
            recognition.start();
        }
    };
}

navigator.mediaDevices.getUserMedia({
        audio: true
    })
    .then(function (stream) {
        console.log('You have given me access to your mic, starting recognition...');
        try {
            soundProcessing(stream)
        } catch (e) {
            console.log(e);
        }
    })
    .catch(function (err) {
        console.log('No mic for you!')
    });
    