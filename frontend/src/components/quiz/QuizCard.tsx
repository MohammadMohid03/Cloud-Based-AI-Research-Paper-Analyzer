import React from 'react';
import { CheckCircle, XCircle } from 'lucide-react';
import { QuizQuestion } from '../../types';

interface QuizCardProps {
  question: QuizQuestion;
  questionNumber: number;
  totalQuestions: number;
  selectedAnswer: number | null;
  onAnswer: (answerIndex: number) => void;
  showResult: boolean;
}

const optionLabels = ['A', 'B', 'C', 'D'];

const QuizCard: React.FC<QuizCardProps> = ({
  question,
  questionNumber,
  totalQuestions,
  selectedAnswer,
  onAnswer,
  showResult,
}) => {
  return (
    <div className="glass-card p-6 sm:p-8 animate-fade-in">
      {/* Progress */}
      <div className="flex items-center justify-between mb-6">
        <span className="text-sm font-medium text-slate-500 dark:text-slate-400">
          Question {questionNumber} of {totalQuestions}
        </span>
        <div className="flex gap-1">
          {Array.from({ length: totalQuestions }).map((_, i) => (
            <div
              key={i}
              className={`w-2 h-2 rounded-full transition-colors ${
                i < questionNumber
                  ? 'bg-primary-500'
                  : 'bg-slate-200 dark:bg-slate-700'
              }`}
            />
          ))}
        </div>
      </div>

      {/* Progress bar */}
      <div className="h-1.5 bg-slate-200 dark:bg-slate-700 rounded-full mb-6 overflow-hidden">
        <div
          className="h-full bg-gradient-to-r from-primary-500 to-violet-500 rounded-full transition-all duration-500"
          style={{ width: `${(questionNumber / totalQuestions) * 100}%` }}
        />
      </div>

      {/* Question */}
      <h3 className="text-lg font-semibold text-slate-900 dark:text-slate-100 mb-6 leading-relaxed">
        {question.question}
      </h3>

      {/* Options */}
      <div className="space-y-3">
        {question.options.map((option, index) => {
          const isSelected = selectedAnswer === index;
          const isCorrect = index === question.correctAnswer;
          const showCorrect = showResult && isCorrect;
          const showIncorrect = showResult && isSelected && !isCorrect;

          return (
            <button
              key={index}
              onClick={() => !showResult && onAnswer(index)}
              disabled={showResult}
              className={`w-full flex items-center gap-3 p-4 rounded-xl border-2 text-left transition-all duration-200 ${
                showCorrect
                  ? 'border-emerald-400 bg-emerald-50 dark:bg-emerald-900/20'
                  : showIncorrect
                  ? 'border-rose-400 bg-rose-50 dark:bg-rose-900/20'
                  : isSelected
                  ? 'border-primary-500 bg-primary-50 dark:bg-primary-900/20'
                  : 'border-slate-200 dark:border-slate-600 hover:border-primary-300 dark:hover:border-primary-600 hover:bg-slate-50 dark:hover:bg-slate-800'
              } ${showResult ? 'cursor-default' : 'cursor-pointer'}`}
            >
              <span
                className={`w-8 h-8 rounded-lg flex items-center justify-center text-sm font-semibold flex-shrink-0 ${
                  showCorrect
                    ? 'bg-emerald-500 text-white'
                    : showIncorrect
                    ? 'bg-rose-500 text-white'
                    : isSelected
                    ? 'bg-primary-600 text-white'
                    : 'bg-slate-100 dark:bg-slate-700 text-slate-600 dark:text-slate-300'
                }`}
              >
                {showCorrect ? (
                  <CheckCircle className="w-4 h-4" />
                ) : showIncorrect ? (
                  <XCircle className="w-4 h-4" />
                ) : (
                  optionLabels[index]
                )}
              </span>
              <span
                className={`text-sm font-medium ${
                  showCorrect
                    ? 'text-emerald-700 dark:text-emerald-300'
                    : showIncorrect
                    ? 'text-rose-700 dark:text-rose-300'
                    : 'text-slate-700 dark:text-slate-200'
                }`}
              >
                {option}
              </span>
            </button>
          );
        })}
      </div>

      {/* Explanation */}
      {showResult && question.explanation && (
        <div className="mt-6 p-4 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-xl animate-fade-in">
          <p className="text-sm font-medium text-blue-800 dark:text-blue-200 mb-1">
            Explanation
          </p>
          <p className="text-sm text-blue-700 dark:text-blue-300 leading-relaxed">
            {question.explanation}
          </p>
        </div>
      )}
    </div>
  );
};

export default QuizCard;
