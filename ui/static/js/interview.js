let currentIndex = 0;
const answers = [];
const questionEl = document.getElementById("question-area");
const inputField = document.getElementById("input-field");
const form = document.getElementById("answer-form");
const finishArea = document.getElementById("finish-area");
const voiceBtn = document.getElementById("start-voice");

console.log("Loaded Questions:", questions);

if (questions.length === 0) {
  console.warn("No questions available for the interview");
}

function renderQuestion() {
  const q = questions[currentIndex];
  if (!q) return;

  questionEl.innerHTML = `<p>${q.text}</p>`;
  inputField.innerHTML = "";

  if (q.type === "text") {
    inputField.innerHTML = `<textarea name="answer" rows="4" required></textarea>`;
  } else if (q.type === "radio" || q.type === "scale") {
    inputField.innerHTML = q.options.map(opt =>
      `<label><input type="radio" name="answer" value="${opt}" required> ${opt}</label><br>`
    ).join('');
  } else if (q.type === "checkbox") {
    inputField.innerHTML = q.options.map(opt =>
      `<label><input type="checkbox" name="answer" value="${opt}"> ${opt}</label><br>`
    ).join('');
  } else if (q.type === "file") {
    inputField.innerHTML = `<input type="file" name="answer">`;
  }

  speak(q.text);
}

function speak(text) {
  if ('speechSynthesis' in window) {
    const utter = new SpeechSynthesisUtterance(text);
    utter.lang = 'en-US';
    window.speechSynthesis.speak(utter);
  }
}

form.addEventListener("submit", e => {
  e.preventDefault();
  const formData = new FormData(form);
  let value = formData.getAll("answer").join(", ");

  if (!value.trim()) {
    alert("Please provide an answer.");
    return;
  }

  answers.push({ question_id: questionEl[currentIndex].id, answer: value });
  currentIndex++;

  if (currentIndex < questionEl.length) {
    renderQuestion();
  } else {
    form.style.display = "none";
    finishArea.style.display = "block";
  }
});

document.getElementById("submit-all").addEventListener("click", async () => {
  const res = await fetch("/submit-interview", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ answers })
  });
  if (res.ok) {
    alert("Interview submitted!");
    window.location.href = "/thank-you";
  } else {
    alert("Error submitting responses.");
  }
});

if ('webkitSpeechRecognition' in window || 'SpeechRecognition' in window) {
  const Recognition = window.SpeechRecognition || window.webkitSpeechRecognition;
  const recog = new Recognition();

  voiceBtn.onclick = () => {
    recog.start();
    voiceBtn.disabled = true;
  };

  recog.onresult = e => {
    const speech = e.results[0][0].transcript;
    const textarea = document.querySelector("textarea");
    const radios = document.querySelectorAll("input[type=radio]");
    if (textarea) textarea.value = speech;
    if (radios.length) {
      radios.forEach(r => {
        if (r.value.toLowerCase() === speech.toLowerCase()) r.checked = true;
      });
    }
    voiceBtn.disabled = false;
  };

  recog.onerror = () => { voiceBtn.disabled = false; };
}

renderQuestion();
