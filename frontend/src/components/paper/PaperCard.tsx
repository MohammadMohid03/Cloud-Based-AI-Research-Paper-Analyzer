import React from 'react';
import { useNavigate } from 'react-router-dom';
import {
  FileText,
  Clock,
  Eye,
  Sparkles,
  Trash2,
  MessageSquare,
  Brain,
} from 'lucide-react';
import { Paper } from '../../types';
import { formatDate, getStatusColor, getStatusLabel } from '../../utils/helpers';

interface PaperCardProps {
  paper: Paper;
  onDelete?: (id: string) => void;
}

const PaperCard: React.FC<PaperCardProps> = ({ paper, onDelete }) => {
  const navigate = useNavigate();

  return (
    <div className="glass-card p-5 hover:shadow-xl transition-all duration-300 hover:-translate-y-0.5 group">
      {/* Header */}
      <div className="flex items-start justify-between mb-3">
        <div className="w-10 h-10 bg-primary-50 dark:bg-primary-900/20 rounded-xl flex items-center justify-center flex-shrink-0">
          <FileText className="w-5 h-5 text-primary-600 dark:text-primary-400" />
        </div>
        <span className={getStatusColor(paper.status)}>
          {getStatusLabel(paper.status)}
        </span>
      </div>

      {/* Title */}
      <h3 className="font-semibold text-slate-900 dark:text-slate-100 mb-2 line-clamp-2 group-hover:text-primary-600 dark:group-hover:text-primary-400 transition-colors">
        {paper.title || paper.fileName}
      </h3>

      {/* Meta */}
      <div className="flex items-center gap-1.5 text-xs text-slate-500 dark:text-slate-400 mb-4">
        <Clock className="w-3.5 h-3.5" />
        <span>{formatDate(paper.uploadedAt)}</span>
        {paper.fileName && (
          <>
            <span className="mx-1">·</span>
            <span className="truncate max-w-[120px]">{paper.fileName}</span>
          </>
        )}
      </div>

      {/* Actions */}
      <div className="flex items-center gap-2 pt-3 border-t border-slate-100 dark:border-slate-700/50">
        <button
          onClick={() => navigate(`/papers/${paper.id}`)}
          className="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-primary-600 dark:text-primary-400 hover:bg-primary-50 dark:hover:bg-primary-900/20 rounded-lg transition-colors"
        >
          <Eye className="w-3.5 h-3.5" />
          View
        </button>

        {paper.status === 'completed' && (
          <>
            <button
              onClick={() => navigate(`/papers/${paper.id}/chat`)}
              className="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-violet-600 dark:text-violet-400 hover:bg-violet-50 dark:hover:bg-violet-900/20 rounded-lg transition-colors"
            >
              <MessageSquare className="w-3.5 h-3.5" />
              Chat
            </button>
            <button
              onClick={() => navigate(`/papers/${paper.id}/quiz`)}
              className="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-emerald-600 dark:text-emerald-400 hover:bg-emerald-50 dark:hover:bg-emerald-900/20 rounded-lg transition-colors"
            >
              <Brain className="w-3.5 h-3.5" />
              Quiz
            </button>
          </>
        )}

        {paper.status === 'pending' && (
          <button
            onClick={() => navigate(`/papers/${paper.id}`)}
            className="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-amber-600 dark:text-amber-400 hover:bg-amber-50 dark:hover:bg-amber-900/20 rounded-lg transition-colors"
          >
            <Sparkles className="w-3.5 h-3.5" />
            Analyze
          </button>
        )}

        <div className="flex-1" />

        {onDelete && (
          <button
            onClick={() => onDelete(paper.id)}
            className="p-1.5 text-slate-400 hover:text-rose-500 hover:bg-rose-50 dark:hover:bg-rose-900/20 rounded-lg transition-colors"
          >
            <Trash2 className="w-3.5 h-3.5" />
          </button>
        )}
      </div>
    </div>
  );
};

export default PaperCard;
