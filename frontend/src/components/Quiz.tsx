import React, { useState, useEffect, useRef, useCallback } from 'react';
import { Question, AnswerResponse, Difficulty, apiService } from '../services/api';
import Timer from './Timer';

interface QuizProps {
  onAnswerSubmitted: (correct: boolean) => void;
  isActive: boolean;
  onTimeUp: () => void;
  difficulty: Difficulty;
}

const Quiz: React.FC<QuizProps> = ({ onAnswerSubmitted, isActive, onTimeUp, difficulty }) => {
  const [question, setQuestion] = useState<Question | null>(null);
  const [userAnswer, setUserAnswer] = useState<string>('');
  const [feedback, setFeedback] = useState<AnswerResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [showFeedback, setShowFeedback] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  const loadNewQuestion = useCallback(async () => {
    try {
      setLoading(true);
      setFeedback(null);
      setShowFeedback(false);
      setUserAnswer('');
      const newQuestion = await apiService.getQuiz(difficulty);
      setQuestion(newQuestion);
      
      // Auto-focus the input field after question loads
      setTimeout(() => {
        if (inputRef.current) {
          requestAnimationFrame(() => {
            inputRef.current?.focus();
            inputRef.current?.scrollIntoView({ behavior: 'smooth', block: 'center' });
          });
        }
      }, 100);
    } catch (error) {
      console.error('Failed to load question:', error);
    } finally {
      setLoading(false);
    }
  }, [difficulty]);

  useEffect(() => {
    if (isActive) {
      loadNewQuestion();
    }
  }, [isActive, loadNewQuestion]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!question || userAnswer.trim() === '') return;

    try {
      setLoading(true);
      const response = await apiService.submitAnswer({
        userAnswer: parseInt(userAnswer),
        questionId: question.id // or include operands/operator if needed
      });
      
      
      setFeedback(response);
      setShowFeedback(true);
      onAnswerSubmitted(response.correct);
      
      // Auto-load next question after showing feedback
      setTimeout(() => {
        loadNewQuestion();
      }, 2000);
      
    } catch (error) {
      console.error('Failed to submit answer:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !loading) {
      handleSubmit(e as any);
    }
  };

  if (!isActive) {
    return (
      <div className="text-center py-8">
        <div className="text-gray-500">Quiz not active</div>
      </div>
    );
  }

  if (loading && !question) {
    return (
      <div className="text-center py-8">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
        <div className="mt-2 text-gray-500">Loading question...</div>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-lg shadow-md p-4">
      <div className="mb-4">
        <Timer isActive={isActive && !showFeedback} onTimeUp={onTimeUp} duration={30} />
      </div>

      {question && (
        <div className="text-center">
          <div className="text-3xl font-bold text-gray-800 mb-4">
            {question.operand1} {question.operator} {question.operand2} = ?
          </div>

          {showFeedback && feedback ? (
            <div className={`p-3 rounded-lg mb-3 ${
              feedback.correct ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
            }`}>
              <div className="text-base font-semibold">
                {feedback.correct ? '✓ Correct!' : '✗ Incorrect'}
              </div>
              {!feedback.correct && (
                <div className="text-sm mt-1">
                  The correct answer was {feedback.correctAnswer}
                </div>
              )}
            </div>
          ) : (
            <form onSubmit={handleSubmit} className="space-y-3">
              <input
                ref={inputRef}
                type="number"
                value={userAnswer}
                onChange={(e) => setUserAnswer(e.target.value)}
                onKeyPress={handleKeyPress}
                placeholder="Enter your answer"
                className="w-full px-3 py-2 text-lg text-center border-2 border-gray-300 rounded-lg focus:border-blue-500 focus:outline-none"
                disabled={loading || showFeedback}
                autoFocus
              />
              
              <button
                type="submit"
                disabled={loading || userAnswer.trim() === '' || showFeedback}
                className="w-full bg-blue-600 text-white py-2 px-4 rounded-lg font-semibold hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
              >
                {loading ? 'Submitting...' : 'Submit Answer'}
              </button>
            </form>
          )}
        </div>
      )}
    </div>
  );
};

export default Quiz;
