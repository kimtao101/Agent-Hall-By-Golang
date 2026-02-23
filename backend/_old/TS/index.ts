// 首先加载环境变量
import './env';

import express from 'express';
import cors from 'cors';
import rateLimit from 'express-rate-limit';
import { Agent } from './agent';
import { xiaohongshuService } from './xiaohongshuService';
import logger from './logger';

const app = express();
const PORT = process.env.PORT || 8015;


// 确保logs目录存在
import fs from 'fs';
import path from 'path';

const logsDir = path.join(__dirname, '..', 'logs');
if (!fs.existsSync(logsDir)) {
  fs.mkdirSync(logsDir, { recursive: true });
}

// 配置请求限流
const limiter = rateLimit({
  windowMs: 15 * 60 * 1000, // 15分钟
  max: 100, // 每个IP限制100个请求
  message: {
    error: 'Too many requests from this IP, please try again later.'
  },
  standardHeaders: true,
  legacyHeaders: false,
  skip: (req) => {
    // 跳过健康检查等特定路径
    return req.path === '/health';
  }
});

// Middleware
app.use(cors());
app.use(express.json());
app.use(limiter);

// Create agent instance
const agent = new Agent();

// 健康检查端点
app.get('/health', (req, res) => {
  res.json({ status: 'ok', timestamp: new Date().toISOString() });
});

// API endpoints
/**
 * @swagger
 * /chat: 
 *   post: 
 *     summary: 发送聊天消息
 *     description: 发送消息给AI助手并获取响应
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             properties:
 *               message: 
 *                 type: string
 *                 description: 用户消息
 *     responses:
 *       200: 
 *         description: 成功获取AI响应
 *       400: 
 *         description: 请求参数错误
 *       500: 
 *         description: 服务器内部错误
 */
app.post('/chat', async (req, res) => {
  try {
    const { message } = req.body;
    
    if (!message || typeof message !== 'string') {
      logger.warn('Invalid message parameter', { body: req.body });
      res.status(400).json({ error: 'Invalid message' });
      return;
    }

    logger.info('Received chat request', { message: message.substring(0, 100) + '...' });

    // Set headers for streaming response
    res.setHeader('Content-Type', 'text/plain');
    res.setHeader('Transfer-Encoding', 'chunked');

    // Generate response with streaming
    await agent.generateResponse(message, (chunk) => {
      res.write(chunk);
    });

    logger.info('Chat response completed');
    res.end();
  } catch (error) {
    logger.error('Error in /chat endpoint:', {
      error: error instanceof Error ? error.message : String(error),
      stack: error instanceof Error ? error.stack : undefined
    });
    res.status(500).json({ error: 'Internal server error' });
  }
});

/**
 * @swagger
 * /xiaohongshu/copy: 
 *   post: 
 *     summary: 生成小红书文案
 *     description: 根据场景和配置生成小红书风格的文案
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             properties:
 *               scene: 
 *                 type: string
 *                 description: 文案场景
 *               config: 
 *                 type: object
 *                 description: 场景配置参数
 *     responses:
 *       200: 
 *         description: 成功生成文案
 *       400: 
 *         description: 请求参数错误
 *       500: 
 *         description: 服务器内部错误
 */
app.post('/xiaohongshu/copy', async (req, res) => {
  try {
    const { scene, config } = req.body;
    
    if (!scene || typeof scene !== 'string') {
      logger.warn('Invalid scene parameter', { body: req.body });
      res.status(400).json({ error: 'Invalid scene' });
      return;
    }

    if (!config || typeof config !== 'object') {
      logger.warn('Invalid config parameter', { body: req.body });
      res.status(400).json({ error: 'Invalid config' });
      return;
    }

    logger.info('Received Xiaohongshu copy request', { scene });

    const result = await xiaohongshuService.generateCopy({ scene, config });
    logger.info('Generated Xiaohongshu copy successfully');

    res.json(result);
  } catch (error) {
    logger.error('Error in /xiaohongshu/copy endpoint:', {
      error: error instanceof Error ? error.message : String(error),
      stack: error instanceof Error ? error.stack : undefined
    });
    res.status(500).json({ error: 'Internal server error' });
  }
});

/**
 * @swagger
 * /history:
 *   get:
 *     summary: 获取聊天历史
 *     description: 获取当前对话的历史记录
 *     responses:
 *       200:
 *         description: 成功获取聊天历史
 *       500:
 *         description: 服务器内部错误
 */
app.get('/history', (req, res) => {
  try {
    const history = agent.getHistory();
    logger.info('Retrieved chat history', { length: history.length });
    res.json(history);
  } catch (error) {
    logger.error('Error in /history endpoint:', {
      error: error instanceof Error ? error.message : String(error)
    });
    res.status(500).json({ error: 'Internal server error' });
  }
});

/**
 * @swagger
 * /clear:
 *   post:
 *     summary: 清除聊天历史
 *     description: 清除当前对话的历史记录
 *     responses:
 *       200:
 *         description: 成功清除聊天历史
 *       500:
 *         description: 服务器内部错误
 */
app.post('/clear', (req, res) => {
  try {
    agent.clearHistory();
    logger.info('Cleared chat history');
    res.json({ success: true });
  } catch (error) {
    logger.error('Error in /clear endpoint:', {
      error: error instanceof Error ? error.message : String(error)
    });
    res.status(500).json({ error: 'Internal server error' });
  }
});

// 错误处理中间件
app.use((err: any, req: express.Request, res: express.Response, next: express.NextFunction) => {
  logger.error('Unhandled error:', {
    error: err.message,
    stack: err.stack,
    path: req.path
  });
  
  res.status(500).json({
    error: 'Internal server error',
    timestamp: new Date().toISOString()
  });
});

// Start server
app.listen(PORT, () => {
  logger.info(`Server running on port ${PORT}`);
  console.log(`Server running on port ${PORT}`);
});
