import React, { useState } from 'react';
import { submissionApi } from '../../shared/api';

const FILE_TYPES = [
  { value: '.cpp', label: 'C++ (.cpp)' },
  { value: '.java', label: 'Java (.java)' },
  { value: '.js', label: 'JavaScript (.js)' },
  { value: '.kt', label: 'Kotlin (.kt)' },
  { value: '.py', label: 'Python (.py)' },
];

export const SubmitCodeForm = ({ onSubmissionComplete, onLoadingChange }) => {
  const [selectedFileType, setSelectedFileType] = useState('');
  const [selectedFile, setSelectedFile] = useState(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState('');

  const handleFileTypeChange = (e) => {
    setSelectedFileType(e.target.value);
    setSelectedFile(null);
    setError('');
  };

  const handleFileChange = (e) => {
    const file = e.target.files[0];
    if (!file) return;

    const fileExtension = '.' + file.name.split('.').pop().toLowerCase();
    if (fileExtension !== selectedFileType) {
      setError(`Выбранный файл не соответствует типу ${selectedFileType}`);
      return;
    }

    setSelectedFile(file);
    setError('');
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!selectedFileType || !selectedFile) {
      setError('Пожалуйста, выберите тип файла и загрузите файл');
      return;
    }    setIsSubmitting(true);
    setError('');
    onLoadingChange?.(true);

    try {
      const fileContent = await readFileContent(selectedFile);
      
      const result = await submissionApi.submitCode(
        selectedFile.name,
        selectedFileType,
        fileContent
      );

      onSubmissionComplete(result);
      
      setSelectedFileType('');
      setSelectedFile(null);
      
    } catch (err) {
      setError('Ошибка при отправке файла: ' + (err.response?.data?.error || err.message));    } finally {
      setIsSubmitting(false);
      onLoadingChange?.(false);
    }
  };

  const readFileContent = (file) => {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = (e) => resolve(e.target.result);
      reader.onerror = (e) => reject(e);
      reader.readAsText(file);
    });
  };

  return (
    <div className="card">
      <h2>Сдать задание</h2>
      
      {error && (
        <div className="alert alert-error">
          {error}
        </div>
      )}

      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label className="form-label">Тип файла:</label>
          <select
            className="form-select"
            value={selectedFileType}
            onChange={handleFileTypeChange}
            disabled={isSubmitting}
          >
            <option value="">Выберите тип файла</option>
            {FILE_TYPES.map((type) => (
              <option key={type.value} value={type.value}>
                {type.label}
              </option>
            ))}
          </select>
        </div>

        {selectedFileType && (
          <div className="form-group">
            <label className="form-label">
              Выберите файл {selectedFileType}:
            </label>
            <input
              type="file"
              className="form-input"
              accept={selectedFileType}
              onChange={handleFileChange}
              disabled={isSubmitting}
            />
          </div>
        )}

        {selectedFile && (
          <div className="form-group">
            <p><strong>Выбранный файл:</strong> {selectedFile.name}</p>
            <p><strong>Размер:</strong> {(selectedFile.size / 1024).toFixed(2)} KB</p>
          </div>
        )}

        <button
          type="submit"
          className="btn btn-primary"
          disabled={!selectedFileType || !selectedFile || isSubmitting}
        >
          {isSubmitting ? 'Отправка...' : 'Сдать задание'}
        </button>
      </form>
    </div>
  );
};
