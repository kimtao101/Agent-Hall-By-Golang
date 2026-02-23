import OpenAI from "openai";
import logger from './logger';


// 创建OpenAI实例
const apiKey = process.env.DEEPSEEK_API_KEY || '';
if (!apiKey) {
  console.error('DEEPSEEK_API_KEY environment variable is not set');
  throw new Error('DEEPSEEK_API_KEY environment variable is required');
}

// 直接创建OpenAI实例，确保在模块顶层正确初始化
const openai = new OpenAI({
  apiKey: apiKey,
  baseURL: 'https://api.deepseek.com'
});

console.log('OpenAI instance created successfully with DeepSeek API');
console.log('Base URL:', 'https://api.deepseek.com');
console.log('API Key configured:', apiKey ? 'Yes' : 'No');

interface ChatRequest {
  messages: Array<{
    role: 'system' | 'user' | 'assistant';
    content: string;
  }>;
  model?: string;
  temperature?: number;
  stream?: boolean;
}

interface ChatResponse {
  id: string;
  object: string;
  created: number;
  model: string;
  choices: Array<{
    index: number;
    message: {
      role: 'assistant';
      content: string | null;
    };
    finish_reason: string;
  }>;
  usage?: {
    prompt_tokens: number;
    completion_tokens: number;
    total_tokens: number;
  };
}

export class OpenAIService {
  /**
   * 发送聊天请求
   * @param request 聊天请求参数
   * @returns 聊天响应
   */
  public async createChatCompletion(request: ChatRequest): Promise<ChatResponse> {
    try {
      const startTime = Date.now();
      logger.info('Creating chat completion', {
        request: {
          ...request,
          messages: request.messages.map(msg => ({
            role: msg.role,
            content: msg.role === 'user' ? msg.content.substring(0, 100) + '...' : msg.content
          }))
        }
      });

      const completion = await openai.chat.completions.create({
        messages: request.messages,
        model: request.model || 'deepseek-chat',
        temperature: request.temperature || 0.7,
        stream: false
      }) as OpenAI.ChatCompletion;

      const endTime = Date.now();
      logger.info('Chat completion created successfully', {
        response: {
          id: completion.id,
          model: completion.model,
          choices: completion.choices.map((choice: any) => ({
            index: choice.index,
            message: {
              role: choice.message.role,
              content: choice.message.content.substring(0, 100) + '...'
            },
            finish_reason: choice.finish_reason
          })),
          usage: completion.usage
        },
        duration: endTime - startTime
      });
      console.log(" openai.chat.completions",completion)
      return completion;
    } catch (error) {
      logger.error('Error creating chat completion', {
        error: error instanceof Error ? error.message : String(error),
        stack: error instanceof Error ? error.stack : undefined
      });
      throw error;
    }
  }

  /**
   * 流式发送聊天请求
   * @param request 聊天请求参数
   * @param onChunk  chunks回调函数
   * @returns 完整响应
   */
  public async createChatCompletionStream(
    request: ChatRequest,
    onChunk: (chunk: string) => void
  ): Promise<string> {
    try {
      const startTime = Date.now();
      logger.info('Creating streaming chat completion', {
        request: {
          ...request,
          messages: request.messages.map(msg => ({
            role: msg.role,
            content: msg.role === 'user' ? msg.content.substring(0, 100) + '...' : msg.content
          }))
        }
      });

      const stream = await openai.chat.completions.create({
        messages: request.messages,
        model: request.model || 'deepseek-chat',
        temperature: request.temperature || 0.7,
        stream: true
      });
      console.log("ai--response--stream",stream)

      let fullResponse = '';

      for await (const chunk of stream) {
        const content = chunk.choices[0]?.delta?.content;
        if (content) {
          fullResponse += content;
          onChunk(content);
        }
      }

      const endTime = Date.now();
      logger.info('Streaming chat completion finished', {
        response: {
          content: fullResponse.substring(0, 100) + '...'
        },
        duration: endTime - startTime
      });

      return fullResponse;
    } catch (error) {
      logger.error('Error creating streaming chat completion', {
        error: error instanceof Error ? error.message : String(error),
        stack: error instanceof Error ? error.stack : undefined
      });
      throw error;
    }
  }
}

// 导出单例实例
export const openaiService = new OpenAIService();
export { openai };
