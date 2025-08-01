import React, { useState, useEffect } from 'react';
import Quiz from '../components/Quiz';
import Score from '../components/Score';
import { ScoreResponse, Difficulty, apiService } from '../services/api';

const Home: React.FC = () => {
  const [isQuizActive, setIsQuizActive] = useState(false);
  const [difficulty, setDifficulty] = useState<Difficulty>('normal');
  const [score, setScore] = useState<ScoreResponse>({
    currentScore: 0,
    highScores: { easy: 0, normal: 0, hard: 0 },
    isNewRecord: false
  });

  const loadScore = async () => {
    try {
      const currentScore = await apiService.getScore();
      setScore(currentScore);
    } catch (error) {
      console.error('Failed to load score:', error);
    }
  };

  useEffect(() => {
    loadScore();
  }, []);

  const handleStartQuiz = () => {
    setIsQuizActive(true);
  };

  const handleStopQuiz = async () => {
    setIsQuizActive(false);
    // Reset score when quiz ends
    try {
      const resetResponse = await apiService.resetScore(difficulty);
      setScore(resetResponse);
    } catch (error) {
      console.error('Failed to reset score:', error);
    }
  };

  const handleAnswerSubmitted = (correct: boolean) => {
    // Reload score after each answer
    loadScore();
  };

  const handleTimeUp = async () => {
    setIsQuizActive(false);
    // Reset score when time is up
    try {
      const resetResponse = await apiService.resetScore(difficulty);
      setScore(resetResponse);
    } catch (error) {
      console.error('Failed to reset score:', error);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
      <div className="container mx-auto px-4 py-4">
        <header className="text-center mb-4">
          <h1 className="text-3xl font-bold text-gray-800 mb-1">
            Mental Math Trainer
          </h1>
          <p className="text-sm text-gray-600">
            Improve your mental math skills with timed exercises
          </p>
        </header>

        <div className="max-w-2xl mx-auto space-y-4">
          <Score score={score} />

          {!isQuizActive ? (
            <div className="bg-white rounded-lg shadow-md p-4 text-center">
              <h2 className="text-xl font-bold text-gray-800 mb-3">
                Ready to Practice?
              </h2>
              <p className="text-sm text-gray-600 mb-4">
                Test your mental math skills with randomly generated questions.
                You'll have 30 seconds per question!
              </p>
              
              <div className="mb-4">
                <label className="block text-xs font-medium text-gray-700 mb-2">
                  Choose Difficulty:
                </label>
                <div className="flex justify-center space-x-2">
                  <button
                    onClick={() => setDifficulty('easy')}
                    className={`px-3 py-1.5 text-sm rounded-lg font-medium transition-colors ${
                      difficulty === 'easy'
                        ? 'bg-green-600 text-white'
                        : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
                    }`}
                  >
                    Easy (1-10, +)
                  </button>
                  <button
                    onClick={() => setDifficulty('normal')}
                    className={`px-3 py-1.5 text-sm rounded-lg font-medium transition-colors ${
                      difficulty === 'normal'
                        ? 'bg-blue-600 text-white'
                        : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
                    }`}
                  >
                    Normal (1-20, +/-)
                  </button>
                  <button
                    onClick={() => setDifficulty('hard')}
                    className={`px-3 py-1.5 text-sm rounded-lg font-medium transition-colors ${
                      difficulty === 'hard'
                        ? 'bg-red-600 text-white'
                        : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
                    }`}
                  >
                    Hard (10-99, +/-/*)
                  </button>
                </div>
              </div>
              
              <button
                onClick={handleStartQuiz}
                className="bg-green-600 text-white py-2 px-6 rounded-lg font-semibold hover:bg-green-700 transition-colors"
              >
                Start Quiz
              </button>
            </div>
          ) : (
            <div className="space-y-3">
              <Quiz
                onAnswerSubmitted={handleAnswerSubmitted}
                isActive={isQuizActive}
                onTimeUp={handleTimeUp}
                difficulty={difficulty}
              />
              <div className="text-center">
                <button
                  onClick={handleStopQuiz}
                  className="bg-red-600 text-white py-1.5 px-4 text-sm rounded-lg font-semibold hover:bg-red-700 transition-colors"
                >
                  Stop Quiz
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default Home;
