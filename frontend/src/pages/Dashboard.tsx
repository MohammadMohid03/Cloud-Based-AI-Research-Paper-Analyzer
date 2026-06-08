import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  FileText,
  Upload,
  CheckCircle,
  Clock,
  Search,
  Plus,
  BarChart3,
  Sparkles,
} from 'lucide-react';
import { useAuth } from '../hooks/useAuth';
import { paperService } from '../services/paperService';
import { Paper } from '../types';
import PaperCard from '../components/paper/PaperCard';
import { Loader, SkeletonCard, SkeletonStat } from '../components/common/Loader';
import { ErrorMessage } from '../components/common/ErrorMessage';
import { getErrorMessage } from '../utils/helpers';

const Dashboard: React.FC = () => {
  const { user, addToast } = useAuth();
  const navigate = useNavigate();
  const [papers, setPapers] = useState<Paper[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [searchQuery, setSearchQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState<string>('all');

  const fetchPapers = async () => {
    setLoading(true);
    setError('');
    try {
      const data = await paperService.getAllPapers();
      setPapers(data);
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPapers();
  }, []);

  const handleDelete = async (id: string) => {
    if (!window.confirm('Are you sure you want to delete this paper?')) return;
    try {
      await paperService.deletePaper(id);
      setPapers((prev) => prev.filter((p) => p.id !== id));
      addToast('success', 'Paper deleted successfully');
    } catch (err) {
      addToast('error', getErrorMessage(err));
    }
  };

  // Stats
  const totalPapers = papers.length;
  const analyzedPapers = papers.filter((p) => p.status === 'completed').length;
  const pendingPapers = papers.filter((p) => p.status === 'pending' || p.status === 'processing').length;

  // Filtered papers
  const filteredPapers = papers.filter((p) => {
    const matchesSearch =
      !searchQuery ||
      p.title?.toLowerCase().includes(searchQuery.toLowerCase()) ||
      p.fileName?.toLowerCase().includes(searchQuery.toLowerCase());
    const matchesStatus = statusFilter === 'all' || p.status === statusFilter;
    return matchesSearch && matchesStatus;
  });

  const stats = [
    {
      label: 'Total Papers',
      value: totalPapers,
      icon: FileText,
      color: 'from-primary-500 to-primary-600',
      bg: 'bg-primary-50 dark:bg-primary-900/20',
    },
    {
      label: 'Analyzed',
      value: analyzedPapers,
      icon: CheckCircle,
      color: 'from-emerald-500 to-emerald-600',
      bg: 'bg-emerald-50 dark:bg-emerald-900/20',
    },
    {
      label: 'Pending',
      value: pendingPapers,
      icon: Clock,
      color: 'from-amber-500 to-amber-600',
      bg: 'bg-amber-50 dark:bg-amber-900/20',
    },
  ];

  return (
    <div className="space-y-8 animate-fade-in">
      {/* Welcome Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl sm:text-3xl font-bold text-slate-900 dark:text-white">
            Welcome back, <span className="text-gradient">{user?.name?.split(' ')[0] || 'Researcher'}</span>
          </h1>
          <p className="text-slate-500 dark:text-slate-400 mt-1">
            Here's an overview of your research papers
          </p>
        </div>
        <button
          onClick={() => navigate('/upload')}
          className="btn-primary"
        >
          <Plus className="w-4 h-4" />
          Upload Paper
        </button>
      </div>

      {/* Stats Cards */}
      {loading ? (
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
          {[1, 2, 3].map((i) => <SkeletonStat key={i} />)}
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
          {stats.map((stat) => (
            <div key={stat.label} className="stat-card">
              <div className={`w-12 h-12 ${stat.bg} rounded-xl flex items-center justify-center`}>
                <stat.icon className={`w-6 h-6 bg-gradient-to-br ${stat.color} bg-clip-text text-transparent`} style={{ color: stat.color.includes('primary') ? '#6366f1' : stat.color.includes('emerald') ? '#10b981' : '#f59e0b' }} />
              </div>
              <div>
                <p className="text-sm text-slate-500 dark:text-slate-400">{stat.label}</p>
                <p className="text-2xl font-bold text-slate-900 dark:text-white">{stat.value}</p>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Filters */}
      {!loading && papers.length > 0 && (
        <div className="flex flex-col sm:flex-row gap-3">
          <div className="relative flex-1">
            <Search className="absolute left-3.5 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
            <input
              type="text"
              placeholder="Search papers..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="input-field pl-10 text-sm"
            />
          </div>
          <div className="flex gap-2">
            {['all', 'completed', 'pending', 'processing', 'failed'].map((status) => (
              <button
                key={status}
                onClick={() => setStatusFilter(status)}
                className={`px-4 py-2.5 text-sm font-medium rounded-lg transition-all duration-200 capitalize ${
                  statusFilter === status
                    ? 'bg-primary-600 text-white shadow-md'
                    : 'bg-white dark:bg-slate-800 text-slate-600 dark:text-slate-300 border border-slate-200 dark:border-slate-700 hover:border-primary-300 dark:hover:border-primary-600'
                }`}
              >
                {status}
              </button>
            ))}
          </div>
        </div>
      )}

      {/* Papers Grid */}
      {loading ? (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {[1, 2, 3, 4, 5, 6].map((i) => <SkeletonCard key={i} />)}
        </div>
      ) : error ? (
        <ErrorMessage message={error} onRetry={fetchPapers} />
      ) : filteredPapers.length === 0 ? (
        <div className="text-center py-16">
          <div className="w-20 h-20 bg-slate-100 dark:bg-slate-800 rounded-2xl flex items-center justify-center mx-auto mb-4">
            {papers.length === 0 ? (
              <Sparkles className="w-10 h-10 text-slate-400" />
            ) : (
              <Search className="w-10 h-10 text-slate-400" />
            )}
          </div>
          <h3 className="text-lg font-semibold text-slate-900 dark:text-slate-100 mb-2">
            {papers.length === 0 ? 'No papers yet' : 'No matching papers'}
          </h3>
          <p className="text-slate-500 dark:text-slate-400 max-w-sm mx-auto mb-6">
            {papers.length === 0
              ? 'Upload your first research paper to get started with AI-powered analysis.'
              : 'Try adjusting your search or filter criteria.'}
          </p>
          {papers.length === 0 && (
            <button onClick={() => navigate('/upload')} className="btn-primary">
              <Upload className="w-4 h-4" />
              Upload Your First Paper
            </button>
          )}
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {filteredPapers.map((paper) => (
            <PaperCard key={paper.id} paper={paper} onDelete={handleDelete} />
          ))}
        </div>
      )}
    </div>
  );
};

export default Dashboard;
