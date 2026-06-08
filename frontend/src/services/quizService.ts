import api from './api';
import { QuizQuestion } from '../types';

export const quizService = {
  async getQuiz(paperId: string): Promise<QuizQuestion[]> {
    const response = await api.get(`/papers/${paperId}/quiz`);
    return response.data.questions || response.data || [];
  },

  async generateQuiz(paperId: string): Promise<QuizQuestion[]> {
    const response = await api.post(`/papers/${paperId}/quiz`);
    return response.data.questions || response.data || [];
  },
};
