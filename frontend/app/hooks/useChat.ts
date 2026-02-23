import { useState, useCallback } from 'react';

type Message = {
  role: 'user' | 'assistant';
  content: string;
};

export function useChat() {
  const [messages, setMessages] = useState<Message[]>([]);
  const [loading, setLoading] = useState(false);

  const sendMessage = useCallback(async (text: string) => {
    if (!text.trim()) return;

    setLoading(true);
    setMessages(prev => [...prev, { role: 'user', content: text }]);
    setMessages(prev => [...prev, { role: 'assistant', content: '' }]);

    try {
      const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8015';
      const response = await fetch(`${API_URL}/chat`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ message: text }),
      });
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      // 处理流式传输：以消费流式响应而不是等待完整文本。
      const reader = response.body?.getReader();
      if (!reader) {
        throw new Error('No response body');
      }
      
      while (true) {
        const { done, value } = await reader.read();
        if (done) break;
        if (value) {
          const chunk = new TextDecoder().decode(value);
          setMessages(prev => {
            const newMessages = [...prev];
            const lastIndex = newMessages.length - 1;
            if (lastIndex >= 0) {
              newMessages[lastIndex] = {
                ...newMessages[lastIndex],
                content: newMessages[lastIndex].content + chunk
              };
            }
            return newMessages;
          });
        }
      }

    } catch (error) {
      console.error('Chat error:', error);
      setMessages(prev => [...prev, {
        role: 'assistant',
        content: 'Sorry, something went wrong. Please try again later.'
      }]);
    } finally {
      setLoading(false);
    }
  }, []);

  return { messages, sendMessage, loading };
}
