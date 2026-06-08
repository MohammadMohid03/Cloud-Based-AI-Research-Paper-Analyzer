// User
export interface User {
  id: string;
  name: string;
  email: string;
  createdAt?: string;
}

// Auth
export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  name: string;
  email: string;
  password: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

// Paper
export interface Paper {
  id: string;
  title: string;
  fileName: string;
  fileSize?: number;
  status: 'pending' | 'processing' | 'completed' | 'failed';
  uploadedAt: string;
  analyzedAt?: string;
  analysis?: PaperAnalysis;
}

export interface PaperAnalysis {
  summary: string;
  keyFindings: string[];
  methodology: string;
  limitations: string[];
  futureScope: string[];
  keywords: string[];
}

// Chat
export interface ChatMessage {
  id: string;
  role: 'user' | 'assistant';
  content: string;
  timestamp: string;
}

export interface ChatRequest {
  question: string;
}

export interface ChatResponse {
  answer: string;
}

// Quiz
export interface QuizQuestion {
  id: string;
  question: string;
  options: string[];
  correctAnswer: number;
  explanation: string;
}

export interface QuizAttempt {
  questionId: string;
  selectedAnswer: number;
  isCorrect: boolean;
}

// API Response wrapper
export interface ApiResponse<T> {
  data: T;
  message?: string;
  success: boolean;
}

// Toast
export interface Toast {
  id: string;
  type: 'success' | 'error' | 'info' | 'warning';
  message: string;
}
