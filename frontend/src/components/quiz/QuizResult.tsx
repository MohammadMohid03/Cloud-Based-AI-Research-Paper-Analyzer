import React from 'react';
import { Trophy, Target, RotateCcw, ArrowLeft, CheckCircle, XCircle } from 'lucide-react';
import { QuizAttempt } from '../../types';
import { Button } from '../common/Button';

interface QuizResultProps {
  attempts: QuizAttempt[];
  totalQuestions: number;
  onRetry: () => void;
  onBack: () => void;
}

const QuizResult: React.FC<QuizResultProps> = ({
  attempts,
  totalQuestions,
  onRetry,
  onBack,
}) => {
  const correctCount = attempts.filter((a) => a.isCorrect).length;
  const percentage = Math.round((correctCount / totalQuestions) * 100);

  const getGrade = () => {
    if (percentage >= 90) return { label: 'Excellent!', color: 'text-emerald-500', emoji: '🎉' };
    if (percentage >= 70) return { label: 'Great Job!', color: 'text-blue-500', emoji: '👏' };
    if (percentage >= 50) return { label: 'Good Effort!', color: 'text-amber-500', emoji: '💪' };
    return { label: 'Keep Learning!', color: 'text-rose-500', emoji: '📚' };
  };

  const grade = getGrade();

  return (
    <div className="glass-card p-8 sm:p-10 text-center animate-slide-up">
      {/* Trophy/Score icon */}
      <div className="mb-6">
        <div className="w-20 h-20 mx-auto bg-gradient-to-br from-primary-100 to-violet-100 dark:from-primary-900/30 dark:to-violet-900/30 rounded-2xl flex items-center justify-center">
          <Trophy className="w-10 h-10 text-primary-600 dark:text-primary-400" />
        </div>
      </div>

      {/* Grade */}
      <p className="text-4xl mb-2">{grade.emoji}</p>
      <h2 className={`text-2xl font-bold mb-2 ${grade.color}`}>{grade.label}</h2>

      {/* Score */}
      <div className="flex items-center justify-center gap-2 mb-8">
        <Target className="w-5 h-5 text-slate-400" />
        <span className="text-3xl font-bold text-slate-900 dark:text-slate-100">
          {correctCount}
        </span>
        <span className="text-lg text-slate-400">/</span>
        <span className="text-lg text-slate-500">{totalQuestions}</span>
        <span className="ml-2 text-sm text-slate-500">({percentage}%)</span>
      </div>

      {/* Score ring */}
      <div className="relative w-32 h-32 mx-auto mb-8">
        <svg className="w-full h-full -rotate-90" viewBox="0 0 36 36">
          <circle
            cx="18"
            cy="18"
            r="16"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            className="text-slate-200 dark:text-slate-700"
          />
          <circle
            cx="18"
            cy="18"
            r="16"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeDasharray={`${percentage} 100`}
            strokeLinecap="round"
            className={
              percentage >= 70
                ? 'text-emerald-500'
                : percentage >= 50
                ? 'text-amber-500'
                : 'text-rose-500'
            }
          />
        </svg>
        <span className="absolute inset-0 flex items-center justify-center text-2xl font-bold text-slate-900 dark:text-slate-100">
          {percentage}%
        </span>
      </div>

      {/* Answer breakdown */}
      <div className="flex items-center justify-center gap-6 mb-8">
        <div className="flex items-center gap-2">
          <CheckCircle className="w-5 h-5 text-emerald-500" />
          <span className="text-sm text-slate-600 dark:text-slate-400">
            {correctCount} correct
          </span>
        </div>
        <div className="flex items-center gap-2">
          <XCircle className="w-5 h-5 text-rose-500" />
          <span className="text-sm text-slate-600 dark:text-slate-400">
            {totalQuestions - correctCount} incorrect
          </span>
        </div>
      </div>

      {/* Actions */}
      <div className="flex items-center justify-center gap-3">
        <Button variant="secondary" onClick={onBack} icon={<ArrowLeft className="w-4 h-4" />}>
          Back to Paper
        </Button>
        <Button variant="primary" onClick={onRetry} icon={<RotateCcw className="w-4 h-4" />}>
          Try Again
        </Button>
      </div>
    </div>
  );
};

export default QuizResult;
