import React, { useEffect, useState } from 'react';
import { ScoreResponse } from '../services/api';

interface ScoreProps {
  score: ScoreResponse;
}

const Score: React.FC<ScoreProps> = ({ score }) => {
  const [showCelebration, setShowCelebration] = useState(false);

  useEffect(() => {
    if (score.isNewRecord) {
      setShowCelebration(true);
      // Hide celebration after 3 seconds
      const timer = setTimeout(() => {
        setShowCelebration(false);
      }, 3000);
      return () => clearTimeout(timer);
    }
  }, [score.isNewRecord]);

  return (
    <div className="bg-white rounded-lg shadow-md p-4 relative overflow-hidden">
      {/* Celebration Animation */}
      {showCelebration && (
        <div className="absolute inset-0 bg-gradient-to-r from-yellow-400 via-pink-500 to-purple-600 opacity-20 animate-pulse z-0">
          <div className="absolute inset-0 flex items-center justify-center">
            <div className="text-6xl animate-bounce">🎉</div>
          </div>
        </div>
      )}
      
      <div className="relative z-10">
        <h2 className="text-xl font-bold text-gray-800 mb-3">
          {showCelebration ? '🏆 NEW HIGH SCORE! 🏆' : 'Your Score'}
        </h2>
        
        <div className="grid grid-cols-2 gap-3">
          {/* Current Score */}
          <div className="text-center">
            <div className="bg-blue-50 rounded-lg p-3">
              <div className={`text-2xl font-bold ${
                showCelebration ? 'text-yellow-600 animate-pulse' : 'text-blue-600'
              }`}>
                {score.currentScore}
              </div>
              <div className="text-sm text-gray-600">Current</div>
            </div>
          </div>
          
          {/* High Scores */}
          <div className="text-center">
            <div className="bg-green-50 rounded-lg p-3">
              <div className="text-lg font-bold text-green-600">
              {score.highScores
                ? Math.max(...Object.values(score.highScores))
                : 0}
              </div>
              <div className="text-sm text-gray-600">Best Overall</div>
            </div>
          </div>
        </div>
        
        {/* Difficulty High Scores */}
        <div className="mt-3">
          <div className="text-sm font-medium text-gray-700 mb-2">High Scores by Difficulty:</div>
          <div className="grid grid-cols-3 gap-2 text-xs">
            <div className="bg-green-100 rounded p-2 text-center">
              <div className="font-bold text-green-700">{score.highScores?.easy || 0}</div>
              <div className="text-green-600">Easy</div>
            </div>
            <div className="bg-blue-100 rounded p-2 text-center">
              <div className="font-bold text-blue-700">{score.highScores?.normal || 0}</div>
              <div className="text-blue-600">Normal</div>
            </div>
            <div className="bg-red-100 rounded p-2 text-center">
              <div className="font-bold text-red-700">{score.highScores?.hard || 0}</div>
              <div className="text-red-600">Hard</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Score;
