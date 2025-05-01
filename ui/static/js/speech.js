// static/js/speech.js

document.addEventListener("DOMContentLoaded", function () {
  const speakButton = document.getElementById("speak-current");
  const voiceInputButton = document.getElementById("voice-fill");

  // Text-to-speech: Read the question out loud
  speakButton.addEventListener("click", function () {
    const questionTextElement = document.querySelector(".question-text");
    if (!questionTextElement) {
      console.error("❌ Could not find .question-text element on the page.");
      alert("Unable to find the question to read aloud.");
      return;
    }
    const questionText = questionTextElement.innerText;
    if ('speechSynthesis' in window) {
      const utterance = new SpeechSynthesisUtterance(questionText);
      utterance.lang = 'en-US';
      window.speechSynthesis.speak(utterance);
    } else {
      console.warn("🔇 This browser does not support speech synthesis.");
      alert("Your browser does not support text-to-speech.");
    }
  });

  // Speech-to-text: Fill the answer field by voice
  if ('webkitSpeechRecognition' in window || 'SpeechRecognition' in window) {
    const SpeechRecognition = window.SpeechRecognition || window.webkitSpeechRecognition;
    const recognition = new SpeechRecognition();
    recognition.lang = 'en-US';
    recognition.interimResults = false;

    voiceInputButton.addEventListener("click", function () {
      try {
        recognition.start();
        voiceInputButton.disabled = true;
      } catch (error) {
        console.error("❌ Failed to start speech recognition:", error);
        alert("Error starting voice input. Please try again.");
      }
    });

    recognition.onresult = function (event) {
      const transcript = event.results[0][0].transcript.trim().toLowerCase();
      const textarea = document.querySelector("textarea[name='answer']");
      const radios = document.querySelectorAll("input[type='radio'][name='answer']");

      if (textarea) {
        textarea.value = transcript;
      } else if (radios.length > 0) {
        let matched = false;

        radios.forEach((radio) => {
          const label = radio.value.toLowerCase();
          const radioNum = label.match(/\d+/);
          const transcriptNum = transcript.match(/\d+/);

          if (
            transcript.includes(label) ||
            (radioNum && transcriptNum && radioNum[0] === transcriptNum[0])
          ) {
            radio.checked = true;
            matched = true;
          }
        });

        if (!matched) {
          console.warn("🎙️ Voice input did not match any radio options.");
          alert("We couldn't match your spoken answer to the available options.");
        }
      } else {
        console.warn("❗ No supported input field found to fill.");
        alert("No appropriate input field to insert your voice response.");
      }

      voiceInputButton.disabled = false;
    };

    recognition.onerror = function (event) {
      console.error("🛑 Speech recognition error:", event.error);
      alert("Voice input error: " + event.error);
      voiceInputButton.disabled = false;
    };
  } else {
    voiceInputButton.disabled = true;
    console.warn("Speech recognition not supported in this browser.");
    alert("Your browser does not support voice input.");
  }
});
