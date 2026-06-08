import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  ArrowLeft,
  FileText,
  Sparkles,
  MessageSquare,
  Brain,
  Download,
  BookOpen,
  Target,
  Microscope,
  AlertTriangle,
  Telescope,
  Tag,
  Clock,
  Loader2,
} from 'lucide-react';
import { useAuth } from '../hooks/useAuth';
import { paperService } from '../services/paperService';
import { Paper } from '../types';
import { Loader, SkeletonText } from '../components/common/Loader';
import { ErrorMessage } from '../components/common/ErrorMessage';
import { Button } from '../components/common/Button';
import SummaryCard from '../components/paper/SummaryCard';
import { formatDate, getStatusColor, getStatusLabel, getErrorMessage } from '../utils/helpers';

type TabKey = 'summary' | 'findings' | 'methodology' | 'limitations' | 'future' | 'keywords';

const tabs: { key: TabKey; label: string; icon: React.ReactNode }[] = [
  { key: 'summary', label: 'Summary', icon: <BookOpen className="w-4 h-4" /> },
  { key: 'findings', label: 'Key Findings', icon: <Target className="w-4 h-4" /> },
  { key: 'methodology', label: 'Methodology', icon: <Microscope className="w-4 h-4" /> },
  { key: 'limitations', label: 'Limitations', icon: <AlertTriangle className="w-4 h-4" /> },
  { key: 'future', label: 'Future Scope', icon: <Telescope className="w-4 h-4" /> },
  { key: 'keywords', label: 'Keywords', icon: <Tag className="w-4 h-4" /> },
];

const PaperAnalysis: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { addToast } = useAuth();
  const [paper, setPaper] = useState<Paper | null>(null);
  const [loading, setLoading] = useState(true);
  const [analyzing, setAnalyzing] = useState(false);
  const [error, setError] = useState('');
  const [activeTab, setActiveTab] = useState<TabKey>('summary');

  const fetchPaper = async () => {
    if (!id) return;
    setLoading(true);
    setError('');
    try {
      const data = await paperService.getPaper(id);
      setPaper(data);
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPaper();
  }, [id]);

  const handleAnalyze = async () => {
    if (!id) return;
    setAnalyzing(true);
    try {
      await paperService.analyzePaper(id);
      addToast('success', 'Analysis started! This may take a moment...');
      // Poll for results
      const pollInterval = setInterval(async () => {
        try {
          const updatedPaper = await paperService.getPaper(id);
          setPaper(updatedPaper);
          if (updatedPaper.status === 'completed' || updatedPaper.status === 'failed') {
            clearInterval(pollInterval);
            setAnalyzing(false);
            if (updatedPaper.status === 'completed') {
              addToast('success', 'Analysis complete!');
            } else {
              addToast('error', 'Analysis failed. Please try again.');
            }
          }
        } catch {
          clearInterval(pollInterval);
          setAnalyzing(false);
        }
      }, 3000);
      // Stop polling after 2 minutes
      setTimeout(() => {
        clearInterval(pollInterval);
        setAnalyzing(false);
      }, 120000);
    } catch (err) {
      setAnalyzing(false);
      addToast('error', getErrorMessage(err));
    }
  };

  const handleDownloadReport = async () => {
    if (!id) return;
    try {
      const data = await paperService.getReport(id);
      const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = `${paper?.title || 'report'}-analysis.json`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
      addToast('success', 'Report downloaded!');
    } catch (err) {
      addToast('error', getErrorMessage(err));
    }
  };

  if (loading) return <Loader text="Loading paper..." />;
  if (error) return <ErrorMessage message={error} onRetry={fetchPaper} />;
  if (!paper) return <ErrorMessage message="Paper not found." />;

  const analysis = paper.analysis;
  const isAnalyzed = paper.status === 'completed' && analysis;

  return (
    <div className="space-y-6 animate-fade-in">
      {/* Header */}
      <div className="flex items-start gap-3">
        <button
          onClick={() => navigate('/dashboard')}
          className="p-2 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors mt-0.5"
        >
          <ArrowLeft className="w-5 h-5 text-slate-600 dark:text-slate-300" />
        </button>
        <div className="flex-1">
          <h1 className="text-2xl font-bold text-slate-900 dark:text-white">
            {paper.title || paper.fileName}
          </h1>
          <div className="flex flex-wrap items-center gap-3 mt-2">
            <span className={getStatusColor(paper.status)}>
              {getStatusLabel(paper.status)}
            </span>
            <span className="text-sm text-slate-500 dark:text-slate-400 flex items-center gap-1">
              <Clock className="w-3.5 h-3.5" />
              {formatDate(paper.uploadedAt)}
            </span>
            {paper.fileName && (
              <span className="text-sm text-slate-500 dark:text-slate-400 flex items-center gap-1">
                <FileText className="w-3.5 h-3.5" />
                {paper.fileName}
              </span>
            )}
          </div>
        </div>
      </div>

      {/* Action Buttons */}
      <div className="flex flex-wrap gap-3">
        {!isAnalyzed && paper.status !== 'processing' && (
          <Button
            variant="primary"
            onClick={handleAnalyze}
            loading={analyzing}
            icon={<Sparkles className="w-4 h-4" />}
          >
            {analyzing ? 'Analyzing...' : 'Analyze Paper'}
          </Button>
        )}
        {paper.status === 'processing' && (
          <Button variant="secondary" disabled loading>
            Processing...
          </Button>
        )}
        {isAnalyzed && (
          <>
            <Button
              variant="secondary"
              onClick={() => navigate(`/papers/${id}/chat`)}
              icon={<MessageSquare className="w-4 h-4" />}
            >
              Chat with Paper
            </Button>
            <Button
              variant="secondary"
              onClick={() => navigate(`/papers/${id}/quiz`)}
              icon={<Brain className="w-4 h-4" />}
            >
              Take Quiz
            </Button>
            <Button
              variant="ghost"
              onClick={handleDownloadReport}
              icon={<Download className="w-4 h-4" />}
            >
              Download Report
            </Button>
          </>
        )}
      </div>

      {/* Analyzing state */}
      {(analyzing || paper.status === 'processing') && (
        <div className="glass-card p-8 text-center">
          <Loader2 className="w-10 h-10 text-primary-600 animate-spin mx-auto mb-4" />
          <h3 className="text-lg font-semibold text-slate-900 dark:text-slate-100 mb-2">
            Analyzing your paper...
          </h3>
          <p className="text-sm text-slate-500 dark:text-slate-400">
            This may take a minute. We're extracting key insights from your paper.
          </p>
        </div>
      )}

      {/* Analysis Results */}
      {isAnalyzed && analysis && (
        <div className="space-y-6">
          {/* Tabs */}
          <div className="flex overflow-x-auto gap-1 bg-white dark:bg-slate-800 p-1.5 rounded-xl border border-slate-200 dark:border-slate-700">
            {tabs.map((tab) => (
              <button
                key={tab.key}
                onClick={() => setActiveTab(tab.key)}
                className={`flex items-center gap-2 px-4 py-2.5 rounded-lg text-sm font-medium whitespace-nowrap transition-all duration-200 ${
                  activeTab === tab.key
                    ? 'bg-primary-600 text-white shadow-md'
                    : 'text-slate-600 dark:text-slate-300 hover:bg-slate-50 dark:hover:bg-slate-700'
                }`}
              >
                {tab.icon}
                {tab.label}
              </button>
            ))}
          </div>

          {/* Tab Content */}
          <div className="animate-fade-in">
            {activeTab === 'summary' && (
              <SummaryCard
                icon={<BookOpen className="w-5 h-5 text-primary-600 dark:text-primary-400" />}
                title="Summary"
              >
                <p className="whitespace-pre-wrap">{analysis.summary || 'No summary available.'}</p>
              </SummaryCard>
            )}

            {activeTab === 'findings' && (
              <SummaryCard
                icon={<Target className="w-5 h-5 text-emerald-600 dark:text-emerald-400" />}
                title="Key Findings"
              >
                {analysis.keyFindings && analysis.keyFindings.length > 0 ? (
                  <ul className="space-y-3">
                    {analysis.keyFindings.map((finding, i) => (
                      <li key={i} className="flex gap-3">
                        <span className="w-6 h-6 bg-emerald-100 dark:bg-emerald-900/30 rounded-lg flex items-center justify-center flex-shrink-0 text-xs font-bold text-emerald-600 dark:text-emerald-400">
                          {i + 1}
                        </span>
                        <span>{finding}</span>
                      </li>
                    ))}
                  </ul>
                ) : (
                  <p>No key findings available.</p>
                )}
              </SummaryCard>
            )}

            {activeTab === 'methodology' && (
              <SummaryCard
                icon={<Microscope className="w-5 h-5 text-blue-600 dark:text-blue-400" />}
                title="Methodology"
              >
                <p className="whitespace-pre-wrap">
                  {analysis.methodology || 'No methodology information available.'}
                </p>
              </SummaryCard>
            )}

            {activeTab === 'limitations' && (
              <SummaryCard
                icon={<AlertTriangle className="w-5 h-5 text-amber-600 dark:text-amber-400" />}
                title="Limitations"
              >
                {analysis.limitations && analysis.limitations.length > 0 ? (
                  <ul className="space-y-2">
                    {analysis.limitations.map((lim, i) => (
                      <li key={i} className="flex gap-2">
                        <span className="text-amber-500 mt-1">•</span>
                        <span>{lim}</span>
                      </li>
                    ))}
                  </ul>
                ) : (
                  <p>No limitations identified.</p>
                )}
              </SummaryCard>
            )}

            {activeTab === 'future' && (
              <SummaryCard
                icon={<Telescope className="w-5 h-5 text-violet-600 dark:text-violet-400" />}
                title="Future Scope"
              >
                {analysis.futureScope && analysis.futureScope.length > 0 ? (
                  <ul className="space-y-2">
                    {analysis.futureScope.map((scope, i) => (
                      <li key={i} className="flex gap-2">
                        <span className="text-violet-500 mt-1">→</span>
                        <span>{scope}</span>
                      </li>
                    ))}
                  </ul>
                ) : (
                  <p>No future scope suggestions available.</p>
                )}
              </SummaryCard>
            )}

            {activeTab === 'keywords' && (
              <SummaryCard
                icon={<Tag className="w-5 h-5 text-rose-600 dark:text-rose-400" />}
                title="Keywords"
              >
                {analysis.keywords && analysis.keywords.length > 0 ? (
                  <div className="flex flex-wrap gap-2">
                    {analysis.keywords.map((keyword, i) => (
                      <span
                        key={i}
                        className="px-3 py-1.5 bg-slate-100 dark:bg-slate-700 text-slate-700 dark:text-slate-300 rounded-lg text-sm font-medium"
                      >
                        {keyword}
                      </span>
                    ))}
                  </div>
                ) : (
                  <p>No keywords extracted.</p>
                )}
              </SummaryCard>
            )}
          </div>
        </div>
      )}

      {/* Not analyzed yet state */}
      {!isAnalyzed && paper.status !== 'processing' && !analyzing && (
        <div className="glass-card p-8 text-center">
          <div className="w-16 h-16 bg-primary-50 dark:bg-primary-900/20 rounded-2xl flex items-center justify-center mx-auto mb-4">
            <Sparkles className="w-8 h-8 text-primary-600 dark:text-primary-400" />
          </div>
          <h3 className="text-lg font-semibold text-slate-900 dark:text-slate-100 mb-2">
            Ready for Analysis
          </h3>
          <p className="text-sm text-slate-500 dark:text-slate-400 max-w-md mx-auto">
            Click the "Analyze Paper" button above to generate an AI-powered summary, extract key
            findings, methodology, and more.
          </p>
        </div>
      )}
    </div>
  );
};

export default PaperAnalysis;
