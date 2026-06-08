import api from './api';
import { Paper } from '../types';

export const paperService = {
  async getAllPapers(): Promise<Paper[]> {
    const response = await api.get('/papers');
    return response.data.papers || response.data || [];
  },

  async getPaper(id: string): Promise<Paper> {
    const response = await api.get(`/papers/${id}`);
    return response.data.paper || response.data;
  },

  async uploadPaper(file: File, title: string, onProgress?: (progress: number) => void): Promise<Paper> {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('title', title);

    const response = await api.post('/papers/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      onUploadProgress: (progressEvent) => {
        if (onProgress && progressEvent.total) {
          const progress = Math.round((progressEvent.loaded * 100) / progressEvent.total);
          onProgress(progress);
        }
      },
    });
    return response.data.paper || response.data;
  },

  async analyzePaper(id: string): Promise<{ message: string }> {
    const response = await api.post(`/papers/${id}/analyze`);
    return response.data;
  },

  async deletePaper(id: string): Promise<void> {
    await api.delete(`/papers/${id}`);
  },

  async downloadReport(id: string): Promise<Blob> {
    const response = await api.get(`/papers/${id}/report`, {
      responseType: 'blob',
    });
    return response.data;
  },

  async getReport(id: string): Promise<any> {
    const response = await api.get(`/papers/${id}/report`);
    return response.data;
  },
};
