'use client';
import { useEffect, useState } from 'react';
import './globals.css';

export default function RootLayout({ children }: { children: React.ReactNode }) {
  const [isTransitioning, setIsTransitioning] = useState(false);
  const [transitionDirection, setTransitionDirection] = useState<'in' | 'out'>('in');

  // 处理页面进入动画
  useEffect(() => {
    // 初始加载和路由变化时的进入动画
    setIsTransitioning(true);
    setTransitionDirection('in');

    // 动画结束后重置状态
    const timer = setTimeout(() => {
      setIsTransitioning(false);
    }, 500);

    return () => clearTimeout(timer);
  }, []); // 空依赖数组，只在初始加载时运行

  return (
    <html lang="zh-CN">
      <body>
        <div
          style={{
            opacity: isTransitioning && transitionDirection === 'out' ? 0 : 1,
            transition: 'opacity 0.3s ease-in-out, transform 0.3s ease-in-out',
            transform: isTransitioning && transitionDirection === 'out' 
              ? 'translateY(20px)' 
              : 'translateY(0)',
            minHeight: '100vh'
          }}
        >
          {children}
        </div>
      </body>
    </html>
  );
}