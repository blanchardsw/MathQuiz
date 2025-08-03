import axios, { AxiosInstance } from 'axios';
import {jwtDecode} from 'jwt-decode';

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

interface DecodedToken {
  exp: number;
}

function isTokenExpired(token: string): boolean {
  try {
    const decoded = jwtDecode<DecodedToken>(token);
    const now = Math.floor(Date.now() / 1000);
    return decoded.exp < now;
  } catch {
    return true;
  }
}

class ApiService {
  private axiosInstance: AxiosInstance;
  private baseURL: string;

  constructor() {
    this.baseURL = process.env.REACT_APP_API_URL || 'http://localhost:4000/api';
    this.axiosInstance = axios.create({
      baseURL: this.baseURL,
      headers: {
        'Content-Type': 'application/json',
      },
      withCredentials: true
    });
  
    // ✅ Add interceptor to inject JWT
    this.axiosInstance.interceptors.request.use(async config => {
      const token = sessionStorage.getItem("jwt");
    
      if (token) {
        if (isTokenExpired(token)) {
          const newToken = await this.refreshToken();
          if (newToken) {
            config.headers.Authorization = `Bearer ${newToken}`;
          }
        } else {
          config.headers.Authorization = `Bearer ${token}`;
        }
      }
    
      return config;
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

  async initSession(): Promise<void> {
    const res = await this.axiosInstance.post(`/init-session`, {}, { withCredentials: true });
    const { token } = res.data;
    if (token) {
      sessionStorage.setItem("jwt", token);
    }
  }

  async refreshToken(): Promise<string | null> {
    try {
      const res = await this.axiosInstance.post<{ token: string }>("/refresh-token", {}, { withCredentials: true });
      const { token } = res.data;
      if (token) {
        sessionStorage.setItem("jwt", token);
        return token;
      }
      this.logout();
      return null;
    } catch (err) {
      console.error("Failed to refresh token", err);
      this.logout();
      return null;
    }
  }

  logout(): void {
    sessionStorage.removeItem("jwt");
    if (window.location.pathname !== "/") {
      window.location.href = "/";
    }
  }
}

export const apiService = new ApiService();
