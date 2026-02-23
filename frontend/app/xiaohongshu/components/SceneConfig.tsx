import React, { useState, useEffect } from 'react';

interface SceneConfigProps {
  sceneId: string | null;
  onConfigChange: (config: any) => void;
}

interface FormField {
  name: string;
  label: string;
  type: 'text' | 'textarea' | 'number';
  placeholder: string;
  required: boolean;
}

const sceneFields: Record<string, FormField[]> = {
  beauty: [
    { name: 'productName', label: '产品名称', type: 'text', placeholder: '例如：XX品牌精华液', required: true },
    { name: 'brand', label: '品牌', type: 'text', placeholder: '例如：XX护肤', required: true },
    { name: 'price', label: '价格', type: 'number', placeholder: '例如：299', required: false },
    { name: 'skinType', label: '适合肤质', type: 'text', placeholder: '例如：干性/油性/混合性', required: false },
    { name: 'texture', label: '质地', type: 'text', placeholder: '例如：轻薄/滋润/清爽', required: false },
    { name: 'keyIngredients', label: '核心成分', type: 'text', placeholder: '例如：玻尿酸、烟酰胺', required: false },
    { name: 'usageFeel', label: '使用感受', type: 'textarea', placeholder: '描述使用时的感受和效果', required: true },
    { name: 'effect', label: '效果', type: 'textarea', placeholder: '描述使用后的效果', required: true },
    { name: 'recommendation', label: '推荐理由', type: 'textarea', placeholder: '为什么推荐这款产品', required: true }
  ],
  fashion: [
    { name: 'clothingType', label: '服装类型', type: 'text', placeholder: '例如：连衣裙、外套', required: true },
    { name: 'style', label: '风格', type: 'text', placeholder: '例如：休闲、正式、复古', required: true },
    { name: 'brand', label: '品牌', type: 'text', placeholder: '例如：XX服饰', required: false },
    { name: 'price', label: '价格', type: 'number', placeholder: '例如：399', required: false },
    { name: 'color', label: '颜色', type: 'text', placeholder: '例如：黑色、白色、粉色', required: false },
    { name: 'material', label: '材质', type: 'text', placeholder: '例如：棉、羊毛、丝绸', required: false },
    { name: 'fit', label: '版型', type: 'text', placeholder: '例如：宽松、修身', required: false },
    { name: 'matchingTips', label: '搭配建议', type: 'textarea', placeholder: '如何搭配其他服饰', required: true },
    { name: 'scenario', label: '适合场景', type: 'text', placeholder: '例如：日常、约会、聚会', required: true },
    { name: 'usageFeel', label: '穿着感受', type: 'textarea', placeholder: '描述穿着时的舒适度', required: true }
  ],
  travel: [
    { name: 'destination', label: '目的地', type: 'text', placeholder: '例如：东京、三亚', required: true },
    { name: 'duration', label: '行程天数', type: 'number', placeholder: '例如：5', required: true },
    { name: 'bestTime', label: '最佳时间', type: 'text', placeholder: '例如：春季、夏季', required: false },
    { name: 'budget', label: '预算', type: 'number', placeholder: '例如：5000', required: false },
    { name: 'attractions', label: '主要景点', type: 'textarea', placeholder: '列出主要景点', required: true },
    { name: 'food', label: '特色美食', type: 'textarea', placeholder: '推荐当地美食', required: true },
    { name: 'accommodation', label: '住宿推荐', type: 'textarea', placeholder: '推荐住宿地点', required: false },
    { name: 'transportation', label: '交通方式', type: 'text', placeholder: '例如：地铁、打车', required: false },
    { name: 'tips', label: '旅行贴士', type: 'textarea', placeholder: '实用建议和注意事项', required: true },
    { name: 'experience', label: '个人体验', type: 'textarea', placeholder: '分享你的旅行感受', required: true }
  ],
  food: [
    { name: 'restaurantName', label: '餐厅名称', type: 'text', placeholder: '例如：XX餐厅', required: true },
    { name: 'location', label: '位置', type: 'text', placeholder: '例如：XX商圈、XX街道', required: true },
    { name: 'cuisineType', label: '菜系', type: 'text', placeholder: '例如：川菜、日料', required: true },
    { name: 'priceRange', label: '价格区间', type: 'text', placeholder: '例如：人均100-200', required: false },
    { name: 'environment', label: '环境', type: 'text', placeholder: '例如：温馨、高端、文艺', required: false },
    { name: 'service', label: '服务', type: 'text', placeholder: '例如：热情、专业', required: false },
    { name: 'signatureDishes', label: '招牌菜', type: 'textarea', placeholder: '推荐几道招牌菜', required: true },
    { name: 'taste', label: '口味', type: 'textarea', placeholder: '描述菜品的味道', required: true },
    { name: 'recommendation', label: '推荐理由', type: 'textarea', placeholder: '为什么推荐这家餐厅', required: true }
  ],
  home: [
    { name: 'productName', label: '产品名称', type: 'text', placeholder: '例如：XX收纳盒', required: true },
    { name: 'category', label: '类别', type: 'text', placeholder: '例如：收纳、厨房、卧室', required: true },
    { name: 'brand', label: '品牌', type: 'text', placeholder: '例如：XX家居', required: false },
    { name: 'price', label: '价格', type: 'number', placeholder: '例如：99', required: false },
    { name: 'material', label: '材质', type: 'text', placeholder: '例如：塑料、木质、金属', required: false },
    { name: 'size', label: '尺寸', type: 'text', placeholder: '例如：30x20x15cm', required: false },
    { name: 'usageScenario', label: '使用场景', type: 'text', placeholder: '例如：客厅、厨房、卧室', required: true },
    { name: 'functionality', label: '功能', type: 'textarea', placeholder: '描述产品的功能', required: true },
    { name: 'usageFeel', label: '使用感受', type: 'textarea', placeholder: '描述使用时的体验', required: true },
    { name: 'spaceSaving', label: '节省空间', type: 'text', placeholder: '例如：是/否', required: false },
    { name: 'recommendation', label: '推荐理由', type: 'textarea', placeholder: '为什么推荐这款产品', required: true }
  ],
  fitness: [
    { name: 'workoutType', label: '运动类型', type: 'text', placeholder: '例如：瑜伽、跑步、力量训练', required: true },
    { name: 'duration', label: '运动时长', type: 'number', placeholder: '例如：30', required: false },
    { name: 'frequency', label: '频率', type: 'text', placeholder: '例如：每周3次', required: false },
    { name: 'equipment', label: '所需装备', type: 'text', placeholder: '例如：瑜伽垫、哑铃', required: false },
    { name: 'difficulty', label: '难度', type: 'text', placeholder: '例如：初级、中级、高级', required: false },
    { name: 'benefits', label: '好处', type: 'textarea', placeholder: '描述这项运动的好处', required: true },
    { name: 'experience', label: '个人体验', type: 'textarea', placeholder: '分享你的运动感受', required: true },
    { name: 'tips', label: '注意事项', type: 'textarea', placeholder: '运动时的注意事项', required: true },
    { name: 'results', label: '效果', type: 'textarea', placeholder: '运动后的效果和变化', required: false }
  ],
  parenting: [
    { name: 'babyAge', label: '宝宝年龄', type: 'text', placeholder: '例如：6个月、2岁', required: true },
    { name: 'topic', label: '主题', type: 'text', placeholder: '例如：辅食、睡眠、教育', required: true },
    { name: 'productName', label: '产品名称（如有）', type: 'text', placeholder: '例如：XX品牌奶瓶', required: false },
    { name: 'brand', label: '品牌（如有）', type: 'text', placeholder: '例如：XX母婴', required: false },
    { name: 'price', label: '价格（如有）', type: 'number', placeholder: '例如：199', required: false },
    { name: 'problem', label: '问题描述', type: 'textarea', placeholder: '描述遇到的问题', required: true },
    { name: 'solution', label: '解决方案', type: 'textarea', placeholder: '分享你的解决方法', required: true },
    { name: 'experience', label: '经验分享', type: 'textarea', placeholder: '分享你的育儿经验', required: true },
    { name: 'tips', label: '小贴士', type: 'textarea', placeholder: '给其他父母的建议', required: true }
  ],
  tech: [
    { name: 'productName', label: '产品名称', type: 'text', placeholder: '例如：XX手机、XX耳机', required: true },
    { name: 'brand', label: '品牌', type: 'text', placeholder: '例如：XX科技', required: true },
    { name: 'price', label: '价格', type: 'number', placeholder: '例如：5999', required: false },
    { name: 'releaseDate', label: '上市时间', type: 'text', placeholder: '例如：2024年1月', required: false },
    { name: 'specs', label: '主要配置', type: 'textarea', placeholder: '描述产品的主要配置', required: true },
    { name: 'design', label: '设计', type: 'textarea', placeholder: '描述产品的外观设计', required: false },
    { name: 'performance', label: '性能', type: 'textarea', placeholder: '描述产品的性能表现', required: true },
    { name: 'battery', label: '续航', type: 'text', placeholder: '例如：10小时', required: false },
    { name: 'camera', label: '相机（如有）', type: 'textarea', placeholder: '描述相机性能', required: false },
    { name: 'userExperience', label: '用户体验', type: 'textarea', placeholder: '描述使用体验', required: true },
    { name: 'pros', label: '优点', type: 'textarea', placeholder: '产品的优点', required: true },
    { name: 'cons', label: '缺点', type: 'textarea', placeholder: '产品的缺点', required: false },
    { name: 'recommendation', label: '推荐理由', type: 'textarea', placeholder: '为什么推荐这款产品', required: true }
  ]
};

export default function SceneConfig({ sceneId, onConfigChange }: SceneConfigProps) {
  const [formData, setFormData] = useState<any>({});

  useEffect(() => {
    if (sceneId) {
      const initialData: any = {};
      sceneFields[sceneId].forEach(field => {
        initialData[field.name] = '';
      });
      setFormData(initialData);
      onConfigChange(initialData);
    }
  }, [sceneId, onConfigChange]);

  const handleChange = (fieldName: string, value: any) => {
    const newData = { ...formData, [fieldName]: value };
    setFormData(newData);
    onConfigChange(newData);
  };

  if (!sceneId) {
    return (
      <div style={{ padding: '2rem', border: '2px dashed #334155', borderRadius: '1rem', textAlign: 'center' }}>
        <p style={{ color: '#94a3b8' }}>请先选择一个文案场景</p>
      </div>
    );
  }

  const fields = sceneFields[sceneId] || [];

  return (
    <div style={{ marginBottom: '2rem' }}>
      <h2 style={{ fontSize: '1.5rem', fontWeight: 600, marginBottom: '1rem' }}>
        场景参数配置
      </h2>
      <div style={{ 
        display: 'grid', 
        gridTemplateColumns: 'repeat(auto-fit, minmax(350px, 1fr))', 
        gap: '2rem',
        boxSizing: 'border-box',
        padding: '0.5rem',
        maxWidth: '100%'
      }}>
        {fields?.map((field) => (
          <div key={field.name} style={{ 
            display: 'flex', 
            flexDirection: 'column', 
            gap: '0.5rem',
            boxSizing: 'border-box'
          }}>
            <label style={{ 
              fontWeight: 500, 
              display: 'flex', 
              alignItems: 'center',
              boxSizing: 'border-box'
            }}>
              {field.label}
              {field.required && <span style={{ color: '#ef4444', marginLeft: '0.25rem' }}>*</span>}
            </label>
            {field.type === 'textarea' ? (
              <textarea
                name={field.name}
                placeholder={field.placeholder}
                value={formData[field.name] || ''}
                onChange={(e) => handleChange(field.name, e.target.value)}
                style={{
                  padding: '0.75rem',
                  border: '1px solid #60a5fa',
                  borderRadius: '0.5rem',
                  backgroundColor: 'var(--secondary)',
                  color: 'var(--foreground)',
                  resize: 'vertical',
                  minHeight: '100px',
                  fontFamily: 'inherit',
                  outline: 'none',
                  transition: 'all 0.2s ease',
                  width: '100%',
                  boxSizing: 'border-box'
                }}
                onFocus={(e) => {
                  e.currentTarget.style.borderColor = '#a78bfa';
                  e.currentTarget.style.boxShadow = '0 0 0 2px rgba(96, 165, 250, 0.2)';
                }}
                onBlur={(e) => {
                  e.currentTarget.style.borderColor = '#60a5fa';
                  e.currentTarget.style.boxShadow = 'none';
                }}
                required={field.required}
              />
            ) : (
              <input
                type={field.type}
                name={field.name}
                placeholder={field.placeholder}
                value={formData[field.name] || ''}
                onChange={(e) => handleChange(field.name, e.target.value)}
                style={{
                  padding: '0.75rem',
                  border: '1px solid #60a5fa',
                  borderRadius: '0.5rem',
                  backgroundColor: 'var(--secondary)',
                  color: 'var(--foreground)',
                  outline: 'none',
                  transition: 'all 0.2s ease',
                  width: '100%',
                  boxSizing: 'border-box'
                }}
                onFocus={(e) => {
                  e.currentTarget.style.borderColor = '#a78bfa';
                  e.currentTarget.style.boxShadow = '0 0 0 2px rgba(96, 165, 250, 0.2)';
                }}
                onBlur={(e) => {
                  e.currentTarget.style.borderColor = '#60a5fa';
                  e.currentTarget.style.boxShadow = 'none';
                }}
                required={field.required}
              />
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
