import axios, { AxiosInstance } from 'axios';

export interface Question {
  operand1: number;
  operand2: number;
  operator: string;
  id: string;
}

export interface AnswerRequest {
  userAnswer: number;
  questionId: string;
}

export interface AnswerResponse {
  correct: boolean;
  correctAnswer: number;
}

export interface ScoreResponse {
  currentScore: number;
  highScores: { [key: string]: number };
  isNewRecord: boolean;
}

export type Difficulty = 'easy' | 'normal' | 'hard';

class ApiService {
  private axiosInstance: AxiosInstance;

  constructor() {
    const baseURL = process.env.REACT_APP_API_URL || 'http://localhost:4000/api';
    this.axiosInstance = axios.create({
      baseURL,
      headers: {
        'Content-Type': 'application/json',
      },
      withCredentials: true // ✅ this sends cookies
    });
  }

  async getQuiz(difficulty: Difficulty = 'normal'): Promise<Question> {
    const response = await this.axiosInstance.get<Question>(`/question?difficulty=${difficulty}`);
    return response.data;
  }  

  async submitAnswer(answerRequest: AnswerRequest): Promise<AnswerResponse> {
    const response = await this.axiosInstance.post<AnswerResponse>('/answer', answerRequest);
    return response.data;
  }

  async getScore(): Promise<ScoreResponse> {
    const response = await this.axiosInstance.get<ScoreResponse>('/score');
    return response.data;
  }

  async resetScore(difficulty: Difficulty): Promise<ScoreResponse> {
    const response = await this.axiosInstance.post<ScoreResponse>('/reset-score', {
      difficulty: difficulty
    });
    return response.data;
  }
}

export const apiService = new ApiService();
