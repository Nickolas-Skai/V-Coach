//file: static/js/interview.js
const questionContainer = document.getElementById('question-container');
const inputContainer = document.getElementById('input-container');

fetch('/api/question')
  .then(response => response.json())
  .then(data => {
    const question = data;
    renderInput(question);
  })
  .catch(error => console.error('Error fetching question:', error));

function renderInput(question) {
  questionContainer.innerHTML = `<p>${question.text}</p>`;
  inputContainer.innerHTML = '';

  let inputHTML = '';
  switch (question.type) {
    case 'text':
      inputHTML = `<textarea name="answer" rows="4" cols="50" required></textarea>`;
      break;
    case 'radio':
      question.options.forEach((opt, i) => {
        inputHTML += `<label><input type="radio" name="answer" value="${opt}" required> ${opt}</label><br>`;
      });
      break;
    case 'checkbox':
      question.options.forEach((opt, i) => {
        inputHTML += `<label><input type="checkbox" name="answer" value="${opt}"> ${opt}</label><br>`;
      });
      break;
    case 'file':
      inputHTML = `<input type="file" name="answer" accept="image/*" required>`;
      break;
    case 'scale':
      for (let i = 1; i <= 5; i++) {
        inputHTML += `<label><input type="radio" name="answer" value="${i}" required> ${i}</label>`;
      }
      break;
  }

  inputContainer.innerHTML = inputHTML;
}

renderInput(question);
