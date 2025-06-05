async function loadQuestions() {
    const res = await fetch('/questions');
    const questions = await res.json();
    const container = document.getElementById('quiz');
    container.innerHTML = '';
  
    questions.forEach(q => {
      const div = document.createElement('div');
      div.className = 'question';
      div.innerHTML = `<h3>${q.question}</h3>`;
  
      if (q.image) {
        div.innerHTML += `<img src="${q.image}" alt="Image"/>`;
      }
  
      q.answers.forEach(ans => {
        const btn = document.createElement('button');
        btn.textContent = ans.text;
        btn.onclick = async () => {
          const response = await fetch(`/answer/${ans.id}`, { method: 'POST' });
          const text = await response.text();
          resultEl.textContent = text;
          resultEl.className = 'result ' + (text.includes('Молодец') ? 'correct' : 'wrong');
        };
        div.appendChild(btn);
      });
  
      const resultEl = document.createElement('div');
      resultEl.className = 'result';
      div.appendChild(resultEl);
      container.appendChild(div);
    });
  }
  
  loadQuestions();
  