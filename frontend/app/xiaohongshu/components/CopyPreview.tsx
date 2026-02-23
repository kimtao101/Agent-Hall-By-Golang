import React from 'react';

interface CopyPreviewProps {
  content: string;
  sceneId: string | null;
}

const sceneIcons: Record<string, string> = {
  beauty: '💄',
  fashion: '👗',
  travel: '✈️',
  food: '🍔',
  home: '🏠',
  fitness: '🏋️',
  parenting: '👶',
  tech: '📱'
};

const sceneNames: Record<string, string> = {
  beauty: '美妆护肤评测',
  fashion: '穿搭搭配分享',
  travel: '旅行打卡攻略',
  food: '美食探店体验',
  home: '家居好物推荐',
  fitness: '健身运动记录',
  parenting: '母婴育儿心得',
  tech: '数码产品测评'
};

export default function CopyPreview({ content, sceneId }: CopyPreviewProps) {
  const copyToClipboard = async () => {
    try {
      await navigator.clipboard.writeText(content);
      alert('文案已复制到剪贴板！');
    } catch (err) {
      console.error('复制失败:', err);
      alert('复制失败，请手动复制');
    }
  };

  if (!content) {
    return (
      <div style={{ padding: '2rem', border: '2px dashed #334155', borderRadius: '1rem', textAlign: 'center' }}>
        <p style={{ color: '#94a3b8' }}>生成文案后将在此处预览</p>
      </div>
    );
  }

  return (
    <div style={{ marginBottom: '2rem' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem' }}>
        <h2 style={{ fontSize: '1.5rem', fontWeight: 600 }}>
          文案预览
        </h2>
        <button
          onClick={copyToClipboard}
          style={{
            padding: '0.5rem 1rem',
            backgroundColor: '#60a5fa',
            color: 'white',
            border: 'none',
            borderRadius: '0.5rem',
            cursor: 'pointer',
            fontSize: '0.9rem'
          }}
        >
          复制文案
        </button>
      </div>
      
      <div style={{
        padding: '2rem',
        border: '1px solid #334155',
        borderRadius: '1rem',
        backgroundColor: 'rgba(30, 41, 59, 0.5)',
        minHeight: '400px'
      }}>
        {/* 小红书风格预览 */}
        <div style={{
          maxWidth: '600px',
          margin: '0 auto',
          fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif'
        }}>
          {/* 顶部信息 */}
          <div style={{ display: 'flex', alignItems: 'center', marginBottom: '1.5rem', paddingBottom: '1rem', borderBottom: '1px solid #e2e8f0' }}>
            <div style={{
              width: '40px',
              height: '40px',
              borderRadius: '50%',
              backgroundColor: '#60a5fa',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              marginRight: '1rem'
            }}>
              {sceneId && sceneIcons[sceneId]}
            </div>
            <div>
              <div style={{ fontWeight: 600, fontSize: '1rem', marginBottom: '0.25rem' }}>
                小红书用户
              </div>
              <div style={{ fontSize: '0.8rem', color: '#94a3b8' }}>
                {sceneId && sceneNames[sceneId]}
              </div>
            </div>
          </div>
          
          {/* 文案内容 */}
          <div style={{ marginBottom: '1.5rem', lineHeight: '1.8' }}>
            {content.split('\n\n').map((paragraph, index) => (
              <div key={index} style={{ marginBottom: '1rem' }}>
                {paragraph.startsWith('#') ? (
                  <div style={{ fontSize: '0.9rem', color: '#60a5fa', marginTop: '1.5rem' }}>
                    {paragraph}
                  </div>
                ) : (
                  <p style={{ margin: 0 }}>{paragraph}</p>
                )}
              </div>
            ))}
          </div>
          
          {/* 互动栏 */}
          <div style={{ display: 'flex', alignItems: 'center', gap: '2rem', paddingTop: '1rem', borderTop: '1px solid #e2e8f0' }}>
            <button style={{
              display: 'flex',
              alignItems: 'center',
              gap: '0.5rem',
              background: 'none',
              border: 'none',
              color: '#94a3b8',
              cursor: 'pointer',
              padding: '0.5rem 0'
            }}>
              <span>👍</span>
              <span>点赞</span>
            </button>
            <button style={{
              display: 'flex',
              alignItems: 'center',
              gap: '0.5rem',
              background: 'none',
              border: 'none',
              color: '#94a3b8',
              cursor: 'pointer',
              padding: '0.5rem 0'
            }}>
              <span>💬</span>
              <span>评论</span>
            </button>
            <button style={{
              display: 'flex',
              alignItems: 'center',
              gap: '0.5rem',
              background: 'none',
              border: 'none',
              color: '#94a3b8',
              cursor: 'pointer',
              padding: '0.5rem 0'
            }}>
              <span>↗️</span>
              <span>分享</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
