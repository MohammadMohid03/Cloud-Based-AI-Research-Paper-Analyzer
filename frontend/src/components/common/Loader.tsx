import React from 'react';
import { Loader2 } from 'lucide-react';

interface LoaderProps {
  size?: 'sm' | 'md' | 'lg';
  text?: string;
  fullScreen?: boolean;
}

export const Loader: React.FC<LoaderProps> = ({ size = 'md', text, fullScreen = false }) => {
  const sizeClasses = {
    sm: 'w-5 h-5',
    md: 'w-8 h-8',
    lg: 'w-12 h-12',
  };

  const content = (
    <div className="flex flex-col items-center justify-center gap-3">
      <Loader2 className={`${sizeClasses[size]} text-primary-600 animate-spin`} />
      {text && (
        <p className="text-sm text-slate-500 dark:text-slate-400 animate-pulse">{text}</p>
      )}
    </div>
  );

  if (fullScreen) {
    return (
      <div className="fixed inset-0 z-50 flex items-center justify-center bg-white/80 dark:bg-slate-900/80 backdrop-blur-sm">
        {content}
      </div>
    );
  }

  return <div className="flex items-center justify-center py-12">{content}</div>;
};

// Skeleton loaders for different content types
export const SkeletonCard: React.FC = () => (
  <div className="glass-card p-6 animate-pulse">
    <div className="h-4 bg-slate-200 dark:bg-slate-700 rounded w-3/4 mb-4" />
    <div className="h-3 bg-slate-200 dark:bg-slate-700 rounded w-1/2 mb-3" />
    <div className="h-3 bg-slate-200 dark:bg-slate-700 rounded w-1/3 mb-4" />
    <div className="flex gap-2">
      <div className="h-8 bg-slate-200 dark:bg-slate-700 rounded w-20" />
      <div className="h-8 bg-slate-200 dark:bg-slate-700 rounded w-20" />
    </div>
  </div>
);

export const SkeletonText: React.FC<{ lines?: number }> = ({ lines = 3 }) => (
  <div className="animate-pulse space-y-3">
    {Array.from({ length: lines }).map((_, i) => (
      <div
        key={i}
        className="h-3 bg-slate-200 dark:bg-slate-700 rounded"
        style={{ width: `${Math.random() * 40 + 60}%` }}
      />
    ))}
  </div>
);

export const SkeletonStat: React.FC = () => (
  <div className="glass-card p-6 animate-pulse flex items-center gap-4">
    <div className="w-12 h-12 bg-slate-200 dark:bg-slate-700 rounded-xl" />
    <div className="flex-1">
      <div className="h-3 bg-slate-200 dark:bg-slate-700 rounded w-20 mb-2" />
      <div className="h-6 bg-slate-200 dark:bg-slate-700 rounded w-12" />
    </div>
  </div>
);
