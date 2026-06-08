import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { ArrowLeft, ArrowRight, Brain, RefreshCw } from 'lucide-react';
import QuizCard from '../components/quiz/QuizCard';
import QuizResult from '../components/quiz/QuizResult';
import { Loader } from '../components/common/Loader';
import { ErrorMessage } from '../components/common/ErrorMessage';
import { Button } from '../components/common/Button';
import { quizService } from '../services/quizService';
import { paperService } from '../services/paperService';
import { QuizQuestion, QuizAttempt } from '../types';
import { getErrorMessage } from '../utils/helpers';
import { useAuth } from '../hooks/useAuth';

const Quiz: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { addToast } = useAuth();
  const [paperTitle, setPaperTitle] = useState('');
  const [questions, setQuestions] = useState<QuizQuestion[]>([]);
  const [currentIndex, setCurrentIndex] = useState(0);
  const [selectedAnswer, setSelectedAnswer] = useState<number | null>(null);
  const [showResult, setShowResult] = useState(false);
  const [attempts, setAttempts] = useState<QuizAttempt[]>([]);
  const [quizComplete, setQuizComplete] = useState(false);
  const [loading, setLoading] = useState(true);
  const [generating, setGenerating] = useState(false);
  const [error, setError] = useState('');

  const fetchQuiz = async () => {
    if (!id) return;
    setLoading(true);
    setError('');
    try {
      const paper = await paperService.getPaper(id);
      setPaperTitle(paper.title || paper.fileName);
      const data = await quizService.getQuiz(id);
      if (Array.isArray(data) && data.length > 0) {
        setQuestions(data.map((q, i) => ({ ...q, id: q.id || String(i) })));
      }
    } catch (err) {
      // Quiz might not exist yet, that's fine
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchQuiz();
  }, [id]);

  const handleGenerateQuiz = async () => {
    if (!id) return;
    setGenerating(true);
    setError('');
    try {
      const data = await quizService.generateQuiz(id);
      if (Array.isArray(data) && data.length > 0) {
        setQuestions(data.map((q, i) => ({ ...q, id: q.id || String(i) })));
        resetQuiz();
        addToast('success', 'Quiz generated successfully!');
      } else {
        addToast('warning', 'No questions were generated. Make sure the paper is analyzed first.');
      }
    } catch (err) {
      addToast('error', getErrorMessage(err));
    } finally {
      setGenerating(false);
    }
  };

  const handleAnswer = (answerIndex: number) => {
    if (showResult) return;
    setSelectedAnswer(answerIndex);
    setShowResult(true);

    const isCorrect = answerIndex === questions[currentIndex].correctAnswer;
    setAttempts((prev) => [
      ...prev,
      {
        questionId: questions[currentIndex].id,
        selectedAnswer: answerIndex,
        isCorrect,
      },
    ]);
  };

  const handleNext = () => {
    if (currentIndex < questions.length - 1) {
      setCurrentIndex((prev) => prev + 1);
      setSelectedAnswer(null);
      setShowResult(false);
    } else {
      setQuizComplete(true);
    }
  };

  const resetQuiz = () => {
    setCurrentIndex(0);
    setSelectedAnswer(null);
    setShowResult(false);
    setAttempts([]);
    setQuizComplete(false);
  };

  if (loading) return <Loader text="Loading quiz..." />;

  return (
    <div className="max-w-2xl mx-auto space-y-6 animate-fade-in">
      {/* Header */}
      <div className="flex items-center gap-3">
        <button
          onClick={() => navigate(`/papers/${id}`)}
          className="p-2 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
        >
          <ArrowLeft className="w-5 h-5 text-slate-600 dark:text-slate-300" />
        </button>
        <div className="flex-1">
          <h1 className="text-2xl font-bold text-slate-900 dark:text-white flex items-center gap-2">
            <Brain className="w-6 h-6 text-primary-600 dark:text-primary-400" />
            Paper Quiz
          </h1>
          <p className="text-sm text-slate-500 dark:text-slate-400 truncate">
            {paperTitle}
          </p>
        </div>
        <Button
          variant="secondary"
          size="sm"
          onClick={handleGenerateQuiz}
          loading={generating}
          icon={<RefreshCw className="w-4 h-4" />}
        >
          {questions.length > 0 ? 'New Quiz' : 'Generate'}
        </Button>
      </div>

      {/* No questions state */}
      {questions.length === 0 && !generating && (
        <div className="glass-card p-8 text-center">
          <div className="w-16 h-16 bg-primary-50 dark:bg-primary-900/20 rounded-2xl flex items-center justify-center mx-auto mb-4">
            <Brain className="w-8 h-8 text-primary-600 dark:text-primary-400" />
          </div>
          <h3 className="text-lg font-semibold text-slate-900 dark:text-slate-100 mb-2">
            No Quiz Available
          </h3>
          <p className="text-sm text-slate-500 dark:text-slate-400 max-w-sm mx-auto mb-6">
            Generate a quiz based on the paper content to test your understanding.
          </p>
          <Button
            variant="primary"
            onClick={handleGenerateQuiz}
            loading={generating}
            icon={<Brain className="w-4 h-4" />}
          >
            Generate Quiz
          </Button>
        </div>
      )}

      {/* Quiz Complete */}
      {quizComplete && (
        <QuizResult
          attempts={attempts}
          totalQuestions={questions.length}
          onRetry={resetQuiz}
          onBack={() => navigate(`/papers/${id}`)}
        />
      )}

      {/* Quiz Card */}
      {questions.length > 0 && !quizComplete && (
        <>
          <QuizCard
            question={questions[currentIndex]}
            questionNumber={currentIndex + 1}
            totalQuestions={questions.length}
            selectedAnswer={selectedAnswer}
            onAnswer={handleAnswer}
            showResult={showResult}
          />

          {/* Next button */}
          {showResult && (
            <div className="flex justify-end animate-fade-in">
              <Button variant="primary" onClick={handleNext} icon={<ArrowRight className="w-4 h-4" />}>
                {currentIndex < questions.length - 1 ? 'Next Question' : 'See Results'}
              </Button>
            </div>
          )}
        </>
      )}
    </div>
  );
};

export default Quiz;
