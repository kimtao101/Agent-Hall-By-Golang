import React, { useState, useEffect } from 'react';

interface CopyEditorProps {
  initialContent: string;
  onContentChange: (content: string) => void;
  loading: boolean;
}

const recommendedTags: Record<string, string[]> = {
  beauty: ['美妆', '护肤', '测评', '种草', '好物推荐', '护肤品', '化妆品', '护肤心得', '美妆分享', '护肤日常'],
  fashion: ['穿搭', '时尚', '搭配', 'OOTD', '穿搭分享', '时尚搭配', '日常穿搭', '服装推荐', '风格穿搭', '时尚好物'],
  travel: ['旅行', '旅游', '打卡', '攻略', '旅行日记', '旅游攻略', '景点推荐', '旅行体验', '旅游打卡', '出行攻略'],
  food: ['美食', '探店', '吃播', '美食推荐', '餐厅评测', '美食打卡', '吃货日常', '美食分享', '餐厅推荐', '美食攻略'],
  home: ['家居', '收纳', '装修', '家居好物', '生活技巧', '家居收纳', '家居装修', '生活好物', '家居布置', '收纳技巧'],
  fitness: ['健身', '运动', '减肥', '健身日记', '运动打卡', '健身心得', '减肥日记', '运动分享', '健身计划', '运动日常'],
  parenting: ['育儿', '母婴', '宝宝', '育儿经验', '母婴好物', '育儿心得', '宝宝日常', '母婴分享', '育儿知识', '母婴推荐'],
  tech: ['数码', '科技', '测评', '数码产品', '科技产品', '产品评测', '数码好物', '科技测评', '数码分享', '科技好物']
};

export default function CopyEditor({ initialContent, onContentChange, loading }: CopyEditorProps) {
  const [content, setContent] = useState(initialContent);
  const [selectedTags, setSelectedTags] = useState<string[]>([]);
  const [customTag, setCustomTag] = useState('');

  // 监听 initialContent 变化，当后端返回新文案时更新
  useEffect(() => {
    setContent(initialContent);
  }, [initialContent]);

  const handleContentChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const newContent = e.target.value;
    setContent(newContent);
    onContentChange(newContent);
  };

  const addTag = (tag: string) => {
    if (!selectedTags.includes(tag)) {
      const newTags = [...selectedTags, tag];
      setSelectedTags(newTags);
      updateContentWithTags(newTags);
    }
  };

  const removeTag = (tag: string) => {
    const newTags = selectedTags.filter(t => t !== tag);
    setSelectedTags(newTags);
    updateContentWithTags(newTags);
  };

  const addCustomTag = () => {
    if (customTag.trim() && !selectedTags.includes(customTag.trim())) {
      const newTags = [...selectedTags, customTag.trim()];
      setSelectedTags(newTags);
      setCustomTag('');
      updateContentWithTags(newTags);
    }
  };

  const updateContentWithTags = (tags: string[]) => {
    if (tags.length === 0) {
      onContentChange(content);
      return;
    }

    const tagsString = tags.map(tag => `#${tag}`).join(' ');
    const contentWithoutTags = content.replace(/\s*#\w+\s*/g, '').trim();
    const newContent = `${contentWithoutTags}\n\n${tagsString}`;
    setContent(newContent);
    onContentChange(newContent);
  };

  const formatContent = () => {
    // 简单的格式化：添加适当的换行和表情符号
    let formatted = content
      .replace(/\n\n+/g, '\n\n')
      .replace(/(。|！|？)([^\n])/g, '$1\n$2');

    // 在段落开头添加表情符号
    const emojis = ['✨', '🌟', '💖', '🎉', '🔥', '💯', '⭐', '💫', '🌈', '🌸'];
    const paragraphs = formatted.split('\n\n');
    const formattedParagraphs = paragraphs.map((para, index) => {
      if (para.trim() && !para.startsWith('#') && !para.match(/^\s*[✨🌟💖🎉🔥💯⭐💫🌈🌸]/)) {
        return `${emojis[index % emojis.length]} ${para}`;
      }
      return para;
    });

    const newContent = formattedParagraphs.join('\n\n');
    setContent(newContent);
    onContentChange(newContent);
  };

  return (
    <div style={{ marginBottom: '2rem' }}>
      <h2 style={{ fontSize: '1.5rem', fontWeight: 600, marginBottom: '1rem' }}>
        文案编辑
      </h2>
      
      <div style={{ marginBottom: '1.5rem' }}>
        <textarea
          value={content}
          onChange={handleContentChange}
          placeholder="编辑你的文案..."
          disabled={loading}
          style={{
            width: '100%',
            minHeight: '300px',
            padding: '1rem',
            border: '1px solid #334155',
            borderRadius: '0.5rem',
            backgroundColor: 'var(--secondary)',
            color: 'var(--foreground)',
            resize: 'vertical',
            fontFamily: 'inherit',
            fontSize: '1rem',
            lineHeight: '1.5'
          }}
        />
      </div>

      <div style={{ marginBottom: '1.5rem' }}>
        <h3 style={{ fontSize: '1.1rem', fontWeight: 600, marginBottom: '0.5rem' }}>
          推荐标签
        </h3>
        <div style={{ display: 'flex', flexWrap: 'wrap', gap: '0.5rem', marginBottom: '1rem' }}>
          {selectedTags.map((tag) => (
            <span
              key={tag}
              style={{
                padding: '0.25rem 0.75rem',
                backgroundColor: '#60a5fa',
                color: 'white',
                borderRadius: '1rem',
                fontSize: '0.9rem',
                display: 'flex',
                alignItems: 'center',
                gap: '0.5rem'
              }}
            >
              #{tag}
              <button
                onClick={() => removeTag(tag)}
                style={{
                  background: 'none',
                  border: 'none',
                  color: 'white',
                  cursor: 'pointer',
                  fontSize: '1rem',
                  padding: 0,
                  width: '16px',
                  height: '16px',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center'
                }}
              >
                ×
              </button>
            </span>
          ))}
        </div>
        
        <div style={{ display: 'flex', gap: '0.5rem', marginBottom: '1rem' }}>
          <input
            type="text"
            value={customTag}
            onChange={(e) => setCustomTag(e.target.value)}
            placeholder="添加自定义标签"
            style={{
              flex: 1,
              padding: '0.5rem',
              border: '1px solid #334155',
              borderRadius: '0.5rem',
              backgroundColor: 'var(--secondary)',
              color: 'var(--foreground)'
            }}
          />
          <button
            onClick={addCustomTag}
            style={{
              padding: '0.5rem 1rem',
              backgroundColor: '#60a5fa',
              color: 'white',
              border: 'none',
              borderRadius: '0.5rem',
              cursor: 'pointer'
            }}
          >
            添加
          </button>
        </div>
      </div>

      <div style={{ display: 'flex', gap: '1rem', marginBottom: '1.5rem' }}>
        <button
          onClick={formatContent}
          disabled={loading}
          style={{
            padding: '0.75rem 1.5rem',
            backgroundColor: '#60a5fa',
            color: 'white',
            border: 'none',
            borderRadius: '0.5rem',
            cursor: loading ? 'not-allowed' : 'pointer',
            opacity: loading ? 0.5 : 1
          }}
        >
          格式化文案
        </button>
        
        <button
          onClick={() => onContentChange('')}
          disabled={loading}
          style={{
            padding: '0.75rem 1.5rem',
            backgroundColor: '#ef4444',
            color: 'white',
            border: 'none',
            borderRadius: '0.5rem',
            cursor: loading ? 'not-allowed' : 'pointer',
            opacity: loading ? 0.5 : 1
          }}
        >
          清空
        </button>
      </div>

      {loading && (
        <div style={{ color: '#94a3b8', fontStyle: 'italic' }}>
          生成文案中...
        </div>
      )}
    </div>
  );
}
