import React from 'react';

export const SubmissionResult = ({ result, onReset }) => {
  if (!result) return null;

  const getGradeColor = (grade) => {
    if (grade >= 5) return '#28a745';
    if (grade >= 4) return '#ffc107';
    return '#dc3545';
  };

  return (
    <div className="card">
      <h2>Результат проверки</h2>
      
      <div className="alert alert-success">
        ✅ Задание успешно отправлено!
      </div>

      <div className="grade" style={{ color: getGradeColor(result.grade) }}>
        Оценка: {result.grade}/5
      </div>

      {result.feedback && (
        <div>
          <h3>Обратная связь:</h3>
          <div className="feedback">
            {result.feedback}
          </div>
        </div>
      )}

      <button
        className="btn btn-primary"
        onClick={onReset}
        style={{ marginTop: '20px' }}
      >
        Сдать новое задание
      </button>
    </div>
  );
};
