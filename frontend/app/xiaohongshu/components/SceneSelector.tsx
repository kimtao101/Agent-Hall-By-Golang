import React from 'react';

interface Scene {
  id: string;
  name: string;
  icon: string;
  description: string;
}

interface SceneSelectorProps {
  selectedScene: string | null;
  onSelectScene: (sceneId: string) => void;
}

const scenes: Scene[] = [
  {
    id: 'beauty',
    name: '美妆护肤评测',
    icon: '💄',
    description: '产品评测、使用感受、效果分享'
  },
  {
    id: 'fashion',
    name: '穿搭搭配分享',
    icon: '👗',
    description: '服装搭配、风格推荐、购物指南'
  },
  {
    id: 'travel',
    name: '旅行打卡攻略',
    icon: '✈️',
    description: '景点推荐、行程规划、旅行体验'
  },
  {
    id: 'food',
    name: '美食探店体验',
    icon: '🍔',
    description: '餐厅评测、美食推荐、用餐体验'
  },
  {
    id: 'home',
    name: '家居好物推荐',
    icon: '🏠',
    description: '家居用品、装修灵感、生活技巧'
  },
  {
    id: 'fitness',
    name: '健身运动记录',
    icon: '🏋️',
    description: '运动计划、健身心得、成果分享'
  },
  {
    id: 'parenting',
    name: '母婴育儿心得',
    icon: '👶',
    description: '育儿经验、产品推荐、成长记录'
  },
  {
    id: 'tech',
    name: '数码产品测评',
    icon: '📱',
    description: '产品评测、使用体验、技术分析'
  }
];

export default function SceneSelector({ selectedScene, onSelectScene }: SceneSelectorProps) {
  return (
    <div style={{ marginBottom: '2rem' }}>
      <h2 style={{ fontSize: '1.5rem', fontWeight: 600, marginBottom: '1rem' }}>
        选择文案场景
      </h2>
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))', gap: '1rem' }}>
        {scenes.map((scene) => (
          <div
            key={scene.id}
            onClick={() => onSelectScene(scene.id)}
            style={{
              padding: '1.5rem',
              border: `2px solid ${selectedScene === scene.id ? '#60a5fa' : '#334155'}`,
              borderRadius: '1rem',
              cursor: 'pointer',
              transition: 'all 0.2s ease',
              backgroundColor: selectedScene === scene.id ? 'rgba(96, 165, 250, 0.1)' : 'var(--secondary)',
              textAlign: 'center'
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.transform = 'translateY(-2px)';
              e.currentTarget.style.boxShadow = '0 4px 6px -1px rgba(0, 0, 0, 0.1)';
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.transform = 'translateY(0)';
              e.currentTarget.style.boxShadow = 'none';
            }}
          >
            <div style={{ fontSize: '2.5rem', marginBottom: '0.5rem' }}>
              {scene.icon}
            </div>
            <h3 style={{ fontSize: '1.1rem', fontWeight: 600, marginBottom: '0.5rem' }}>
              {scene.name}
            </h3>
            <p style={{ fontSize: '0.9rem', color: '#94a3b8', lineHeight: '1.4' }}>
              {scene.description}
            </p>
          </div>
        ))}
      </div>
    </div>
  );
}
