import api from './api';
import { ChatResponse } from '../types';

export const chatService = {
  async sendMessage(paperId: string, question: string): Promise<ChatResponse> {
    const response = await api.post<ChatResponse>(`/papers/${paperId}/chat`, { question });
    return response.data;
  },

  async getChatHistory(paperId: string): Promise<any[]> {
    try {
      const response = await api.get(`/papers/${paperId}/chat`);
      return response.data.messages || response.data || [];
    } catch {
      return [];
    }
  },
};
