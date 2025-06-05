async function loadEditList() {
    const res = await fetch('/questions');
    const questions = await res.json();
    const container = document.getElementById('edit-list');
    container.innerHTML = '';
  
    for (const q of questions) {
      const div = document.createElement('div');
      div.className = 'question';
  
      const answerInputs = q.answers.map((ans, i) => {
        return `
          <input type="text" value="${ans.text}" data-index="${i}" class="ans-input"/>
          <label>
            <input type="radio" name="correct-${q.id}" value="${i}" ${ans.is_correct ? 'checked' : ''}/> Правильный
          </label><br/>
        `;
      }).join('');
  
      div.innerHTML = `
        <form onsubmit="return updateQuestion(event, ${q.id})">
          <textarea name="question">${q.question}</textarea><br/>
          ${answerInputs}
          <input type="file" name="image"/><br/>
          <button type="submit">Сохранить</button>
          <button type="button" onclick="deleteQuestion(${q.id})">Удалить</button>
        </form>
      `;
  
      container.appendChild(div);
    }
  }
  
  async function updateQuestion(e, id) {
    e.preventDefault();
    const form = e.target;
    const question = form.querySelector('textarea[name="question"]').value;
    const answers = [...form.querySelectorAll('.ans-input')].map((input, i) => ({
      text: input.value,
      is_correct: form.querySelector(`input[name="correct-${id}"]:checked`).value == i
    }));
  
    const res = await fetch(`/questions/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ question, answers })
    });
  
    if (!res.ok) {
      alert('Ошибка при обновлении');
      return;
    }
  
    const imageInput = form.querySelector('input[type="file"]');
    if (imageInput.files.length > 0) {
      const formData = new FormData();
      formData.append('image', imageInput.files[0]);
  
      await fetch(`/upload/${id}`, {
        method: 'POST',
        body: formData
      });
    }
  
    alert('Обновлено!');
    loadEditList();
  }
  
  async function deleteQuestion(id) {
    if (!confirm('Удалить вопрос?')) return;
    await fetch(`/questions/${id}`, { method: 'DELETE' });
    loadEditList();
  }
  
  loadEditList();
  