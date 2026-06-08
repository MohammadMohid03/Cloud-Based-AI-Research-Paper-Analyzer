import React, { useState } from 'react';
import {
  User,
  Mail,
  Lock,
  Moon,
  Sun,
  Shield,
  Save,
  Loader2,
} from 'lucide-react';
import { useAuth } from '../hooks/useAuth';
import { Button } from '../components/common/Button';
import api from '../services/api';
import { getErrorMessage } from '../utils/helpers';

const Settings: React.FC = () => {
  const { user, darkMode, toggleDarkMode, addToast } = useAuth();
  const [currentPassword, setCurrentPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmNewPassword, setConfirmNewPassword] = useState('');
  const [changingPassword, setChangingPassword] = useState(false);

  const handleChangePassword = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!currentPassword || !newPassword) {
      addToast('warning', 'Please fill in all password fields.');
      return;
    }
    if (newPassword.length < 6) {
      addToast('warning', 'New password must be at least 6 characters.');
      return;
    }
    if (newPassword !== confirmNewPassword) {
      addToast('error', 'New passwords do not match.');
      return;
    }
    setChangingPassword(true);
    try {
      await api.put('/auth/password', {
        currentPassword,
        newPassword,
      });
      addToast('success', 'Password changed successfully!');
      setCurrentPassword('');
      setNewPassword('');
      setConfirmNewPassword('');
    } catch (err) {
      addToast('error', getErrorMessage(err));
    } finally {
      setChangingPassword(false);
    }
  };

  return (
    <div className="max-w-2xl mx-auto space-y-6 animate-fade-in">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-slate-900 dark:text-white">Settings</h1>
        <p className="text-slate-500 dark:text-slate-400 mt-1">
          Manage your account settings and preferences
        </p>
      </div>

      {/* Profile Information */}
      <div className="glass-card p-6">
        <h2 className="text-lg font-semibold text-slate-900 dark:text-white mb-4 flex items-center gap-2">
          <User className="w-5 h-5 text-primary-600 dark:text-primary-400" />
          Profile Information
        </h2>
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1.5">
              Name
            </label>
            <div className="relative">
              <User className="absolute left-3.5 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-400" />
              <input
                type="text"
                value={user?.name || ''}
                readOnly
                className="input-field pl-11 bg-slate-50 dark:bg-slate-700/50 cursor-not-allowed"
              />
            </div>
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1.5">
              Email
            </label>
            <div className="relative">
              <Mail className="absolute left-3.5 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-400" />
              <input
                type="email"
                value={user?.email || ''}
                readOnly
                className="input-field pl-11 bg-slate-50 dark:bg-slate-700/50 cursor-not-allowed"
              />
            </div>
          </div>
          {user?.createdAt && (
            <p className="text-xs text-slate-500 dark:text-slate-400">
              Member since {new Date(user.createdAt).toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })}
            </p>
          )}
        </div>
      </div>

      {/* Appearance */}
      <div className="glass-card p-6">
        <h2 className="text-lg font-semibold text-slate-900 dark:text-white mb-4 flex items-center gap-2">
          {darkMode ? (
            <Moon className="w-5 h-5 text-primary-600 dark:text-primary-400" />
          ) : (
            <Sun className="w-5 h-5 text-primary-600 dark:text-primary-400" />
          )}
          Appearance
        </h2>
        <div className="flex items-center justify-between">
          <div>
            <p className="font-medium text-slate-900 dark:text-white">Dark Mode</p>
            <p className="text-sm text-slate-500 dark:text-slate-400">
              Toggle between light and dark theme
            </p>
          </div>
          <button
            onClick={toggleDarkMode}
            className={`relative inline-flex h-7 w-12 items-center rounded-full transition-colors duration-300 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 ${
              darkMode ? 'bg-primary-600' : 'bg-slate-300'
            }`}
          >
            <span
              className={`inline-block h-5 w-5 transform rounded-full bg-white shadow-md transition-transform duration-300 ${
                darkMode ? 'translate-x-6' : 'translate-x-1'
              }`}
            />
          </button>
        </div>
      </div>

      {/* Change Password */}
      <div className="glass-card p-6">
        <h2 className="text-lg font-semibold text-slate-900 dark:text-white mb-4 flex items-center gap-2">
          <Shield className="w-5 h-5 text-primary-600 dark:text-primary-400" />
          Change Password
        </h2>
        <form onSubmit={handleChangePassword} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1.5">
              Current Password
            </label>
            <div className="relative">
              <Lock className="absolute left-3.5 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-400" />
              <input
                type="password"
                value={currentPassword}
                onChange={(e) => setCurrentPassword(e.target.value)}
                placeholder="••••••••"
                className="input-field pl-11"
              />
            </div>
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1.5">
              New Password
            </label>
            <div className="relative">
              <Lock className="absolute left-3.5 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-400" />
              <input
                type="password"
                value={newPassword}
                onChange={(e) => setNewPassword(e.target.value)}
                placeholder="••••••••"
                className="input-field pl-11"
              />
            </div>
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1.5">
              Confirm New Password
            </label>
            <div className="relative">
              <Lock className="absolute left-3.5 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-400" />
              <input
                type="password"
                value={confirmNewPassword}
                onChange={(e) => setConfirmNewPassword(e.target.value)}
                placeholder="••••••••"
                className="input-field pl-11"
              />
            </div>
          </div>
          <Button
            type="submit"
            variant="primary"
            loading={changingPassword}
            icon={<Save className="w-4 h-4" />}
          >
            Update Password
          </Button>
        </form>
      </div>

      {/* Account Info */}
      <div className="glass-card p-6">
        <h2 className="text-lg font-semibold text-slate-900 dark:text-white mb-3">
          About
        </h2>
        <div className="text-sm text-slate-500 dark:text-slate-400 space-y-1">
          <p>Research Paper Analyzer v1.0</p>
          <p>AI-Powered Analysis · Chat · Quiz</p>
          <p className="text-xs mt-2">Built with React, TypeScript & Tailwind CSS</p>
        </div>
      </div>
    </div>
  );
};

export default Settings;
