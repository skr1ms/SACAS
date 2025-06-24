import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

export const submissionApi = {
  async submitCode(fileName, fileType, content) {
    const response = await api.post('/api/submit', {
      file_name: fileName,
      file_type: fileType,
      content: content,
    });
    return response.data;
  },

  async getSubmissions() {
    const response = await api.get('/api/submissions');
    return response.data;
  },

  async getSubmission(id) {
    const response = await api.get(`/api/submissions/${id}`);
    return response.data;
  },
};

export default api;
