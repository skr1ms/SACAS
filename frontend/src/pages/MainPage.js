import React, { useState } from 'react';
import { SubmitCodeForm } from '../features/submit-code';
import { SubmissionResult } from '../widgets/submission-result';

export const MainPage = () => {
  const [submissionResult, setSubmissionResult] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  const handleSubmissionComplete = (result) => {
    setSubmissionResult(result);
  };

  const handleLoadingChange = (loading) => {
    setIsLoading(loading);
  };

  const handleReset = () => {
    setSubmissionResult(null);
  };

  return (
    <div className="container">
      <header style={{ textAlign: 'center', marginBottom: '40px' }}>
        <h1>CodeGrader</h1>
        <p>Система автоматической проверки и оценки кода студентов</p>
      </header>

      {isLoading && (
        <div className="loading">
          <h3>Анализ кода...</h3>
          <p>Пожалуйста, подождите. Ваш код анализируется с помощью ИИ.</p>
        </div>
      )}      {!isLoading && !submissionResult && (
        <SubmitCodeForm 
          onSubmissionComplete={handleSubmissionComplete}
          onLoadingChange={handleLoadingChange}
        />
      )}

      {!isLoading && submissionResult && (
        <SubmissionResult result={submissionResult} onReset={handleReset} />
      )}
    </div>
  );
};
