import React, { useState } from 'react';
import { ChevronDown, ChevronUp } from 'lucide-react';

interface SummaryCardProps {
  icon: React.ReactNode;
  title: string;
  children: React.ReactNode;
  defaultExpanded?: boolean;
}

const SummaryCard: React.FC<SummaryCardProps> = ({
  icon,
  title,
  children,
  defaultExpanded = true,
}) => {
  const [expanded, setExpanded] = useState(defaultExpanded);

  return (
    <div className="glass-card overflow-hidden transition-all duration-300">
      <button
        onClick={() => setExpanded(!expanded)}
        className="w-full flex items-center justify-between p-5 hover:bg-slate-50 dark:hover:bg-slate-700/30 transition-colors"
      >
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 bg-primary-50 dark:bg-primary-900/20 rounded-xl flex items-center justify-center flex-shrink-0">
            {icon}
          </div>
          <h3 className="font-semibold text-slate-900 dark:text-slate-100">{title}</h3>
        </div>
        {expanded ? (
          <ChevronUp className="w-5 h-5 text-slate-400" />
        ) : (
          <ChevronDown className="w-5 h-5 text-slate-400" />
        )}
      </button>
      {expanded && (
        <div className="px-5 pb-5 animate-fade-in">
          <div className="pl-[52px] text-slate-600 dark:text-slate-300 leading-relaxed text-sm">
            {children}
          </div>
        </div>
      )}
    </div>
  );
};

export default SummaryCard;
