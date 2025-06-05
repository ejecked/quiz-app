document.getElementById('questionForm').addEventListener('submit', async function (e) {
    e.preventDefault();
  
    const question = document.getElementById('questionText').value;
    const answers = [
      document.getElementById('ans1').value,
      document.getElementById('ans2').value,
      document.getElementById('ans3').value,
    ];
    const correctIndex = parseInt(document.getElementById('correctAnswer').value);
  
    const payload = {
      question,
      answers: answers.map((text, i) => ({ text, is_correct: i === correctIndex }))
    };
  
    const res = await fetch('/questions', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    });
  
    if (!res.ok) {
      alert('Ошибка при добавлении');
      return;
    }
  
    if (document.getElementById('image').files.length > 0) {
      const questionsRes = await fetch('/questions');
      const questionsList = await questionsRes.json();
      const lastQuestion = questionsList[questionsList.length - 1];
  
      const formData = new FormData();
      formData.append('image', document.getElementById('image').files[0]);
  
      await fetch(`/upload/${lastQuestion.id}`, {
        method: 'POST',
        body: formData
      });
    }
  
    alert('Вопрос добавлен!');
    e.target.reset();
  });
  