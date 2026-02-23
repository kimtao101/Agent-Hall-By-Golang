'use client';
import { useState, useCallback, useEffect, useRef } from 'react';
import { useRouter } from 'next/navigation';

interface AgentItem {
  id: string;
  name: string;
  description: string;
  icon: string;
  route?: string;
}

export default function AgentHall() {
  const router = useRouter();
  
  // 防抖定时器
  const debounceTimerRef = useRef<NodeJS.Timeout | null>(null);

  // Agent数据，包含小红书agent作为第一个元素
  const agents: AgentItem[] = [
    {
      id: 'xiaohongshu',
      name: '小红书Agent',
      description: '智能生成小红书平台优质文案',
      icon: '📝',
      route: '/xiaohongshu'
    },
    {
      id: 'marketing',
      name: '营销Agent',
      description: '专业营销文案与策略生成',
      icon: '📈'
    },
    {
      id: 'research',
      name: '研究Agent',
      description: '深度调研与分析报告生成',
      icon: '🔍'
    },
    {
      id: 'travel',
      name: '旅行Agent',
      description: '旅行规划与攻略生成',
      icon: '✈️'
    },
    {
      id: 'hr',
      name: 'HR Agent',
      description: '人力资源管理与招聘',
      icon: '👥'
    },
    {
      id: 'media',
      name: '媒体Agent',
      description: '媒体内容与新闻稿生成',
      icon: '📺'
    },
    {
      id: 'ecommerce',
      name: '电商Agent',
      description: '电商运营与产品描述',
      icon: '🛒'
    },
    {
      id: 'social',
      name: '社交Agent',
      description: '社交媒体内容与互动',
      icon: '🌐'
    },
    {
      id: 'language',
      name: '语言Agent',
      description: '多语言翻译与内容优化',
      icon: '🌍'
    },
    {
      id: 'product',
      name: '产品Agent',
      description: '产品设计与用户研究',
      icon: '📱'
    }
  ];

  // 搜索相关状态
  const [searchQuery, setSearchQuery] = useState('');
  const [searchResults, setSearchResults] = useState<AgentItem[]>([]);
  const [isSearching, setIsSearching] = useState(false);
  const [searchCache, setSearchCache] = useState<Record<string, AgentItem[]>>({});
  const [displayedAgents, setDisplayedAgents] = useState<AgentItem[]>(agents);

  const handleAgentClick = (agent: AgentItem) => {
    if (agent.route) {
      // 直接导航，依赖布局文件的动画效果
      router.push(agent.route);
    } else {
      // 对于未实现的Agent，显示提示
      alert(`${agent.name} 正在开发中，敬请期待！`);
    }
  };

  // 搜索处理函数
  const handleSearch = useCallback((query: string) => {
    setSearchQuery(query);
    
    // 清除之前的防抖定时器
    if (debounceTimerRef.current) {
      clearTimeout(debounceTimerRef.current);
    }
    
    if (!query.trim()) {
      setSearchResults([]);
      setDisplayedAgents(agents);
      return;
    }
    
    // 检查缓存
    if (searchCache[query]) {
      setSearchResults(searchCache[query]);
      setDisplayedAgents(searchCache[query]);
      return;
    }
    
    // 设置防抖定时器
    debounceTimerRef.current = setTimeout(() => {
      setIsSearching(true);
      
      // 模拟搜索延迟
      setTimeout(() => {
        // 执行搜索逻辑
        const results = agents.filter(agent => {
          const searchTerm = query.toLowerCase();
          return (
            agent.name.toLowerCase().includes(searchTerm) ||
            agent.description.toLowerCase().includes(searchTerm) ||
            agent.id.toLowerCase().includes(searchTerm)
          );
        }).slice(0, 10); // 最多显示10条结果
        
        // 更新缓存
        setSearchCache(prev => ({
          ...prev,
          [query]: results
        }));
        
        setSearchResults(results);
        setDisplayedAgents(results);
        setIsSearching(false);
      }, 200); // 模拟搜索响应时间
    }, 300); // 防抖时间300ms
  }, [agents, searchCache]);

  // 处理搜索提交
  const handleSearchSubmit = useCallback((e: React.FormEvent) => {
    e.preventDefault();
    // 搜索提交时，搜索结果已经通过handleSearch更新到displayedAgents
  }, []);

  // 重置搜索结果
  useEffect(() => {
    // 键盘事件处理：按ESC键清除搜索
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        setSearchQuery('');
        setDisplayedAgents(agents);
      }
    };

    document.addEventListener('keydown', handleKeyDown);

    return () => {
      document.removeEventListener('keydown', handleKeyDown);
    };
  }, []);

  return (
    <main className="container" style={{ alignItems: 'center' }}>
      {/* 页面标题 */}
      <header style={{ 
        marginBottom: '3rem', 
        textAlign: 'center',
        width: '100%'
      }}>
        <h1 style={{ 
          fontSize: '3rem', 
          fontWeight: 700, 
          backgroundImage: 'linear-gradient(to right, #60a5fa, #a78bfa)', 
          WebkitBackgroundClip: 'text', 
          color: 'transparent',
          marginBottom: '1rem'
        }}>
          Agent大厅
        </h1>
        <p style={{ 
          color: '#94a3b8',
          fontSize: '1.1rem'
        }}>
          探索AI驱动的专业助手，赋能各类业务场景
        </p>
      </header>

      {/* 搜索组件 */}
      <div style={{
        width: '100%',
        maxWidth: '800px',
        marginBottom: '3rem',
        position: 'relative'
      }}>
        <div style={{
          display: 'flex',
          width: '100%',
          position: 'relative'
        }}>
          <input
            type="text"
            placeholder="搜索智能体..."
            value={searchQuery}
            onChange={(e) => handleSearch(e.target.value)}
            style={{
              flex: 1,
              padding: '1rem 4rem 1rem 1.5rem',
              fontSize: '1rem',
              border: '1px solid #334155',
              borderRadius: '0.5rem 0 0 0.5rem',
              backgroundColor: 'var(--secondary)',
              color: 'var(--foreground)',
              outline: 'none',
              transition: 'all 0.2s ease'
            }}
          />
          
          {/* 删除按钮 */}
          {searchQuery.length > 0 && (
            <button
              type="button"
              onClick={() => {
                setSearchQuery('');
                setDisplayedAgents(agents);
              }}
              style={{
                position: 'absolute',
                right: '8rem',
                top: '50%',
                transform: 'translateY(-50%)',
                backgroundColor: 'transparent',
                color: '#94a3b8',
                border: 'none',
                borderRadius: '50%',
                width: '2rem',
                height: '2rem',
                cursor: 'pointer',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                fontSize: '1.2rem',
                transition: 'all 0.2s ease',
                zIndex: 10
              }}
              onMouseEnter={(e) => {
                e.currentTarget.style.color = '#60a5fa';
              }}
              onMouseLeave={(e) => {
                e.currentTarget.style.color = '#94a3b8';
              }}
            >
              ×
            </button>
          )}
          
          <button
            type="submit"
            onClick={handleSearchSubmit}
            style={{
              padding: '1rem 2rem',
              fontSize: '1rem',
              backgroundColor: '#60a5fa',
              color: 'white',
              border: 'none',
              borderRadius: '0 0.5rem 0.5rem 0',
              cursor: 'pointer',
              transition: 'background 0.2s ease',
              minWidth: '8rem'
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.backgroundColor = '#3b82f6';
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.backgroundColor = '#60a5fa';
            }}
          >
            {isSearching ? '搜索中...' : '搜索'}
          </button>
        </div>
        
        {/* 搜索结果提示 */}
        {searchQuery.trim() && (
          <div style={{
            marginTop: '0.75rem',
            textAlign: 'center',
            fontSize: '0.875rem',
            color: '#94a3b8'
          }}>
            {isSearching 
              ? '正在搜索智能体...' 
              : searchResults.length > 0 
                ? `找到 ${searchResults.length} 个相关智能体` 
                : '未找到相关智能体'}
          </div>
        )}
      </div>

      {/* Agent矩阵布局 */}
      <div style={{
        display: 'grid',
        gridTemplateColumns: 'repeat(auto-fill, minmax(250px, 1fr))',
        gridTemplateRows: 'repeat(6, auto)',
        gap: '1.5rem',
        maxWidth: '1400px',
        width: '100%',
        marginBottom: '3rem'
      }}>
        {displayedAgents.map((agent) => (
          <div
            key={agent.id}
            onClick={() => handleAgentClick(agent)}
            style={{
              background: 'var(--secondary)',
              border: agent.id === 'xiaohongshu' ? '2px solid #60a5fa' : '1px solid #334155',
              borderRadius: '1rem',
              padding: '1.5rem',
              textAlign: 'center',
              cursor: 'pointer',
              transition: 'all 0.3s ease',
              boxShadow: agent.id === 'xiaohongshu' ? '0 0 20px rgba(96, 165, 250, 0.3)' : 'none'
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.transform = 'translateY(-5px)';
              e.currentTarget.style.boxShadow = agent.id === 'xiaohongshu' 
                ? '0 10px 30px rgba(96, 165, 250, 0.4)' 
                : '0 5px 15px rgba(0, 0, 0, 0.2)';
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.transform = 'translateY(0)';
              e.currentTarget.style.boxShadow = agent.id === 'xiaohongshu' 
                ? '0 0 20px rgba(96, 165, 250, 0.3)' 
                : 'none';
            }}
          >
            <div style={{ fontSize: '3rem', marginBottom: '1rem' }}>
              {agent.icon}
            </div>
            <h3 style={{ 
              fontSize: '1.2rem', 
              fontWeight: 600, 
              marginBottom: '0.5rem',
              color: agent.id === 'xiaohongshu' ? '#60a5fa' : 'var(--foreground)'
            }}>
              {agent.name}
            </h3>
            <p style={{ 
              color: '#94a3b8', 
              fontSize: '0.9rem',
              lineHeight: '1.4'
            }}>
              {agent.description}
            </p>
            {agent.route && (
              <div style={{ 
                marginTop: '1rem',
                fontSize: '0.8rem',
                color: '#60a5fa'
              }}>
                点击进入 →
              </div>
            )}
          </div>
        ))}
      </div>

      {/* 页脚 */}
      <footer style={{ 
        marginTop: '3rem', 
        textAlign: 'center', 
        color: '#64748b',
        fontSize: '0.9rem',
        width: '100%'
      }}>
        <p>© 2026 Agent系统 | 智能赋能未来</p>
      </footer>
    </main>
  );
}