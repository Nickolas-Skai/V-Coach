
// static/js/interview.js
let currentQuestionIndex = 0;

let answers = [];

const questionArea = document.getElementById('question-area');
const inputField = document.getElementById('input-field');
const answerForm = document.getElementById('answer-form');
const finishArea = document.getElementById('finish-area');
const voiceBtn = document.getElementById('start-voice');

function renderQuestion() {
  const q = questions[currentQuestionIndex];
  console.log("Rendering question:", q);

 
 if (!q ) {
  console.error("❌ Invalid question data at index", currentQuestionIndex);
  questionArea.innerHTML = "<p>Invalid question data.</p>";
  return;
}

  questionArea.innerHTML = `<p>${q.text}</p>`;
  inputField.innerHTML = '';

  if (q.type === 'text') {
    inputField.innerHTML = '<textarea name="answer" rows="4" required></textarea>';
  } else if (q.type === 'radio' || q.type === 'scale') {
    if (q.options && Array.isArray(q.options)) {
      inputField.innerHTML = q.options.map(opt => 
        `<label><input type="radio" name="answer" value="${opt}" required> ${opt}</label><br>`
      ).join('');
    } else {
      inputField.innerHTML = '<p>No options available.</p>';
    }
  } else if (q.type === 'checkbox') {
    if (q.options && Array.isArray(q.options)) {
      inputField.innerHTML = q.options.map(opt => 
        `<label><input type="checkbox" name="answer" value="${opt}"> ${opt}</label><br>`
      ).join('');
    } else {
      inputField.innerHTML = '<p>No options available.</p>';
    }
  } else if (q.type === 'file') {
    inputField.innerHTML = '<input type="file" name="answer" required>';
  }

  speakQuestion(q.text);
}

function speakQuestion(text) {
  if ('speechSynthesis' in window) {
    const utter = new SpeechSynthesisUtterance(text);
    utter.lang = 'en-US';
    window.speechSynthesis.speak(utter);
  }
}

if ('webkitSpeechRecognition' in window || 'SpeechRecognition' in window) {
  const SpeechRecognition = window.SpeechRecognition || window.webkitSpeechRecognition;
  const recognition = new SpeechRecognition();

  voiceBtn.onclick = () => {
    recognition.start();
    voiceBtn.disabled = true;
  };

  recognition.onresult = (event) => {
    const spokenText = event.results[0][0].transcript.trim();
    const textarea = document.querySelector('textarea');
    const radios = document.querySelectorAll('input[type=radio]');

    if (textarea) {
      textarea.value = spokenText;
    } else if (radios.length) {
      radios.forEach(radio => {
        if (radio.value.toLowerCase() === spokenText.toLowerCase() || radio.value.includes(spokenText)) {
          radio.checked = true;
        }
      });
    }
    voiceBtn.disabled = false;
  };

  recognition.onerror = () => {
    voiceBtn.disabled = false;
  };
}

answerForm.addEventListener('submit', function(e) {
  e.preventDefault();
  const formData = new FormData(answerForm);
  let answer = '';

  if (formData.getAll('answer').length > 1) {
    answer = formData.getAll('answer').join(', ');
  } else {
    answer = formData.get('answer');
  }

  if (!answer || answer.trim() === '') {
    alert('Please answer the question before continuing.');
    return;
  }

  answers.push({
    question_id: questions[currentQuestionIndex].id,
    answer: answer
  });

  currentQuestionIndex++;
  if (currentQuestionIndex >= questions.length) {
    document.getElementById('answer-form').style.display = 'none';
    finishArea.style.display = 'block';
  } else {
    renderQuestion();
  }
});

// Final submission
const submitAll = document.getElementById('submit-all');
submitAll.addEventListener('click', async () => {
  const res = await fetch('/submit-interview', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ answers: answers })
  });

  if (res.ok) {
    alert('✅ Interview submitted successfully!');
    window.location.href = '/thank-you';
  } else {
    alert('❌ Failed to submit. Please try again.');
  }
});

renderQuestion();