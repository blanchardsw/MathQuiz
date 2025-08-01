import React, { useState, useEffect } from 'react';

interface TimerProps {
  isActive: boolean;
  onTimeUp: () => void;
  duration?: number; // in seconds
}

const Timer: React.FC<TimerProps> = ({ isActive, onTimeUp, duration = 30 }) => {
  const [timeLeft, setTimeLeft] = useState(duration);

  useEffect(() => {
    setTimeLeft(duration);
  }, [duration]);

  useEffect(() => {
    let interval: NodeJS.Timeout | null = null;

    if (isActive && timeLeft > 0) {
      interval = setInterval(() => {
        setTimeLeft((time) => {
          if (time <= 1) {
            onTimeUp();
            return 0;
          }
          return time - 1;
        });
      }, 1000);
    }

    return () => {
      if (interval) clearInterval(interval);
    };
  }, [isActive, timeLeft, onTimeUp]);

  const formatTime = (seconds: number): string => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  const getTimerColor = (): string => {
    const percentage = (timeLeft / duration) * 100;
    if (percentage > 50) return 'text-green-600';
    if (percentage > 20) return 'text-yellow-600';
    return 'text-red-600';
  };

  return (
    <div className="text-center">
      <div className={`text-lg font-bold ${getTimerColor()}`}>
        {formatTime(timeLeft)}
      </div>
      <div className="w-full bg-gray-200 rounded-full h-1.5 mt-1">
        <div
          className={`h-1.5 rounded-full transition-all duration-1000 ${
            timeLeft > duration * 0.5 ? 'bg-green-600' :
            timeLeft > duration * 0.2 ? 'bg-yellow-600' : 'bg-red-600'
          }`}
          style={{ width: `${(timeLeft / duration) * 100}%` }}
        ></div>
      </div>
    </div>
  );
};

export default Timer;
