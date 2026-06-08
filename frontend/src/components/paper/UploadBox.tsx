import React, { useState, useRef, useCallback } from 'react';
import { Upload, FileText, X, AlertTriangle } from 'lucide-react';
import { formatFileSize } from '../../utils/helpers';

interface UploadBoxProps {
  onFileSelect: (file: File) => void;
  uploading?: boolean;
  progress?: number;
}

const MAX_FILE_SIZE = 50 * 1024 * 1024; // 50MB

const UploadBox: React.FC<UploadBoxProps> = ({ onFileSelect, uploading = false, progress = 0 }) => {
  const [isDragging, setIsDragging] = useState(false);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [error, setError] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const validateFile = useCallback((file: File): boolean => {
    setError(null);
    if (!file.name.toLowerCase().endsWith('.pdf')) {
      setError('Only PDF files are accepted.');
      return false;
    }
    if (file.size > MAX_FILE_SIZE) {
      setError(`File size must be less than ${formatFileSize(MAX_FILE_SIZE)}.`);
      return false;
    }
    return true;
  }, []);

  const handleFile = useCallback(
    (file: File) => {
      if (validateFile(file)) {
        setSelectedFile(file);
        onFileSelect(file);
      }
    },
    [validateFile, onFileSelect]
  );

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
    const file = e.dataTransfer.files[0];
    if (file) handleFile(file);
  };

  const handleClick = () => {
    fileInputRef.current?.click();
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) handleFile(file);
  };

  const clearFile = () => {
    setSelectedFile(null);
    setError(null);
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  return (
    <div className="w-full">
      <div
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
        onClick={handleClick}
        className={`relative border-2 border-dashed rounded-2xl p-8 sm:p-12 text-center cursor-pointer transition-all duration-300 ${
          isDragging
            ? 'border-primary-500 bg-primary-50 dark:bg-primary-900/20 scale-[1.02]'
            : selectedFile
            ? 'border-emerald-400 bg-emerald-50/50 dark:bg-emerald-900/10'
            : 'border-slate-300 dark:border-slate-600 hover:border-primary-400 hover:bg-slate-50 dark:hover:bg-slate-800/50'
        } ${uploading ? 'pointer-events-none' : ''}`}
      >
        <input
          ref={fileInputRef}
          type="file"
          accept=".pdf"
          onChange={handleInputChange}
          className="hidden"
        />

        {selectedFile ? (
          <div className="flex flex-col items-center gap-3">
            <div className="w-16 h-16 bg-emerald-100 dark:bg-emerald-900/30 rounded-2xl flex items-center justify-center">
              <FileText className="w-8 h-8 text-emerald-600 dark:text-emerald-400" />
            </div>
            <div>
              <p className="font-semibold text-slate-900 dark:text-slate-100">
                {selectedFile.name}
              </p>
              <p className="text-sm text-slate-500 dark:text-slate-400">
                {formatFileSize(selectedFile.size)}
              </p>
            </div>
            {!uploading && (
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  clearFile();
                }}
                className="flex items-center gap-1 text-sm text-rose-500 hover:text-rose-600 transition-colors"
              >
                <X className="w-4 h-4" />
                Remove
              </button>
            )}
          </div>
        ) : (
          <div className="flex flex-col items-center gap-4">
            <div
              className={`w-16 h-16 rounded-2xl flex items-center justify-center transition-colors ${
                isDragging
                  ? 'bg-primary-100 dark:bg-primary-800/30'
                  : 'bg-slate-100 dark:bg-slate-800'
              }`}
            >
              <Upload
                className={`w-8 h-8 transition-colors ${
                  isDragging
                    ? 'text-primary-600 dark:text-primary-400'
                    : 'text-slate-400 dark:text-slate-500'
                }`}
              />
            </div>
            <div>
              <p className="font-semibold text-slate-900 dark:text-slate-100 mb-1">
                {isDragging ? 'Drop your PDF here' : 'Drag & drop your PDF here'}
              </p>
              <p className="text-sm text-slate-500 dark:text-slate-400">
                or click to browse · Max 50MB
              </p>
            </div>
          </div>
        )}

        {/* Upload progress bar */}
        {uploading && (
          <div className="mt-6 w-full max-w-xs mx-auto">
            <div className="flex items-center justify-between text-sm mb-2">
              <span className="text-slate-600 dark:text-slate-400">Uploading...</span>
              <span className="font-medium text-primary-600">{progress}%</span>
            </div>
            <div className="h-2 bg-slate-200 dark:bg-slate-700 rounded-full overflow-hidden">
              <div
                className="h-full bg-gradient-to-r from-primary-500 to-violet-500 rounded-full transition-all duration-300"
                style={{ width: `${progress}%` }}
              />
            </div>
          </div>
        )}
      </div>

      {/* Error message */}
      {error && (
        <div className="flex items-center gap-2 mt-3 px-3 py-2 bg-rose-50 dark:bg-rose-900/20 border border-rose-200 dark:border-rose-800 rounded-lg">
          <AlertTriangle className="w-4 h-4 text-rose-500 flex-shrink-0" />
          <p className="text-sm text-rose-600 dark:text-rose-400">{error}</p>
        </div>
      )}
    </div>
  );
};

export default UploadBox;
