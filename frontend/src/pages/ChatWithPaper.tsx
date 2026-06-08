import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { ArrowLeft } from 'lucide-react';
import ChatBox from '../components/chat/ChatBox';
import { Loader } from '../components/common/Loader';
import { ErrorMessage } from '../components/common/ErrorMessage';
import { paperService } from '../services/paperService';
import { chatService } from '../services/chatService';
import { ChatMessage } from '../types';
import { generateId, getErrorMessage } from '../utils/helpers';
import { useAuth } from '../hooks/useAuth';

const ChatWithPaper: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { addToast } = useAuth();
  const [paperTitle, setPaperTitle] = useState('');
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [loading, setLoading] = useState(true);
  const [sending, setSending] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    const init = async () => {
      if (!id) return;
      setLoading(true);
      try {
        const paper = await paperService.getPaper(id);
        setPaperTitle(paper.title || paper.fileName);

        // Try to load chat history
        try {
          const history = await chatService.getChatHistory(id);
          if (Array.isArray(history) && history.length > 0) {
            setMessages(
              history.map((msg: any) => ({
                id: msg.id || generateId(),
                role: msg.role || (msg.sender === 'user' ? 'user' : 'assistant'),
                content: msg.content || msg.message || msg.text || '',
                timestamp: msg.timestamp || msg.createdAt || new Date().toISOString(),
              }))
            );
          }
        } catch {
          // No history available, that's fine
        }
      } catch (err) {
        setError(getErrorMessage(err));
      } finally {
        setLoading(false);
      }
    };
    init();
  }, [id]);

  const handleSend = async (content: string) => {
    if (!id) return;

    // Add user message immediately
    const userMessage: ChatMessage = {
      id: generateId(),
      role: 'user',
      content,
      timestamp: new Date().toISOString(),
    };
    setMessages((prev) => [...prev, userMessage]);
    setSending(true);

    try {
      const response = await chatService.sendMessage(id, content);
      const aiMessage: ChatMessage = {
        id: generateId(),
        role: 'assistant',
        content: response.answer || response.toString(),
        timestamp: new Date().toISOString(),
      };
      setMessages((prev) => [...prev, aiMessage]);
    } catch (err) {
      addToast('error', getErrorMessage(err));
      // Add error message in chat
      const errorMessage: ChatMessage = {
        id: generateId(),
        role: 'assistant',
        content: 'Sorry, I encountered an error processing your question. Please try again.',
        timestamp: new Date().toISOString(),
      };
      setMessages((prev) => [...prev, errorMessage]);
    } finally {
      setSending(false);
    }
  };

  if (loading) return <Loader text="Loading chat..." />;
  if (error) return <ErrorMessage message={error} onRetry={() => window.location.reload()} />;

  return (
    <div className="h-[calc(100vh-8rem)] flex flex-col animate-fade-in">
      {/* Header */}
      <div className="flex items-center gap-3 mb-4">
        <button
          onClick={() => navigate(`/papers/${id}`)}
          className="p-2 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
        >
          <ArrowLeft className="w-5 h-5 text-slate-600 dark:text-slate-300" />
        </button>
        <div>
          <h1 className="text-xl font-bold text-slate-900 dark:text-white">
            Chat with Paper
          </h1>
          <p className="text-sm text-slate-500 dark:text-slate-400 truncate max-w-md">
            {paperTitle}
          </p>
        </div>
      </div>

      {/* Chat */}
      <div className="flex-1 min-h-0">
        <ChatBox
          messages={messages}
          onSend={handleSend}
          loading={sending}
          paperTitle={paperTitle}
        />
      </div>
    </div>
  );
};

export default ChatWithPaper;
