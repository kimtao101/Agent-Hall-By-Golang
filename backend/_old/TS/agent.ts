// 首先加载环境变量
import './env';

import { openaiService } from './openaiService';

interface Message {
  role: 'user' | 'assistant' | 'system';
  content: string;
}

export class Agent {
  private messages: Message[];

  constructor() {
    this.messages = [];
    this.initializeSystemPrompt();
  }

  private initializeSystemPrompt(): void {
    this.messages.push({
      role: 'system',
      content: 'You are a helpful research assistant. Provide detailed and accurate responses to user queries.'
    });
  }

  public addMessage(role: 'user' | 'assistant', content: string): void {
    this.messages.push({ role, content });
    this.trimHistory();
  }

  private trimHistory(): void {
    if (this.messages.length > 20) {
      this.messages = [this.messages[0], ...this.messages.slice(-19)];
    }
  }

  public async generateResponse(
    userInput: string,
    onChunk: (chunk: string) => void
  ): Promise<void> {
    this.addMessage('user', userInput);

    try {
      const fullResponse = await openaiService.createChatCompletionStream(
        {
          messages: this.messages,
          model: 'deepseek-chat',
          temperature: 0.7,
          stream: true
        },
        onChunk
      );

      this.addMessage('assistant', fullResponse);
    } catch (error) {
      console.error('Error generating response:', error);
      const errorMessage = 'I apologize, but I encountered an error while processing your request. Please try again later.';
      onChunk(errorMessage);
      this.addMessage('assistant', errorMessage);
    }
  }

  public getHistory(): Message[] {
    return [...this.messages];
  }

  public clearHistory(): void {
    this.messages = [];
    this.initializeSystemPrompt();
  }
}