import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { ArrowLeft, Upload, CheckCircle, FileText } from 'lucide-react';
import UploadBox from '../components/paper/UploadBox';
import { Button } from '../components/common/Button';
import { useAuth } from '../hooks/useAuth';
import { paperService } from '../services/paperService';
import { getErrorMessage } from '../utils/helpers';

const UploadPaper: React.FC = () => {
  const navigate = useNavigate();
  const { addToast } = useAuth();
  const [file, setFile] = useState<File | null>(null);
  const [title, setTitle] = useState('');
  const [uploading, setUploading] = useState(false);
  const [progress, setProgress] = useState(0);
  const [uploadedPaperId, setUploadedPaperId] = useState<string | null>(null);

  const handleFileSelect = (selectedFile: File) => {
    setFile(selectedFile);
    if (!title) {
      // Auto-fill title from filename
      const name = selectedFile.name.replace(/\.pdf$/i, '').replace(/[-_]/g, ' ');
      setTitle(name);
    }
  };

  const handleUpload = async () => {
    if (!file) {
      addToast('warning', 'Please select a file to upload');
      return;
    }
    if (!title.trim()) {
      addToast('warning', 'Please enter a title for the paper');
      return;
    }

    setUploading(true);
    setProgress(0);
    try {
      const paper = await paperService.uploadPaper(file, title.trim(), (p) => setProgress(p));
      setUploadedPaperId(paper.id);
      addToast('success', 'Paper uploaded successfully!');
    } catch (err) {
      addToast('error', getErrorMessage(err));
    } finally {
      setUploading(false);
    }
  };

  return (
    <div className="max-w-2xl mx-auto space-y-6 animate-fade-in">
      {/* Header */}
      <div className="flex items-center gap-3">
        <button
          onClick={() => navigate(-1)}
          className="p-2 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
        >
          <ArrowLeft className="w-5 h-5 text-slate-600 dark:text-slate-300" />
        </button>
        <div>
          <h1 className="text-2xl font-bold text-slate-900 dark:text-white">Upload Paper</h1>
          <p className="text-sm text-slate-500 dark:text-slate-400">
            Upload a research paper (PDF) for AI analysis
          </p>
        </div>
      </div>

      {uploadedPaperId ? (
        /* Success state */
        <div className="glass-card p-8 text-center">
          <div className="w-16 h-16 bg-emerald-100 dark:bg-emerald-900/30 rounded-2xl flex items-center justify-center mx-auto mb-4">
            <CheckCircle className="w-8 h-8 text-emerald-500" />
          </div>
          <h2 className="text-xl font-bold text-slate-900 dark:text-white mb-2">
            Upload Successful!
          </h2>
          <p className="text-slate-500 dark:text-slate-400 mb-6">
            Your paper has been uploaded and is ready for analysis.
          </p>
          <div className="flex items-center justify-center gap-3">
            <Button
              variant="secondary"
              onClick={() => {
                setFile(null);
                setTitle('');
                setProgress(0);
                setUploadedPaperId(null);
              }}
              icon={<Upload className="w-4 h-4" />}
            >
              Upload Another
            </Button>
            <Button
              variant="primary"
              onClick={() => navigate(`/papers/${uploadedPaperId}`)}
              icon={<FileText className="w-4 h-4" />}
            >
              View Paper
            </Button>
          </div>
        </div>
      ) : (
        /* Upload form */
        <div className="space-y-6">
          <div className="glass-card p-6">
            <UploadBox onFileSelect={handleFileSelect} uploading={uploading} progress={progress} />
          </div>

          {file && (
            <div className="glass-card p-6 animate-slide-up">
              <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-2">
                Paper Title
              </label>
              <input
                type="text"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                placeholder="Enter a title for the paper"
                className="input-field"
              />
              <div className="mt-4">
                <Button
                  variant="primary"
                  size="lg"
                  onClick={handleUpload}
                  loading={uploading}
                  className="w-full"
                  icon={<Upload className="w-4 h-4" />}
                >
                  {uploading ? 'Uploading...' : 'Upload Paper'}
                </Button>
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default UploadPaper;
