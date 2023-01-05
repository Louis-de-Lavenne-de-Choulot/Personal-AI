function getAllIndexes(arr, val) {
    var indexes = [];
    var i = -1;
    indexes.push(arr.indexOf(val, 0));
    while ((i = arr.indexOf(val, i + 1)) != -1) {
        indexes.push(arr.indexOf(val, i + 1));
    }
    //remove last element if it is -1
    indexes.pop();
    return indexes;
}
const ia = document.querySelector('.ia');
const ia_a = document.querySelector('.ia-answer');
let lang = "";
const activation_word = document.querySelector('.activation_word');
const speech_recognition = document.querySelector('.speech_recognition');


async function SpeechSend(transcript) {
    //replace activation word by "" in the transcript
    transcript = transcript.replace(activation_word, "");
    let fl = "";
    if (lang) {
        fl = "&lang="+lang.value;
    }
    let fl2 = "";
    if (speech_recognition) {
        // remove after ,
        fl2 = "&fromlang=" + speech_recognition.value.substring(0, speech_recognition.value.indexOf(","));
    }
    //get answer from the server bestiaever.ml/ia/answer with the transcript, lang
    const response = await fetch('http://localhost:5019/ai/answer?text=' + transcript + fl+fl2);
    const answer = await response.text();

    //say answer
    window.speechSynthesis.speak(new SpeechSynthesisUtterance(answer));
    // write answer in the text with class ia-answer
    ia_a.innerHTML += `<div class="ia-answer">${answer}</div>`;
}